package judgetools

import (
	"errors"
	"fmt"
	"os/exec"
)

func ExecCp(first string, second string) error {
	cmd := exec.Command("cp", "-f", first, second)
	err := cmd.Run()
	if err != nil {
		fmt.Println("cp file err: ", first, second)
		fmt.Println(err)
		return errors.New("cp file err")
	}
	return nil
}

func ExecGcc(first string, second string) error {
	myExec := exec.Command("g++", "-o", first, second)
	execErr := myExec.Run()
	if execErr != nil {
		fmt.Println("g++ " + second + " error: ", execErr)
		return errors.New("g++ " + second + " error")
	}
	return nil
}
