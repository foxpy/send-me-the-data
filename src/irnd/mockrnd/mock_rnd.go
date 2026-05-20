package mockrnd

import "github.com/foxpy/send-me-the-data/src/irnd"

type MockRND struct {
	publicIDResponses     []string
	sessionTokenResponses []string
}

var _ irnd.Random = &MockRND{}

func NewMockRND() *MockRND {
	return &MockRND{}
}
