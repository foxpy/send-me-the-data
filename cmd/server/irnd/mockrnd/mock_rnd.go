package mockrnd

import "github.com/foxpy/send-me-the-data/cmd/server/irnd"

type MockRND struct {
	publicIDResponses     []string
	sessionTokenResponses []string
}

var _ irnd.Random = &MockRND{}

func NewMockRND() *MockRND {
	return &MockRND{}
}
