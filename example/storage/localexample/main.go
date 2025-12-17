package main

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/syntaxfa/quick-connect/adapter/storage/local"
)

func main() {
	cfg := local.Config{
		RootPath: "./uploads",
		BaseURL:  "http://localhost:2560/files",
	}

	ctx := context.Background()
	key := "example/test.txt"

	storage, nErr := local.New(cfg, slog.Default())
	if nErr != nil {
		panic(nErr)
	}

	dataBytes, rErr := os.ReadFile("example/storage/localexample/quick-conn-exam.txt")
	if rErr != nil {
		panic(rErr)
	}

	file := bytes.NewReader(dataBytes)

	key, uErr := storage.Upload(ctx, file, 1, key, "", false)
	if uErr != nil {
		panic(uErr)
	}

	fmt.Println(key)

	fmt.Println("GetURL:")

	url, gErr := storage.GetURL(ctx, key)
	if gErr != nil {
		panic(gErr)
	}

	fmt.Println(url)

	fmt.Println("GetPresignedURL:")
	url, gErr = storage.GetPresignedURL(ctx, key)
	if gErr != nil {
		panic(gErr)
	}

	fmt.Println(url)

	//fmt.Println("Delete:")
	//if dErr := storage.Delete(ctx, key); dErr != nil {
	//	panic(dErr)
	//} else {
	//	fmt.Printf("file %s deleted\n", key)
	//}
}
