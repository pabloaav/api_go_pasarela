package mockservice

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockStoreService struct {
	mock.Mock
}

func (mock *MockStoreService) DeleteObject(ctx context.Context, key string) error {
	args := mock.Called(ctx, key)
	return args.Error(0)
}

func (mock *MockStoreService) GetObject(ctx context.Context, key string) (string, error) {
	args := mock.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (mock *MockStoreService) PutObject(ctx context.Context, data []byte, filename, fileType string) error {
	args := mock.Called(ctx, data, filename, fileType)
	return args.Error(0)
}

func (mock *MockStoreService) GetObjectS3(ctx context.Context, key string) (result []byte, erro error) {
	return
}
