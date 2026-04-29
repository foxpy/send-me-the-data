package mockdb

func (d *MockDB) GenerateRandomPublicID() string {
	if len(d.randomPublicIDResponses) == 0 {
		panic("must mock GenerateRandomPublicID() response")
	}

	response := d.randomPublicIDResponses[0]
	d.randomPublicIDResponses = d.randomPublicIDResponses[1:]
	return response
}

func (d *MockDB) MockGenerateRandomPublicIDResponse(response string) {
	d.randomPublicIDResponses = append(d.randomPublicIDResponses, response)
}
