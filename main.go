package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Version number constant.
const Version = "0.0.2"

var pathExp = regexp.MustCompile(`s3:\/\/(?P<bucket>\w+)\/(?P<key>.+)`)

func main() {
	numArgs := len(os.Args[1:])
	if numArgs >= 1 && os.Args[1] == "-v" {
		fmt.Printf("Version: v%s %s\n", Version, runtime.Version())
		os.Exit(0)
	}
	if numArgs != 2 {
		fmt.Printf("Usage: %s [-v] source-file s3://bucket/key/path\n", os.Args[0])
		os.Exit(0)
	}
	sourcefile := os.Args[1]
	destfile := os.Args[2]

	info, err := os.Stat(sourcefile)
	if err != nil {
		log.Panic("Couldn't stat file: " + err.Error())
	} else if info.IsDir() {
		log.Panic("Source is not a file.")
	}
	f, err := os.Open(sourcefile)
	if err != nil {
		log.Panic(err.Error())
	}

	match := pathExp.FindStringSubmatch(destfile)
	results := make(map[string]string)
	for i, name := range match {
		results[pathExp.SubexpNames()[i]] = name
	}
	if results["bucket"] == "" || results["key"] == "" {
		log.Panic("Couldn't parse destination: ", results)
		os.Exit(1)
	}

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	client := s3.NewFromConfig(cfg)

	uploader := manager.NewUploader(client)
	result, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(results["bucket"]),
		Key:    aws.String(results["key"]),
		Body:   f,
	})

	if err != nil {
		panic(err)
	}

	log.Println("File Uploaded Successfully, URL : ", result.Location)
}
