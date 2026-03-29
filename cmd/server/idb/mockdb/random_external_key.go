package mockdb

func (d *MockDB) GenerateRandomExternalKey() string {
	if len(d.randomExternalKeyResponses) == 0 {
		panic("must mock GenerateRandomExternalKey() response")
	}

	response := d.randomExternalKeyResponses[0]
	d.randomExternalKeyResponses = d.randomExternalKeyResponses[1:]
	return response
}

func (d *MockDB) MockGenerateRandomExternalKeyResponse(response string) {
	d.randomExternalKeyResponses = append(d.randomExternalKeyResponses, response)
}
