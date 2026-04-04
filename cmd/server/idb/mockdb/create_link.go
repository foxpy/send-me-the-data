package mockdb

import (
	"fmt"
	"reflect"

	"github.com/foxpy/send-me-the-data/cmd/server/idb"
)

func (d *MockDB) CreateLink(name, externalKey string, userDownloadable bool, maxFileSize uint64) error {
	if len(d.expectedCreateLinkCalls) == 0 {
		panic("must mock expected CreateLink() call")
	}

	mockedCall := d.expectedCreateLinkCalls[0]
	d.expectedCreateLinkCalls = d.expectedCreateLinkCalls[1:]

	actual := idb.Link{
		Name:             name,
		ExternalKey:      externalKey,
		UserDownloadable: userDownloadable,
		MaxFileSize:      maxFileSize,
	}

	if !reflect.DeepEqual(mockedCall.link, actual) {
		return fmt.Errorf("expected CreateLink(%v), got CreateLink(%v)", mockedCall.link, actual)
	}

	if mockedCall.resultFunc == nil {
		return nil
	}

	return mockedCall.resultFunc()
}

func (d *MockDB) MockExpectedCreateLinkCall(name, externalKey string, userDownloadable bool, maxFileSize uint64, mockedResult func() error) {
	link := idb.Link{
		Name:             name,
		ExternalKey:      externalKey,
		UserDownloadable: userDownloadable,
		MaxFileSize:      maxFileSize,
	}
	d.expectedCreateLinkCalls = append(d.expectedCreateLinkCalls, CreateLinkCall{
		link:       link,
		resultFunc: mockedResult,
	})
}
