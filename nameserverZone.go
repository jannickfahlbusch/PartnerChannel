package partnerchannel

import (
	"errors"
	"net/http"

	"time"

	"gitlab.com/jannickfahlbusch/partnerchannel/recordType"
)

type NameserverZone interface {
	AddRecord(*AddRecordRequest, *MetadataRequest) (*MetadataResponse, error)
	RemoveRecord(int, *MetadataRequest) (*MetadataResponse, error)
	Info(int, *MetadataRequest) (*InfoResponse, *MetadataResponse, error)
	//UpdateRecord(*UpdateRecordRequest, *MetadataRequest) (*MetadataResponse, error)
}

type nameserverZoneEndpoint struct {
	client *Client
}

type AddRecordRequest struct {
	Aux              int
	Data             string
	Name             string
	NameserverZoneID int
	TTL              int
	Type             recordType.RecordType
}

type InfoResponse struct {
	ExpirationTime  time.Time
	ExpireTime      int
	InsertTime      time.Time
	LastCheckTime   time.Time
	MinimumTTL      int
	NameserverSetID int
	ObjectID        int
	ObjectState     string
	RecordList      []Record
	RefreshTime     int
	Reseller        IdName
	RetryTime       int
	SerialNumber    int
	TTL             int
	XFER            []string
	ZoneName        string
	ZoneNameACE     string
}

type Record struct {
	Aux      int
	Data     string
	Name     string
	ObjectID string
	TTL      int
	Type     recordType.RecordType
}

func (endpoint *nameserverZoneEndpoint) AddRecord(addRecordRequest *AddRecordRequest, metadataRequest *MetadataRequest) (*MetadataResponse, error) {
	request, err := endpoint.client.NewRequest(http.MethodPost, "nameserverZone/addRecord", nil, &genericRequest{
		Data:     addRecordRequest,
		Metadata: metadataRequest,
	})

	if err != nil {
		return nil, err
	}

	response := &genericResponse{}

	_, err = endpoint.client.Do(request, response)
	if err != nil {
		return nil, err
	}

	if !response.Metadata.Check() {
		return nil, errors.New(response.Metadata.ResponseCode)
	}

	return response.Metadata, err
}

func (endpoint *nameserverZoneEndpoint) Info(objectID int, metadataRequest *MetadataRequest) (*InfoResponse, *MetadataResponse, error) {
	request, err := endpoint.client.NewRequest(http.MethodGet, "nameserverZone/info", nil, &genericRequest{
		Metadata: metadataRequest,
		ObjectID: objectID,
	})
	if err != nil {
		return nil, nil, err
	}

	response := &struct {
		genericResponse
		Data *InfoResponse
	}{}

	_, err = endpoint.client.Do(request, response)
	if err != nil {
		return nil, nil, err
	}

	if !response.Metadata.Check() {
		return nil, nil, errors.New(response.Metadata.ResponseCode)
	}

	return response.Data, response.Metadata, err
}

func (endpoint *nameserverZoneEndpoint) RemoveRecord(objectID int, metadataRequest *MetadataRequest) (*MetadataResponse, error) {
	request, err := endpoint.client.NewRequest(http.MethodPost, "nameserverZone/removeRecord", nil, &struct {
		ObjectID int
		Metadata *MetadataRequest
	}{
		ObjectID: objectID,
		Metadata: metadataRequest,
	})

	if err != nil {
		return nil, err
	}

	response := &genericResponse{}

	_, err = endpoint.client.Do(request, response)
	if err != nil {
		return nil, err
	}

	if !response.Metadata.Check() {
		return nil, errors.New(response.Metadata.ResponseCode)
	}

	return response.Metadata, err
}
