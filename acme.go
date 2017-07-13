package partnerchannel

// Acme represents the /acme endpoint
type Acme interface {
	AddChallenge()
	RemoveChallenge()
}

type acmeEndpoint struct {
	client *Client
}

func (endpoint *acmeEndpoint) AddChallenge() {

}

func (endpoint *acmeEndpoint) RemoveChallenge() {

}
