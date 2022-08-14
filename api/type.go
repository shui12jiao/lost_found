package api

import (
	"database/sql"
	"lost_found/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

//添加大地点
type addTypeWideRequest struct {
	Name string `json:"name" binding:"required"`
}

func (server *Server) addTypeWide(ctx *gin.Context) {
	var request addTypeWideRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	wide, err := server.store.AddTypeWide(ctx, request.Name)
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
	ctx.JSON(http.StatusOK, wide)
}

//添加小地点
type addTypeNarrowRequest struct {
	Name   string `json:"name" binding:"required"`
	WideID int16  `json:"wide_id" binding:"required"`
}

func (server *Server) addTypeNarrow(ctx *gin.Context) {
	var request addTypeNarrowRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	param := sqlc.AddTypeNarrowParams{
		Name:   request.Name,
		WideID: request.WideID,
	}
	narrow, err := server.store.AddTypeNarrow(ctx, param)
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
	ctx.JSON(http.StatusOK, narrow)
}

//展示地点
func (server *Server) listType(ctx *gin.Context) {
	narrow := make(map[string][]sqlc.TypeNarrow, 10)

	wide, err := server.store.ListTypeWide(ctx)
	if err != nil && err != sql.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	for _, w := range wide {
		n, err := server.store.ListTypeNarrowByWide(ctx, w.ID)
		if err != nil && err != sql.ErrNoRows {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		narrow[w.Name] = n
	}
	response := listTypeResponse{
		Wide:   wide,
		Narrow: narrow,
	}
	ctx.JSON(http.StatusOK, response)
}

type listTypeResponse struct {
	Wide   []sqlc.TypeWide              `json:"wide"`
	Narrow map[string][]sqlc.TypeNarrow `json:"narrow"`
}

//删除大地点
type deleteTypeWideRequest struct {
	ID int16 `json:"id" binding:"required"`
}

func (server *Server) deleteTypeWide(ctx *gin.Context) {
	var request deleteTypeWideRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := server.store.DeleteTypeWide(ctx, request.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

//删除小地点
type deleteTypeNarrowRequest struct {
	ID int16 `json:"id" binding:"required"`
}

func (server *Server) deleteTypeNarrow(ctx *gin.Context) {
	var request deleteTypeNarrowRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := server.store.DeleteTypeNarrow(ctx, request.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, nil)
}
