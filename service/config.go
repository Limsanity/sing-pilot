package service

import (
	"fmt"
	"log"
	"os"

	"github.com/limsanity/sing-pilot/model"
	"gorm.io/gorm"
)

type Config struct {
	file string
	db   *gorm.DB
}

func NewConfigService(file string, db *gorm.DB) Config {
	return Config{
		file: file,
		db:   db,
	}
}

func (c *Config) GetFile() string {
	return c.file
}

func (c *Config) UseFile(id uint) error {
	config := model.Config{}
	if result := c.db.First(&config); result.Error != nil {
		return result.Error
	}

	file := fmt.Sprint("tmp/", id, ".json")
	fd, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0644)
	defer func() {
		if err := fd.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if err != nil {
		return err
	}

	_, err = fd.WriteString(config.Content)

	if err != nil {
		return err
	}

	c.file = file

	return nil
}
