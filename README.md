# youdaoOCR
有道OCR Go SDK
用于请求有道OCR
支持使用文件名和go的 Image.image对象获取OCR识别结果
提供Word和Line结构体用于适配按行解析和按字解析的结果

# example
```go
  // 创建OCR对象
	o := ydocr.NewOcr(appKey, appSecret)
  // 请求OCR识别结果
	rsp, err := o.FileOcr(filename)
  if err != nil{
		log.Fatal("base post err: %s", err)
	}
  // 解析结果
	l, err := ydocr.Response2Line(rsp)
	if err != nil{
		log.Fatal("response to line err: %s", err)
	}
	t.Logf("rsp: %+v", l)
```