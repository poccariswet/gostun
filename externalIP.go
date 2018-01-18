package main

import (
	"io"
	"net/http"
	"os"
)

func main() {
	resp, err := http.Get("http://myexternalip.com/raw")
	//	resp, err := http.Get("http://httpbin.org/ip")
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
}
