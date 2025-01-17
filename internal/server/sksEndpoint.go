package server

import (
	"encoding/json"
	"fmt"
	"github.com/quan-to/chevron/internal/keymagic"
	"github.com/quan-to/chevron/internal/models"
	"github.com/quan-to/chevron/pkg/interfaces"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/quan-to/slog"
)

type SKSEndpoint struct {
	sm  interfaces.SecretsManager
	gpg interfaces.PGPManager
	log slog.Instance
}

// MakeSKSEndpoint creates a handler for SKS Server Endpoint
func MakeSKSEndpoint(log slog.Instance, sm interfaces.SecretsManager, gpg interfaces.PGPManager) *SKSEndpoint {
	if log == nil {
		log = slog.Scope("SKS")
	} else {
		log = log.SubScope("SKS")
	}

	return &SKSEndpoint{
		sm:  sm,
		gpg: gpg,
		log: log,
	}
}

func (sks *SKSEndpoint) AttachHandlers(r *mux.Router) {
	r.HandleFunc("/getKey", sks.getKey).Methods("GET")
	r.HandleFunc("/searchByName", sks.searchByName).Methods("GET")
	r.HandleFunc("/searchByFingerPrint", sks.searchByFingerPrint).Methods("GET")
	r.HandleFunc("/searchByEmail", sks.searchByEmail).Methods("GET")
	r.HandleFunc("/search", sks.search).Methods("GET")
	r.HandleFunc("/addKey", sks.addKey).Methods("POST")
}

func (sks *SKSEndpoint) getKey(w http.ResponseWriter, r *http.Request) {
	ctx := wrapContextWithRequestID(r)
	log := wrapLogWithRequestID(sks.log, r)
	InitHTTPTimer(log, r)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	q := r.URL.Query()

	fingerPrint := q.Get("fingerPrint")
	key, _ := sks.gpg.GetPublicKeyASCII(ctx, fingerPrint)

	if key == "" {
		NotFound("fingerPrint", fmt.Sprintf("Key with fingerPrint %s was not found", fingerPrint), w, r, log)
		return
	}

	w.Header().Set("Content-Type", models.MimeText)
	w.WriteHeader(200)
	n, _ := w.Write([]byte(key))
	LogExit(log, r, 200, n)
}

func (sks *SKSEndpoint) searchByName(w http.ResponseWriter, r *http.Request) {
	log := wrapLogWithRequestID(sks.log, r)
	InitHTTPTimer(log, r)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	q := r.URL.Query()
	name := q.Get("name")
	pageStartS := q.Get("pageStart")
	pageEndS := q.Get("pageEnd")

	if name == "" {
		InvalidFieldData("name", "you should provide a name", w, r, log)
		return
	}

	pageStart, err := strconv.ParseInt(pageStartS, 10, 32)
	if err != nil {
		pageStart = models.DefaultPageStart
	}

	pageEnd, err := strconv.ParseInt(pageEndS, 10, 32)
	if err != nil {
		pageEnd = models.DefaultPageEnd
	}

	gpgKeys, err := keymagic.PKSSearchByName(name, int(pageStart), int(pageEnd))

	if err != nil {
		NotFound("name", err.Error(), w, r, log)
		return
	}

	bodyData, err := json.Marshal(gpgKeys)

	if err != nil {
		InternalServerError("There was an internal server error. Please try again", nil, w, r, log)
		return
	}

	w.Header().Set("Content-Type", models.MimeJSON)
	w.WriteHeader(200)
	n, _ := w.Write(bodyData)
	LogExit(log, r, 200, n)
}

func (sks *SKSEndpoint) searchByFingerPrint(w http.ResponseWriter, r *http.Request) {
	log := wrapLogWithRequestID(sks.log, r)
	InitHTTPTimer(log, r)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	q := r.URL.Query()
	fingerPrint := q.Get("fingerPrint")
	pageStartS := q.Get("pageStart")
	pageEndS := q.Get("pageEnd")

	if fingerPrint == "" {
		InvalidFieldData("fingerPrint", "you should provide a fingerPrint", w, r, log)
		return
	}

	pageStart, err := strconv.ParseInt(pageStartS, 10, 32)
	if err != nil {
		pageStart = models.DefaultPageStart
	}

	pageEnd, err := strconv.ParseInt(pageEndS, 10, 32)
	if err != nil {
		pageEnd = models.DefaultPageEnd
	}

	gpgKeys, err := keymagic.PKSSearchByFingerPrint(fingerPrint, int(pageStart), int(pageEnd))

	if err != nil {
		NotFound("fingerPrint", err.Error(), w, r, log)
		return
	}

	bodyData, err := json.Marshal(gpgKeys)

	if err != nil {
		InternalServerError("There was an internal server error. Please try again", nil, w, r, log)
		return
	}

	w.Header().Set("Content-Type", models.MimeJSON)
	w.WriteHeader(200)
	n, _ := w.Write(bodyData)
	LogExit(log, r, 200, n)
}

func (sks *SKSEndpoint) searchByEmail(w http.ResponseWriter, r *http.Request) {
	log := wrapLogWithRequestID(sks.log, r)
	InitHTTPTimer(log, r)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	q := r.URL.Query()
	email := q.Get("email")
	pageStartS := q.Get("pageStart")
	pageEndS := q.Get("pageEnd")

	if email == "" {
		InvalidFieldData("email", "you should provide a email", w, r, log)
		return
	}

	pageStart, err := strconv.ParseInt(pageStartS, 10, 32)
	if err != nil {
		pageStart = models.DefaultPageStart
	}

	pageEnd, err := strconv.ParseInt(pageEndS, 10, 32)
	if err != nil {
		pageEnd = models.DefaultPageEnd
	}

	gpgKeys, err := keymagic.PKSSearchByEmail(email, int(pageStart), int(pageEnd))

	if err != nil {
		NotFound("email", err.Error(), w, r, log)
		return
	}

	bodyData, err := json.Marshal(gpgKeys)

	if err != nil {
		InternalServerError("There was an internal server error. Please try again", nil, w, r, log)
		return
	}

	w.Header().Set("Content-Type", models.MimeJSON)
	w.WriteHeader(200)
	n, _ := w.Write(bodyData)
	LogExit(log, r, 200, n)
}

func (sks *SKSEndpoint) search(w http.ResponseWriter, r *http.Request) {
	log := wrapLogWithRequestID(sks.log, r)
	InitHTTPTimer(log, r)

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	q := r.URL.Query()
	valueData := q.Get("valueData")
	pageStartS := q.Get("pageStart")
	pageEndS := q.Get("pageEnd")

	if valueData == "" {
		InvalidFieldData("email", "you should provide a valueData", w, r, log)
		return
	}

	pageStart, err := strconv.ParseInt(pageStartS, 10, 32)
	if err != nil {
		pageStart = models.DefaultPageStart
	}

	pageEnd, err := strconv.ParseInt(pageEndS, 10, 32)
	if err != nil {
		pageEnd = models.DefaultPageEnd
	}

	gpgKeys, err := keymagic.PKSSearch(valueData, int(pageStart), int(pageEnd))

	if err != nil {
		NotFound("valueData", err.Error(), w, r, log)
		return
	}

	bodyData, err := json.Marshal(gpgKeys)

	if err != nil {
		InternalServerError("There was an internal server error. Please try again", nil, w, r, log)
		return
	}

	w.Header().Set("Content-Type", models.MimeJSON)
	w.WriteHeader(200)
	n, _ := w.Write(bodyData)
	LogExit(log, r, 200, n)
}

func (sks *SKSEndpoint) addKey(w http.ResponseWriter, r *http.Request) {
	ctx := wrapContextWithRequestID(r)
	log := wrapLogWithRequestID(sks.log, r)
	InitHTTPTimer(log, r)

	var data models.SKSAddKey

	if !UnmarshalBodyOrDie(&data, w, r, log) {
		return
	}

	defer func() {
		if rec := recover(); rec != nil {
			CatchAllError(rec, w, r, log)
		}
	}()

	status := keymagic.PKSAdd(ctx, data.PublicKey)

	if status != "OK" {
		InvalidFieldData("PublicKey", "Invalid Public Key specified. Check if its in ASCII Armored Format", w, r, log)
		return
	}

	w.Header().Set("Content-Type", models.MimeText)
	w.WriteHeader(200)
	n, _ := w.Write([]byte("OK"))
	LogExit(log, r, 200, n)
}
