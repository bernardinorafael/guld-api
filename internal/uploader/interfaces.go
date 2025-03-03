package uploader

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
)

type UploaderInterface interface {
	UploadFile(ctx context.Context, file io.Reader, key string) (*manager.UploadOutput, error)
}
