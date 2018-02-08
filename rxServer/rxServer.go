package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func retrieveFilename(s string) (string, error) {
	tokens := strings.Split(s, ";")
	for _, v := range tokens {
		st := strings.TrimSpace(v)
		if strings.HasPrefix(st, "filename") {
			nameValue := strings.Split(st, "=")
			if len(nameValue) == 2 {
				return strings.Trim(nameValue[1], "\""), nil
			}
		}
	}
	return "", errors.New("Error locating filename in content-disposition header")
}

func saveUploadedFile(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	filename, err := retrieveFilename(r.Header.Get("Content-Disposition"))
	if err != nil {
		log.Fatal("Error retrieving filename", err)
	}

	file, err := os.Create("downloads/" + filename)
	if err != nil {
		log.Fatal("Error retrieving filename", err)
	}
	defer file.Close()

	_, err = io.Copy(file, r.Body)
	if err != nil {
		log.Fatal("Error saving data", err)
	}
	fmt.Println("Successully saved file", filename)
}

func printLocalIpAddr() {
	ifaces, err := net.Interfaces()
	if err != nil {
		panic("Error getting interfaces from os.")
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			panic("Error retrieving address from interface.")
		}
		fmt.Println("Interface:", i.Name)

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			fmt.Println("\tIP Address:", ip)
		}
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/upload", saveUploadedFile)

	fmt.Println("rxServer is running on ip address:")
	//print local ip addressess
	printLocalIpAddr()

	err := http.ListenAndServe(":8080", mux) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
