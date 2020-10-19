package main

import (
	"database/sql"
	"fmt"
	"judge/controlroutine"
	"judge/core"
	"judge/db"
)

func main() {
	var Db *sql.DB
	var err error
	Db, err = sql.Open("mysql", db.DbMysql)
	defer Db.Close()
	if err != nil {
		fmt.Println("conn database error: ", err)
		return
	}

	chanRoutine := controlroutine.NewChanRoutine(2)

	for {
		cid, flag, err := db.GetCodeId(Db, chanRoutine) //查到未运行的代码
		if err != nil {
		} else {
			go core.StartRunCode(Db, cid, flag, chanRoutine)
		}
	}
}