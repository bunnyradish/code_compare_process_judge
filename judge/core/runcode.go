package core

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"judge/controlroutine"
	"judge/db"
	"judge/dockerexec"
	"judge/environment"
	"judge/judgetools"
	"judge/zapconf"
	"os"
	"strconv"
	"strings"
	"time"
)

func StartRunCode(Db *sql.DB, cid string, flag string, cr *controlroutine.ChanRoutine) {
	codeGoPath := environment.CodeRunRootPath + cid + "/"
	runPath := codeGoPath + ServerRunCodePath
	ioPath := codeGoPath + ServerIOPath
	judgetools.FastCreateFile(codeGoPath)
	judgetools.FastCreateFile(runPath)
	judgetools.FastCreateFile(ioPath)

	myCode, err := db.GetEvaCodeWithCodeId(Db, cid)
	if err != nil {
		return
	}

	fmt.Printf("\n---code running---\n")
	myRunCode, err := db.GetRunCodeWithCodeId(Db, strconv.Itoa(myCode.CodeId), strconv.Itoa(myCode.UserId), flag)
	if err != nil {
		return
	}
	outputFile, msgFile, err := RunCode(myCode.Path, judgetools.GetMillisecond(), runPath, strconv.Itoa(myRunCode.CodeId), strconv.Itoa(myRunCode.UserId), myRunCode.InputPath, ioPath, codeGoPath)
	if err != nil {
		zapconf.GetWarnLog().Warn("run code err: " + err.Error())
		return
	}
	//把输出文件以及运行信息文件的路径返回来。应该用数据库来接收此信息
	if !db.UpdateRunCodeEndMsg(Db, strconv.Itoa(myRunCode.CodeId), strconv.Itoa(myRunCode.UserId), outputFile, msgFile, flag) {
		zapconf.GetInfoLog().Info("update err last msg")
	}
	cr.DelGoRoutine()
}

func RunCode(codePath string, nowDate string, runPath string, codeId string, userId string, inputPath string, ioPath string, codeGoPath string) (string, string, error) {
	runName := runPath + codeId + nowDate
	runCodeName := runName + ".cpp"
	err := judgetools.ExecCp(codePath, runCodeName)
	if err != nil {
		zapconf.GetWarnLog().Warn("cp error")
		return "", "", errors.New("cp error")
	}
	err = judgetools.ExecCp(environment.CodeJudgePath, codeGoPath)
	if err != nil {
		zapconf.GetWarnLog().Warn("cp error")
		return "", "", errors.New("cp error")
	}
	err = judgetools.ExecGcc(runName, runCodeName)
	if err != nil {
		zapconf.GetWarnLog().Warn("g++ error")
		return "", "", errors.New("g++ " + runCodeName + " error")
	}

	inputData := inputPath //输入数据文件是从数据库读
	outputData := ioPath + codeId + "_" + userId + "_" + nowDate + "output.txt" //输出结果文件还需要构造 是没有的 空的
	msgFile :=  runPath + codeId + "_" + userId + "_" + nowDate + "msg.txt"               //此文件保存运行情况 也是需要构造的
	os.Create(msgFile)

	hostPath := codeGoPath
	rlyRunName := strings.Replace(runName, codeGoPath, "", -1)
	rlyInputData := strings.Replace(inputData, codeGoPath, "", -1)
	rlyOutputData := strings.Replace(outputData, codeGoPath, "", -1)
	rlyMsgFile := strings.Replace(msgFile, codeGoPath, "", -1)
	containerId := dockerexec.GoDocker(hostPath, "/work", "/work", rlyRunName, "15000", "65535", rlyInputData, rlyOutputData, rlyMsgFile, nowDate+codeId+userId+judgetools.GetMillisecond())
	fmt.Println(containerId)

	codeMsg := ""
	startTime := time.Now().Unix()
	//至此上面的运行结果与运行情况都已经获取到了 下面要对两种结果的获取并删除过程文件
	for codeMsg == "" {
		fmt.Println("msgfile path:", msgFile)
		msgF, err := ioutil.ReadFile(msgFile)
		if err != nil {
			zapconf.GetWarnLog().Warn("read failed: " + err.Error())
		}
		fmt.Println("msgf:", msgF)
		codeMsg = string(msgF)
		endTime := time.Now().Unix()
		if endTime-startTime > 60 {
			codeMsg = "{\"status\":\"100\",\"timeUsed\":\"15000\",\"memoryUsed\":\"0\"}"
			break
		}
		time.Sleep(time.Duration(1) * time.Second)
	}

	res := db.CodeStatusMsg{}
	fmt.Println("msgf string", codeMsg)
	fmt.Println("unmarshal")
	if err := json.Unmarshal([]byte(codeMsg), &res); err != nil {
		panic(err)
	}
	fmt.Println("output:")

	fmt.Println("res --")
	fmt.Println(res)

	fmt.Println("string(outputdataF) --")

	fmt.Println("codemsg --", codeMsg)
	fmt.Println("res in  --")
	fmt.Println(res.Status)
	fmt.Println(res.TimeUsed + "ms")
	fmt.Println(res.MemoryUsed + "kb")
	dockerexec.StopDocker(containerId)
	dockerexec.DelDocker(containerId)
	judgetools.FastDelFile(runName)
	judgetools.FastDelFile(runCodeName)
	return outputData, msgFile, nil
}
