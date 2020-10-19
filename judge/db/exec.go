package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"judge/controlroutine"
	"judge/judgetools"
	_ "github.com/go-sql-driver/mysql"
)

func GetCompareId(conn *sql.DB, cr *controlroutine.ChanRoutine) (string, error) {
	var cid string
	sqlData := "select compare_id from run_compare where version = 0 ORDER BY compare_id limit 1 for update"
	err := conn.QueryRow(sqlData).Scan(&cid)
	if err != nil {
		return "", errors.New("no record need run")
	}
	fmt.Println(cid)

	cr.AddGoRoutine()
	if !updateRunCompareVersion(conn, cid) {
		panic("update version error")
	}
	return cid, nil
}

func updateRunCompareVersion(conn *sql.DB, cid string) bool {
	sqlData := "update run_compare set version = 1 where compare_id = " + cid
	fmt.Println(sqlData)
	_, upErr := conn.Exec(sqlData)
	if upErr != nil {
		fmt.Println("update error:", upErr)
		return false
	}
	return true
}

func GetEvaCompareWithCompareId(conn *sql.DB, cid string) (*EvaCompare, error) {
	takeCodeIdSql := "select * from eva_compare where compare_id = " + cid
	fmt.Println(takeCodeIdSql)
	myCompare := &EvaCompare{}
	qErr := conn.QueryRow(takeCodeIdSql).Scan(&myCompare.CompareId, &myCompare.CompareName, &myCompare.UserId, &myCompare.FirstCodeId, &myCompare.SecondCodeId, &myCompare.InputDataPath, &myCompare.MaxInputGroup, &myCompare.CreateTime, &myCompare.UpdateTime, &myCompare.Remarks)
	if qErr != nil {
		fmt.Println("select eva_compare error")
		fmt.Println(qErr)
		return &EvaCompare{}, errors.New("select eva_compare error")
	}
	fmt.Println(myCompare)
	return myCompare, nil
}

func GetEvaCodePathWithCodeId(conn *sql.DB, cid string) (string, error) {
	takeCodePathSql := "select path from eva_code where code_id = '" + cid + "'"
	path := ""
	fcpErr := conn.QueryRow(takeCodePathSql).Scan(&path)
	if fcpErr != nil {
		fmt.Println("select code path error: ", fcpErr)
		return "", fcpErr
	}
	return path, nil
}

func InsertData(firstCodeMsg CodeRunData, secondCodeMsg CodeRunData, compareId string, i string, maxx string, ioPath string, Db *sql.DB) bool {
	flag := false
	selectData := "select compare_data from run_compare where compare_id = " + compareId
	selectDataString := ""
	Db.QueryRow(selectData).Scan(&selectDataString)
	fmt.Println(selectDataString)
	tmp := make(map[string]CodeMsg)
	json.Unmarshal([]byte(selectDataString), &tmp)
	fmt.Println(tmp)
	fmt.Println("--------")
	makeNew := make(map[string]string)
	makeNew["first_runTime"] = firstCodeMsg.TimeUsed
	makeNew["first_runMemory"] = firstCodeMsg.MemoryUsed
	makeNew["second_runTime"] = secondCodeMsg.TimeUsed
	makeNew["second_runMemory"] = secondCodeMsg.MemoryUsed
	fmt.Println("make new : ", makeNew)
	newAddMsg := CodeMsg{}
	if !judgetools.Compare(firstCodeMsg.Output, secondCodeMsg.Output) {
		makeNew["input_path"] = ioPath + i + "input.txt"
		newAddMsg.Input_path = makeNew["input_path"]
		flag = true
	}
	fmt.Println(makeNew)
	fmt.Println("-------- add new")
	newAdd, _ := json.Marshal(makeNew)
	fmt.Println(string(newAdd))
	fmt.Println(i)
	newAddMsg.First_runTime = firstCodeMsg.TimeUsed
	newAddMsg.First_runMemory = firstCodeMsg.MemoryUsed
	newAddMsg.Second_runTime = secondCodeMsg.TimeUsed
	newAddMsg.Second_runMemory = secondCodeMsg.MemoryUsed
	tmp[i] = newAddMsg
	fmt.Println(tmp)
	fmt.Println("--------")
	insData, _ := json.Marshal(tmp)
	fmt.Println(string(insData))
	addStr := ""
	if i == maxx || flag == true {
		addStr = " , version=2"
	}
	updateSql := "update run_compare set compare_data='" + string(insData) + "'" + addStr + " where compare_id=" + compareId
	fmt.Println(updateSql)

	//把这条数据存入 或者判断不对等就flag=true 跳过后面的对拍
	_, inErr := Db.Exec(updateSql)
	if inErr != nil {
		fmt.Println("update error")
		panic(inErr)
	}

	return flag
}

func GetCodeId(conn *sql.DB, cr *controlroutine.ChanRoutine) (string, string, error) {
	var cid string
	var uid string
	var flag string
	sqlData := "select code_id, user_id, select_flag from run_code where version = 0 ORDER BY code_id limit 1 for update"
	err := conn.QueryRow(sqlData).Scan(&cid, &uid, &flag)
	if err != nil {
		return "", "", errors.New("no record need run")
	}
	cr.AddGoRoutine()
	if !updateRunCodeVersion(conn, cid, uid, flag) {
		panic("update version error")
	}
	return cid, flag, nil
}

func updateRunCodeVersion(conn *sql.DB, cid string, uid string, flag string) bool {
	sqlData := "update run_code set version = 1 where code_id = " + cid + " and user_id = " + uid + " and select_flag = '" + flag + "'"
	_, upErr := conn.Exec(sqlData)
	if upErr != nil {
		fmt.Println("update error:", upErr)
		return false
	}
	return true
}

func GetEvaCodeWithCodeId(conn *sql.DB, cid string) (*EvaCode, error) {
	takeCodeIdSql := "select * from eva_code where code_id = " + cid
	fmt.Println(takeCodeIdSql)
	myCode := &EvaCode{}
	qErr := conn.QueryRow(takeCodeIdSql).Scan(&myCode.CodeId, &myCode.CodeName, &myCode.UserId, &myCode.CodeText, &myCode.Path, &myCode.CreateTime, &myCode.UpdateTime)
	if qErr != nil {
		fmt.Println("select eva_code error")
		fmt.Println(qErr)
		return &EvaCode{}, errors.New("select eva_code error")
	}
	fmt.Println(myCode)
	return myCode, nil
}

func GetRunCodeWithCodeId(conn *sql.DB, cid string, uid string, flag string) (*runCode, error) {
	takeRunCodeSql := "select code_id, user_id, input_path from run_code where code_id = " + cid + " and user_id = " + uid + " and select_flag = '" + flag + "'"
	fmt.Println(takeRunCodeSql)
	myRunCode := &runCode{}
	rcErr := conn.QueryRow(takeRunCodeSql).Scan(&myRunCode.CodeId, &myRunCode.UserId, &myRunCode.InputPath)
	if rcErr != nil {
		fmt.Println("select run_code error")
		fmt.Println(rcErr)
		return &runCode{}, errors.New("select run_code error")
	}
	return myRunCode, nil
}

func UpdateRunCodeEndMsg(conn *sql.DB, code_id string, user_id string, outputFile string, msgFile string, flag string) (bool) {
	sqlData := "update run_code set run_data='" + outputFile + "', msg_data='" + msgFile + "' where code_id = " + code_id + " and user_id = " + user_id + " and select_flag = '" + flag + "'"
	fmt.Println(sqlData)
	//defer conn.Close()
	_, upErr := conn.Exec(sqlData)
	if upErr != nil {
		fmt.Println("update error: ", upErr)
		return false
	}
	return true
}