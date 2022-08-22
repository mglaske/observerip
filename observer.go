package observerip

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

/*
USAGE:

x := observerip.New(5080)
x.SetPassthrough(true)
x.ListenAndServe()

*/

// Since our router is just forwarding all port 80 requests from ObserverIP to
// us, we don't know where to send the data if passthrough is turned on.
// This is the URL that we prepend the response URI with.
var WU_URL string = "http://rtupdate.wunderground.com"
var WU_PASSTHROUGH bool = false

// ObserverIP Handler
type ObHandler struct {
	stationResponse  StationResponse
	endpointResponse EndpointResponse
	server           *http.Server
	port             int
	ctx              context.Context
}

// This one just calls ListenAndServe (does not return)
func (me *ObHandler) ListenAndServe() {
	me.server.ListenAndServe()
}

// This starts the listenAndServe, and returns
func (me *ObHandler) Start() {
	go me.server.ListenAndServe()
}

func (me *ObHandler) SetPassthrough(b bool) {
	WU_PASSTHROUGH = b
}

func (me *ObHandler) SetPassthroughURL(u string) {
	WU_URL = u
}

func (me *ObHandler) Close() {
	me.server.Close()
}

func New(port int) (*ObHandler, error) {
	ctx := context.Background()
	return NewWithContext(port, ctx)
}

func NewWithContext(port int, ctx context.Context) (*ObHandler, error) {
	h := ObHandler{
		port: port,
		ctx:  ctx,
	}
	mux := http.NewServeMux()
	h.server = &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: mux}

	// Actual Weather Underground Endpoints SENT by ObserverIP
	mux.HandleFunc("/weatherstation/updateweatherstation.php", h.stationResponse.Parse)
	mux.HandleFunc("/endpoint", h.endpointResponse.Parse)

	// Our Endpoints for polling stored data
	mux.HandleFunc("/endpoints", h.GetEndpoints)
	mux.HandleFunc("/info", h.GetInfo)

	return &h, nil
}

func (me *ObHandler) GetEndpoints(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(me.endpointResponse)
}

func (me *ObHandler) GetInfo(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(me.stationResponse)
}
