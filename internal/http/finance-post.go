package http

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	//gsheets "github.com/panikkuo/finance-controller/integration/google_sheets"
)

func HandlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Can't parse request body: %v", err)
		http.Error(w, "Can't parse request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		log.Printf("Can't parse JSON: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	//gsheets.AdjustBalance(data["type"], data["total"])

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
