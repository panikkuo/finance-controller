package http

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	gsheets "github.com/panikkuo/finance-controller/integration/google_sheets"
	utils "github.com/panikkuo/finance-controller/utils"
)

func parse(data map[string]interface{}) (string, string, string, error) {
	operationType, err := utils.GetFieldFromMapAsString(data, "type", true)
	if err != nil {
		return "", "", "", err
	}
	totalSum, err := utils.GetFieldFromMapAsString(data, "total", true)

	if err != nil {
		return "", "", "", err
	}

	category, err := utils.GetFieldFromMapAsString(data, "category", false)

	return totalSum, operationType, category, nil
}

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

	totalSum, operationType, category, err := parse(data)

	if err != nil {
		log.Printf("err.Error(): %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = gsheets.AdjustBalance(operationType, totalSum, category)

	if err != nil {
		log.Printf("err.Error(): %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
