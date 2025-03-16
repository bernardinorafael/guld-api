package uploader

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	. "github.com/bernardinorafael/internal/_shared/errors"

	"github.com/bernardinorafael/pkg/logger"
)

const (
	S3_KEY    = ""
	S3_SECRET = ""

	S3_REGION = "us-east-2"
	S3_BUCKET = "gulg-profile-picture"

	MAX_FILE_SIZE = 3 << 20 // 3mb
)

type Uploader struct {
	log    logger.Logger
	Client *s3.Client
}

func NewUploader(ctx context.Context, log logger.Logger) UploaderInterface {
	credential := credentials.NewStaticCredentialsProvider(S3_KEY, S3_SECRET, "")

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithCredentialsProvider(credential),
		config.WithRegion(S3_REGION),
	)

	if err != nil {
		log.Criticalw(ctx, "Error loading default config: %v", logger.Err(err))
		panic(err)
	}

	return &Uploader{
		log:    log,
		Client: s3.NewFromConfig(cfg),
	}
}

func (svc Uploader) UploadFile(ctx context.Context, file io.Reader, key string) (*manager.UploadOutput, error) {
	if svc.Client == nil {
		return nil, NewBadRequestError("S3 client not initialized", nil)
	}
	if file == nil {
		return nil, NewBadRequestError("file reader is nil", nil)
	}

	uploader := manager.NewUploader(svc.Client)
	object := &s3.PutObjectInput{
		Bucket:      aws.String(S3_BUCKET),
		Key:         aws.String(key),
		ContentType: aws.String("image/webp"),
		Body:        file,
	}

	out, err := uploader.Upload(ctx, object)
	if err != nil {
		svc.log.Errorw(ctx, "error uploading file to s3", logger.Err(err))
		return nil, NewBadRequestError("error uploading file to s3", err)
	}

	return out, nil
}
