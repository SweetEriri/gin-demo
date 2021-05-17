package controller

import (
	"gin.vue.demo/ginessential/common"
	"gin.vue.demo/ginessential/model"
	"gin.vue.demo/ginessential/util"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
	hasedPassWord, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code" : 500, "msg": "加密错误"})
		return
	}
	newUser := model.User{
		Name: name,
		Telephone: telephone,
		Password: string(hasedPassWord),
	}

	DB.Create(&newUser)


	// 返回结果
	ctx.JSON(200, gin.H{
		"code" : "200",
		"msg" : "注册成功",
	})
}

func Login (ctx *gin.Context)  {
	DB := common.GetDB()

	// 获取参数
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")

	//数据验证
	if len(telephone) != 11 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code" : 422, "msg" : "手机号必须为11位"})
		return
	}
	if len(password) < 6 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code" : 422, "msg": "密码不能少于6位"})
		return
	}

	//判断手机号是否存在
	var user model.User
	DB.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code" : 422, "msg" : "用户不存在"})
		return
	}

	// 判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) ; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code" : 400, "msg" : "密码错误"})
		return
	}

	// 发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code" : 500, "msg" : "系统异常"})
		log.Printf("token generate error : %v", err)
		return
	}

	//发放结果
	ctx.JSON(200, gin.H{
		"code" : "200",
		"data" : gin.H{"token" : token},
		"msg" : "登陆成功",
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