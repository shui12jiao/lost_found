package api

import (
	"database/sql"
	"lost_found/db/sqlc"
	"lost_found/middleware"
	"lost_found/middleware/session"
	"lost_found/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

//-------拾取物品-------
//添加拾取物
type addFoundRequst struct {
	FoundDate      time.Time           `json:"foundDate" binding:"required"`
	TimeBucket     sqlc.TimeBucket     `json:"timeBucket" binding:"required"`
	LocationID     int16               `json:"locationID" binding:"required"`
	LocationInfo   string              `json:"locationInfo" binding:"required"`
	LocationStatus sqlc.LocationStatus `json:"locationStatus" binding:"required"`
	TypeID         int16               `json:"typeID" binding:"required"`
	ItemInfo       string              `json:"itemInfo" binding:"required"`
	Image          []byte              `json:"image" binding:"required"`
	ImageKey       string              `json:"imageKey" binding:"required"`
	OwnerInfo      string              `json:"ownerInfo"`
	AddtionalInfo  string              `json:"addtionalInfo"`
}

func (server *Server) addFound(ctx *gin.Context) {
	var request addFoundRequst
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	userSession := ctx.MustGet(middleware.SessionHeaderKey).(session.Session)

	param := sqlc.AddFoundParams{
		PickerOpenid:   userSession.ID.String(),
		FoundDate:      request.FoundDate,
		TimeBucket:     request.TimeBucket,
		LocationID:     request.LocationID,
		LocationInfo:   request.LocationInfo,
		LocationStatus: request.LocationStatus,
		TypeID:         request.TypeID,
		ItemInfo:       request.ItemInfo,
		Image:          request.Image,
		ImageKey:       request.ImageKey,
	}
	found, err := server.store.AddFound(ctx, param)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, found)
}

//通过id获取拾取物
type getFoundRequest struct {
	ID int32 `json:"id" binding:"required"`
}

func (server *Server) getFound(ctx *gin.Context) {
	var request getFoundRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	found, err := server.store.GetFound(ctx, request.ID)
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, found)
}

//获取我的拾取物
func (server *Server) listMyFound(ctx *gin.Context) {
	openid := ctx.MustGet(middleware.SessionHeaderKey).(session.Session).ID.String()

	found, err := server.store.ListFoundByPicker(ctx, openid)
	if err != nil && err != sql.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, found)
}

//拾取物品已归还
type completeMyFoundRequest struct {
	OwnerOpenid string `json:"ownerOpenid" binding:"required"`
	Comment     string `json:"comment" binding:"required"`
}

func (server *Server) completeMyFound(ctx *gin.Context) {
	var itemId int32
	if err := ctx.ShouldBindUri(&itemId); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var request completeMyFoundRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	openid := ctx.MustGet(middleware.SessionHeaderKey).(session.Session).ID.String()

	found, err := server.store.GetFound(ctx, itemId)
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if found.PickerOpenid != openid {
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	param := service.CompleteFoundTxParams{
		AddMatchParam: sqlc.AddMatchParams{
			PickerOpenid: openid,
			OwnerOpenid:  request.OwnerOpenid,
			FoundDate:    time.Now(),
			LostDate:     found.FoundDate,
			TypeID:       found.TypeID,
			ItemInfo:     found.ItemInfo,
			Image:        found.Image,
			ImageKey:     found.ImageKey,
			Comment:      request.Comment,
		},
		ID: itemId,
	}
	match, err := server.store.CompleteFoundTx(ctx, param)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, match)
}

//展示拾取物
type listFoundRequest struct {
	PageID   int32 `json:"pageID" binding:"required,gte=1"`
	PageSize int32 `json:"pageSize" binding:"required,min=5,max=30"`
}

func (server *Server) listFound(ctx *gin.Context) {
	var request listFoundRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	param := sqlc.ListFoundParams{
		Limit:  request.PageSize,
		Offset: (request.PageID - 1) * request.PageSize,
	}
	found, err := server.store.ListFound(ctx, param)
	if err != nil && err != sql.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, found)
}

//删除拾取物
type deleteFoundRequest struct {
	ID int32 `json:"id" binding:"required"`
}

func (server *Server) deleteFound(ctx *gin.Context) {
	var request deleteFoundRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	openid := ctx.MustGet(middleware.SessionHeaderKey).(session.Session).ID.String()
	manager, ok := ctx.Get(middleware.ManagerHeaderKey)

	found, err := server.store.GetFound(ctx, request.ID)
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if found.PickerOpenid != openid {
		if ok {
			per := manager.(sqlc.Manager).Permission
			if per != sqlc.PermissionLevel2 && per != sqlc.PermissionLevel3 {
				ctx.JSON(http.StatusForbidden, errorResponse(ErrPermissionDenied))
				return
			}
		} else {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
	}

	err = server.store.DeleteFound(ctx, request.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

//-------遗失物品-------
//添加遗失物
type addLostRequst struct {
	LostDate    time.Time       `json:"lostDate" binding:"required"`
	TimeBucket  sqlc.TimeBucket `json:"timeBucket" binding:"required"`
	TypeID      int16           `json:"typeID" binding:"required"`
	ItemInfo    string          `json:"itemInfo" binding:"required"`
	Image       []byte          `json:"image"`    //遗物可能没有图片
	ImageKey    string          `json:"imageKey"` //同上
	LocationID  int16           `json:"locationID" binding:"required"`
	LocationId1 int16           `json:"locationId1"`
	LocationId2 int16           `json:"locationId2"`
}

//添加遗失物
func (server *Server) addLost(ctx *gin.Context) {
	var request addLostRequst
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	openid := ctx.MustGet(middleware.SessionHeaderKey).(session.Session).ID.String()

	param := sqlc.AddLostParams{
		OwnerOpenid: openid,
		LostDate:    request.LostDate,
		TimeBucket:  request.TimeBucket,
		TypeID:      request.TypeID,
		ItemInfo:    request.ItemInfo,
		Image:       request.Image,
		ImageKey:    request.ImageKey,
		LocationID:  request.LocationID,
		LocationId1: request.LocationId1,
		LocationId2: request.LocationId2,
	}
	lost, err := server.store.AddLost(ctx, param)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, lost)
}

//通过id获取遗失物
type getLostRequest struct {
	ID int32 `json:"id" binding:"required"`
}

func (server *Server) getLost(ctx *gin.Context) {
	var request getLostRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	lost, err := server.store.GetLost(ctx, request.ID)
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, lost)
}

//获取我的遗失物
func (server *Server) listMyLost(ctx *gin.Context) {
	openid := ctx.MustGet(middleware.SessionHeaderKey).(session.Session).ID.String()

	found, err := server.store.ListLostByOwner(ctx, openid)
	if err != nil && err != sql.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, found)
}

//遗失物品已寻回
type completeMyLostRequest struct {
	PickerOpenid string `json:"pickerOpenid" binding:"required"`
	Comment      string `json:"comment" binding:"required"`
}

func (server *Server) completeMyLost(ctx *gin.Context) {
	var itemId int32
	if err := ctx.ShouldBindUri(&itemId); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var request completeMyLostRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	openid := ctx.MustGet(middleware.SessionHeaderKey).(session.Session).ID.String()

	lost, err := server.store.GetLost(ctx, itemId)
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if lost.OwnerOpenid != openid {
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	param := service.CompleteLostTxParams{
		AddMatchParam: sqlc.AddMatchParams{
			PickerOpenid: request.PickerOpenid,
			OwnerOpenid:  lost.OwnerOpenid,
			FoundDate:    time.Now(),
			LostDate:     lost.LostDate,
			TypeID:       lost.TypeID,
			ItemInfo:     lost.ItemInfo,
			Image:        lost.Image,
			ImageKey:     lost.ImageKey,
			Comment:      request.Comment,
		},
		ID: itemId,
	}
	match, err := server.store.CompleteLostTx(ctx, param)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, match)
}

//展示遗失物
type listLostRequest struct {
	PageID   int32 `json:"pageID" binding:"required,gte=1"`
	PageSize int32 `json:"pageSize" binding:"required,min=5,max=30"`
}

func (server *Server) listLost(ctx *gin.Context) {
	var request listLostRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	param := sqlc.ListLostParams{
		Limit:  request.PageSize,
		Offset: (request.PageID - 1) * request.PageSize,
	}
	lost, err := server.store.ListLost(ctx, param)
	if err != nil && err != sql.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, lost)
}

//删除遗失物
type deleteLostRequest struct {
	ID int32 `json:"id" binding:"required"`
}

func (server *Server) deleteLost(ctx *gin.Context) {
	var request deleteLostRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	openid := ctx.MustGet(middleware.SessionHeaderKey).(session.Session).ID.String()
	manager, ok := ctx.Get(middleware.ManagerHeaderKey)

	lost, err := server.store.GetLost(ctx, request.ID)
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if lost.OwnerOpenid != openid {
		if ok {
			per := manager.(sqlc.Manager).Permission
			if per != sqlc.PermissionLevel2 && per != sqlc.PermissionLevel3 {
				ctx.JSON(http.StatusForbidden, errorResponse(ErrPermissionDenied))
				return
			}
		} else {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
	}

	err = server.store.DeleteFound(ctx, request.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

//-------归还物品-------
//添加已归还物品
type addMatchRequst struct {
	PickerOpenid string    `json:"pickerOpenid"`
	OwnerOpenid  string    `json:"ownerOpenid"`
	FoundDate    time.Time `json:"foundDate" binding:"required"`
	LostDate     time.Time `json:"lostDate" binding:"required"`
	TypeID       int16     `json:"typeID" binding:"required"`
	ItemInfo     string    `json:"itemInfo" binding:"required"`
	Image        []byte    `json:"image" binding:"required"`
	ImageKey     string    `json:"imageKey" binding:"required"`
	Comment      string    `json:"comment"`
}

func (server *Server) addMatch(ctx *gin.Context) {
	var request addMatchRequst
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	param := sqlc.AddMatchParams{
		PickerOpenid: request.PickerOpenid,
		OwnerOpenid:  request.OwnerOpenid,
		FoundDate:    request.FoundDate,
		LostDate:     request.LostDate,
		TypeID:       request.TypeID,
		ItemInfo:     request.ItemInfo,
		Image:        request.Image,
		ImageKey:     request.ImageKey,
		Comment:      request.Comment,
	}
	match, err := server.store.AddMatch(ctx, param)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, match)
}

//通过id获取已归还物品
type getMatchRequest struct {
	ID int32 `json:"id" binding:"required"`
}

func (server *Server) getMatch(ctx *gin.Context) {
	var request getMatchRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	match, err := server.store.GetMatch(ctx, request.ID)
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, match)
}

//获取我的归还物品，包括拾取的和遗失的
func (server *Server) listMyMatch(ctx *gin.Context) {
	openid := ctx.MustGet(middleware.SessionHeaderKey).(session.Session).ID.String()

	found, err := server.store.ListMatchByPicker(ctx, openid)
	if err != nil && err != sql.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	lost, err := server.store.ListMatchByOwner(ctx, openid)
	if err != nil && err != sql.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := listMyMatchResponse{
		Found: found,
		Lost:  lost,
	}
	ctx.JSON(http.StatusOK, response)
}

type listMyMatchResponse struct {
	Found []sqlc.Match `json:"found"`
	Lost  []sqlc.Match `json:"lost"`
}

//展示已归还物品
type listMatchRequest struct {
	PageID   int32 `json:"pageID" binding:"required,gte=1"`
	PageSize int32 `json:"pageSize" binding:"required,min=5,max=30"`
}

func (server *Server) listMatch(ctx *gin.Context) {
	var request listMatchRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	param := sqlc.ListMatchParams{
		Limit:  request.PageSize,
		Offset: (request.PageID - 1) * request.PageSize,
	}
	match, err := server.store.ListMatch(ctx, param)
	if err != nil && err != sql.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, match)
}

//删除归还物品
type deleteMatchRequest struct {
	ID int32 `json:"id" binding:"required"`
}

func (server *Server) deleteMatch(ctx *gin.Context) {
	var request deleteMatchRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	per := ctx.MustGet(middleware.ManagerHeaderKey).(sqlc.Manager).Permission
	if per != sqlc.PermissionLevel2 && per != sqlc.PermissionLevel3 {
		ctx.JSON(http.StatusForbidden, errorResponse(ErrPermissionDenied))
		return
	}
}
