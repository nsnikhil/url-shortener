package reporters

import "github.com/stretchr/testify/mock"

type MockStatsDClient struct {
	mock.Mock
}

func (msc *MockStatsDClient) ReportAttempt(bucket string) {
	msc.Called(bucket)
}

func (msc *MockStatsDClient) ReportSuccess(bucket string) {
	msc.Called(bucket)
}

func (msc *MockStatsDClient) ReportFailure(bucket string) {
	msc.Called(bucket)
}
