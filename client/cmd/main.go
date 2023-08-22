package main

import (
	"flag"
	"fmt"
	"github.com/bitmyth/upload/client"
	"log"
	"os"
)

func main() {
	var path string
	var host string
	var port string
	var scheme string
	var base string
	flag.StringVar(&path, "f", "", "filepath")
	flag.StringVar(&host, "h", "localhost", "host")
	flag.StringVar(&port, "P", "443", "port")
	flag.StringVar(&scheme, "s", "https", "scheme")
	flag.StringVar(&base, "b", "", "api url base")
	flag.Parse()

	client.BaseUrl = fmt.Sprintf("%s://%s:%s", scheme, host, port)
	if base != "" {
		client.BaseUrl += "/" + base
	}

	hostEnv := os.Getenv("HOST")
	if hostEnv != "" {
		client.BaseUrl = hostEnv
	}
	log.Println(client.BaseUrl)

	if path == "" {
		log.Println("Please provide file path by -f option")
		return
	}

	stat, err := os.Stat(path)
	if err != nil {
		log.Fatalln(err)
		return
	}

	c := client.Client{
		FileInfo: stat,
		Filepath: path,
	}

	err = c.NewUpload()
	if err != nil {
		log.Fatalln(err)
		return
	}

	err = c.Upload()
	if err != nil {
		log.Fatalln(err)
		return
	}

	log.Println("URL:", c.DownloadLink())
}
