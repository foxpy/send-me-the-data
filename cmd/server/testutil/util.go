package testutil

import (
	"time"

	"github.com/foxpy/send-me-the-data/cmd/server/ifs"
)

type LinkFiles struct {
	Name  string
	Files []ifs.File
}

var MockTime = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
var MockTimeMilli = uint64(MockTime.UTC().UnixMilli())
