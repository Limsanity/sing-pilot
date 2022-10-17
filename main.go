package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/limsanity/sing-pilot/model"
	"github.com/limsanity/sing-pilot/service"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	DB_FILE = "sing_pilot.db"
)

func main() {
	sb := service.SingBox{}
	sb.Start()

	// initialize db
	db, err := gorm.Open(sqlite.Open(DB_FILE), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&model.Config{})

	// initialize http server
	router := gin.Default()

	// create config
	router.POST("/config", func(ctx *gin.Context) {
		config := model.Config{}
		err := ctx.ShouldBindJSON(&config)
		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
			return
		}

		if result := db.Create(&config); result.Error != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "success"})
	})

	// get all config
	router.GET("/config", func(ctx *gin.Context) {
		var configList []model.Config
		if result := db.Find(&configList); result.Error != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"result": configList})
	})

	// restart sing-box
	router.POST("/sing_box/restart", func(ctx *gin.Context) {
		sb.Stop()
		sb.Start()
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	router.POST("/sing_box/start", func(ctx *gin.Context) {
		sb.Start()
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	router.POST("/sing_box/stop", func(ctx *gin.Context) {
		sb.Stop()
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	log.Fatal(router.Run(":8080"))
}
