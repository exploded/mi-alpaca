package main

// alpacaResponse contains the common ASCOM Alpaca response fields
type alpacaResponse struct {
	ClientTransactionID uint32 `json:"ClientTransactionID"`
	ServerTransactionID uint32 `json:"ServerTransactionID"`
	ErrorNumber         int32  `json:"ErrorNumber"`
	ErrorMessage        string `json:"ErrorMessage"`
}

// stringResponse returns a string value
type stringResponse struct {
	Value string `json:"Value"`
	alpacaResponse
}

// stringlistResponse returns an array of strings
type stringlistResponse struct {
	Value []string `json:"Value"`
	alpacaResponse
}

// booleanResponse returns a boolean value
type booleanResponse struct {
	Value bool `json:"Value"`
	alpacaResponse
}

// int32Response returns an int32 value
type int32Response struct {
	Value int32 `json:"Value"`
	alpacaResponse
}

// doubleResponse returns a double (int64) value
type doubleResponse struct {
	Value int64 `json:"Value"`
	alpacaResponse
}

// uint32listResponse returns an array of uint32 values
type uint32listResponse struct {
	Value []uint32 `json:"Value"`
	alpacaResponse
}

// putResponse is the standard response for PUT operations
type putResponse struct {
	alpacaResponse
}

// managementDevicesListResponse returns the list of configured devices
type managementDevicesListResponse struct {
	Value []DeviceConfiguration `json:"Value"`
	alpacaResponse
}

// DeviceConfiguration describes a single ASCOM Alpaca device
type DeviceConfiguration struct {
	DeviceName   string `json:"DeviceName"`
	DeviceType   string `json:"DeviceType"`
	DeviceNumber uint32 `json:"DeviceNumber"`
	UniqueID     string `json:"UniqueID"`
}

// managementDescriptionResponse returns server description information
type managementDescriptionResponse struct {
	Value ServerDescription `json:"Value"`
	alpacaResponse
}

// ServerDescription contains information about the ASCOM Alpaca server
type ServerDescription struct {
	ServerName          string `json:"ServerName"`
	Manufacturer        string `json:"Manufacturer"`
	ManufacturerVersion string `json:"ManufacturerVersion"`
	Location            string `json:"Location"`
}
