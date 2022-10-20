package service

import (
	"bufio"
	"log"
	"os/exec"
)

const (
	SING_BOX = "sing-box"
)

type SingBox struct {
	cmd *exec.Cmd
}

func (sb *SingBox) Stop() {
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

func (sb *SingBox) Start(cf string) {
	if sb.cmd != nil {
		return
	}

	// run sing-box
	args := []string{"run", "-c", cf}
	sb.cmd = exec.Command(SING_BOX, args...)

	// log sing-box stdout
	stdoutReader, err := sb.cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		scanner := bufio.NewScanner(stdoutReader)
		for scanner.Scan() {
			log.Print("[sing-box] ", scanner.Text())
		}
	}()

	// log sing-box stderr
	stderrReader, err := sb.cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		scanner := bufio.NewScanner(stderrReader)
		for scanner.Scan() {
			log.Print("[sing-box] ", scanner.Text())
		}
	}()

	err = sb.cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
}
