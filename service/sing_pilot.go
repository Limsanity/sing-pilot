package service

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/limsanity/sing-pilot/model"
	"gorm.io/gorm"
)

const (
	SING_BOX       = "sing-box"
	DEFAULT_CONFIG = "default.json"
)

type SingPilot struct {
	cmd  *exec.Cmd
	db   *gorm.DB
	file string
}

func NewSingPilotService(db *gorm.DB) SingPilot {
	sp := SingPilot{
		db:   db,
		file: DEFAULT_CONFIG,
	}

	var userConfig model.UserConfig
	if result := db.First(&userConfig); result.Error == nil {
		sp.UseFile(userConfig.ConfigId)
	}

	return sp
}

func (sp *SingPilot) UseFile(id uint) error {
	config := model.Config{}
	if result := sp.db.First(&config, id); result.Error != nil {
		return result.Error
	}

	file := fmt.Sprint("tmp/", id, ".json")
	fd, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	defer func() {
		if err := fd.Close(); err != nil {
			panic(err)
		}
	}()

	_, err = fd.WriteString(config.Content)

	if err != nil {
		return err
	}

	sp.file = file

	return nil
}

func (sb *SingPilot) Stop() {
	if sb.cmd == nil {
		return
	}

	// stop sing-box
	err := sb.cmd.Process.Kill()
	if err != nil {
		log.Fatal("kill sing-box error: ", err)
	}

	sb.cmd = nil
}

func (sp *SingPilot) Start() {
	if sp.cmd != nil {
		return
	}

	// run sing-box
	args := []string{"run", "-c", sp.file}
	sp.cmd = exec.Command(SING_BOX, args...)

	// log sing-box stdout
	stdoutReader, err := sp.cmd.StdoutPipe()
	if err != nil {
		log.Fatal("stdout pipe sing-box error: ", err)
	}
	go func() {
		scanner := bufio.NewScanner(stdoutReader)
		for scanner.Scan() {
			log.Print("[sing-box] ", scanner.Text())
		}
	}()

	// log sing-box stderr
	stderrReader, err := sp.cmd.StderrPipe()
	if err != nil {
		log.Fatal("stderr pipe sing-box error: ", err)
	}
	go func() {
		scanner := bufio.NewScanner(stderrReader)
		for scanner.Scan() {
			log.Print("[sing-box] ", scanner.Text())
		}
	}()

	err = sp.cmd.Start()
	if err != nil {
		log.Fatal("start sing-box error: ", err)
	}
}
