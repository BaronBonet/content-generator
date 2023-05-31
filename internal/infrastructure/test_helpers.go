package infrastructure

import "github.com/stretchr/testify/mock"

// TearDownAdapters resets the mocks so that they can be reused
func TearDownAdapters(adapters ...*mock.Mock) {
	for _, adapter := range adapters {
		adapter.ExpectedCalls = []*mock.Call{}
		adapter.Calls = []mock.Call{}
	}
}
