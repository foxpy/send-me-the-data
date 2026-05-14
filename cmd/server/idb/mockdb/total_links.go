package mockdb

func (d *MockDB) TotalLinks() (uint64, error) {
	if d.totalLinksResponse == nil {
		panic("must mock TotalLinks() response")
	}

	return *d.totalLinksResponse, nil
}

func (d *MockDB) MockTotalLinksResponse(response uint64) {
	d.totalLinksResponse = &response
}
