package mockrnd

func (r *MockRND) PublicID() string {
	if len(r.publicIDResponses) == 0 {
		panic("must mock PublicID() response")
	}

	response := r.publicIDResponses[0]
	r.publicIDResponses = r.publicIDResponses[1:]
	return response
}

func (r *MockRND) MockPublicIDResponse(response string) {
	r.publicIDResponses = append(r.publicIDResponses, response)
}
