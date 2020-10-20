package main

import (
	"database/sql"
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
		zapconf.GetWarnLog().Warn("conn database error: " + err.Error())
		return
	}
	if err = Db.Ping(); err != nil {
		zapconf.GetWarnLog().Warn("conn database error: " + err.Error())
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