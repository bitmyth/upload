package upload

import "time"

type CreateUploadResponse struct {
	UploadId string
}

type UploadChunkResponse struct {
	ID   string
	Hash string
}

type FileChunk struct {
	Hash        string
	ChunkNumber int
}

type ReassembleChunksRequest struct {
	UploadId string
	Chunks   []FileChunk
	Filename string
}

type ReassembleChunksResponse struct {
	Hash     string
	Filename string
}

type DownloadRequest struct {
	UploadId string
}

type Chunk struct {
	ID        string
	Hash      string
	CreatedAt time.Time
}

type File struct {
	UploadId string
	Filename string
	Size     uint64
}
