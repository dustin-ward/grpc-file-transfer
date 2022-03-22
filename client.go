package main

import (
	"context"
	"grpc-demo/file-transfer/fileStreaming"
	"io"
	"log"
	"os"

	"google.golang.org/grpc"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Dial: %v", err)
	}
	defer conn.Close()

	// Open file
	filename := "./client_files/important_file.txt"
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("unable to open file")
	}

	c := fileStreaming.NewFileUploadServiceClient(conn)
	stream, err := c.Upload(context.Background())
	if err != nil {
		log.Fatalf("error establishing upload")
	}

	buf := make([]byte, fileStreaming.CHUNK_SIZE)
	for {
		n, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Fatalf("error reading data from file")
		}

		err = stream.Send(&fileStreaming.Chunk{
			Content: buf[:n],
		})
		if err != nil {
			log.Fatalf("error sending data during transfer")
		}
	}

	status, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error retrieving status")
	}
	if status.Code != fileStreaming.UploadStatusCode_Ok {
		log.Fatalf("Error, file transfer finished with non-ok status")
	} else {
		println("File Transfer Complete!")
	}
}
