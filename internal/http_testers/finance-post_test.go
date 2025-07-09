package http_test

import (
	"bytes"
	"net/http"
	"testing"
)

func TestIncomeOkTest(t *testing.T) {
	jsonData := `{"type": "income", "total": "100.50"}`
	reqBody := bytes.NewBuffer([]byte(jsonData))

	resp, err := http.Post("http://127.0.0.1:8080/finance", "application/json", reqBody)
	if err != nil {
		t.Fatalf("Ошибка запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался 200, получил %d\nТело запроса: %v", resp.StatusCode, resp.Body)
	}
}

func TestExpenseOkTest(t *testing.T) {
	jsonData := `{"type": "expense", "total": "100.50"}`
	reqBody := bytes.NewBuffer([]byte(jsonData))

	resp, err := http.Post("http://127.0.0.1:8080/finance", "application/json", reqBody)
	if err != nil {
		t.Fatalf("Ошибка запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался 200, получил %d\nТело запроса: %v", resp.StatusCode, resp.Body)
	}
}

func TestExpenseWithCategoryOkTest(t *testing.T) {
	jsonData := `{"type": "expense", "total": "200", "category" : "Шлюхи"}`
	reqBody := bytes.NewBuffer([]byte(jsonData))

	resp, err := http.Post("http://127.0.0.1:8080/finance", "application/json", reqBody)
	if err != nil {
		t.Fatalf("Ошибка запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался 200, получил %d\nТело запроса: %v", resp.StatusCode, resp.Body)
	}
}
