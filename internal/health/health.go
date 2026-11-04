package health

import (
	"encoding/json"
	"net/http"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (hh *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	res, _ := json.Marshal(map[string]string{"status": "ok"})
	w.Write(res)
}
