package api

import (
	"database/sql"
	"errors"
	"lost_found/db/sqlc"
	"lost_found/middleware"
	"lost_found/middleware/session"
	"lost_found/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrNoSession = errors.New("no session")
)

type userResponse struct {
	Name      string `json:"name"`
	StudentID string `json:"studentID"`
	Phone     string `json:"phone"`
	AvatarUrl string `json:"avatar_url"`
	Avatar    []byte `json:"avatar"`
}

func newUserResponse(user sqlc.Usr) userResponse {
	return userResponse{
		Name:      user.Name,
		StudentID: user.StudentID,
		Phone:     user.Phone,
		AvatarUrl: user.AvatarUrl,
		Avatar:    user.Avatar,
	}
}

//用户获取自己的信息
func (server *Server) getUserSelf(ctx *gin.Context) {
	userSession, ok := ctx.Get(middleware.SessionHeaderKey)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrNoSession))
		return
	}
	openid := userSession.(session.Session).ID.String()

	user, err := server.store.GetUsr(ctx, openid)
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	response := newUserResponse(user)
	ctx.JSON(http.StatusOK, response)
}

//用户永久注销自己账户
func (server *Server) deleteUserSelf(ctx *gin.Context) {
	userSession, ok := ctx.Get(middleware.SessionHeaderKey)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrNoSession))
		return
	}
	openid := userSession.(session.Session).ID.String()

	err := server.store.DeleteUsr(ctx, openid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	ctx.JSON(http.StatusOK, nil)
}

//用户更新自己账户信息
type updateUserSelfRequest struct {
	Name      string `json:"name" binding:"required,min=4,max=24"`
	AvatarUrl string `json:"avatar_url" binding:"required,url"`
	Avatar    []byte `json:"avatar" binding:"required"`
}

func (server *Server) updateUserSelf(ctx *gin.Context) {
	userSession, ok := ctx.Get(middleware.SessionHeaderKey)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, errorResponse(ErrNoSession))
		return
	}
	openid := userSession.(session.Session).ID.String()

	var request updateUserSelfRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUsr(ctx, openid)
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	if user.Name != request.Name {
		param := sqlc.UpdateUsrNameParams{
			Openid: openid,
			Name:   request.Name,
		}
		user, err = server.store.UpdateUsrName(ctx, param)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	if request.Avatar != nil {
		param := sqlc.UpdateUsrAvatarParams{
			Openid:    openid,
			AvatarUrl: request.AvatarUrl,
			Avatar:    request.Avatar,
		}
		user, err = server.store.UpdateUsrAvatar(ctx, param)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	response := newUserResponse(user)
	ctx.JSON(http.StatusOK, response)
}

//注册添加微信用户
type addUserRequest struct {
	Code      string `json:"code" binding:"required,alphanum"`
	PhoneCode string `json:"phone_code" binding:"required,alphanum"`
}

func (server *Server) addUser(ctx *gin.Context) {
	var addUserRequest addUserRequest
	if err := ctx.ShouldBindJSON(&addUserRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	auth := server.wxMini.GetAuth()
	code2SessionRes, err := auth.Code2Session(addUserRequest.Code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	getPhoneNumberRes, err := auth.GetPhoneNumber(addUserRequest.PhoneCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	param := sqlc.AddUsrParams{
		Openid:    code2SessionRes.OpenID,
		Name:      "微信用户",
		Phone:     getPhoneNumberRes.PhoneInfo.PurePhoneNumber,
		StudentID: "",
		AvatarUrl: "",
		Avatar:    nil,
	}
	user, err := server.store.AddUsr(ctx, param)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	response := newUserResponse(user)
	ctx.JSON(http.StatusOK, response)
}

type addUserResponse = userResponse

//用户登录
type loginUserRequest struct {
	Code      string `json:"code" binding:"required,alphanum"`
	SessionId string `json:"session_id" binding:"required,uuid"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var loginUserRequest loginUserRequest
	if err := ctx.ShouldBindJSON(&loginUserRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	session, err := server.sessionManager.ReadSession(loginUserRequest.SessionId)
	if session != nil && err == nil { //用户已登录
		ctx.Redirect(http.StatusFound, "/") //重定向至首页
		return
	}

	//用户登录
	//获取用户OpenId和SessionKey
	auth := server.wxMini.GetAuth()
	res, err := auth.Code2Session(loginUserRequest.Code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	sessionValue := service.FormSessionValue(res)

	//添加新session
	sessionId, err := server.sessionManager.AddSession(sessionValue)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	//通过openid获取用户
	user, err := server.store.GetUsr(ctx, res.OpenID)
	if err == sql.ErrNoRows { //未注册，添加该用户
		ctx.Redirect(http.StatusFound, "/users/add")
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.SetCookie(cookieName, sessionId, int(server.config.SessionLifeTime.Seconds()), "/", domain, false, true)
	response := loginUserResponse{
		User: newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, response)
}

type loginUserResponse struct {
	User userResponse
}

//用户获取
type getUserRequest struct {
	Openid string `json:"openid" binding:"alphanum"`
}

func (server *Server) getUser(ctx *gin.Context) {

}

//显示所有用户
type listUserRequest struct{}

func (server *Server) listUser(ctx *gin.Context) {

}

type listUserResponse struct{}

//更新用户信息
type updateUserRequest struct{}

func (server *Server) updateUser(ctx *gin.Context) {

}

type updateUserResponse struct{}

//删除用户
type deleteUserRequest struct{}

func (server *Server) deleteUser(ctx *gin.Context) {

}

type deleteUserResponse struct{}
