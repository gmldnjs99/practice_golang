# Go 언어를 활용한 데이터 수집 및 파이프라인 개발

## Introduction

**실시간 및 과거 주식 데이터를 수집하고, 이를 저장 및 분석하여 투자 참고 정보를 제공하는 데이터 파이프라인 구축 관련 repository입니다.** 

## 기본 정보

- **Base URL**: `http://localhost:8080`
- **Content-Type**: `application/json`

---

### API

| API 종류 | HTTP Method |        기능        |       url      |                                          설명                                          |
|:--------:|:--------:|:------------------:|:-------------------:|:--------------------------------------------------------------------------------------:|
| 루트 페이지 |   GET   |  환영 메시지 반환 |        /       | 기본 엔드포인트, 서비스 정상 동작 여부 확인용. Welcome to Stock Data Pipeline! 메시지 제공. |
| 헬스 체크 |  GET  |    헬스 체크   |      /health     | 서버가 정상적으로 동작하는지 상태 확인. { "status": "ok" } 반환.                           |
| 파비콘 무시         |  GET  |      파비콘 요청 무시     |     /favicon.ico    | 불필요한 파비콘 요청을 무시하고, 204 No Content 반환.       |
| 삼성전자 주식 저장         |   POST   |   삼성전자 및 삼성전자 우선주 주식 저장   | /samsung-stocks	 | 삼성전자(005930), 삼성전자우(005935) 주식 가격을 네이버 금융에서 크롤링 후 DB 저장.                |
| S&P 500 주식 저장         |   POST    |      Tiger 미국 S&P 500 주식 저장     |        /snp-stocks        | Tiger 미국 S&P 500 (360750) 주식 가격을 크롤링 후 DB에 저장.                                                                       |
|  개별 주식 데이터 저장 |   POST   |      특정 종목 주식 데이터 저장      |         /stock/:symbol        | 요청한 종목코드(symbol)의 주식 가격을 크롤링 후 DB 저장.                           |
| 저장된 주식 데이터 조회         |   GET   |     저장된 주식 데이터 전체 조회    |      /saved-stocks-info     | DB에 저장된 최근 주식 데이터 전체 조회.                                                             |

## API 목록

### 1. Root Endpoint

- **Method**: `GET`
- **URL**: `/`
- **Description**: API의 기본 페이지입니다. 서비스의 간단한 소개를 반환합니다.

#### 예시
```bash
curl -X GET http://localhost:8080/
```

#### Response
```json
{
  "message": "Welcome to Stock Data Pipeline!"
}
```

---

### 2. Health Check

- **Method**: `GET`
- **URL**: `/health`
- **Description**: 서버의 상태를 확인하는 엔드포인트입니다. 서버가 정상 동작 중인지 확인할 수 있습니다.

#### 예시
```bash
curl -X GET http://localhost:8080/health
```

#### Response
```json
{
  "status": "ok"
}
```

---

### 3. 삼성전자 주식 가격 저장

- **Method**: `GET`
- **URL**: `/samsung-stocks`
- **Description**: 삼성전자(005930)와 삼성전자 우선주(005935)의 주식 데이터를 네이버에서 가져와 MySQL에 저장합니다.

#### 예시
```bash
curl -X GET http://localhost:8080/samsung-stocks
```

#### Response
```json
{
  "message": "Samsung Stock data saved to DB"
}
```

---

### 4. Tiger 미국 S&P 주식 가격 저장

- **Method**: `GET`
- **URL**: `/snp-stocks`
- **Description**: Tiger 미국 S&P 주식 데이터(360750)를 네이버에서 가져와 MySQL에 저장합니다.

#### 예시
```bash
curl -X GET http://localhost:8080/snp-stocks
```

#### Response
```json
{
  "message": "Tiger 미국 S&P Stock data saved to DB"
}
```

---

### 5. 특정 주식 데이터 저장

- **Method**: `POST`
- **URL**: `/stock`
- **Description**: 특정 주식의 주식 데이터를 네이버에서 가져와 MySQL에 저장합니다. 요청 본문에서 주식의 `symbol`을 전달해야 합니다.

#### Request Body
```json
{
  "symbol": "005930"
}
```

#### 예시
```bash
curl -X POST http://localhost:8080/stock -H "Content-Type: application/json" -d '{"symbol": "005930"}'
```

#### Response
```json
{
  "message": "Stock data saved successfully",
  "data": {
    "symbol": "005930",
    "name": "Samsung Electronics",
    "price": 82500
  }
}
```

---

### 6. 최근 저장된 주식 정보 조회

- **Method**: `GET`
- **URL**: `/saved-stocks-info`
- **Description**: MySQL에 저장된 최근 주식 데이터 목록을 반환합니다.

#### 예시
```bash
curl -X GET http://localhost:8080/saved-stocks-info
```

#### Response
```json
{
  "data": [
    {
      "symbol": "005930",
      "name": "Samsung Electronics",
      "price": 82500,
      "timestamp": "2025-03-18T10:00:00Z"
    },
    {
      "symbol": "360750",
      "name": "Tiger US S&P",
      "price": 34500,
      "timestamp": "2025-03-18T10:05:00Z"
    }
  ]
}
```

---

## 에러 처리

모든 API는 에러가 발생할 경우 적절한 HTTP 상태 코드와 함께 에러 메시지를 반환합니다.

### 예시

#### 1. 잘못된 주식 심볼을 요청한 경우
```bash
curl -X POST http://localhost:8080/stock -H "Content-Type: application/json" -d '{"symbol": "INVALID"}'
```

#### Response
```json
{
  "error": "Failed to fetch stock data"
}
```

---

## 서버 실행

서버는 `8080` 포트에서 실행됩니다. 서버를 시작하려면 다음 명령어를 실행하십시오:

```bash
go run main.go
```

서버가 정상적으로 실행되면, `http://localhost:8080`에서 API를 사용할 수 있습니다.