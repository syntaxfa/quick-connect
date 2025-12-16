package aws

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/syntaxfa/quick-connect/app/storageapp/service"
)

var _ service.Storage = (*Adapter)(nil)

type Adapter struct {
	cfg           Config
	client        *s3.Client
	presignClient *s3.PresignClient
	bucketName    string
	publicBaseURL string
}

func New(ctx context.Context, cfg Config) (*Adapter, error) {
	awsCfg, lcErr := config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, "")),
	)
	if lcErr != nil {
		return nil, fmt.Errorf("load s3 config failed: %s", lcErr.Error())
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
		}

		o.UsePathStyle = cfg.UsePathStyle
	})

	_, hbErr := client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(cfg.BucketName),
	})
	if hbErr != nil {
		return nil, fmt.Errorf("connection check failed for bucket %s: %s", cfg.BucketName, hbErr.Error())
	}

	return &Adapter{
		cfg:           cfg,
		client:        client,
		presignClient: s3.NewPresignClient(client),
		bucketName:    cfg.BucketName,
		publicBaseURL: cfg.Endpoint,
	}, nil
}

func (a *Adapter) Upload(ctx context.Context, file io.Reader, size int64, key string, contentType string, isPublic bool) (string, error) {
	input := &s3.PutObjectInput{
		Bucket:        aws.String(a.bucketName),
		Key:           aws.String(key),
		Body:          file,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(contentType),
	}

	if isPublic {
		input.ACL = types.ObjectCannedACLPublicRead
	} else {
		input.ACL = types.ObjectCannedACLPrivate
	}

	_, err := a.client.PutObject(ctx, input)
	if err != nil {
		return "", fmt.Errorf("s3 upload failed: %w", err)
	}

	return key, nil
}

func (a *Adapter) Delete(ctx context.Context, key string) error {
	_, dErr := a.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(a.bucketName),
		Key:    aws.String(key),
	})
	if dErr != nil {
		return fmt.Errorf("s3 delete failed: %s", dErr.Error())
	}

	return nil
}

func (a *Adapter) GetURL(ctx context.Context, key string) (string, error) {
	if a.cfg.SupportObjectACL {
		return fmt.Sprintf("%s/%s/%s", a.publicBaseURL, a.bucketName, key), nil
	}

	req, err := a.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(a.bucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = a.cfg.PresignPublicExpire
	})

	if err != nil {
		return "", fmt.Errorf("failed to generate public url: %w", err)
	}

	return req.URL, nil
}

func (a *Adapter) GetPresignedURL(ctx context.Context, key string) (string, error) {
	req, pErr := a.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(a.bucketName),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = a.cfg.PresignPrivateExpire
	})
	if pErr != nil {
		return "", fmt.Errorf("presign url failed: %s", pErr.Error())
	}

	return req.URL, nil
}

func (a *Adapter) Exists(ctx context.Context, key string) (bool, error) {
	_, err := a.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(a.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return false, err
	}

	return true, nil
}
