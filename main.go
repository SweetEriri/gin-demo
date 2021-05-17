package main

import (
	"gin.vue.demo/ginessential/common"
	"gin.vue.demo/ginessential/model"
	"github.com/gin-gonic/gin"
)



func main(){
	db := common.InitDB()
	db.AutoMigrate(&model.User{})
	r := gin.Default()
	r = CollectRoute(r)
	panic(r.Run()) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}





