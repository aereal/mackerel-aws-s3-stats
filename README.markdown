# mackerel-aws-s3-stats

## Usage

```sh
MACKEREL_APIKEY='...' mackerel-aws-s3-stats --region $REGION --bucket $BUCKET_NAME --service [--bucket $ANOTHER_BUCKET_NAME]
```

Sends total size of objects (bytes) and number of objects to [Mackerel][].

or only calculating stats, do not request to Mackerel.

```sh
mackerel-aws-s3-stats --region $REGION --bucket $BUCKET_NAME --no-post [--bucket $ANOTHER_BUCKET_NAME]
```

See also how to specify the AWS credentials:
https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html

## Build

```sh
go build ./...
GOOS=linux GOARCH=amd64 go build -o mackerel-aws-s3-stats.linux.amd64 ./...
```

[Mackerel]: https://mackerel.io/
