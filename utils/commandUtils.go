package utils

import (
	"errors"
	"os/exec"
	"runtime"
	"terraform-updater/constants"
)

type CommandService struct {
	ExecCommand func (string, ...string) *exec.Cmd
}

func (cs *CommandService) FindCommand(command string) (string, error) {

	os := runtime.GOOS
	var whereCmd string

	if os == constants.OS_LINUX {
		whereCmd = "whereis"
	} else if os == constants.OS_WINDOWS {
		whereCmd = "where"
	} else {
		return "", errors.New("not supported")
	}

	out, err := cs.ExecCommand(whereCmd,command).Output()

	if err != nil {
		return "", err
	}

	return string(out), nil


}
