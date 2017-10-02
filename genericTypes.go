package partnerchannel

type IdName struct {
	ID   int
	Name string
}

// genericRequest represents the structure of the default request format
type genericRequest struct {
	Data     interface{}
	ObjectID int
	Metadata *MetadataRequest
}

// genericResponse represents the structure of the default response format
type genericResponse struct {
	Metadata *MetadataResponse
}
