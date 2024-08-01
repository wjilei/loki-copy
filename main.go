package main

import (
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

func main() {
	err := ConfigInit()
	if err != nil {
		panic(err)
	}
	err = InitDb()
	if err != nil {
		panic(err)
	}
	defer CloseDb()
	lokiCli := NewLokiClient("http://iotgw.poersmart.com:3100")
	queries := viper.Get("queries")
	start := time.Now().Add(time.Minute * -30).UTC().UnixNano()
	end := time.Now().UTC().UnixNano()
	for _, v := range queries.([]interface{}) {
		query := v.(string)
		log.Println(query)
		res, err := lokiCli.QueryRange(query, start, end)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println(res)
	}
}
