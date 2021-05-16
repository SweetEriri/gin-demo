package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"net/http"
)

type User struct {
	gorm.Model
	Name string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"type:varchar(11)";not null;unique`
	Password string `gorm:"size:255";not null`
}

func main(){
	r := gin.Default()
	r.POST("/api/auth/register", func(ctx *gin.Context) {
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
			name = RandomString(10)
		}

		log.Println(name, telephone, password)
		// 判断手机号是否存在
		if isTelephoneExist(InitDB(), telephone){
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code" : 422, "msg": "用户已存在"})
			return
		}
		// 创建用户
		newUser := User{
			Name: name,
			Telephone: telephone,
			Password: password,
		}

		InitDB().Create(&newUser)


		// 返回结果
		ctx.JSON(200, gin.H{
			"msg": "注册成功",
		})
	})
	panic(r.Run()) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}

	return false
}

func RandomString(n int) string {
	var letter = []byte("asdfghjklqwertyuiopASDFGHJKLQWERTYUIOP")
	result :=make([]byte, n)
	for i := range result{
		result[i] = letter[rand.Intn(len(letter))]
	}

	return string(result)
}

func InitDB() *gorm.DB {
	//driverName := "mysql"
	host := "localhost"
	port := "3306"
	database := "ginessential"
	username := "root"
	password := ""
	charset := "utf8"
	args := fmt.Sprintf("%s%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)

	db, err := gorm.Open(mysql.Open(args), &gorm.Config{})
	if err != nil {
		panic("fial to connect database, err:" + err.Error())
	}
	db.AutoMigrate(&User{})

	return db
}