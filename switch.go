package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (srv *ApiServer) configureSwitchAPI(router *httprouter.Router) {
	// ASCOM methods specific to the Switch API
	router.GET("/setup/v1/switch/1/setup", srv.handleSwitchSetup)
	router.GET("/api/v1/switch/1/maxswitch", srv.handleMaxSwitch)
	router.GET("/api/v1/switch/1/canwrite", srv.handleCanWrite)
	router.GET("/api/v1/switch/1/getswitch", srv.handleGetSwitch)
	router.GET("/api/v1/switch/1/getswitchdescription", srv.handleGetSwitchDescription)
	router.GET("/api/v1/switch/1/getswitchname", srv.handleGetSwitchName)
	router.GET("/api/v1/switch/1/getswitchvalue", srv.handleGetSwitchValue)
	router.GET("/api/v1/switch/1/minswitchvalue", srv.handleMinSwitchValue)
	router.GET("/api/v1/switch/1/maxswitchvalue", srv.handleMaxSwitchValue)
	router.PUT("/api/v1/switch/1/setswitch", srv.handleSetSwitch)
	router.PUT("/api/v1/switch/1/setswitchname", srv.handleSetSwitchName)
	router.PUT("/api/v1/switch/1/setswitchvalue", srv.handleSetSwitchValue)
	router.GET("/api/v1/switch/1/switchstep", srv.handleSwitchStep)
}

// handleSwitchSetup returns the setup page
func (srv *ApiServer) handleSwitchSetup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintln(w, "Alpaca Mi switch server")
}

// handleMaxSwitch returns the number of switches (devices numbered from 0 to MaxSwitch - 1)
func (srv *ApiServer) handleMaxSwitch(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	resp := int32Response{Value: NumOnOffSwitch + NumVarSwitch}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// handleCanWrite reports if the switch can be written to
func (srv *ApiServer) handleCanWrite(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sn, err := getIdFromRequest(r)
	if err != nil {
		resp := stringResponse{Value: err.Error()}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := booleanResponse{Value: MiGetCanWrite(sn)}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// handleGetSwitch returns the state of the specified switch as a boolean
func (srv *ApiServer) handleGetSwitch(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sn, err := getIdFromRequest(r)
	if err != nil {
		resp := stringResponse{Value: err.Error()}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	result, err := MiGetOnOff(sn)
	if err != nil {
		resp := stringResponse{Value: err.Error()}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := booleanResponse{Value: result}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// handleGetSwitchDescription returns the description of the specified switch
func (srv *ApiServer) handleGetSwitchDescription(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sn, err := getIdFromRequest(r)
	if err != nil {
		resp := stringResponse{Value: err.Error()}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := stringResponse{Value: MiGetName(sn)}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// handleGetSwitchName returns the name of the specified switch
func (srv *ApiServer) handleGetSwitchName(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sn, err := getIdFromRequest(r)
	if err != nil {
		resp := stringResponse{Value: err.Error()}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := stringResponse{Value: MiGetName(sn)}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// handleGetSwitchValue returns the value of the specified switch
func (srv *ApiServer) handleGetSwitchValue(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sn, err := getIdFromRequest(r)
	if err != nil {
		resp := stringResponse{Value: err.Error()}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := doubleResponse{Value: MiGetValue(sn)}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// handleMinSwitchValue returns the minimum value of the specified switch
func (srv *ApiServer) handleMinSwitchValue(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sn, err := getIdFromRequest(r)
	if err != nil {
		resp := stringResponse{Value: err.Error()}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := doubleResponse{Value: MiGetMin(sn)}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// handleMaxSwitchValue returns the maximum value of the specified switch
func (srv *ApiServer) handleMaxSwitchValue(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sn, err := getIdFromRequest(r)
	if err != nil {
		resp := stringResponse{Value: err.Error()}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := doubleResponse{Value: MiGetMax(sn)}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// handleSetSwitch sets the specified switch to the given state
func (srv *ApiServer) handleSetSwitch(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Print("Set Switch called")
	sn, err := getIdFromRequest(r)
	if err != nil {
		resp := stringResponse{Value: err.Error()}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	log.Printf("Set Switch called for switch number %d", sn)
	sv, err := getSwitchStateFromRequest(r)
	if err != nil {
		resp := stringResponse{Value: err.Error()}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	err = MiSetOnOff(sn, sv)
	if err != nil {
		resp := stringResponse{Value: err.Error()}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := stringResponse{Value: "Success"}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// handleSetSwitchName sets the name of the specified switch
func (srv *ApiServer) handleSetSwitchName(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sn, err := getIdFromRequest(r)
	if err != nil {
		resp := stringResponse{Value: err.Error()}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	sna, err := getSwitchNameFromRequest(r)
	if err != nil {
		resp := stringResponse{Value: err.Error()}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	err = MiSetName(sn, sna)
	if err != nil {
		resp := stringResponse{Value: err.Error()}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := putResponse{}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// handleSetSwitchValue sets the value of the specified switch
func (srv *ApiServer) handleSetSwitchValue(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Print("Set Switch value called")
	sn, err := getIdFromRequest(r)
	if err != nil {
		resp := stringResponse{Value: err.Error()}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	log.Printf("Set Switch value for switch number %d", sn)
	sv, err := getValueFromRequest(r)
	if err != nil {
		resp := stringResponse{Value: err.Error()}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	newState := sv != 0
	err = MiSetOnOff(sn, newState)
	if err != nil {
		resp := stringResponse{Value: err.Error()}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := stringResponse{Value: "Success"}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// handleSwitchStep returns the step size for the specified switch
func (srv *ApiServer) handleSwitchStep(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sn, err := getIdFromRequest(r)
	if err != nil {
		resp := stringResponse{Value: err.Error()}
		srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := doubleResponse{Value: MiGetStep(sn)}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
