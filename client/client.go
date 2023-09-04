package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
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

var (
	BaseUrl = "http://localhost:8080"
)

const (
	ChunkSize = 1024 * 1024 * 4

	PathNewUpload = "uploads"
	PathDownload  = "file"
)

func fullPath(path string) string {
	return BaseUrl + "/" + path
}

type Client struct {
	Transport http.RoundTripper
	UploadId  string
	Chunks    []upload.FileChunk
	FileInfo  os.FileInfo
	Filepath  string
}

func (c *Client) NewUpload() error {
	url := fullPath(PathNewUpload)

	data := ""
	body := strings.NewReader(data)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	httpClient := c.client()
	setCommonHeaders(req)

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode > 400 {
		log.Println(resp.StatusCode)
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

func (c *Client) DownloadLink() string {
	url := fmt.Sprintf("%s/%s%s", fullPath(PathDownload), c.UploadId, c.FileInfo.Name())
	return url
}

func (c *Client) Upload() error {
	size := c.FileInfo.Size()

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

	httpClient := c.client()

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

	httpClient := c.client()

	url := fullPath(filepath.Join(PathNewUpload, c.UploadId))

	payload := upload.ReassembleChunksRequest{
		UploadId: c.UploadId,
		Chunks:   c.Chunks,
		Filename: c.FileInfo.Name(),
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

func (c *Client) client() *http.Client {
	var tr http.RoundTripper
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	if c.Transport != nil {
		tr = c.Transport
		log.Println("using custom http transport")
	}

	return &http.Client{Transport: tr}
}

func setCommonHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")
}
