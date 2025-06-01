package models

import "github.com/minio/minio-go/v7"

type Bucket struct {
	EntityType string `json:"entity_type"`
	Name       string `json:"name"`
}

type MinIOInfo struct {
	Url       string   `json:"url"`
	AccessKey string   `json:"access_key"`
	SecretKey string   `json:"secret_key"`
	Api       string   `json:"api"`
	Path      string   `json:"path"`
	UseSSL    bool     `json:"use_ssl"`
	Buckets   []Bucket `json:"buckets"`
}

type MinIOPutModel struct {
	BucketName    string
	ObjectName    string
	EncodedString string
	Metadata      map[string]string
	MinIOClient   *minio.Client
	PutOptions    minio.PutObjectOptions
}

type MinIOUploadModel struct {
	ObjectName    string
	EncodedString string
	PutOptions    minio.PutObjectOptions
	EntityType    string
	MinIOInfo     MinIOInfo
	Module        string
}

type MinIOGetOrDeleteModel struct {
	ObjectName string
	EntityType string
	MinIOInfo  MinIOInfo
}

type MinIOModel struct {
	BucketName  string
	ObjectName  string
	MinIOClient *minio.Client
}

type MinIOMigrateModel struct {
	SourceBucketEntityType      string
	DestinationBucketEntityType string
	ObjectName                  string
	MinIOClient                 *minio.Client
	SourceGetOptions            minio.GetObjectOptions
	DestinationPutOptions       minio.PutObjectOptions
	SourceRemoveOptions         minio.RemoveObjectOptions
	DestinationMetadata         map[string]string
	MinIOInfo                   MinIOInfo
}

type ErrorFileUploadModel struct {
	SrcFileName                 string
	SourceBucketEntityType      string
	DestinationBucketEntityType string
	Module                      string
}
