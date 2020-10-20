package core

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"judge/zapconf"
	"os"
	"judge/db"
	"judge/environment"
	"judge/judgetools"
	"strconv"
	"time"
)

//生成随机输入数据的编译
/*运行代码路径，io路径，时间戳，当前对拍的root路径，对拍的各项信息
 */
func GenerateInputData(runPath string, ioPath string, nowDate string, runRootPath string, compareData *db.EvaCompare) {
	maxInputGroup := compareData.MaxInputGroup
	runName := runPath + strconv.Itoa(compareData.CompareId) + nowDate
	runCodeName := runPath + strconv.Itoa(compareData.CompareId) + nowDate + ".cpp"
	fmt.Println(compareData)
	err := judgetools.ExecCp(compareData.InputDataPath, runCodeName)
	if err != nil {
		zapconf.GetWarnLog().Warn("cp error")
		return
	}
	err = judgetools.ExecCp(environment.CompareJudgePath, runRootPath)
	if err != nil {
		zapconf.GetWarnLog().Warn("cp error")
		return
	}
	err = judgetools.ExecGcc(runName, runCodeName)
	if err != nil {
		zapconf.GetWarnLog().Warn("g++ error")
		return
	}

	for i := 1; i <= maxInputGroup; i++ {
		checkNum := 0
		/*这里也有重试机制，若运行失败有十次的重试，还不行就panic*/
		if !GenerateInputRun(runName, strconv.Itoa(i), ioPath, nowDate, runRootPath, nowDate+strconv.Itoa(compareData.CompareId)+strconv.Itoa(i)) {
			i--
			checkNum++
			if checkNum > 10 {
				zapconf.GetErrorLog().Error("cant GenerateInputRun right")
				panic("cant GenerateInputRun right")
			}
		} //生成随机输入数据的运行
	}

	judgetools.FastDelFile(runCodeName)
	judgetools.FastDelFile(runName)
}


//生成随机输入数据的运行 其中这个里面的输出其实是为了给真正跑的代码做输入数据
func GenerateInputRun(runName string, inputData string, ioPath string, nowDate string, runRootPath string, dockerName string) bool { //对参数做个解释：runname顾名思义就是运行生成随机输入数据的可执行文件的名字，inputdata就是为了给循环的各个文件做个区分同时将组数传递给随机生成代码作为参数，iopath就是输入输出文件的地方，nowdate也为了做个区分，唯一性
	inputDataName := ioPath + "generate" + inputData + nowDate//io的根路径
	inputDataPath := inputDataName + "input.txt"//输入数据文件
	outputDataPath := inputDataName + "output.txt"//输出数据文件
	msgPath := runName + "msg.txt"
	fmt.Println("generatemsg: ===========", msgPath)
	file, err := os.OpenFile(inputDataPath, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		zapconf.GetWarnLog().Warn("open " + inputDataPath + "error")
		panic(err)
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(inputData)
	write.Flush()

	//用户可执行文件，输入数据文件，输出数据文件，运行信息路径，运行这id的root路径，容器名
	Run(runName, inputDataPath, outputDataPath, msgPath, runRootPath, dockerName) //跑代码的核心函数，这一步要跑代码得到随机生成的输入数据

	/*有十秒的等待期等待docker跑完，完了后读取数据放入到之后 对拍要用的输入文件中*/
	startTime := time.Now().Unix()
	for !judgetools.IsExist(outputDataPath) {
		endTime := time.Now().Unix()
		if endTime-startTime > 10 {
			return false
		}
		time.Sleep(time.Duration(1) * time.Second)
		fmt.Println("have no outputpath wait...")
	}
	myOutputFile, err := os.OpenFile(outputDataPath, os.O_RDONLY, 0777)
	if err != nil {
		zapconf.GetWarnLog().Warn("open " + outputDataPath + " error")
		panic(err)
	}
	defer myOutputFile.Close()
	codeInputData, err := ioutil.ReadAll(myOutputFile)
	if err != nil {
		zapconf.GetWarnLog().Warn("read file " + outputDataPath + " error")
		return false
	}
	judgetools.FastDelFile(inputDataPath)
	judgetools.FastDelFile(outputDataPath)
	judgetools.FastDelFile(msgPath)

	rlyInputPath := ioPath + inputData + "input.txt"
	fmt.Println("rlyinputpath", rlyInputPath)
	rlyInputFile, err := os.OpenFile(rlyInputPath, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		zapconf.GetWarnLog().Warn("open " + rlyInputPath + " error")
		return false
	}
	defer rlyInputFile.Close()
	writer := bufio.NewWriter(rlyInputFile)
	writer.WriteString(string(codeInputData))
	writer.Flush()
	return true
}

