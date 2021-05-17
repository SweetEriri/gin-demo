package main

import (
	"gin.vue.demo/ginessential/controller"
	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine)  *gin.Engine {
	r.POST("/api/auth/register", controller.Register)
	return r
}