package timezone

import (
	"encoding/json"
	"fmt"
)

// 表示全体
type OutputRoot struct {
	Jobnetworks []*OutputJobNet `json:"jobnetworks"`
}

// 表示用のジョブネットワーク構造体
type OutputJobNet struct {
	Id         int          `json:"id"`
	Jobnetwork string       `json:"jobnetwork"`
	StartDate  string       `json:"startdate"`
	EndDate    string       `json:"enddate"`
	Status     int          `json:"status"`
	Detail     string       `json:"detail"`
	CreateDate string       `json:"createdate"`
	UpdateDate string       `json:"updatedate"`
	Jobs       []*OutputJob `json:"jobs"`
}

// 表示用のジョブ構造体
type OutputJob struct {
	JobId      string `json:"jobid"`
	Jobname    string `json:"jobname"`
	StartDate  string `json:"startdate"`
	EndDate    string `json:"enddate"`
	Status     int    `json:"status"`
	Detail     string `json:"detail"`
	Rc         int    `json:"rc"`
	Node       string `json:"node`
	Port       int    `json:"port"`
	Variable   string `json:"variable"`
	CreateDate string `json:"createdate"`
	UpdateDate string `json:"updatedate"`
}

func parseShowOutput(showOutput string) (*OutputRoot, error) {
	jsData := new(OutputRoot)
	err := json.Unmarshal([]byte(showOutput), jsData)
	if err != nil {
		return nil, fmt.Errorf("Json Parse failed. : %v", err)
	}
	if len(jsData.Jobnetworks) == 0 {
		return nil, fmt.Errorf("Output of show utility is empty.")
	}

	return jsData, nil
}
