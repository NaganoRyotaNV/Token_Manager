package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Token struct {
	ProjectName string `json:"projectName"`
	Token       string `json:"token"`
	Permission  string `json:"permission"`
	UserID      string `json:"userId"`
	UserName    string `json:"userName"`
	ExpiryDate  string `json:"expiryDate"`
}

var filePath = "uploaded.csv"

func main() {
	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)

	http.HandleFunc("/api/tokens", cors(handleTokens))
	http.HandleFunc("/api/upload", cors(handleUpload))
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("サーバーはポート%sで開始されました\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
func cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "無効なリクエストメソッドです", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Printf("ファイルの取得中にエラーが発生しました: %v", err)
		http.Error(w, "ファイルの取得中にエラーが発生しました", http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Printf("アップロードされたファイル: %+v\n", header.Filename)
	log.Printf("ファイルサイズ: %+v\n", header.Size)
	log.Printf("MIMEヘッダー: %+v\n", header.Header)

	tempFile, err := os.Create(filePath)
	if err != nil {
		log.Printf("ファイルの作成中にエラーが発生しました: %v", err)
		http.Error(w, "ファイルの作成中にエラーが発生しました", http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("ファイルの読み取り中にエラーが発生しました: %v", err)
		http.Error(w, "ファイルの読み取り中にエラーが発生しました", http.StatusInternalServerError)
		return
	}

	tempFile.Write(fileBytes)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "ファイルが正常にアップロードされました"})
}

func handleTokens(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		projectName := r.URL.Query().Get("projectName")
		if projectName != "" {
			showProjectTokens(filePath, projectName, w)
		} else {
			showData(filePath, w)
		}
	case http.MethodPost:
		addData(filePath, r, w)
	case http.MethodPut:
		updateProjectTokens(filePath, r, w)
	case http.MethodDelete:
		deleteData(filePath, r, w)
	default:
		http.Error(w, "サポートされていないメソッドです", http.StatusMethodNotAllowed)
	}
}

func showData(filePath string, w http.ResponseWriter) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("ファイルが存在しません: %v", err)
		http.Error(w, "ファイルが見つかりません", http.StatusNotFound)
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("ファイルのオープン中にエラーが発生しました: %v", err)
		http.Error(w, "ファイルのオープン中にエラーが発生しました", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("CSVデータの読み取り中にエラーが発生しました: %v", err)
		http.Error(w, "CSVデータの読み取り中にエラーが発生しました", http.StatusInternalServerError)
		return
	}

	var tokens []Token
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

func showProjectTokens(filePath, projectName string, w http.ResponseWriter) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("ファイルが存在しません: %v", err)
		http.Error(w, "ファイルが見つかりません", http.StatusNotFound)
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("ファイルのオープン中にエラーが発生しました: %v", err)
		http.Error(w, "ファイルのオープン中にエラーが発生しました", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("CSVデータの読み取り中にエラーが発生しました: %v", err)
		http.Error(w, "CSVデータの読み取り中にエラーが発生しました", http.StatusInternalServerError)
		return
	}

	var tokens []Token
	for _, record := range records {
		if len(record) == 6 && record[0] == projectName {
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

func addData(filePath string, r *http.Request, w http.ResponseWriter) {
	var token Token
	if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
		log.Printf("リクエストボディの解析中にエラーが発生しました: %v", err)
		http.Error(w, "リクエストボディの解析中にエラーが発生しました", http.StatusBadRequest)
		return
	}

	records, err := readCSV(filePath)
	if err != nil {
		log.Printf("CSVデータの読み取り中にエラーが発生しました: %v", err)
		http.Error(w, "CSVデータの読み取り中にエラーが発生しました", http.StatusInternalServerError)
		return
	}

	records = append(records, []string{token.ProjectName, token.Token, token.Permission, token.UserID, token.UserName, token.ExpiryDate})
	err = writeCSV(filePath, records)
	if err != nil {
		http.Error(w, "データの追加中にエラーが発生しました", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func updateProjectTokens(filePath string, r *http.Request, w http.ResponseWriter) {
	var data struct {
		ProjectName string  `json:"projectName"`
		Tokens      []Token `json:"tokens"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("リクエストボディの解析中にエラーが発生しました: %v", err)
		http.Error(w, "リクエストボディの解析中にエラーが発生しました", http.StatusBadRequest)
		return
	}

	records, err := readCSV(filePath)
	if err != nil {
		log.Printf("CSVデータの読み取り中にエラーが発生しました: %v", err)
		http.Error(w, "CSVデータの読み取り中にエラーが発生しました", http.StatusInternalServerError)
		return
	}

	newRecords := [][]string{}
	for _, record := range records {
		if record[0] != data.ProjectName {
			newRecords = append(newRecords, record)
		}
	}

	for _, token := range data.Tokens {
		newRecords = append(newRecords, []string{token.ProjectName, token.Token, token.Permission, token.UserID, token.UserName, token.ExpiryDate})
	}

	err = writeCSV(filePath, newRecords)
	if err != nil {
		http.Error(w, "データの更新中にエラーが発生しました", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteData(filePath string, r *http.Request, w http.ResponseWriter) {
	lineNumStr := r.URL.Query().Get("line")
	lineNum, err := strconv.Atoi(lineNumStr)
	if err != nil {
		log.Printf("無効な行番号です: %v", err)
		http.Error(w, "無効な行番号です", http.StatusBadRequest)
		return
	}

	records, err := readCSV(filePath)
	if err != nil {
		log.Printf("CSVデータの読み取り中にエラーが発生しました: %v", err)
		http.Error(w, "CSVデータの読み取り中にエラーが発生しました", http.StatusInternalServerError)
		return
	}

	if lineNum < 0 || lineNum >= len(records) {
		log.Printf("行番号が範囲外です: %d", lineNum)
		http.Error(w, "行番号が範囲外です", http.StatusBadRequest)
		return
	}

	log.Printf("削除前のレコード: %v", records) 
	records = append(records[:lineNum], records[lineNum+1:]...)
	log.Printf("削除後のレコード: %v", records) 

	err = writeCSV(filePath, records)
	if err != nil {
		log.Printf("CSVデータの書き込み中にエラーが発生しました: %v", err)
		http.Error(w, "CSVデータの書き込み中にエラーが発生しました", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func readCSV(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("ファイルのオープン中にエラーが発生しました: %v", err)
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("CSVデータの読み取り中にエラーが発生しました: %v", err)
		return nil, err
	}

	return records, nil
}

func writeCSV(filePath string, records [][]string) error {
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("ファイルの作成中にエラーが発生しました: %v", err)
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	err = writer.WriteAll(records)
	if err != nil {
		log.Printf("CSVファイルへの書き込み中にエラーが発生しました: %v", err)
		return err
	}

	return nil
}
