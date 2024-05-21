package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
)

func main() {
	url := ""
	filePath := "C:\\Users\\upload\\test.png"
	mchntCd := ""
	message := ""

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// 创建一个新的表单数据
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加mchnt_cd字段
	_ = writer.WriteField("mchnt_cd", mchntCd)

	// 添加message字段
	_ = writer.WriteField("message", message)

	// 创建文件字段
	originfileName := filepath.Base(file.Name())
	part, err := createFormFileWithContentType("file", originfileName, "image/png", writer)
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return
	}

	// 将文件内容复制到表单数据中
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Error copying file to form:", err)
		return
	}

	// 完成表单数据
	err = writer.Close()
	if err != nil {
		fmt.Println("Error closing writer:", err)
		return
	}

	// 创建一个HTTP POST请求
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// 设置请求头中的Content-Type
	request.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer response.Body.Close()

	// 打印响应
	fmt.Println("Response Status:", response.Status)

	resp_body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.Status)
	fmt.Println(string(resp_body))
}

func createFormFileWithContentType(fieldName, fileName, contentType string, writer *multipart.Writer) (io.Writer, error) {
	// 创建 MIME 头部
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldName, fileName))
	h.Set("Content-Type", contentType)

	// 创建文件部分
	return writer.CreatePart(h)
}
