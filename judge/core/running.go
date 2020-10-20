package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"judge/db"
	"judge/dockerexec"
	"judge/environment"
	"judge/judgetools"
	"judge/zapconf"
	"os"
	"strings"
	"time"
)

func Run(runName string, inputData string, outputData string, msgPath string, hostPath string, dockerName string) { //运行函数的核心代码，参数解释：runname显然是运行代码的可执行文件名称，inputdata是存放输入数据的文件路径
	rlyRunName := strings.Replace(runName, hostPath, "", -1)
	rlyInputData := strings.Replace(inputData, hostPath, "", -1)
	rlyOutputData := strings.Replace(outputData, hostPath, "", -1)
	rlyMsgFile := strings.Replace(msgPath, hostPath, "", -1)
	fmt.Println(hostPath)
	fmt.Println(rlyRunName)
	fmt.Println(rlyInputData)
	fmt.Println(rlyOutputData)
	/*
		hostpath是挂载主机的路径
		containerpath是挂载到容器中路径
		workdir是容器中工作路径
		rlyrunname是可执行文件路径 应该对应到挂载主机路径下相对路径
		rlyinputdata 输入文件路径
		rlyoutputdata 输出文件路径
		msg在这里不需要其实 因为这里是构造输入数据
	*/
	containerId := dockerexec.GoDocker(hostPath, "/work", "/work", rlyRunName, "15000", "65535", rlyInputData, rlyOutputData, rlyMsgFile, dockerName)
	dockerexec.StopDocker(containerId)
	dockerexec.DelDocker(containerId)
}

/*代码文件存放路径，时间戳，代码运行路径，代码id，输入文件路径，输出文件路径，运行信息文件路径，代码运行root路径(挂载主机路径)
*/
func RunCodeGetResult(codePath string, nowDate string, runPath string, codeId string, inputPath string, outputPath string, msgPath string, runRootPath string, num string) (string, string, string, string, string) {
	runName := runPath + codeId + nowDate
	runCodeName := runName + ".cpp"
	err := judgetools.ExecCp(codePath, runCodeName)
	if err != nil {
		zapconf.GetWarnLog().Warn("cp error")
		return "", "", "", "", ""
	}
	err = judgetools.ExecCp(environment.CompareJudgePath, runRootPath)
	if err != nil {
		zapconf.GetWarnLog().Warn("cp error")
		return "", "", "", "", ""
	}
	err = judgetools.ExecGcc(runName, runCodeName)
	if err != nil {
		zapconf.GetWarnLog().Warn("g++ error")
		return "", "", "", "", ""
	}

	os.Create(msgPath)
	hostPath := runRootPath
	rlyRunName := strings.Replace(runName, runRootPath, "", -1)
	rlyInputData := strings.Replace(inputPath, runRootPath, "", -1)
	rlyOutputData := strings.Replace(outputPath, runRootPath, "", -1)
	rlyMsgFile := strings.Replace(msgPath, runRootPath, "", -1)
	dockerName := ""
	if !(num == "0") {
		dockerName = nowDate+codeId+num+judgetools.GetMillisecond()
	} else {
		dockerName = nowDate+codeId+judgetools.GetMillisecond()
	}
	containerId := dockerexec.GoDocker(hostPath, "/work", "/work", rlyRunName, "15000", "65535", rlyInputData, rlyOutputData, rlyMsgFile, dockerName)
	fmt.Println(containerId)

	codeMsg := ""
	startTime := time.Now().Unix()
	//至此上面的运行结果与运行情况都已经获取到了 下面要对两种结果的获取并删除过程文件
	for codeMsg == "" {
		fmt.Println("msgfile path:", msgPath)
		msgF, err := ioutil.ReadFile(msgPath)
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
	println("unmarshal")
	if err := json.Unmarshal([]byte(codeMsg), &res); err != nil {
		panic(err)
		return "", "", "", "", ""
	}
	fmt.Println("output:")

	fmt.Println("res --")
	fmt.Println(res)

	fmt.Println("string(outputdataF) --")
	//fmt.Println(string(outputdataF))

	fmt.Println("codemsg --", codeMsg)

	fmt.Println("res in  --")
	fmt.Println(string(res.Status))
	fmt.Println(string(res.TimeUsed + "ms"))
	fmt.Println(string(res.MemoryUsed + "kb"))
	dockerexec.StopDocker(containerId)
	dockerexec.DelDocker(containerId)
	judgetools.FastDelFile(runName)
	judgetools.FastDelFile(runCodeName)
	return outputPath, res.Status, res.TimeUsed+"ms", res.MemoryUsed+"kb", "ok"
}
