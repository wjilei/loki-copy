package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"xorm.io/xorm"
)

var engine *xorm.Engine

func InitDb() error {
	var err error
	engine, err = xorm.NewEngine("sqlite3", viper.GetString("db"))
	if err != nil {
		return err
	}
	err = engine.Ping()
	if err != nil {
		return err
	}
	engine.ShowSQL(true)
	engine.Sync2(new(Position))
	return nil
}

func CloseDb() {
	if engine != nil {
		engine.Close()
	}
}

type Position struct {
	Id    int64  `xorm:"pk autoincr"`
	Query string `xorm:"index"`
	Pos   int64
}
