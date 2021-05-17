package controller

import (
	"gin.vue.demo/ginessential/common"
	"gin.vue.demo/ginessential/model"
	"gin.vue.demo/ginessential/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func Register(ctx *gin.Context) {

	// 获取数据库连接
	DB := common.GetDB()

	// 获取参数
	name := ctx.PostForm("name")
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")
	// 数据验证
	if len(telephone) != 11 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code" : 422, "msg" : "手机号必须为11位"})
		return
	}
	if len(password) < 6 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code" : 422, "msg": "密码不能少于6位"})
		return
	}

	// 名称可为空 如果名称不传 给予一个随机的10位字符串
	if len(name) == 0 {
		name = util.RandomString(10)
	}

	log.Println(name, telephone, password)
	// 判断手机号是否存在
	if isTelephoneExist(DB, telephone){
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code" : 422, "msg": "用户已存在"})
		return
	}
	// 创建用户
	newUser := model.User{
		Name: name,
		Telephone: telephone,
		Password: password,
	}

	DB.Create(&newUser)


	// 返回结果
	ctx.JSON(200, gin.H{
		"msg": "注册成功",
	})
}



func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}

	return false
}