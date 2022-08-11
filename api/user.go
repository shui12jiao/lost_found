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
	ErrNoSession         = errors.New("no session")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrNoPermission      = errors.New("no permission")
	ErrPermissionDenied  = errors.New("permission denied")
)

type userResponse struct {
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	StudentID string `json:"studentID"`
	AvatarUrl string `json:"avatarUrl"`
	Avatar    []byte `json:"avatar"`
}

func newUserResponse(user sqlc.Usr) userResponse {
	return userResponse{
		Name:      user.Name,
		Phone:     user.Phone,
		StudentID: user.StudentID,
		AvatarUrl: user.AvatarUrl,
		Avatar:    user.Avatar,
	}
}

//用户获取自己的信息
func (server *Server) getUserSelf(ctx *gin.Context) {
	openid := ctx.MustGet(middleware.SessionHeaderKey).(session.Session).ID.String()

	user, err := server.store.GetUsr(ctx, openid)
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusBadRequest, errorResponse(ErrUserNotFound))
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	response := newUserResponse(user)
	ctx.JSON(http.StatusOK, response)
}

//用户永久注销自己账户
func (server *Server) deleteUserSelf(ctx *gin.Context) {
	openid := ctx.MustGet(middleware.SessionHeaderKey).(session.Session).ID.String()

	err := server.store.DeleteUsr(ctx, openid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	ctx.JSON(http.StatusOK, nil)
}

//用户更新自己账户信息
type updateUserSelfRequest struct {
	Name      string `json:"name" binding:"required,min=4,max=24"`
	AvatarUrl string `json:"avatarUrl" binding:"required,url"`
	Avatar    []byte `json:"avatar" binding:"required"`
}

func (server *Server) updateUserSelf(ctx *gin.Context) {
	openid := ctx.MustGet(middleware.SessionHeaderKey).(session.Session).ID.String()

	var request updateUserSelfRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUsr(ctx, openid)
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusBadRequest, errorResponse(ErrUserNotFound))
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
	PhoneCode string `json:"phoneCode" binding:"required,alphanum"`
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

	_, err = server.store.GetUsr(ctx, code2SessionRes.OpenID)
	if err != sql.ErrNoRows {
		ctx.JSON(http.StatusBadRequest, errorResponse(ErrUserAlreadyExists))
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
	SessionID string `json:"sessionID" binding:"required,uuid"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var loginUserRequest loginUserRequest
	if err := ctx.ShouldBindJSON(&loginUserRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	session, err := server.sessionManager.ReadSession(loginUserRequest.SessionID)
	if session != nil && err == nil { //用户已登录
		ctx.Redirect(http.StatusFound, "/") //重定向至首页
		return
	}

	//用户登录
	//获取用户OpenID和SessionKey
	auth := server.wxMini.GetAuth()
	res, err := auth.Code2Session(loginUserRequest.Code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	sessionValue := service.FormSessionValue(res)

	//添加新session
	sessionID, err := server.sessionManager.AddSession(sessionValue)
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

	ctx.SetCookie(cookieName, sessionID, int(server.config.SessionLifeTime.Seconds()), "/", domain, false, true)
	response := loginUserResponse{
		User: newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, response)
}

type loginUserResponse struct {
	User userResponse `json:"user"`
}

//获取用户
type getUserRequest struct {
	Openid string `json:"openid" binding:"required,alphanum,min=28"`
}

func (server *Server) getUser(ctx *gin.Context) {
	var request getUserRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUsr(ctx, request.Openid)
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusBadRequest, errorResponse(ErrUserNotFound))
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	response := newUserResponse(user)
	ctx.JSON(http.StatusOK, response)
}

//搜索用户
type searchUserRequest struct {
	Query string `json:"query" binding:"required,alphanum"`
}

func (server *Server) searchUser(ctx *gin.Context) {
	var request searchUserRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	users, err := server.store.SearchUsr(ctx, request.Query)
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusOK, nil)
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	//返回值包括openid
	ctx.JSON(http.StatusOK, users)
}

//显示所有用户
type listUserRequest struct {
	PageID   int32 `json:"pageID" binding:"required,gte=1"`
	PageSize int32 `json:"pageSize" binding:"required,min=5,max=30"`
}

func (server *Server) listUser(ctx *gin.Context) {
	var request listUserRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	param := sqlc.ListUsrParams{
		Limit:  request.PageSize,
		Offset: (request.PageID - 1) * request.PageSize,
	}
	users, err := server.store.ListUsr(ctx, param)
	if err != nil && err != sql.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	//返回值包括用户openid
	ctx.JSON(http.StatusOK, users)
}

//更新用户信息
type updateUserRequest struct {
	Openid    string `json:"openid" binding:"required,alphanum,min=28"`
	Name      string `json:"name" binding:"required,min=4,max=24"`
	Phone     string `json:"phone" binding:"required,len=11"`
	StudentID string `json:"studentID" binding:"required"`
	AvatarUrl string `json:"avatarUrl" binding:"required,url"`
	Avatar    []byte `json:"avatar" binding:"required"`
}

func (server *Server) updateUser(ctx *gin.Context) {
	var request updateUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//低于PermissionLevel2的管理员无法更新用户信息
	manager := ctx.MustGet(middleware.ManagerHeaderKey).(sqlc.Manager)
	if per := manager.Permission; per != sqlc.PermissionLevel2 && per != sqlc.PermissionLevel3 {
		ctx.JSON(http.StatusForbidden, errorResponse(ErrPermissionDenied))
		return
	}

	param := sqlc.UpdateUsrParams{
		Openid:    request.Openid,
		Name:      request.Name,
		Phone:     request.Phone,
		StudentID: request.StudentID,
		AvatarUrl: request.AvatarUrl,
		Avatar:    request.Avatar,
	}
	user, err := server.store.UpdateUsr(ctx, param)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, user)
}

//删除用户
type deleteUserRequest struct {
	Openid string `json:"openid" binding:"required,alphanum,min=28"`
}

func (server *Server) deleteUser(ctx *gin.Context) {
	var request deleteUserRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//低于PermissionLevel2的管理员无法删除用户
	manager := ctx.MustGet(middleware.ManagerHeaderKey).(sqlc.Manager)
	if per := manager.Permission; per != sqlc.PermissionLevel2 && per != sqlc.PermissionLevel3 {
		ctx.JSON(http.StatusForbidden, errorResponse(ErrPermissionDenied))
		return
	}

	err := server.store.DeleteUsr(ctx, request.Openid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
