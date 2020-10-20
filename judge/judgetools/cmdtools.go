package judgetools

import (
	"errors"
	"judge/zapconf"
	"os/exec"
)

func ExecCp(first string, second string) error {
	cmd := exec.Command("cp", "-f", first, second)
	err := cmd.Run()
	if err != nil {
		zapconf.GetWarnLog().Warn("cp file err: " + first + " " + second)
		return errors.New("cp file err")
	}
	return nil
}

func ExecGcc(first string, second string) error {
	myExec := exec.Command("g++", "-o", first, second)
	execErr := myExec.Run()
	if execErr != nil {
		zapconf.GetWarnLog().Warn("g++ " + second + " error: " + execErr.Error())
		return errors.New("g++ " + second + " error")
	}
	return nil
}
