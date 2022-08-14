package api

import (
	"database/sql"
	"lost_found/db/sqlc"
	"lost_found/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type managerResponse struct {
	ID         int16           `json:"id"`
	Permission sqlc.Permission `json:"permission"`
}

func newManagerResponse(manager sqlc.Manager) managerResponse {
	return managerResponse{
		ID:         manager.ID,
		Permission: manager.Permission,
	}
}

//添加管理员
type addManagerRequst struct {
	UserOpenid string          `json:"userOpenid" binding:"required"`
	Permission sqlc.Permission `json:"permission" binding:"required"`
}

func (server *Server) addManager(ctx *gin.Context) {
	var request addManagerRequst
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if !request.Permission.Valid() {
		ctx.JSON(http.StatusBadRequest, errorResponse(ErrInvalidPermission))
		return
	}

	per := ctx.MustGet(middleware.ManagerHeaderKey).(sqlc.Manager).Permission
	if per != sqlc.PermissionLevel3 { //仅三级管理员有权限添加管理员
		ctx.JSON(http.StatusForbidden, errorResponse(ErrPermissionDenied))
		return
	}

	param := sqlc.AddManagerParams{
		UsrOpenid:  request.UserOpenid,
		Permission: request.Permission,
	}
	manager, err := server.store.AddManager(ctx, param)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, manager)
}

//通过id获取管理员
type getManagerRequest struct {
	ID int16 `json:"id" binding:"required"`
}

func (server *Server) getManager(ctx *gin.Context) {
	var requst getManagerRequest
	if err := ctx.ShouldBindUri(&requst); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	manager, err := server.store.GetManager(ctx, requst.ID)
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, manager)
}

//更改管理员权限
type updateManagerRequest struct {
	ID         int16           `json:"id" binding:"required"`
	Permission sqlc.Permission `json:"permission" binding:"required"`
}

func (server *Server) updateManager(ctx *gin.Context) {
	var request updateManagerRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if !request.Permission.Valid() {
		ctx.JSON(http.StatusBadRequest, errorResponse(ErrInvalidPermission))
		return
	}

	per := ctx.MustGet(middleware.ManagerHeaderKey).(sqlc.Manager).Permission
	if per != sqlc.PermissionLevel3 { //仅三级管理员有权限修改管理员等级
		ctx.JSON(http.StatusForbidden, errorResponse(ErrPermissionDenied))
		return
	}

	param := sqlc.UpdateManagerParams{
		ID:         request.ID,
		Permission: request.Permission,
	}
	manager, err := server.store.UpdateManager(ctx, param)
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, manager)
}

//展示管理员
func (server *Server) listManager(ctx *gin.Context) {
	managers, err := server.store.ListManager(ctx)
	if err != nil && err != sql.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	response := make([]managerResponse, 10)
	for _, m := range managers {
		response = append(response, newManagerResponse(m))
	}

	ctx.JSON(http.StatusOK, response)
}

//删除管理员
type deleteManagerRequest struct {
	ID int16 `json:"id" binding:"required"`
}

func (server *Server) deleteManager(ctx *gin.Context) {
	var request deleteManagerRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	per := ctx.MustGet(middleware.ManagerHeaderKey).(sqlc.Manager).Permission
	if per != sqlc.PermissionLevel3 { //仅三级管理员有权限删除管理员
		ctx.JSON(http.StatusForbidden, errorResponse(ErrPermissionDenied))
	}

	err := server.store.DeleteManager(ctx, request.ID)
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}
