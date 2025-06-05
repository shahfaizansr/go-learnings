package minio_utils

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/shahfaizansr/models"
)

const (
	DATE_TIME_FORMAT                   = "09-07-2017 17:06:06"
	NO_SUCH_BUCKET                     = "no bucket found"
	FAILED_TO_DECODE                   = "failed to decode the string"
	FAILED_TO_PREPARE_PARAMS           = "failed to prepare the minio params"
	FAILED_TO_STORE_OBJECT             = "failed to store object"
	NO_SUCH_FILE                       = "no such file exists"
	FAILED_TO_READ_THE_FILE_FROM_STORE = "failed to read the file from store"
	OPERATION_FAILED                   = "operation failed"
)

type MinIOClientInterface interface {
	CheckFileExistInBucket(ctx context.Context, bucketName, objectName string, searchOptions minio.StatObjectOptions) (minio.ObjectInfo, error)
	CreateBucket(ctx context.Context, bucketName string) error
	CheckIfBucketExist(ctx context.Context, bucketName string) (bool, error)
	PutFileIntoMinIOBucket(ctx context.Context, minioPutModel models.MinIOPutModel) (minio.UploadInfo, error)
	GetMinIOObject(ctx context.Context, bucketName, objectName string, getOptions minio.GetObjectOptions) (*minio.Object, error)
	RemoveObjectFromMinIO(ctx context.Context, bucketName, objectName string, removeOptions minio.RemoveObjectOptions) error
	MigrateObjectAndDelete(ctx context.Context, minioMigrateModel models.MinIOMigrateModel) (minio.UploadInfo, error)
	GetObjectInByteArray(ctx context.Context, bucketName, objectName string, eventOptions MinIOEventObjectOptions) ([]byte, error)
	PresignedGetObject(ctx context.Context, bucketName, objectName string, expires time.Duration, reqParams map[string][]string) (*url.URL, error)
	PutMultipleObjects(ctx context.Context, bucketName, objectName string, reader io.Reader, size int64, options minio.PutObjectOptions) (minio.UploadInfo, error)
	ListMultipleObjects(ctx context.Context, bucketName string) ([]string, error)
	GetObjectMetadata(ctx context.Context, bucketName, objectName string) (minio.ObjectInfo, error)
	MigrateObject(ctx context.Context, minioMigrateModel models.MinIOMigrateModel) (minio.UploadInfo, error)
}

type MinIOClientModel struct {
	MinIOClient *minio.Client
}

type MinIOEventObjectOptions struct {
	searchOptions minio.StatObjectOptions
	getOptions    minio.GetObjectOptions
	// putOptions    minio.PutObjectOptions
	// removeOptions minio.RemoveObjectOptions
}

var minioC *minio.Client

func newMinIOClientModel(minioInfo models.MinIOInfo) (MinIOClientInterface, error) {
	// Initialize minio client object.
	minioClient, err := minio.New(minioInfo.Url, &minio.Options{
		Creds:  credentials.NewStaticV4(minioInfo.AccessKey, minioInfo.SecretKey, ""),
		Secure: minioInfo.UseSSL,
	})
	if err != nil {
		return nil, err
	}
	minioC = minioClient
	return &MinIOClientModel{
		MinIOClient: minioClient,
	}, nil
}

func GetMinioClient() *minio.Client {
	return minioC
}

func BuildMinIOSetup(ctx context.Context, minioInfo models.MinIOInfo) (MinIOClientInterface, error) {
	minioClientModel, err := newMinIOClientModel(minioInfo)
	if err != nil {
		return nil, err
	}

	// Commented because we should not be create bucket automatically to prevent the unwanted bucket creation.
	// entityTypeAndBucketsMap := GetAllBucketsAndEntityType(minioInfo)
	// for _, bucketName := range entityTypeAndBucketsMap {
	// 	exists, errBucketExists := minioClientModel.CheckIfBucketExist(ctx, bucketName)

	// 	if errBucketExists != nil {
	// 		return nil, errBucketExists
	// 	}

	// 	if !exists {
	// 		err := minioClientModel.CreateBucket(ctx, bucketName)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 	}
	// }
	return minioClientModel, nil
}

func (m MinIOClientModel) CheckFileExistInBucket(ctx context.Context, bucketName, objectName string, searchOptions minio.StatObjectOptions) (minio.ObjectInfo, error) {
	return m.MinIOClient.StatObject(ctx, bucketName, objectName, searchOptions)
}

func (m MinIOClientModel) CreateBucket(ctx context.Context, bucketName string) error {
	return m.MinIOClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
}

func (m MinIOClientModel) CheckIfBucketExist(ctx context.Context, bucketName string) (bool, error) {
	return m.MinIOClient.BucketExists(ctx, bucketName)
}

func (m MinIOClientModel) PutFileIntoMinIOBucket(ctx context.Context, minioPutModel models.MinIOPutModel) (minio.UploadInfo, error) {
	reader, fileLength, err := prepareMinIOParamsForPut(minioPutModel.EncodedString)
	if err != nil {
		return minio.UploadInfo{}, errors.New(FAILED_TO_PREPARE_PARAMS)
	}
	info, err := m.MinIOClient.PutObject(ctx, minioPutModel.BucketName, minioPutModel.ObjectName, reader, int64(fileLength), minioPutModel.PutOptions)
	if err != nil {
		return minio.UploadInfo{}, errors.New(FAILED_TO_STORE_OBJECT)
	}
	return info, nil

}
func (m MinIOClientModel) GetMinIOObject(ctx context.Context, bucketName, objectName string, getOptions minio.GetObjectOptions) (*minio.Object, error) {
	return m.MinIOClient.GetObject(ctx, bucketName, objectName, getOptions)
}
func (m MinIOClientModel) RemoveObjectFromMinIO(ctx context.Context, bucketName, objectName string, removeOptions minio.RemoveObjectOptions) error {
	return m.MinIOClient.RemoveObject(ctx, bucketName, objectName, removeOptions)

}
func (m MinIOClientModel) MigrateObjectAndDelete(ctx context.Context, minioMigrateModel models.MinIOMigrateModel) (minio.UploadInfo, error) {

	// get source bucket name by entity type
	sourceBucketName, err := GetBucketNameByEntityType(minioMigrateModel.SourceBucketEntityType, minioMigrateModel.MinIOInfo)
	if sourceBucketName == "" || err != nil {
		return minio.UploadInfo{}, err
	}

	// get destination bucket name by entity type
	destinationBucketName, err := GetBucketNameByEntityType(minioMigrateModel.DestinationBucketEntityType, minioMigrateModel.MinIOInfo)
	if destinationBucketName == "" || err != nil {
		return minio.UploadInfo{}, err
	}

	object, err := m.GetMinIOObject(ctx, sourceBucketName, minioMigrateModel.ObjectName, minioMigrateModel.SourceGetOptions)
	if err != nil {
		return minio.UploadInfo{}, err
	}
	byteArray, err := convertObjectToByteArray(object)

	if err != nil {
		return minio.UploadInfo{}, err
	}
	fileLength := len(string(byteArray))
	info, err := m.MinIOClient.PutObject(ctx, destinationBucketName, minioMigrateModel.ObjectName, object, int64(fileLength), minioMigrateModel.DestinationPutOptions)

	if err != nil {
		return minio.UploadInfo{}, err
	}
	err = m.RemoveObjectFromMinIO(ctx, sourceBucketName, minioMigrateModel.ObjectName, minioMigrateModel.SourceRemoveOptions)
	if err != nil {
		return minio.UploadInfo{}, err
	}
	return info, nil
}

func (m MinIOClientModel) GetObjectInByteArray(ctx context.Context, bucketName, objectName string, eventOptions MinIOEventObjectOptions) ([]byte, error) {
	_, err := m.CheckFileExistInBucket(ctx, bucketName, objectName, eventOptions.searchOptions)
	if err != nil {
		return nil, errors.New(NO_SUCH_FILE)
	}
	object, err := m.GetMinIOObject(ctx, bucketName, objectName, eventOptions.getOptions)
	if err != nil {
		return nil, errors.New(FAILED_TO_STORE_OBJECT)
	}
	defer object.Close()
	return convertObjectToByteArray(object)

}

func prepareMinIOParamsForPut(encodedString string) (io.Reader, int, error) {
	byteArray, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		return nil, 0, errors.New(FAILED_TO_DECODE)
	}
	reader := bytes.NewReader(byteArray)
	fileLength := len(string(byteArray))
	return reader, fileLength, nil
}

func convertObjectToByteArray(object *minio.Object) ([]byte, error) {
	var buffer bytes.Buffer
	_, err := buffer.ReadFrom(object)
	if err != nil {
		return nil, errors.New(FAILED_TO_READ_THE_FILE_FROM_STORE)
	}
	byteArray := buffer.Bytes()
	return byteArray, nil
}

func GetAllBucketsAndEntityType(minioInfo models.MinIOInfo) map[string]string {
	var entityTypeAndBucketsMap map[string]string = make(map[string]string)
	for _, bucket := range minioInfo.Buckets {
		entityTypeAndBucketsMap[bucket.EntityType] = bucket.Name
	}
	return entityTypeAndBucketsMap
}

func GetBucketNameByEntityType(entityType string, minioInfo models.MinIOInfo) (string, error) {
	for _, bucket := range minioInfo.Buckets {
		if entityType == bucket.EntityType {
			return bucket.Name, nil
		}
	}
	return "", errors.New(NO_SUCH_BUCKET)
}

func (m MinIOClientModel) PresignedGetObject(ctx context.Context, bucketName, objectName string, expires time.Duration, reqParams map[string][]string) (*url.URL, error) {
	// Generate the pre-signed URL.
	presignedURL, err := m.MinIOClient.PresignedGetObject(ctx, bucketName, objectName, expires, reqParams)
	if err != nil {
		return nil, err
	}

	return presignedURL, nil
}

// uploads the object to the MinIO bucket using the PutObject method of the MinIO client without encoding
func (m *MinIOClientModel) PutMultipleObjects(ctx context.Context, bucketName, objectName string, reader io.Reader, size int64, options minio.PutObjectOptions) (minio.UploadInfo, error) {
	info, err := m.MinIOClient.PutObject(ctx, bucketName, objectName, reader, size, options)
	if err != nil {
		return minio.UploadInfo{}, err
	}
	return info, nil
}

// lists objects present in the MinIO bucket
func (m *MinIOClientModel) ListMultipleObjects(ctx context.Context, bucketName string) ([]string, error) {
	objectCh := m.MinIOClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Recursive: true,
	})
	var files []string
	for object := range objectCh {
		if object.Err != nil {
			log.Println(object.Err)
			continue
		}
		files = append(files, object.Key)
	}

	return files, nil
}

// gets metadata of an object
func (m *MinIOClientModel) GetObjectMetadata(ctx context.Context, bucketName, objectName string) (minio.ObjectInfo, error) {
	objectInfo, err := m.MinIOClient.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return minio.ObjectInfo{}, err
	}

	return objectInfo, nil
}

func (m MinIOClientModel) MigrateObject(ctx context.Context, minioMigrateModel models.MinIOMigrateModel) (minio.UploadInfo, error) {

	// get source bucket name by entity type
	sourceBucketName, err := GetBucketNameByEntityType(minioMigrateModel.SourceBucketEntityType, minioMigrateModel.MinIOInfo)
	if sourceBucketName == "" || err != nil {
		return minio.UploadInfo{}, err
	}

	// get destination bucket name by entity type
	destinationBucketName, err := GetBucketNameByEntityType(minioMigrateModel.DestinationBucketEntityType, minioMigrateModel.MinIOInfo)
	if destinationBucketName == "" || err != nil {
		return minio.UploadInfo{}, err
	}

	object, err := m.GetMinIOObject(ctx, sourceBucketName, minioMigrateModel.ObjectName, minioMigrateModel.SourceGetOptions)
	if err != nil {
		return minio.UploadInfo{}, err
	}
	byteArray, err := convertObjectToByteArray(object)

	if err != nil {
		return minio.UploadInfo{}, err
	}
	fileLength := len(string(byteArray))
	info, err := m.MinIOClient.PutObject(ctx, destinationBucketName, minioMigrateModel.ObjectName, object, int64(fileLength), minioMigrateModel.DestinationPutOptions)

	if err != nil {
		return minio.UploadInfo{}, err
	}
	return info, nil
}
