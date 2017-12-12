package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	session "github.com/aws/aws-sdk-go/aws/session"
	s3 "github.com/aws/aws-sdk-go/service/s3"
	mkr "github.com/mackerelio/mackerel-client-go"
)

type options struct {
	bucket string
	region string
}

func parseOptions() (*options, error) {
	bucket := flag.String("bucket", "", "bucket name")
	region := flag.String("region", "", "region name")
	flag.Parse()

	opts := &options{}
	if bucket == nil {
		return nil, fmt.Errorf("bucket required")
	}
	if b := *bucket; b == "" {
		return nil, fmt.Errorf("bucket required")
	}
	opts.bucket = *bucket

	if region == nil {
		return nil, fmt.Errorf("region required")
	}
	if r := *region; r == "" {
		return nil, fmt.Errorf("region required")
	}
	opts.region = *region

	return opts, nil
}

func main() {
	opts, err := parseOptions()
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}

	metrics, err := fetchS3Metrics(opts)
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}
	for _, m := range metrics {
		fmt.Fprintf(os.Stdout, "%s\t%d\t%d\n", m.Name, m.Value, m.Time)
	}
}

func fetchS3Metrics(opts *options) ([]*mkr.MetricValue, error) {
	metrics := make([]*mkr.MetricValue, 0)

	cre := credentials.NewEnvCredentials()
	s, err := session.NewSession(&aws.Config{
		Credentials: cre,
		Region:      &opts.region,
	})
	// opts := session.Options{}
	// s, err := session.NewSessionWithOptions(opts)
	if err != nil {
		return metrics, err
	}

	srv := s3.New(s)
	out, err := srv.ListObjects(&s3.ListObjectsInput{Bucket: &opts.bucket})
	if err != nil {
		return metrics, err
	}

	totalSize := int64(0)
	for _, obj := range out.Contents {
		if size := obj.Size; size != nil {
			totalSize += *size
		}
	}
	ts := time.Now()
	metrics = append(metrics, &mkr.MetricValue{
		Name:  "objects_count." + opts.bucket,
		Value: len(out.Contents),
		Time:  ts.Unix(),
	})
	metrics = append(metrics, &mkr.MetricValue{
		Name:  "total_size." + opts.bucket,
		Value: totalSize,
		Time:  ts.Unix(),
	})
	return metrics, nil
}
