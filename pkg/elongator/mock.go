package elongator

import "github.com/stretchr/testify/mock"

type MockElongator struct {
	mock.Mock
}

func (mock *MockElongator) Elongate(hash string) (string, error) {
	args := mock.Called(hash)
	return args.String(0), args.Error(1)
}
