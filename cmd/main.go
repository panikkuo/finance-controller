package main

import (
	"log"
	"net/http"
	"os"

	gsheets "github.com/panikkuo/finance-controller/integration/google_sheets"
	myhttp "github.com/panikkuo/finance-controller/internal/http"
)

func main() {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Не удалось прочитать файл учетных данных: %v", err)
	}

	gsheets.InitClient(b, "15Bz08E9MvWbsRE5NwAEWlZGsDO2bsZHa2H05X2ald9g", "finance")
	addr := ":8080"
	router := myhttp.NewRouter()

	log.Printf("Сервер запущен на %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
