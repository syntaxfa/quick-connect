package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/syntaxfa/quick-connect/adapter/storage/aws"
	"os"
	"time"
)

func main() {
	cfg := aws.Config{
		Endpoint:             "https://storage.iran.liara.space",
		AccessKeyID:          "tki1ets9amtds2di",
		SecretAccessKey:      "472fc79f-06cd-4109-ac9f-da11196de879",
		BucketName:           "stellarhonar",
		Region:               "default",
		UseSSL:               true,
		UsePathStyle:         true,
		SupportObjectACL:     false,
		PresignPublicExpire:  time.Hour * 24 * 7,
		PresignPrivateExpire: time.Second * 120,
	}

	ctx := context.Background()

	storage, err := aws.New(ctx, cfg)
	if err != nil {
		panic(err)
	}

	dataBytes, rErr := os.ReadFile("quick-conn-exam.txt")
	if rErr != nil {
		panic(rErr)
	}

	file := bytes.NewReader(dataBytes)

	key := "examplegolang.txt"
	k, uErr := storage.Upload(ctx, file, file.Size(), key, "docs", true)
	if uErr != nil {
		panic(uErr)
	}

	fmt.Println("uploaded key:")
	fmt.Println(k)

	fmt.Println("get public key")
	publicURL, getPErr := storage.GetURL(ctx, key)
	if getPErr != nil {
		panic(getPErr)
	}

	fmt.Println(publicURL)

	privateURL, getPrErr := storage.GetPresignedURL(ctx, key)
	if getPrErr != nil {
		panic(getPrErr)
	}

	fmt.Println("private URL:")
	fmt.Println(privateURL)

	fmt.Println("exists:")

	exists, exErr := storage.Exists(ctx, key)
	if exErr != nil {
		fmt.Println("not exists")
	}
	fmt.Println("exists", exists)

	fmt.Println("delete:")
	if dErr := storage.Delete(ctx, key); dErr != nil {
		fmt.Println("error failed", dErr.Error())
	} else if dErr == nil {
		fmt.Println("file deleted")
	}

	fmt.Println("exists:")

	exists, exErr = storage.Exists(ctx, key)
	if exErr != nil {
		fmt.Println("not exists")
	}
	fmt.Println("exists", exists)
}
