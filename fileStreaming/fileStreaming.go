package fileStreaming

import (
	"io"
	"log"
	"math/rand"
	"os"
)

type Server struct{}
type Client struct{}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const CHUNK_SIZE = 1024

func serverError(stream FileUploadService_UploadServer) {
	err := stream.SendAndClose(&UploadStatus{
		Message: "Error uploading file",
		Code:    UploadStatusCode_Failed,
	})
	if err != nil {
		log.Fatalf("unable to send error status")
	}
}

func (s *Server) Upload(stream FileUploadService_UploadServer) (err error) {
	name_rand := make([]byte, 6)
	for i := range name_rand {
		name_rand[i] = letterBytes[rand.Int63()%int64(6)]
	}

	filename := "./server_files/file_" + string(name_rand)
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		serverError(stream)
		log.Fatalf("unable to open file")
	}

	complete := false
	for {
		chunk, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				complete = true
				break
			}

			serverError(stream)
			return err
		}

		f.Write(chunk.Content)
	}

	if complete {
		log.Println("Upload Completed!")
	}

	err = stream.SendAndClose(&UploadStatus{
		Message: "Upload complete",
		Code:    UploadStatusCode_Ok,
	})
	if err != nil {
		return err
	}

	return nil
}
