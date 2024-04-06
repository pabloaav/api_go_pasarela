package mockrepository

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockRepositoryAWSStore struct {
	mock.Mock
}

func (mock *MockRepositoryAWSStore) PutObject(ctx context.Context, data []byte, filename, fileType string) error {
	args := mock.Called(ctx, data, filename, fileType)
	return args.Error(1)
}
func (mock *MockRepositoryAWSStore) GetObject(ctx context.Context, key string) (string, error) {
	args := mock.Called(key)
	resultado := args.String(0)
	return resultado, args.Error(1)
}
