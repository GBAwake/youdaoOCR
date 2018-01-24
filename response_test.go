package ydocr_test

import (
	"testing"
	"mytest/ydocr"
	"fmt"
)

func TestResponseToLine(t *testing.T) {
	r := []byte(`{"errorCode":"0","Result":{"orientation":"UP","regions":[{"boundingBox":"0,173,157,173,157,207,0,207","lines":[{"boundingBox":"0,173,157,173,157,207,0,207","words":[{"boundingBox":"21,169,52,169,52,212,21,212","word":"道"},{"boundingBox":"63,169,84,169,84,212,63,212","word":"云"},{"boundingBox":"94,169,115,169,115,212,94,212","word":"笔"},{"boundingBox":"126,169,157,169,157,212,126,212","word":"记"}],"text":"道云笔记"}]}]}}`)

	w, err := ydocr.Response2Line(r)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", *w)
}

func TestResponseToWord(t *testing.T) {
	r := []byte(`{
	"errorCode": "0",
	"Result": {
		"orientation": "Up",
		"textAngle": 0,
		"language": "en",
		"lines": [{
			"boundingBox": "30,33,25,10",
			"words": "hello"
		}]
	}
	}`)
	w, err := ydocr.Response2Word(r)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", *w)

}
