package api

import (
	"encoding/json"
	"net/http"
)

func (api *API) FetchAllClass(w http.ResponseWriter, r *http.Request) {
	classes, err := api.classService.FetchAll()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(classes)
}
