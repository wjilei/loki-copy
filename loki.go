package main

import (
	"log"

	"github.com/imroc/req/v3"
)

type LokiClient struct {
	client  *req.Client
	baseUrl string
}

func NewLokiClient(baseUrl string) *LokiClient {
	return &LokiClient{
		client:  req.NewClient(),
		baseUrl: baseUrl,
	}
}

type QueryResult struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Stream map[string]string `json:"stream"`
			Values [][]string        `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

func (c *LokiClient) QueryRange(query string, start int64, end int64) (*QueryResult, error) {
	var result QueryResult

	log.Println("QueryRange", query, start, end)
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
	if r.StatusCode != 200 {
		return nil, err
	}
	return &result, nil
}
