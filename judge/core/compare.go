package core

import (
	"database/sql"
	"fmt"
	"judge/controlroutine"
	"judge/db"
	"judge/environment"
	"judge/judgetools"
	"judge/zapconf"
	"strconv"
)

func StartCompare(Db *sql.DB, cid string, cr *controlroutine.ChanRoutine) {
	comparePath := environment.CompareRootPath + cid + "/" //组合对拍的运行环境
	runPath := comparePath + ServerRunCodePath
	ioPath := comparePath + ServerIOPath
	judgetools.FastCreateFile(comparePath)
	judgetools.FastCreateFile(runPath)
	judgetools.FastCreateFile(ioPath)

	myCompare, err := db.GetEvaCompareWithCompareId(Db, cid)
	if err != nil {
		return
	}
	/*运行代码路径，io路径，时间戳，当前对拍的root路径，对拍的各项信息
	*/
	GenerateInputData(runPath, ioPath, judgetools.GetMillisecond(), comparePath, myCompare)

	//随机输入数据的构造已完成！下面是让两个代码在各输入数据下运行的结果
	FirstCodePath, err := db.GetEvaCodePathWithCodeId(Db, strconv.Itoa(myCompare.FirstCodeId))
	if err != nil {
		zapconf.GetWarnLog().Warn("select first code err")
		return
	}
	SecondCodePath, err := db.GetEvaCodePathWithCodeId(Db, strconv.Itoa(myCompare.SecondCodeId))
	if err != nil {
		zapconf.GetWarnLog().Warn("select second code err")
		return
	}

	flag := false
	for i := 1; i <= myCompare.MaxInputGroup; i++ {
		if flag == false {
			// 跑第一个代码在第i项的结果，返回包括运行结果/运行时间/内存消耗
			fmt.Printf("\n---%d first---\n", i)

			firstCodeData := db.CodeRunData{}
			firstCheck := ""
			firstInputFile := ioPath + strconv.Itoa(i) + "input.txt"            //输入数据文件是由之前构造的时候就弄好的，是已经有的
			firstOutputFile := ioPath + strconv.Itoa(i) + strconv.Itoa(myCompare.FirstCodeId) + "output.txt" //输出结果文件是没有的 空的
			firstMsgFile := runPath + strconv.Itoa(i) + "msg.txt"               //此文件保存运行情况
			firstCodeData.Output, firstCodeData.Status, firstCodeData.TimeUsed, firstCodeData.MemoryUsed, firstCheck = RunCodeGetResult(FirstCodePath, judgetools.GetMillisecond(), runPath, strconv.Itoa(i)+cid, firstInputFile, firstOutputFile, firstMsgFile, comparePath, strconv.Itoa(i))//在codeid那个地方要codeid+i
			if firstCheck != "ok" {
				zapconf.GetWarnLog().Warn("run code docker error")
				return
			}

			// 跑第一个代码在第i项的结果，返回包括运行结果/运行时间/内存消耗
			fmt.Printf("\n---%d second---\n", i)

			secondCodeData := db.CodeRunData{}
			secondCheck := ""
			secondInputFile := ioPath + strconv.Itoa(i) + "input.txt"            //输入数据文件是由之前构造的时候就弄好的，是已经有的
			secondOutputFile := ioPath + strconv.Itoa(i) + strconv.Itoa(myCompare.SecondCodeId) + "output.txt" //输出结果文件是没有的 空的
			secondMsgFile := runPath + strconv.Itoa(i) + "msg.txt"               //此文件保存运行情况
			secondCodeData.Output, secondCodeData.Status, secondCodeData.TimeUsed, secondCodeData.MemoryUsed, secondCheck = RunCodeGetResult(SecondCodePath, judgetools.GetMillisecond(), runPath, strconv.Itoa(i)+cid, secondInputFile, secondOutputFile, secondMsgFile, comparePath, strconv.Itoa(i))//在codeid那个地方要codeid+i
			if secondCheck != "ok" {
				zapconf.GetWarnLog().Warn("run code docker error")
				return
			}

			// 把这些数据加上之前的数据插入数据库，如果insert返回true了（即两代码结果不同），那么说明对拍失败没必要继续进行下去了，可以把flag置为true
			if db.InsertData(firstCodeData, secondCodeData, strconv.Itoa(myCompare.CompareId), strconv.Itoa(i), strconv.Itoa(myCompare.MaxInputGroup), ioPath, Db) {
				flag = true
				continue
			}
			fmt.Println("msgfile : ======>", firstMsgFile, secondMsgFile)
		}
		// 删除第i次运行过程中所产生的io过程文件 包括构造的随机输入数据
		inputData := ioPath + strconv.Itoa(i) + "input.txt"                                           //输入数据文件是由之前构造的时候就弄好的，是已经有的
		outputData1 := ioPath + strconv.Itoa(i) + strconv.Itoa(myCompare.FirstCodeId) + "output.txt"  //输出结果文件是没有的 空的
		outputData2 := ioPath + strconv.Itoa(i) + strconv.Itoa(myCompare.SecondCodeId) + "output.txt" //输出结果文件是没有的 空的
		msgFile := runPath + strconv.Itoa(i) + "msg.txt"                                              //此文件保存运行情况
		judgetools.FastDelFile(inputData)
		judgetools.FastDelFile(outputData1)
		judgetools.FastDelFile(outputData2)
		judgetools.FastDelFile(msgFile)
	}

	cr.DelGoRoutine()
}