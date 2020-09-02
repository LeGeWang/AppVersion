package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
)

type AppVersion struct {
	gorm.Model
	AppName string `gorm:"type:varchar(100);unique_index"`
	AppVersion string `gorm:"type:varchar(100)"`
	url string
}
var db *gorm.DB
var err error
func main() {
	webGin := gin.Default()
	//连接mysql
	db, err = gorm.Open("mysql", "zs_app:RJEdD5S8BSstNJsz@(39.102.64.220:3306)/zs_app?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		println(err.Error())
		panic("链接mysql 失败")
	}
	db.AutoMigrate(&AppVersion{})

	//增加一个版本
	webGin.POST("/appVersion/add", addAppV)
	//删除一个版本
	webGin.DELETE("/appVersion/del/:id",delAppV)
	//修改一个版本
	webGin.POST("/appVersion/update",updateAppV)
	//查找一个游戏版本
	webGin.GET("/appVersion",getAppV)

	_ = webGin.Run("0.0.0.0:8084")
}

//增加一个版本
func addAppV(c *gin.Context) {
	appV := AppVersion{AppName: c.Request.FormValue("AppName"),AppVersion: c.Request.FormValue("AppVersion")}
	result := db.Create(&appV)
	c.JSON(http.StatusOK,gin.H{
		"msg": result,
	})

}

//删除一个版本
func delAppV(c *gin.Context) {
	id := c.Param("id")
	println(c,"请求删除的id")
	var appV AppVersion
	db.First(&appV, id)
	if appV.ID == 0 {
		c.JSON(404, gin.H{"message": "user not AppVersion"})
		return
	} else {
		_ = c.BindJSON(&appV)
		db.Delete(&appV)
		c.JSON(200, gin.H{"message": "delete success"})
	}
}

//修改一个版本
func updateAppV(c *gin.Context) {
	id := c.PostForm("id")
	var appV AppVersion
	db.First(&appV, id)

	if appV.ID == 0{
		c.JSON(404, gin.H{"message": "AppVersion not found"})
		return
	}else{
		appV.AppName = c.PostForm("AppName")
		appV.AppVersion = c.PostForm("AppVersion")
		db.Save(&appV)
		c.JSON(200,appV)
	}

}

//获取一个版本
//http://localhost:8084/appVersion?appName=ledada
func getAppV(c *gin.Context)  {
	appName := c.Query("appName")
	var appV AppVersion
	db.Where("app_name = ?",appName).First(&appV)
	c.JSON(http.StatusOK,appV)
}