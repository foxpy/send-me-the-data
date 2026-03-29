package mockdb

import "github.com/foxpy/send-me-the-data/cmd/server/idb"

func (d *MockDB) AllLinks() ([]idb.Link, error) {
	if d.allLinksResponse == nil {
		panic("must mock AllLinks() response")
	}

	return d.allLinksResponse, nil
}

func (d *MockDB) MockAllLinksResponse(response []idb.Link) {
	d.allLinksResponse = response
}
