package handler

import (
	"encoding/json"
	"net/http"
)

func response(w http.ResponseWriter, resp interface{}) {
	r, _ := json.Marshal(resp)
	w.Write(r)
}
