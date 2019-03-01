package http

import (
	"encoding/json"
	"log"
	"net/http"
)

type clusterHandler struct {
	*Server
}

func (h *clusterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	m := h.Members()
	b, e := json.Marshal(m)
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(e)
		return
	}
	w.Write(b)
}

func (s *Server) clusterHandler() http.Handler {
	return &clusterHandler{s}
}
