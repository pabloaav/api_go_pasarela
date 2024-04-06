package mockrepository

import (
	"github.com/stretchr/testify/mock"
)

type MockRepositoryApiLink struct {
	mock.Mock
}

func (mock *MockRepositoryApiLink) PutApilinkCierrelote(listaDebinesId []string) error {
	return nil
}
