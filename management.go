package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/julienschmidt/httprouter"
)

// configureManagementAPI sets up the management API routes
func (srv *ApiServer) configureManagementAPI(router *httprouter.Router) {
	router.GET("/", srv.handleRoot)
	router.GET("/management/apiversions", srv.handleApiVersions)
	router.GET("/management/v1/description", srv.handleDescription)
	router.GET("/management/v1/configureddevices", srv.handleConfiguredDevices)
}

// handleRoot returns the root web page
func (srv *ApiServer) handleRoot(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintln(w, "Alpaca Mi server")
}

// handleApiVersions returns an array of supported Alpaca API version numbers
func (srv *ApiServer) handleApiVersions(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	resp := uint32listResponse{Value: []uint32{1}}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// handleDescription returns server information
func (srv *ApiServer) handleDescription(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	bi, _ := debug.ReadBuildInfo()
	resp := managementDescriptionResponse{
		Value: ServerDescription{
			ServerName:          "Mi Switch Controller",
			Manufacturer:        "https://github.com/exploded/",
			ManufacturerVersion: bi.Main.Version,
			Location:            Location,
		},
	}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// handleConfiguredDevices returns the list of configured devices
func (srv *ApiServer) handleConfiguredDevices(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	resp := managementDevicesListResponse{Value: MiGetInit()}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
