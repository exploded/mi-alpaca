package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type ApiServer struct {
	ApiPort             uint32
	ServerTransactionID uint32
}

func NewApiServer(apiPort uint32) *ApiServer {
	return &ApiServer{
		ApiPort: apiPort,
	}
}

func (srv *ApiServer) Start() {
	router := httprouter.New()
	srv.configureManagementAPI(router)
	srv.configureCommonAPI(router)
	srv.configureSwitchAPI(router)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", srv.ApiPort), router))
}

func (srv *ApiServer) handleNotSupported(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resp := stringResponse{Value: "Command not supported"}
	srv.prepareAlpacaResponse(r, &resp.alpacaResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(resp)
}

func (srv *ApiServer) prepareAlpacaResponse(r *http.Request, resp *alpacaResponse) {
	ctid := getClientTransactionId(r)
	if ctid < 0 {
		ctid = 0
	}
	srv.ServerTransactionID++
	resp.ClientTransactionID = uint32(ctid)
	resp.ServerTransactionID = srv.ServerTransactionID
}

func (srv *ApiServer) validAlpacaRequest(r *http.Request) bool {
	cid := getClientId(r)
	if cid < 0 {
		return false
	}
	ctidv := getClientTransactionId(r)
	return ctidv >= 0
}

func getClientId(r *http.Request) int {
	var cidv string
	if r.Method == "GET" {
		cidv = r.URL.Query().Get("ClientID")
		if cidv == "" {
			cidv = r.URL.Query().Get("clientid")
		}
	} else {
		cidv = r.PostFormValue("ClientID")
		if cidv == "" {
			cidv = r.PostFormValue("clientid")
		}
	}
	if cidv == "" {
		return -1
	}
	cid, err := strconv.Atoi(cidv)
	if err != nil || cid < 0 {
		return -1
	}
	return cid
}

func getClientTransactionId(r *http.Request) int {
	var ctidv string
	if r.Method == "GET" {
		ctidv = r.URL.Query().Get("ClientTransactionID")
		if ctidv == "" {
			ctidv = r.URL.Query().Get("clienttransactionid")
		}
	} else {
		ctidv = r.PostFormValue("ClientTransactionID")
		if ctidv == "" {
			ctidv = r.PostFormValue("clienttransactionid")
		}
	}
	if ctidv == "" {
		return -1
	}
	ctid, err := strconv.Atoi(ctidv)
	if err != nil || ctid < 0 {
		return -1
	}
	return ctid
}

func getIdFromRequest(r *http.Request) (int32, error) {
	var sid string
	if r.Method == "GET" {
		sid = r.URL.Query().Get("Id")
	} else {
		sid = r.PostFormValue("Id")
	}
	if sid == "" {
		return -1, errors.New("id parameter missing")
	}

	iid, err := strconv.ParseInt(sid, 10, 32)
	if err != nil {
		return -1, errors.New("id parameter not numeric")
	}
	if iid < 0 {
		return -1, errors.New("id parameter out of range")
	}
	return int32(iid), nil
}

func getValueFromRequest(r *http.Request) (int64, error) {
	if r.Method != "PUT" {
		return -1, errors.New("expected PUT")
	}
	svalue := r.PostFormValue("Value")
	if svalue == "" {
		return -1, errors.New("value parameter missing")
	}
	ivalue, err := strconv.ParseInt(svalue, 10, 64)
	if err != nil {
		return -1, errors.New("value parameter not numeric")
	}
	if ivalue < 0 {
		return -1, errors.New("value parameter out of range")
	}
	return ivalue, nil
}

func getPositionFromRequest(r *http.Request) (int32, error) {
	if r.Method != "PUT" {
		return -1, errors.New("expected PUT")
	}
	sposition := r.PostFormValue("Position")
	if sposition == "" {
		return -1, errors.New("position parameter missing")
	}
	ivalue, err := strconv.ParseInt(sposition, 10, 32)
	if err != nil {
		return -1, errors.New("position parameter not numeric")
	}
	if ivalue == 0 {
		return -1, errors.New("position parameter out of range")
	}
	return int32(ivalue), nil
}

func getSwitchStateFromRequest(r *http.Request) (bool, error) {
	var sstate string
	if r.Method == "GET" {
		sstate = r.URL.Query().Get("State")
	} else {
		sstate = r.PostFormValue("State")
	}
	if sstate == "" {
		return false, errors.New("state parameter missing")
	}
	return strconv.ParseBool(sstate)
}

func getSwitchNameFromRequest(r *http.Request) (string, error) {
	var sname string
	if r.Method == "GET" {
		sname = r.URL.Query().Get("Name")
	} else {
		sname = r.PostFormValue("Name")
	}
	if sname == "" {
		return "", errors.New("name parameter missing")
	}
	return sname, nil
}

func getConnectedFromRequest(r *http.Request) (bool, error) {
	connect, err := strconv.ParseBool(r.PostFormValue("Connected"))
	if err != nil {
		return false, errors.New("connected parameter missing or invalid")
	}
	return connect, nil
}
