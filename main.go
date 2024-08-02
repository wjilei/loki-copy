package main

import (
	"cmp"
	"errors"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/viper"
)

func main() {
	Init()
	sigchan := make(chan os.Signal, 1)
	stop := false
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigchan
		log.Println("Exiting...")
		stop = true
	}()
	for {
		if stop {
			CloseDb()
			break
		}
		lokiCliSrc := NewLokiClient(viper.GetString("loki-source"))
		lokiCliDest := NewLokiClient(viper.GetString("loki-destination"))
		queries := viper.Get("queries")
		for _, v := range queries.([]interface{}) {
			query := v.(string)
			start, err := GetQueryPos(query)
			if err != nil {
				start = time.Now().Add(time.Second * -30).UnixNano()
			}
			log.Printf("Query: %s, Start: %s\n", query, time.Unix(start/1e9, start%1e9).Format("2006-01-02 15:04:05.999999999"))
			res, err := lokiCliSrc.QueryRange(query, start, time.Now().UnixNano())
			if err != nil {
				log.Println(err)
				continue
			}

			pos, err := getNewReadPos(res)
			if err != nil {
				log.Println(err)
				continue
			}

			pushReq := PushRequest{
				Streams: res.Data.Result,
			}

			err = lokiCliDest.Push(&pushReq)
			if err != nil {
				log.Println("Push Error:", err)
				continue
			}

			if pos > start {
				err = SetQueryPos(query, pos)
				if err != nil {
					log.Println(err)
				}
			}
		}
		time.Sleep(time.Second * 3)
	}
	log.Println("Exited.")
}

func getNewReadPos(res *QueryResult) (int64, error) {
	if len(res.Data.Result) == 0 {
		return 0, errors.New("data result is empty")
	}
	var posStr string
	for _, result := range res.Data.Result {
		for _, v := range result.Values {
			if cmp.Compare(v[0], posStr) > 0 {
				posStr = v[0]
			}
		}
	}
	if posStr == "" {
		return 0, errors.New("no position found")
	}
	pos, err := strconv.ParseInt(posStr, 10, 64)
	if err != nil {
		return 0, errors.New("position format error")
	}
	return pos, nil
}

func Init() {
	err := ConfigInit()
	if err != nil {
		panic(err)
	}
	err = InitDb()
	if err != nil {
		panic(err)
	}
}
