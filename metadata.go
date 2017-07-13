package partnerchannel

import (
	"time"
)

type MetadataRequest struct {
	ClientTransactionID string
}

type MetadataResponse struct {
	MetadataRequest
	EPPCode            int
	ResponseCode       string
	ResponseTime       *time.Time
	TransactionID      int64
	ValidationResponse *ValidationResponse
}

type ValidationResponse struct {
	Metadata  MetadataResponse
	EntryList []ValidationEntry
}

type ValidationEntry struct {
	CheckName              string
	Message                string // ToDo: ENUM
	ParameterItemExtension string
	ParameterItemName      string
	ParameterName          string
	Type                   string // ToDo: ENUM
}

func (metadata *MetadataResponse) Check() bool {
	return metadata.ResponseCode == "success"
}
