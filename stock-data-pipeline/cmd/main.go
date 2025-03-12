package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gmldnjs99/stock-data-pipeline/internal/collector"
	"github.com/gmldnjs99/stock-data-pipeline/internal/storage"
)

func main() {
	// DB 연결
	storage.InitDB()
	defer storage.CloseDB()

	// Gin 엔진 생성
	router := gin.Default()

	// 루트 페이지 엔드포인트 추가
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to Stock Data Pipeline!"})
	})

	// 헬스 체크 엔드포인트
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Favicon 요청 무시
	router.GET("/favicon.ico", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	// 삼성전자 주식 가격 조회 및 MySQL 저장
	router.GET("/samsung-stocks", func(c *gin.Context) {
		// 네이버에서 삼성전자 주식 데이터를 가져오고 MySQL에 저장
		stockData, err := collector.GetStockPrice("005930") // 삼성전자
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stock data"})
			return
		}
		// 주식 데이터를 DB에 저장
		storage.SaveStockData(stockData.Symbol, stockData.Name, stockData.Price)

		stockData, err = collector.GetStockPrice("005935") // 삼성전자 우선주
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stock data"})
			return
		}
		// 주식 데이터를 DB에 저장
		storage.SaveStockData(stockData.Symbol, stockData.Name, stockData.Price)

		c.JSON(http.StatusOK, gin.H{"message": "Samsung Stock data saved to DB"})
	})

	// Tiger 미국 S&P 주식 가격 조회 및 MySQL 저장
	router.GET("/snp-stocks", func(c *gin.Context) {
		// 네이버에서 Tiger 미국 S&P 주식 데이터를 가져오고 MySQL에 저장
		stockData, err := collector.GetStockPrice("360750") // Tiger 미국 S&P
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stock data"})
			return
		}
		// 주식 데이터를 DB에 저장
		storage.SaveStockData(stockData.Symbol, stockData.Name, stockData.Price)

		c.JSON(http.StatusOK, gin.H{"message": "Tiger 미국 S&P Stock data saved to DB"})
	})

	// 특정 주식 데이터를 가져오고 DB에 저장하는 엔드포인트 추가
	router.GET("/stock/:symbol", func(c *gin.Context) {
		symbol := c.Param("symbol") // URL에서 주식 코드 가져오기
		stockData, err := collector.GetStockPrice(symbol)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stock data"})
			return
		}

		err = storage.SaveStockData(stockData.Symbol, stockData.Name, stockData.Price)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save stock data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Stock data saved successfully",
			"data":    stockData,
		})
	})

	router.GET("/saved-stocks-info", func(c *gin.Context) {
		stocks, err := storage.GetRecentStockData()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stock data"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": stocks})
	})

	// 서버 실행
	log.Println("Server is running on port 8080")
	router.Run(":8080")
}
