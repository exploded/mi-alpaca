package main

import (
	"encoding/json"
	"net/http"
	"runtime/debug"

	"github.com/julienschmidt/httprouter"
)

func (srv *ApiServer) configureCommonAPI(router *httprouter.Router) {
	// ASCOM Methods Common To All Devices
	router.PUT("/api/v1/switch/1/action", srv.handleNotSupported)
	router.PUT("/api/v1/switch/1/commandblind", srv.handleNotSupported)
	router.PUT("/api/v1/switch/1/commandbool", srv.handleNotSupported)
	router.PUT("/api/v1/switch/1/commandstring", srv.handleNotSupported)

	router.GET("/api/v1/switch/1/connected", srv.handleConnected)
	router.PUT("/api/v1/switch/1/connected", srv.handleConnect)

	router.GET("/api/v1/switch/1/description", srv.handleDescriptionCommon)
	router.GET("/api/v1/switch/1/driverinfo", srv.handleDriverinfo)
	router.GET("/api/v1/switch/1/driverversion", srv.handleDriverVersion)
	router.GET("/api/v1/switch/1/interfaceversion", srv.handleInterfaceVersion)
	router.GET("/api/v1/switch/1/name", srv.handleName)
	router.GET("/api/v1/switch/1/supportedactions", srv.handleSupportedActions)
}

// ASCOM Common API handlers

// Retrieves the connected state of the device
func (srv *ApiServer) handleConnected(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resp := booleanResponse{
		Value: MiGetConnected(),
	}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Sets the connected state of the device
func (srv *ApiServer) handleConnect(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	connected, err := getConnectedFromRequest(r)
	if err != nil {
		resp := stringResponse{Value: err.Error()}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	MiSetConnect(connected)

	resp := stringResponse{Value: ""}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Returns the description of the device
func (srv *ApiServer) handleDescriptionCommon(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resp := stringResponse{Value: "Xiaomi Mi Smart Plug Switch Controller"}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Returns driver information
func (srv *ApiServer) handleDriverinfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resp := stringResponse{Value: "Xiaomi Mi Alpaca Switch Driver https://github.com/exploded/mi-ascom-alpaca"}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Returns the driver version
func (srv *ApiServer) handleDriverVersion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	bi, _ := debug.ReadBuildInfo()
	resp := stringResponse{Value: bi.Main.Version}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Returns the ASCOM device interface version (Switch v2)
func (srv *ApiServer) handleInterfaceVersion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resp := int32Response{Value: 2}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Returns the device name
func (srv *ApiServer) handleName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resp := stringResponse{Value: "Mi Switch Controller"}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Returns the list of supported actions (none for this driver)
func (srv *ApiServer) handleSupportedActions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resp := stringlistResponse{Value: []string{}}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
