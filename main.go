package main

import (
	"fmt"
	"log"
	"os"
	"time"

	session "github.com/aws/aws-sdk-go/aws/session"
	s3 "github.com/aws/aws-sdk-go/service/s3"
	mkr "github.com/mackerelio/mackerel-client-go"
)

func main() {
	bucket := ""
	metrics, err := fetchS3Metrics(bucket)
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}
	for _, m := range metrics {
		fmt.Fprintf(os.Stdout, "%s\t%d\t%d\n", m.Name, m.Value, m.Time)
	}
}

func fetchS3Metrics(bucket string) ([]*mkr.MetricValue, error) {
	metrics := make([]*mkr.MetricValue, 0)

	opts := session.Options{
		SharedConfigState: session.SharedConfigStateFromEnv,
	}
	s, err := session.NewSessionWithOptions(opts)
	if err != nil {
		return metrics, err
	}

	srv := s3.New(s)
	out, err := srv.ListObjects(&s3.ListObjectsInput{Bucket: &bucket})
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
		Name:  "objects_count." + bucket,
		Value: len(out.Contents),
		Time:  ts.Unix(),
	})
	metrics = append(metrics, &mkr.MetricValue{
		Name:  "total_size." + bucket,
		Value: totalSize,
		Time:  ts.Unix(),
	})
	return metrics, nil
}
