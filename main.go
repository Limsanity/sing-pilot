package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/limsanity/sing-pilot/dto"
	"github.com/limsanity/sing-pilot/model"
	"github.com/limsanity/sing-pilot/service"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var cf = "default.json"

const (
	DB_FILE = "sing_pilot.db"
)

func main() {
	sb := service.SingBox{}
	sb.Start(cf)

	// initialize db
	db, err := gorm.Open(sqlite.Open(DB_FILE), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&model.Config{})

	// initialize http server
	router := gin.Default()

	api := router.Group("/api")

	// create config
	api.POST("/config", func(ctx *gin.Context) {
		dto := dto.CreateConfigDto{}
		err := ctx.ShouldBindJSON(&dto)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		config := model.Config{
			Content: dto.Content,
		}

		if result := db.Create(&config); result.Error != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "success"})
	})

	// update config
	api.PATCH("/config", func(ctx *gin.Context) {
		var dto dto.PatchConfigDto
		err := ctx.ShouldBindJSON(&dto)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		var config model.Config
		if result := db.First(&config, dto.ID); result.Error != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
			return
		}

		config.Content = dto.Content
		db.Save(&config)

		ctx.JSON(http.StatusCreated, gin.H{"message": "success"})
	})

	// get all config
	api.GET("/config", func(ctx *gin.Context) {
		var configList []model.Config
		if result := db.Find(&configList); result.Error != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"result": configList})
	})

	// delete config
	api.DELETE("/config/:id", func(ctx *gin.Context) {
		dto := &dto.DeleteConfigDto{}
		err := ctx.ShouldBindUri(&dto)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if result := db.Delete(&model.Config{}, dto.ID); result.Error != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// restart sing-box
	api.POST("/sing_box/restart", func(ctx *gin.Context) {
		var dto dto.RestartDto
		err := ctx.ShouldBindJSON(&dto)
		if err == nil && dto.ConfigId != nil {
			var config model.Config
			if result := db.First(&config, *dto.ConfigId); result.Error != nil {
				ctx.JSON(http.StatusBadGateway, gin.H{"message": result.Error.Error()})
				return
			}

			file := "tmp/" + fmt.Sprint(config.ID) + ".json"
			fd, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0644)
			defer func() {
				if err := fd.Close(); err != nil {
					log.Fatal(err)
				}
			}()

			if err != nil {
				ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
				return
			}

			fd.WriteString(config.Content)
			cf = file
		}

		sb.Stop()
		sb.Start(cf)
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	api.POST("/sing_box/start", func(ctx *gin.Context) {
		sb.Start(cf)
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	api.POST("/sing_box/stop", func(ctx *gin.Context) {
		sb.Stop()
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	log.Fatal(router.Run(":8080"))
}
