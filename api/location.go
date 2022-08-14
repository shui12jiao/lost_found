package api

import (
	"database/sql"
	"lost_found/db/sqlc"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

//添加大地点
type addLocationWideRequest struct {
	Name   string `json:"name" binding:"required"`
	Campus string `json:"campus" binding:"required"`
}

func (server *Server) addLocationWide(ctx *gin.Context) {
	var request addLocationWideRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	param := sqlc.AddLocationWideParams{
		Name:   request.Name,
		Campus: sqlc.Campus(request.Campus),
	}
	wide, err := server.store.AddLocationWide(ctx, param)
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
type addLocationNarrowRequest struct {
	Name   string `json:"name" binding:"required"`
	WideID int16  `json:"wide_id" binding:"required"`
}

func (server *Server) addLocationNarrow(ctx *gin.Context) {
	var request addLocationNarrowRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	param := sqlc.AddLocationNarrowParams{
		Name:   request.Name,
		WideID: request.WideID,
	}
	narrow, err := server.store.AddLocationNarrow(ctx, param)
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
func (server *Server) listLocation(ctx *gin.Context) {
	narrow := make(map[string][]sqlc.LocationNarrow, 10)

	wide, err := server.store.ListLocationWide(ctx)
	if err != nil && err != sql.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	for _, w := range wide {
		n, err := server.store.ListLocationNarrowByWide(ctx, w.ID)
		if err != nil && err != sql.ErrNoRows {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		}
		narrow[w.Name] = n
	}
	response := listLocationResponse{
		Wide:   wide,
		Narrow: narrow,
	}
	ctx.JSON(http.StatusOK, response)
}

type listLocationResponse struct {
	Wide   []sqlc.LocationWide              `json:"wide"`
	Narrow map[string][]sqlc.LocationNarrow `json:"narrow"`
}

//删除大地点
type deleteLocationWideRequest struct {
	ID int16 `json:"id" binding:"required"`
}

func (server *Server) deleteLocationWide(ctx *gin.Context) {
	var request deleteLocationWideRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := server.store.DeleteLocationWide(ctx, request.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

//删除小地点
type deleteLocationNarrowRequest struct {
	ID int16 `json:"id" binding:"required"`
}

func (server *Server) deleteLocationNarrow(ctx *gin.Context) {
	var request deleteLocationNarrowRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if err := server.store.DeleteLocationNarrow(ctx, request.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, nil)
}
