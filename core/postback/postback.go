package postback

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type PostbackInf interface {
}

type Postback struct {
}

type HttpPostback struct {
	filename  string
	targetUrl string
	Content   []byte
	Guid      string
	rid		  string
}

func (self *HttpPostback) SetGuid(guid string) {
	self.Guid = guid
}

func (self *HttpPostback) SetRid(rid string) {
	self.rid = rid
}

func (self *HttpPostback) SetTargetUrl(targetUrl string) {
	self.targetUrl = targetUrl
}

func (self *HttpPostback) SetFileName(filename string) {
	self.filename = filename
}

func (self *HttpPostback) PostFile() error {
	body_buf := bytes.NewBufferString("")
	body_writer := multipart.NewWriter(body_buf)

	// use the body_writer to write the Part headers to the buffer
	_, err := body_writer.CreateFormFile("uploadfile", self.filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}

	// the file data will be the second part of the body
	fh, err := os.Open(self.filename)
	if err != nil {
		fmt.Println("error opening file")
		return err
	}
	defer fh.Close()
	// need to know the boundary to properly close the part myself.
	boundary := body_writer.Boundary()

	close_buf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

	// use multi-reader to defer the reading of the file data until
	// writing to the socket buffer.
	request_reader := io.MultiReader(body_buf, fh, close_buf)
	fi, err := fh.Stat()
	if err != nil {
		fmt.Printf("Error Stating file: %s", self.filename)
		return err
	}
	req, err := http.NewRequest("POST", self.targetUrl, request_reader)
	if err != nil {
		return err
	}

	// Set headers for multipart, and Content Length
	req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	req.Header.Set("Guid", self.Guid)
	req.ContentLength = fi.Size() + int64(body_buf.Len()) + int64(close_buf.Len())

	_, err = http.DefaultClient.Do(req)
	return err
}

func (self *HttpPostback) PostContent() error {
	// 判断 WebHook 通知
	reader := bytes.NewReader(self.Content)

	timeout := time.Duration(3) * time.Second
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	request, _ := http.NewRequest("POST", self.targetUrl+"/"+self.rid+"/"+self.filename, reader)
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("Guid", self.Guid)
	client := &http.Client{
		Transport: tr,
		Timeout:   timeout,
	}
	resp, err := client.Do(request)

	if err != nil {
		log.Print("上报记录失败.", err)
	} else {
		log.Print("上报记录成功.")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
	}
	log.Print("回应：", string(body))
	return nil
}
