package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Version number constant.
const Version = "0.0.1"

var (
	sourcefile = flag.String("file", "", "File to upload")
	destfile   = flag.String("dest", "", "Destination to upload")
	versDisp   = flag.Bool("version", false, "Display version")
	pathExp    = regexp.MustCompile(`s3:\/\/(?P<bucket>\w+)\/(?P<key>.+)`)
)

func main() {
	flag.Parse()

	if *versDisp {
		fmt.Printf("Version: v%s %s\n", Version, runtime.Version())
		os.Exit(0)
	}

	_, err := os.Stat(*sourcefile)
	if err != nil {
		log.Panic("Couldn't stat file: " + err.Error())
	}
	f, err := os.Open(*sourcefile)
	if err != nil {
		log.Panic("" + err.Error())
	}

	match := pathExp.FindStringSubmatch(*destfile)
	results := make(map[string]string)
	for i, name := range match {
		results[pathExp.SubexpNames()[i]] = name
	}
	if results["bucket"] == "" || results["key"] == "" {
		log.Panic("Couldn't parse destination: ", results)
		os.Exit(1)
	}

	ctx := context.Background()

	s3session := session.Must(session.NewSession(
		&aws.Config{
			Endpoint:         aws.String(os.Getenv("AWS_S3_ENDPOINT")),
			S3ForcePathStyle: aws.Bool(true),
		},
	))

	uploader := s3manager.NewUploader(s3session)

	result, err := uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(results["bucket"]),
		Key:    aws.String(results["key"]),
		Body:   f,
	})

	if err != nil {
		panic(err)
	}

	log.Println("File Uploaded Successfully, URL : ", result.Location)
}
