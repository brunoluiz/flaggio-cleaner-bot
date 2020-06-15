package linear

import "net/http"

type addHeaderTransport struct {
	T     http.RoundTripper
	token string
}

func (adt *addHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", adt.token)
	return adt.T.RoundTrip(req)
}

func newAddHeaderTransport(token string) *addHeaderTransport {
	return &addHeaderTransport{http.DefaultTransport, token}
}
