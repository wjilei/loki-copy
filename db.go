package main

import (
	"errors"

	_ "github.com/logoove/sqlite"
	"github.com/spf13/viper"
	"xorm.io/xorm"
)

var engine *xorm.Engine

func InitDb() error {
	var err error
	engine, err = xorm.NewEngine("sqlite", viper.GetString("db"))
	if err != nil {
		return err
	}
	err = engine.Ping()
	if err != nil {
		return err
	}
	engine.ShowSQL(viper.GetBool("show-sql"))
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

func GetQueryPos(query string) (int64, error) {
	pos := Position{Query: query}
	has, err := engine.Get(&pos)
	if err != nil {
		return 0, err
	}
	if !has {
		return 0, errors.New("no such record")
	}
	return pos.Pos, nil
}

func SetQueryPos(query string, pos int64) error {
	posDb := Position{Query: query}
	exist, err := engine.Exist(&posDb)
	if err != nil {
		return err
	}
	if !exist {
		_, err = engine.Insert(&Position{Query: query, Pos: pos})
	} else {
		_, err = engine.Where("query = ?", query).Update(&Position{Pos: pos})
	}

	return err
}
