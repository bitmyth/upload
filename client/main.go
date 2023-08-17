package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bitmyth/upload"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	var path string
	flag.StringVar(&path, "p", "", "filepath")
	flag.Parse()

	stat, err := os.Stat(path)
	if err != nil {
		log.Fatalln(err)
		return
	}

	c := Client{
		fileInfo: stat,
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

	url := fmt.Sprintf("%s/%s%s", fullPath(PathDownload), c.UploadId, c.fileInfo.Name())
	log.Println("URL:")
	log.Println(url)
}

var (
	baseUrl = "http://localhost:8080/v1/files"
)

const (
	ChunkSize = 1024 * 1024 * 4

	PathNewUpload = "uploads"
	PathDownload  = "file"
)

func fullPath(path string) string {
	return baseUrl + "/" + path
}

type Client struct {
	UploadId string
	Chunks   []upload.FileChunk
	fileInfo os.FileInfo
	Filepath string
}

func (c *Client) NewUpload() error {
	url := fullPath(PathNewUpload)

	data := ""
	body := strings.NewReader(data)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	httpClient := client()
	setCommonHeaders(req)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(all))
	defer resp.Body.Close()

	var result upload.CreateUploadResponse
	json.Unmarshal(all, &result)

	c.UploadId = result.UploadId

	println("upload id:", c.UploadId)

	return nil
}

func (c *Client) Upload() error {
	size := c.fileInfo.Size()

	chunksCount := size/ChunkSize + 1

	fileData, err := os.ReadFile(c.Filepath)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	log.Println("Total chunks:", chunksCount)

	for i := 0; i < int(chunksCount); i++ {
		begin := i * ChunkSize
		end := min(len(fileData), begin+ChunkSize)

		err = c.UploadChunk(fileData[begin:end], i+1)
		if err != nil {
			log.Fatalln(i+1, err)
			return err
		}
	}

	err = c.ReassembleChunks()
	if err != nil {
		log.Fatalln(err)
		return err
	}

	return nil
}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}

func (c *Client) UploadChunk(data []byte, chunkNumber int) error {
	log.Println("Uploading chunk:", chunkNumber, " length:", len(data))

	httpClient := client()

	url := fullPath(filepath.Join(PathNewUpload, c.UploadId, strconv.Itoa(chunkNumber)))

	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	setCommonHeaders(req)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(all))
	defer resp.Body.Close()

	var result upload.UploadChunkResponse
	json.Unmarshal(all, &result)

	c.Chunks = append(c.Chunks, upload.FileChunk{
		Hash:        result.Hash,
		ChunkNumber: chunkNumber,
	})

	return nil
}

func (c *Client) ReassembleChunks() error {
	log.Println("Reassemble chunks")

	httpClient := client()

	url := fullPath(filepath.Join(PathNewUpload, c.UploadId))

	payload := upload.ReassembleChunksRequest{
		UploadId: c.UploadId,
		Chunks:   c.Chunks,
		Filename: c.fileInfo.Name(),
	}

	marshal, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewReader(marshal))
	if err != nil {
		return err
	}
	setCommonHeaders(req)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println(string(all))

	return nil
}

func client() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: tr}
}

func setCommonHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")
}
