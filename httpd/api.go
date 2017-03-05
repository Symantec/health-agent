package httpd

import (
	"fmt"
	"io"
	"net"
	"net/http"
)

type HtmlWriter interface {
	WriteHtml(writer io.Writer)
}

type RequestHtmlWriter interface {
	HtmlWriter
	RequestWriteHtml(writer io.Writer, req *http.Request)
}

var htmlWriters []HtmlWriter

func StartServer(portNum uint) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", portNum))
	if err != nil {
		return err
	}
	http.HandleFunc("/", statusHandler)
	http.HandleFunc("/favicon.ico", func(http.ResponseWriter, *http.Request) {})
	go http.Serve(listener, nil)
	return nil
}

func AddHtmlWriter(htmlWriter HtmlWriter) {
	htmlWriters = append(htmlWriters, htmlWriter)
}
