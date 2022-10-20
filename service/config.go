package service

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	File string
}

func (c *Config) UseFile(id uint, content string) error {
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

	_, err = fd.WriteString(content)

	if err != nil {
		return err
	}

	c.File = file

	return nil
}
