package storage

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DB 연결 변수
var db *sql.DB

// StockData 구조체 정의
type StockData struct {
	Symbol    string    `json:"symbol"`
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"created_at"`  // time.Time으로 선언
}

// DB 연결 함수
func InitDB() {
	var err error
	// MySQL DB 연결 (사용자명:비밀번호@/DB명)
	dsn := "stock_user:gmldnjs5695@tcp(localhost:3306)/stock_data"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("DB 연결 실패:", err)
	}
	// DB 연결 확인
	if err := db.Ping(); err != nil {
		log.Fatal("DB 연결 실패:", err)
	} else {
		log.Println("DB 연결 성공")
	}
}

// DB 종료 함수
func CloseDB() {
	if err := db.Close(); err != nil {
		log.Fatal("DB 종료 실패:", err)
	}
	log.Println("DB 연결 종료")
}

// 주식 데이터 저장 함수
func SaveStockData(symbol, name string, price int) error {
	query := "INSERT INTO stocks (symbol, name, price) VALUES (?, ?, ?)"
	_, err := db.Exec(query, symbol, name, price)
	if err != nil {
		return fmt.Errorf("데이터 저장 실패: %v", err)
	}
	log.Println("주식 데이터 저장 성공:", symbol, name, price)
	return nil
}

// 최근 주식 데이터 가져오기
func GetRecentStockData() ([]StockData, error) {
	var stocks []StockData
	query := "SELECT symbol, name, price, created_at FROM stocks ORDER BY id DESC LIMIT 100"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var stock StockData
		var createdAtStr string  // created_at을 string으로 먼저 읽음
		err := rows.Scan(&stock.Symbol, &stock.Name, &stock.Price, &createdAtStr)
		if err != nil {
			return nil, err
		}

		// 문자열을 time.Time으로 변환
		stock.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("created_at 파싱 실패: %v", err)
		}

		stocks = append(stocks, stock)
	}
	return stocks, nil
}