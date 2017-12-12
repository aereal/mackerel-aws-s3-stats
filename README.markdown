# mackerel-aws-s3-stats

## Usage

```sh
mackerel-aws-s3-stats --region $REGION --bucket $BUCKET_NAME [--bucket $ANOTHER_BUCKET_NAME]
```

Writes total size of objects (bytes) and number of objects to stdout as Sensu-compatible format.

See also how to specify the AWS credentials:
https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html

## Build

```sh
go build ./...
GOOS=linux GOARCH=amd64 go build -o mackerel-aws-s3-stats.linux.amd64 ./...
```
