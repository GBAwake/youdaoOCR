package ydocr

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"crypto/md5"
	"math/rand"
	"strconv"
	"io"
	"net/url"
	"net/http"
	"io/ioutil"
	"os"
	"image"
	"image/jpeg"
)

type Ocr struct {
	appKey string
	appSecret string
	url string
	DetecType string
	ImageType string
	LangType string
}

func NewOcr(appKey, appSecret string) *Ocr{
	// 默认的ocr 使用 按行识别， 使用base64编码， 中英混合
	o := &Ocr{
		appKey:appKey,
		appSecret:appSecret,
		url:"http://openapi.youdao.com/ocrapi",
		DetecType:"10012",
		ImageType:"1",
		LangType:"zh-en",
	}
	return o
}

func (o *Ocr) FileOcr(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	f, err := ioutil.ReadAll(file)
	if err  != nil {
		return nil, err
	}
	return o.BasePost(f)
}

func (o *Ocr) ImageOcr(img *image.Image) ([]byte, error) {
	buff := new(bytes.Buffer)
	err := jpeg.Encode(buff, *img, nil)
	if err != nil {
		return nil, fmt.Errorf("jpeg image err: %s", err)
	}
	return o.BasePost(buff.Bytes())
}



func (o *Ocr) BasePost(img []byte) ([]byte, error){
	buff := new(bytes.Buffer)
	err := mbase64Encode(img, buff)
	if err != nil {
		return nil, err
	}
	imgEncode := buff.String()
	sign, salt, err := countSign(o.appKey, o.appSecret, imgEncode)
	if err != nil {
		return nil, err
	}

	v := url.Values{
		"appKey": []string{o.appKey},
		"img": []string{imgEncode},
		"detectType": []string{o.DetecType},
		"imageType": []string{o.ImageType},
		"langType": []string{o.LangType},
		"salt": []string{salt},
		"sign": []string{sign},
	}
	resp, err := http.PostForm(o.url, v)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func mbase64Encode(src []byte, dst *bytes.Buffer) error {
	en := base64.NewEncoder(base64.StdEncoding, dst)
	defer en.Close()
	lenSrc := len(src)
	for i:= 0; i < lenSrc; {
		n, err := en.Write(src[i:])
		if err != nil {
			return fmt.Errorf("base64 encode err: %s", err)
		}
		i += n
	}
	return nil
}

func mCountMd5(src []byte) string{
	return fmt.Sprintf("%x", md5.Sum(src))
}

func countSign(appKey, appSecret, img string) (sign, salt string, e error) {
	// 生成随机salt
	salt = strconv.Itoa(rand.Intn(65535) + 1)
	c := md5.New()
	_, err := io.WriteString(c, appKey)
	if err != nil {
		e = fmt.Errorf("write to md5 err: %s ", err)
		return
	}
	_, err = io.WriteString(c, img)
	if err != nil {
		e = fmt.Errorf("write to md5 err: %s ", err)
		return
	}
	_, err = io.WriteString(c, salt)
	if err != nil {
		e = fmt.Errorf("write to md5 err: %s ", err)
		return
	}
	_, err = io.WriteString(c, appSecret)
	if err != nil {
		e = fmt.Errorf("write to md5 err: %s ", err)
		return
	}
	sign = fmt.Sprintf("%x", c.Sum(nil))
	return
}


