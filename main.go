package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	session "github.com/aws/aws-sdk-go/aws/session"
	s3 "github.com/aws/aws-sdk-go/service/s3"
	mkr "github.com/mackerelio/mackerel-client-go"
)

type options struct {
	buckets []string
	region  string
}

type bucketFlags []string

func (f *bucketFlags) String() string {
	buf := ""
	for _, b := range *f {
		buf += b
	}
	return buf
}

func (f *bucketFlags) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func parseOptions() (*options, error) {
	var buckets bucketFlags
	flag.Var(&buckets, "bucket", "bucket name")
	region := flag.String("region", "", "region name")
	flag.Parse()

	opts := &options{}
	if len(buckets) == 0 {
		return nil, fmt.Errorf("bucket required")
	}
	opts.buckets = buckets

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

	s, err := session.NewSession(&aws.Config{
		Region: &opts.region,
	})
	if err != nil {
		return metrics, err
	}

	for _, b := range opts.buckets {
		ms, er := fetchS3MetricsByBucket(s, b)
		if er != nil {
			log.Printf("! Error on bucket %s: %s\n", b, er)
			continue
		}
		metrics = append(metrics, ms...)
	}

	return metrics, err
}

func fetchS3MetricsByBucket(s *session.Session, bucket string) ([]*mkr.MetricValue, error) {
	metrics := make([]*mkr.MetricValue, 0)

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
		Name:  "s3_objects_count." + bucket,
		Value: len(out.Contents),
		Time:  ts.Unix(),
	})
	metrics = append(metrics, &mkr.MetricValue{
		Name:  "s3_total_size." + bucket,
		Value: totalSize,
		Time:  ts.Unix(),
	})
	return metrics, nil
}
