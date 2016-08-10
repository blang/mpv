package mpv

import "encoding/json"
import "net/http"

// JSONRequest send to the server.
type JSONRequest struct {
	Command []interface{} `json:"command"`
}

// JSONResponse send from the server.
type JSONResponse struct {
	Err  string      `json:"error"`
	Data interface{} `json:"data"` // May contain float64, bool or string
}

type httpServerHandler struct {
	llclient LLClient
}

// HTTPServerHandler returns a http.Handler to access a client via a lowlevel json-api.
// Register as route on your server:
// 		http.Handle("/mpv", mpv.HTTPHandler(lowlevelclient)
//
// Use api:
// POST http://host/lowlevel `{ "command": ["get_property", "fullscreen"] }`
// Result:`{"error":"success","data":false}`
func HTTPServerHandler(client LLClient) http.Handler {
	return &httpServerHandler{
		llclient: client,
	}
}

func (h *httpServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req JSONRequest
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&req)
	if err != nil {
		http.Error(w, "Can not decode request", http.StatusBadRequest)
		return
	}
	resp, err := h.llclient.Exec(req.Command...)
	if err != nil {
		if err == ErrTimeoutRecv || err == ErrTimeoutSend {
			http.Error(w, "Timeout", http.StatusGatewayTimeout)
			return
		}
		// TODO: Handle error, maybe send json response
		http.Error(w, "Client returned unknown error", http.StatusInternalServerError)
		return
	}
	jsonResp := JSONResponse{
		Err:  resp.Err,
		Data: resp.Data,
	}
	b, err := json.Marshal(jsonResp)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(b)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}
}
