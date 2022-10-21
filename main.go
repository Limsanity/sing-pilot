package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/limsanity/sing-pilot/dto"
	"github.com/limsanity/sing-pilot/model"
	"github.com/limsanity/sing-pilot/service"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	DB_FILE = "sing_pilot.db"
)

func main() {
	// initialize db
	db, err := gorm.Open(sqlite.Open(DB_FILE), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&model.Config{})
	db.AutoMigrate(&model.UserConfig{})

	sp := service.NewSingPilotService(db)
	sp.Start()

	// initialize http server
	router := gin.Default()

	api := router.Group("/api")

	singBoxApi := api.Group("/sing_box")

	// create config
	singBoxApi.POST("/config", func(ctx *gin.Context) {
		var dto dto.CreateConfigDto
		err := ctx.ShouldBindJSON(&dto)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		config := model.Config{
			Content: dto.Content,
		}

		if result := db.Create(&config); result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": result.Error.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "success"})
	})

	// update config
	singBoxApi.PATCH("/config", func(ctx *gin.Context) {
		var dto dto.PatchConfigDto
		err := ctx.ShouldBindJSON(&dto)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		var config model.Config
		if result := db.First(&config, dto.ID); result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": result.Error.Error()})
			return
		}

		config.Content = dto.Content
		db.Save(&config)

		ctx.JSON(http.StatusCreated, gin.H{"message": "success"})
	})

	// get all config
	singBoxApi.GET("/config", func(ctx *gin.Context) {
		var configList []model.Config
		if result := db.Find(&configList); result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": result.Error.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"result": configList})
	})

	// delete config
	singBoxApi.DELETE("/config/:id", func(ctx *gin.Context) {
		var dto dto.DeleteConfigDto
		err := ctx.ShouldBindUri(&dto)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		if result := db.Delete(&model.Config{}, dto.ID); result.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": result.Error.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// restart sing-box
	singBoxApi.POST("/restart", func(ctx *gin.Context) {
		var dto dto.RestartDto
		err := ctx.ShouldBindJSON(&dto)
		if err == nil && dto.ConfigId != nil {
			if err := sp.UseFile(*dto.ConfigId); err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}

			var userConfig model.UserConfig
			db.First(&userConfig)
			userConfig.ConfigId = *dto.ConfigId
			db.Save(&userConfig)
		}

		sp.Stop()
		sp.Start()
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	singBoxApi.POST("/start", func(ctx *gin.Context) {
		sp.Start()
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	singBoxApi.POST("/stop", func(ctx *gin.Context) {
		sp.Stop()
		ctx.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	log.Fatal(router.Run(":8080"))
}
