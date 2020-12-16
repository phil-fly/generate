package postback

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

type PostbackInf interface {

}

type Postback struct {

}



type HttpPostback struct {
	filename string
	targetUrl	string
	Content []byte
}

func (self *HttpPostback)SetTargetUrl(targetUrl string) {
	self.targetUrl = targetUrl
}

func (self *HttpPostback)SetFileName(filename string) {
	self.filename = filename
}

func (self *HttpPostback)PostFile() (error) {
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
	req.ContentLength = fi.Size() + int64(body_buf.Len()) + int64(close_buf.Len())

	_,err = http.DefaultClient.Do(req)
	return err
}

func (self *HttpPostback)PostContent() (error) {
	body_buf := bytes.NewBufferString("")
	body_writer := multipart.NewWriter(body_buf)

	// use the body_writer to write the Part headers to the buffer
	_, err := body_writer.CreateFormFile("uploadfile", self.filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}
	// need to know the boundary to properly close the part myself.
	boundary := body_writer.Boundary()

	close_buf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

	// use multi-reader to defer the reading of the file data until
	// writing to the socket buffer.

	var bytesData bytes.Buffer //Buffer是一个实现了读写方法的可变大小的字节缓冲

	bytesData.Write(body_buf.Bytes())
	bytesData.Write(self.Content)
	bytesData.Write(close_buf.Bytes())
	request_reader := bytes.NewReader(bytesData.Bytes())
	req, err := http.NewRequest("POST", self.targetUrl, request_reader)
	if err != nil {
		return err
	}

	// Set headers for multipart, and Content Length
	req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	req.ContentLength = int64(len(self.Content)) + int64(body_buf.Len()) + int64(close_buf.Len())
	_,err = http.DefaultClient.Do(req)
	if err !=nil {
		log.Print(err)
	}
	return err
}