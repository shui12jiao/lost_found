package api

import "github.com/gin-gonic/gin"

//添加类型
type addTypeRequest struct{}

func (server *Server) addType(ctx *gin.Context) {

}

type addTypeResponse struct{}

//展示类型
type listTypeRequest struct{}

func (server *Server) listType(ctx *gin.Context) {

}

type listTypeResponse struct{}

//删除类型
type deleteTypeRequest struct{}

func (server *Server) deleteType(ctx *gin.Context) {}

type deleteTypeResponse struct{}
