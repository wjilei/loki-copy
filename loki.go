package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"time"

	"github.com/imroc/req/v3"
)

type LokiClient struct {
	client  *req.Client
	baseUrl string
}

func NewLokiClient(baseUrl string) *LokiClient {
	c := req.NewClient()
	c.SetTimeout(10 * time.Second)
	return &LokiClient{
		client:  c,
		baseUrl: baseUrl,
	}
}

type QueryResult struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string       `json:"resultType"`
		Result     []StreamData `json:"result"`
	} `json:"data"`
}

func (c *LokiClient) QueryRange(query string, start int64, end int64) (*QueryResult, error) {
	var result QueryResult
	r, err := c.client.R().
		SetQueryParamsAnyType(map[string]interface{}{
			"query":     query,
			"start":     start,
			"end":       end,
			"direction": "forward",
			"limit":     1000,
		}).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetSuccessResult(&result).
		Get(c.baseUrl + "/loki/api/v1/query_range")
	if err != nil {
		return nil, err
	}
	if !r.IsSuccessState() {
		return nil, errors.New("query failed")
	}
	return &result, nil
}

type StreamData struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

type PushRequest struct {
	Streams []StreamData `json:"streams"`
}

var ErrorPushError = errors.New("push error")

func (c *LokiClient) Push(req *PushRequest) error {
	b, err := c.gzipReq(req)
	if err != nil {
		return err
	}
	r, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(b).
		Post(c.baseUrl + "/loki/api/v1/push")
	if err != nil {
		// log.Println("push error:", err)
		return err
	}
	if r.StatusCode == 200 || r.StatusCode == 204 {
		return nil
	}
	return errors.New(r.String())
}

func (c *LokiClient) gzipReq(req *PushRequest) ([]byte, error) {
	var buf []byte
	buffer := bytes.NewBuffer(buf)
	w := gzip.NewWriter(buffer)
	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	w.Write(b)
	w.Close()

	return buffer.Bytes(), nil
}
