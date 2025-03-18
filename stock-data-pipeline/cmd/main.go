package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gmldnjs99/stock-data-pipeline/internal/collector"
	"github.com/gmldnjs99/stock-data-pipeline/internal/storage"
)

func main() {
	// ========================
	// 데이터베이스 초기화 및 종료
	// ========================
	storage.InitDB()
	defer storage.CloseDB()

	// ========================
	// Gin 엔진 생성 및 미들웨어 설정
	// ========================
	router := gin.Default()

	// CORS 설정 (모든 origin 허용)
	router.Use(cors.Default())

	// ========================
	// 기본 라우터 핸들러
	// ========================

	// 루트 페이지 - 기본 안내 메시지
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to Stock Data Pipeline!"})
	})

	// 헬스 체크 엔드포인트
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// favicon 요청 무시
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	// ========================
	// 주식 데이터 수집 및 저장 엔드포인트 (POST)
	// ========================

	// 삼성전자 및 삼성전자 우선주 데이터를 수집하고 DB에 저장하는 POST 엔드포인트
	router.POST("/samsung-stocks", func(c *gin.Context) {
		// 삼성전자 주식 데이터 수집
		stockData, err := collector.GetStockPrice("005930") // 삼성전자
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Samsung Electronics stock data"})
			return
		}
		// DB 저장
		storage.SaveStockData(stockData.Symbol, stockData.Name, stockData.Price)

		// 삼성전자 우선주 데이터 수집
		stockData, err = collector.GetStockPrice("005935") // 삼성전자 우선주
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Samsung Electronics Preferred stock data"})
			return
		}
		// DB 저장
		storage.SaveStockData(stockData.Symbol, stockData.Name, stockData.Price)

		// 성공 응답
		c.JSON(http.StatusOK, gin.H{"message": "Samsung stock data saved to DB"})
	})

	// TIGER 미국 S&P500 ETF 데이터를 수집하고 DB에 저장하는 POST 엔드포인트
	router.POST("/snp-stocks", func(c *gin.Context) {
		// TIGER 미국 S&P500 ETF 데이터 수집
		stockData, err := collector.GetStockPrice("360750")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get S&P500 ETF stock data"})
			return
		}
		// DB 저장
		storage.SaveStockData(stockData.Symbol, stockData.Name, stockData.Price)

		c.JSON(http.StatusOK, gin.H{"message": "S&P500 ETF stock data saved to DB"})
	})

	// 주식 데이터를 수집하고 DB에 저장하는 POST 엔드포인트 (symbol은 body에서 받음)
	router.POST("/stock", func(c *gin.Context) {
		// JSON 바디에서 symbol을 받음
		var requestBody struct {
			Symbol string `json:"symbol"`  // 요청 본문에서 symbol을 받음
		}

		// 요청 본문 파싱
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// 주식 데이터 수집
		stockData, err := collector.GetStockPrice(requestBody.Symbol)  // 이제 symbol은 requestBody.Symbol로 받음
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stock data"})
			return
		}

		// DB에 주식 데이터 저장
		err = storage.SaveStockData(stockData.Symbol, stockData.Name, stockData.Price)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save stock data"})
			return
		}

		// 성공 응답
		c.JSON(http.StatusOK, gin.H{
			"message": "Stock data saved successfully",
			"data":    stockData,
		})
	})

	// DB에 저장된 최근 주식 데이터를 조회하는 GET 엔드포인트
	router.GET("/saved-stocks-info", func(c *gin.Context) {
		// 최근 주식 데이터 조회
		stocks, err := storage.GetRecentStockData()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch saved stock data"})
			return
		}

		// 성공 응답
		c.JSON(http.StatusOK, gin.H{"data": stocks})
	})

	// ========================
	// 서버 실행
	// ========================
	log.Println("✅ Server is running on port 8080")
	router.Run(":8080")
}
