package main

import (
	"database/sql"
	"fmt"
	"judge/controlroutine"
	"judge/core"
	"judge/db"
	"judge/zapconf"
)

func main() {
	var Db *sql.DB
	var err error

	Db, err = sql.Open("mysql", db.DbMysql)
	defer Db.Close()
	if err != nil {
		zapconf.GetErrorLog().Error("conn database error: " + err.Error())
		return
	}
	if err = Db.Ping(); err != nil {
		zapconf.GetWarnLog().Warn("conn database error: " + err.Error())
		return
	}

	chanRoutine := controlroutine.NewChanRoutine(2)

	for {
		cid, err := db.GetCompareId(Db, chanRoutine) //查到未运行的对拍
		if err != nil {
		} else {
			fmt.Println("lets go : ", cid)
			go core.StartCompare(Db, cid, chanRoutine)
		}
	}
}