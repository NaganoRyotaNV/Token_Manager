package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setupTestFile() (string, error) {
	file, err := os.CreateTemp("", "test-uploaded.csv")
	if err != nil {
		return "", err
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.WriteAll([][]string{
		{"Project Name", "Token", "Permission", "User ID", "User Name", "Expiry Date"},
		{"Project A", "token456", "read,write", "102", "Bob", "2024/6/30"},
		{"Project B", "token456", "write", "102", "Bob", "2024/6/30"},
		{"Project D", "token444", "read,write", "102", "Bob", "2024/3/30"},
	})
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}

func TestHandleTokens(t *testing.T) {
	tempFilePath, err := setupTestFile()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFilePath)

	originalFilePath := filePath
	filePath = tempFilePath
	defer func() { filePath = originalFilePath }()

	req, err := http.NewRequest("GET", "/api/ttokens", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleTokens)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("ハンドラーが返したステータスコードが正しくありません: %v (期待値: %v)", status, http.StatusOK)
	}

	var actual []Token
	err = json.Unmarshal(rr.Body.Bytes(), &actual)
	if err != nil {
		t.Fatalf("レスポンスのパース中にエラーが発生しました: %v", err)
	}

	expected := []Token{
		{"Project Name", "Token", "Permission", "User ID", "User Name", "Expiry Date"},
		{"Project A", "token456", "read,write", "102", "Bob", "2024/6/30"},
		{"Project B", "token456", "write", "102", "Bob", "2024/6/30"},
		{"Project D", "token444", "read,write", "102", "Bob", "2024/3/30"},
	}

	if !compareTokens(actual, expected) {
		t.Errorf("ハンドラーが返したボディが予期しないものでした: %v (期待値: %v)", actual, expected)
	}
}

func compareTokens(a, b []Token) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestHandleUpload(t *testing.T) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.csv")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte("file content"))
	writer.Close()

	req, err := http.NewRequest("POST", "/api/upload", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleUpload)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("ハンドラーが返したステータスコードが正しくありません: %v (期待値: %v)", status, http.StatusCreated)
	}

	var actual map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &actual)
	if err != nil {
		t.Fatalf("レスポンスのパース中にエラーが発生しました: %v", err)
	}

	expected := map[string]string{"message": "ファイルが正常にアップロードされました"}
	if !compareMaps(actual, expected) {
		t.Errorf("ハンドラーが返したボディが予期しないものでした: %v (期待値: %v)", actual, expected)
	}
}

func compareMaps(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for key := range a {
		if a[key] != b[key] {
			return false
		}
	}
	return true
}

func TestAddToken(t *testing.T) {
	tempFilePath, err := setupTestFile()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFilePath)

	originalFilePath := filePath
	filePath = tempFilePath
	defer func() { filePath = originalFilePath }()

	token := Token{
		ProjectName: "Project E",
		Token:       "token789",
		Permission:  "read",
		UserID:      "103",
		UserName:    "Alice",
		ExpiryDate:  "2025/6/30",
	}
	body, _ := json.Marshal(token)

	req, err := http.NewRequest("POST", "/api/tokens", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleTokens)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("ハンドラーが返したステータスコードが正しくありません: %v (期待値: %v)", status, http.StatusCreated)
	}

	var tokens []Token
	file, err := os.Open(tempFilePath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	for _, record := range records {
		if len(record) == 6 {
			tokens = append(tokens, Token{
				ProjectName: record[0],
				Token:       record[1],
				Permission:  record[2],
				UserID:      record[3],
				UserName:    record[4],
				ExpiryDate:  record[5],
			})
		}
	}

	expected := Token{"Project E", "token789", "read", "103", "Alice", "2025/6/30"}
	found := false
	for _, t := range tokens {
		if t == expected {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("トークンが追加されていません: %v (期待値: %v)", tokens, expected)
	}
}

func TestUpdateToken(t *testing.T) {
	tempFilePath, err := setupTestFile()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFilePath)

	originalFilePath := filePath
	filePath = tempFilePath
	defer func() { filePath = originalFilePath }()

	updateData := struct {
		ProjectName string  `json:"projectName"`
		Tokens      []Token `json:"tokens"`
	}{
		ProjectName: "Project A",
		Tokens: []Token{
			{
				ProjectName: "Project A",
				Token:       "updated_token",
				Permission:  "write",
				UserID:      "102",
				UserName:    "Bob",
				ExpiryDate:  "2024/12/31",
			},
		},
	}
	body, _ := json.Marshal(updateData)

	req, err := http.NewRequest("PUT", "/api/tokens", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleTokens)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("ハンドラーが返したステータスコードが正しくありません: %v (期待値: %v)", status, http.StatusOK)
	}

	var tokens []Token
	file, err := os.Open(tempFilePath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	for _, record := range records {
		if len(record) == 6 {
			tokens = append(tokens, Token{
				ProjectName: record[0],
				Token:       record[1],
				Permission:  record[2],
				UserID:      record[3],
				UserName:    record[4],
				ExpiryDate:  record[5],
			})
		}
	}

	expected := Token{"Project A", "updated_token", "write", "102", "Bob", "2024/12/31"}
	found := false
	for _, t := range tokens {
		if t == expected {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("トークンが更新されていません: %v (期待値: %v)", tokens, expected)
	}
}

func TestDeleteToken(t *testing.T) {
    tempFilePath, err := setupTestFile()
    if err != nil {
        t.Fatal(err)
    }
    defer os.Remove(tempFilePath)

    originalFilePath := filePath
    filePath = tempFilePath
    defer func() { filePath = originalFilePath }()

    req, err := http.NewRequest("DELETE", "/api/tokens?line=2", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(handleTokens)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("ハンドラーが返したステータスコードが正しくありません: %v (期待値: %v)", status, http.StatusOK)
    }

    var tokens []Token
    file, err := os.Open(tempFilePath)
    if err != nil {
        t.Fatal(err)
    }
    defer file.Close()
    reader := csv.NewReader(file)
    records, err := reader.ReadAll()
    if err != nil {
        t.Fatal(err)
    }
    for _, record := range records {
        if len(record) == 6 {
            tokens = append(tokens, Token{
                ProjectName: record[0],
                Token:       record[1],
                Permission:  record[2],
                UserID:      record[3],
                UserName:    record[4],
                ExpiryDate:  record[5],
            })
        }
    }
	
    t.Logf("削除後のトークンリスト: %v", tokens)

    expected := Token{"Project B", "token456", "write", "102", "Bob", "2024/6/30"}
    found := false
    for _, t := range tokens {
        if t == expected {
            found = true
            break
        }
    }
    if found {
        t.Errorf("トークンが削除されていません: %v (期待値: 削除)", expected)
    }
}
