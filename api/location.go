package api

import "github.com/gin-gonic/gin"

//添加地点
type addLocationRequest struct{}

func (server *Server) addLocation(ctx *gin.Context) {

}

type addLocationResponse struct{}

//展示地点
type listLocationRequest struct{}

func (server *Server) listLocation(ctx *gin.Context) {

}

type listLocationResponse struct{}

//删除地点
type deleteLocationRequest struct{}

func (server *Server) deleteLocation(ctx *gin.Context) {}

type deleteLocationResponse struct{}
