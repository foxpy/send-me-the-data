package mockdb

import "github.com/foxpy/send-me-the-data/src/idb"

func (d *MockDB) ListLinks(offset, limit uint) ([]idb.Link, error) {
	links, ok := d.listLinksResponses[struct{ offset, limit uint }{offset, limit}]
	if !ok {
		panic("must mock ListLinks() response")
	}

	return links, nil
}

func (d *MockDB) MockListLinksResponse(offset, limit uint, response []idb.Link) {
	d.listLinksResponses[struct{ offset, limit uint }{offset, limit}] = response
}
