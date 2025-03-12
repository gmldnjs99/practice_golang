package collector

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/gmldnjs99/stock-data-pipeline/internal/storage"
	"github.com/gocolly/colly/v2"
)

// StockData 구조체 정의 (종목명 추가)
type StockData struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
	Price  int    `json:"price"` // 소수점 없이 정수로 저장
}

// GetStockPrice 함수 - 네이버 금융에서 특정 주식의 가격을 가져옴
func GetStockPrice(symbol string) (*StockData, error) {
	// 네이버 금융 URL (주식 정보 URL)
	url := fmt.Sprintf("https://finance.naver.com/item/main.nhn?code=%s", symbol)

	// 크롤러 초기화
	c := colly.NewCollector()

	// 주식 가격을 저장할 변수
	var price string
	var name string

	// 주식 가격을 가져오는 부분
	c.OnHTML(".no_today .blind", func(e *colly.HTMLElement) {
		// 주식 가격 가져오기 (쉼표 제거)
		price = e.Text
	})

	// 종목명 가져오는 부분
	c.OnHTML(".wrap_company h2 a", func(e *colly.HTMLElement) {
		name = e.Text
	})

	// 에러 처리
	c.OnError(func(r *colly.Response, err error) {
		log.Println("크롤링 중 오류 발생:", err)
	})

	// 네이버 금융 페이지 크롤링 시작
	err := c.Visit(url)
	if err != nil {
		return nil, fmt.Errorf("주식 가격을 가져오는 중 오류 발생: %v", err)
	}

	// 주식 가격이 없으면 에러 반환
	if price == "" || name == "" {
		return nil, fmt.Errorf("주식 데이터를 찾을 수 없음")
	}

	// 가격에서 쉼표 제거
	price = strings.Replace(price, ",", "", -1)

	// 문자열에서 가격을 정수로 변환
	var priceFloat float64
	_, err = fmt.Sscanf(price, "%f", &priceFloat)
	if err != nil {
		return nil, fmt.Errorf("가격 변환 오류: %v", err)
	}

	// 가격 소수점 없애기
	priceFloat = math.Floor(priceFloat)

	// 결과 반환
	return &StockData{
		Symbol: symbol,
		Name:   name,
		Price:  int(priceFloat), // 소수점 없는 정수로 반환
	}, nil
}

// GetSamsungStockPrices 함수 - 삼성전자와 삼성전자 우선주의 가격을 가져옴
func GetSamsungStockPrices() (map[string]*StockData, error) {
	// 삼성전자와 삼성전자 우선주 주식 코드
	stocks := []string{"005930", "005935"}

	// 주식 데이터 맵
	stockDataMap := make(map[string]*StockData)

	// 각 주식 코드에 대해 가격을 가져옴
	for _, symbol := range stocks {
		stockData, err := GetStockPrice(symbol)
		if err != nil {
			return nil, err
		}
		stockDataMap[symbol] = stockData

		// MySQL에 주식 데이터 저장
		err = storage.SaveStockData(stockData.Symbol, stockData.Name, stockData.Price)
		if err != nil {
			log.Printf("DB 저장 오류: %v", err)
		}
	}

	// 결과 반환
	return stockDataMap, nil
}

// GetTigerSPStockPrices 함수 - Tiger 미국 S&P의 가격을 가져옴
func GetTigersnpStockPrices() (map[string]*StockData, error) {
	// Tiger 미국 S&P 주식 코드
	stocks := []string{"360750"}

	// 주식 데이터 맵
	stockDataMap := make(map[string]*StockData)

	// 각 주식 코드에 대해 가격을 가져옴
	for _, symbol := range stocks {
		stockData, err := GetStockPrice(symbol)
		if err != nil {
			return nil, err
		}
		stockDataMap[symbol] = stockData

		// MySQL에 주식 데이터 저장
		err = storage.SaveStockData(stockData.Symbol, stockData.Name, stockData.Price)
		if err != nil {
			log.Printf("DB 저장 오류: %v", err)
		}
	}

	// 결과 반환
	return stockDataMap, nil
}
