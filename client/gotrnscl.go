package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var destURL = flag.String("dest", "", "URL location of server")
var filePath = flag.String("path", "", "Location of file(s) to transfer.")

func main() {
	flag.Parse()

	fmt.Println(os.Args)

	//Check to see if the required number of parameters are set
	if len(*filePath) == 0 || len(*destURL) == 0 {
		flag.Usage()
		return
	}

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}

	//Open requested file
	//fmt.Println("Selected File:", *filePath)
	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	req, err := http.NewRequest("POST", *destURL, file)
	if err != nil {
		log.Fatal(err)
	}

	//set the content-type and the name of the file
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Content-Disposition", "attachment; filename=\""+*filePath+"\"")

	//turn off chunked transfer encoding
	req.TransferEncoding = []string{"identity"}
	fi, err := file.Stat()
	req.ContentLength = fi.Size()

	req.Header.Set("AS2-To", "\"axway -> interchange\"")
	req.Header.Set("AS2-From", "\"axway-trading\"")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("StatusCode:", resp.StatusCode)
}
