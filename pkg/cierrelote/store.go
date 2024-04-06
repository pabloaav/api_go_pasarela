package cierrelote

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/config"
	"github.com/Corrientes-Telecomunicaciones/api_go_pasarela/internal/storage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Store interface {
	/*
		PutObject almacena un archivo en el store y devuelve:
		nil si la operacion fue correncta.
		error si no se pudo realizar la operación
	*/
	PutObject(ctx context.Context, data []byte, filename, fileType string) error
	/*
		DeleteObject mediante el nombre de acceso de un archivo se lo elimina del store.
	*/
	DeleteObject(ctx context.Context, key string) error
	/*
		GetObject mediante el nombre de acceso devuelve del store un archivo.
	*/
	GetObject(ctx context.Context, key string) (string, error)

	GetObjectS3(ctx context.Context, key string) (result []byte, erro error)
}

type store struct {
	storage storage.Storage
}

func NewStore(st storage.Storage) Store {
	return &store{
		storage: st,
	}
}

func (s *store) GetObject(ctx context.Context, key string) (string, error) {
	inputS3 := &s3.ListObjectsInput{
		Bucket:  aws.String(config.AWS_BUCKET),
		MaxKeys: aws.Int64(1000),
		Prefix:  aws.String(key),
	}
	result1, err := s.storage.S3.ListObjects(inputS3)
	if err != nil {
		return "", fmt.Errorf("s3.GetObjectWithContext: %w", err)
	}
	if len(result1.Contents) <= 0 {
		return "", fmt.Errorf("lista de archivos vacia ")
	}
	directorio := time.Now().Local().Format("02-01-2006") + config.DIR_TEMP_NAME
	direcName, err := ioutil.TempDir(".."+config.DOC_CL, directorio)
	if err != nil {
		return "", fmt.Errorf("error creando directorio: %w", err)
	}
	for _, v := range result1.Contents {
		stringArray := strings.Split(*aws.String(*v.Key), "/")
		if stringArray[2] != "" {

			obj, err := s.storage.S3.GetObject(&s3.GetObjectInput{
				Bucket: aws.String(config.AWS_BUCKET),
				Key:    aws.String(*v.Key),
			})
			if err != nil {
				return "", fmt.Errorf("erro al obtener los archivos: %w", err)
			}

			file, err := os.Create(direcName + "/" + stringArray[2])
			if err != nil {
				return "", fmt.Errorf("error al crear archivo temporal: %w", err)
			}
			io.Copy(file, obj.Body)
		}
	}
	return direcName, nil
}

func (s *store) PutObject(ctx context.Context, data []byte, filename, fileType string) error {
	key := s.getFileKey(filename, fileType)
	_, err := s.storage.S3.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Body:   bytes.NewReader(data),
		Bucket: aws.String(config.AWS_BUCKET),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("s3.PutObjectWithContext: %w", err)
	}
	// logs.Info(fmt.Sprintf("documento: %-v", object.String()))
	return nil
}

func (s *store) DeleteObject(ctx context.Context, key string) error {
	_, err := s.storage.S3.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(config.AWS_BUCKET),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("error al intentar borrar los archivos del minio: %w", err)
	}

	return nil
}

func (s *store) getFileKey(fileID string, fileType string) string {
	return fmt.Sprintf("%s.%s", fileID, fileType)
}

func (s *store) GetObjectS3(ctx context.Context, key string) (result []byte, erro error) {

	if len(key) <= 0 {
		erro = fmt.Errorf("lista de archivos vacia")
		return
	}
	inputS3 := &s3.GetObjectInput{
		Bucket: aws.String(config.AWS_BUCKET),
		Key:    aws.String(key),
	}
	object, err := s.storage.S3.GetObject(inputS3)
	if err != nil {
		erro = fmt.Errorf("s3.GetObjectWithContext: %w", err)
		return
	}
	defer object.Body.Close()

	objtoBytes, err := ioutil.ReadAll(object.Body)
	if err != nil {
		erro = fmt.Errorf("erro al leer archivo: %w", err)
		return
	}
	result = objtoBytes

	return
}
