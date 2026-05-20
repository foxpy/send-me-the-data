package mockrnd

func (r *MockRND) SessionToken() string {
	if len(r.sessionTokenResponses) == 0 {
		panic("must mock SessionToken() response")
	}

	response := r.sessionTokenResponses[0]
	r.sessionTokenResponses = r.sessionTokenResponses[1:]
	return response
}

func (r *MockRND) MockSessionTokenResponse(response string) {
	r.sessionTokenResponses = append(r.sessionTokenResponses, response)
}
