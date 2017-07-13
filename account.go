package partnerchannel

import (
	"net/http"
)

// Account represents the /account endpoint
type Account interface {
	Login(*LoginRequest, *MetadataRequest) (*SessionResponse, *MetadataResponse, error)
	Logout(*MetadataRequest) (*MetadataResponse, error)
}

type accountEndpoint struct {
	client *Client
}

type LoginRequest struct {
	Fingerprint  string
	Password     string
	ResellerName string
	UserName     string
}

type SessionResponse struct {
	SessionID string
}

func (endpoint *accountEndpoint) Login(loginRequest *LoginRequest, metadataRequest *MetadataRequest) (*SessionResponse, *MetadataResponse, error) {

	request, err := endpoint.client.NewRequest(http.MethodPost, "account/login", nil, &genericRequest{
		Data:     loginRequest,
		Metadata: metadataRequest,
	})

	if err != nil {
		return nil, nil, err
	}

	response := &struct {
		genericResponse
		Data *SessionResponse
	}{}

	_, err = endpoint.client.Do(request, response)
	if err != nil {
		return nil, nil, err
	}

	if response.Metadata.Check() {
		endpoint.client.sessionToken = response.Data.SessionID
	}

	return response.Data, response.Metadata, err
}

func (endpoint *accountEndpoint) Logout(*MetadataRequest) (*MetadataResponse, error) {
	request, err := endpoint.client.NewRequest(http.MethodPost, "account/logout", nil, nil)
	if err != nil {
		return nil, err
	}

	response := &genericResponse{}

	_, err = endpoint.client.Do(request, response)

	return response.Metadata, err
}
