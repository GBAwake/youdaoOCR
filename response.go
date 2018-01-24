package ydocr

import (
	"encoding/json"
	"fmt"
)

type Word struct {
	ErrorCode string `json:"errorCode"`
	Result    struct {
		Orientation string `json:"orientation"`
		TextAngle   int    `json:"textAngle"`
		Language    string `json:"language"`
		Lines       []struct {
			BoundingBox string `json:"boundingBox"`
			Words       string `json:"words"`
		} `json:"lines"`
	} `json:"Result"`
}

type Line struct {
	ErrorCode string `json:"errorCode"`
	Result    struct {
		Orientation string `json:"orientation"`
		Regions     []struct {
			BoundingBox string `json:"boundingBox"`
			Lines       []struct {
				BoundingBox string `json:"boundingBox"`
				Words       []struct {
					BoundingBox string `json:"boundingBox"`
					Word        string `json:"word"`
				} `json:"words"`
				Text string `json:"text"`
			} `json:"lines"`
		} `json:"regions"`
	} `json:"Result"`
}

func Response2Word(r []byte) (*Word, error){
	v := &Word{}
	err := json.Unmarshal(r, v)
	if err != nil {
		return nil, fmt.Errorf("response to word json error: %s", err)
	}
	return v, nil
}

func Response2Line(r []byte) (*Line, error){
	v := &Line{}
	err := json.Unmarshal(r, v)
	if err != nil {
		return nil, fmt.Errorf("response to line json error: %s", err)
	}
	return v, nil
}
