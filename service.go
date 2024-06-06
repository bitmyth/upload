package upload

import (
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"
)

var Dir = "upload"

func InitDir() {
	err := os.MkdirAll(Dir+"/chunks", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

type Service struct {
	hash.Hash
}

func NewService(dep dependency) Service {
	var s Service
	s.Hash = dep.Hash()
	return s
}

func (u Service) newUploadId() string {
	return fmt.Sprintf("%d", time.Now().UnixMicro())
}

func (u Service) CreateUpload() CreateUploadResponse {
	return CreateUploadResponse{
		UploadId: u.newUploadId(),
	}
}

func (u Service) chunkFilepath(sum string) string {
	return filepath.Join(Dir+"/chunks", sum)
}

func (u Service) finalFilepath(uploadId string) string {
	return filepath.Join(Dir+"/", uploadId)
}

// UploadChunk upload each 4MB chunk of a file
func (u Service) UploadChunk(req *http.Request) (*UploadChunkResponse, error) {
	d, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	sum := u.Hash.Sum(d)
	sumHex := hex.EncodeToString(sum[:])

	chunk, err := os.Create(u.chunkFilepath(sumHex))
	defer chunk.Close()
	if err != nil {
		return nil, err
	}

	_, err = chunk.Write(d)
	if err != nil {
		return nil, err
	}

	// TODO save chunk record to database, generate ID for each chunk
	// chk := Chunk{
	// 	ID:   "",
	// 	Hash: sumHex,
	// }

	return &UploadChunkResponse{
		ID:   "",
		Hash: sumHex,
	}, nil
}

// ReassembleChunk put chunks together
func (u Service) ReassembleChunk(req ReassembleChunksRequest) (*ReassembleChunksResponse, error) {
	sort.Slice(req.Chunks, func(i, j int) bool {
		return req.Chunks[i].ChunkNumber < req.Chunks[j].ChunkNumber
	})

	finalFile, err := os.Create(u.finalFilepath(req.UploadId + req.Filename))
	if err != nil {
		return nil, err
	}
	defer finalFile.Close()

	fileHash := u.Hash
	for _, chunk := range req.Chunks {
		chunkFilepath := u.chunkFilepath(chunk.Hash)

		chunkData, e := os.ReadFile(chunkFilepath)
		if e != nil {
			return nil, e
		}

		fileHash.Write(chunkData)

		_, err = finalFile.Write(chunkData)
		if err != nil {
			return nil, err
		}

		// Remove chunk data
		_ = os.RemoveAll(chunkFilepath)
	}

	sum := fileHash.Sum(nil)

	resp := &ReassembleChunksResponse{
		Hash:     hex.EncodeToString(sum[:]),
		Filename: req.Filename,
	}

	// TODO persist file record into database

	return resp, nil
}

func (u Service) Download(req DownloadRequest, header http.Header, writer http.ResponseWriter) error {
	// Retrieve file by uploadId
	//fileRecord := File{
	//	UploadId: req.UploadId,
	//	Filename: "",
	//}

	name := u.finalFilepath(req.UploadId)
	stat, err := os.Stat(name)
	if err != nil {
		return err
	}

	filename := stat.Name()
	header.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	header.Set("Content-Length", fmt.Sprintf("%d", stat.Size()))

	f, err := os.Open(name)
	if err != nil {
		return err
	}

	http.ServeContent(writer, req.Req, filename, stat.ModTime(), f)
	//io.Copy(writer, f)
	defer f.Close()

	return nil
}
