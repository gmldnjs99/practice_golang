import React, { useState, useEffect } from "react";
import StockChart from './components/StockChart'; // StockChart 컴포넌트 import
import { Button, TextField, Container, Typography, CircularProgress, Paper, Grid, Accordion, AccordionSummary, AccordionDetails } from '@mui/material'; // Material UI 컴포넌트 import
import ExpandMoreIcon from '@mui/icons-material/ExpandMore'; // Accordion 아이콘
import './App.css'; // CSS 파일 import

const StockData = () => {
  const [savedStocks, setSavedStocks] = useState([]); // 주식 데이터 상태
  const [isLoading, setIsLoading] = useState(true); // 로딩 상태
  const [error, setError] = useState(null); // 에러 상태
  const [symbol, setSymbol] = useState(""); // 사용자로부터 입력받은 주식 코드 상태
  const [postError, setPostError] = useState(null); // POST 요청 에러 상태
  const [newsData, setNewsData] = useState({}); // 주식별 뉴스 데이터 상태 (기본적으로 빈 객체로 초기화)
  const [newsLoading, setNewsLoading] = useState({}); // 뉴스 로딩 상태 (각 주식별로 관리)
  const [newsVisible, setNewsVisible] = useState({}); // 각 주식별 뉴스 가시성 상태 관리
  const [newsExpanded, setNewsExpanded] = useState(false); // 전체 뉴스 섹션 열림 여부 상태

  // 페이지 로드 시 백엔드 API에서 주식 데이터를 가져옴
  useEffect(() => {
    fetch("http://localhost:8080/saved-stocks-info")
      .then((response) => response.json()) // JSON 형식으로 응답 처리
      .then((data) => {
        setSavedStocks(data.data); // 데이터가 있으면 상태에 저장
        setIsLoading(false); // 로딩 상태 종료
      })
      .catch((error) => {
        console.error("Failed to fetch stock data:", error);
        setError("Failed to fetch stock data"); // 에러 발생 시 상태 업데이트
        setIsLoading(false); // 로딩 상태 종료
      });
  }, []); // 컴포넌트가 마운트될 때 한번만 실행

  // POST 요청으로 주식 데이터 저장
  const addStock = () => {
    if (!symbol) {
      setPostError("Please enter a stock symbol.");
      return;
    }

    fetch("http://localhost:8080/stock", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ symbol }), // body에 symbol 포함
    })
      .then((response) => response.json())
      .then((data) => {
        if (data.message) {
          // 성공적인 저장 후, 저장된 데이터를 다시 가져옴
          setSavedStocks((prevStocks) => [...prevStocks, data.data]);
          setPostError(null);
        }
      })
      .catch((error) => {
        console.error("Failed to post stock data:", error);
        setPostError("Failed to post stock data.");
      });
  };

  // 뉴스 데이터 가져오기
  const fetchNews = (stockSymbol) => {
    setNewsLoading((prev) => ({ ...prev, [stockSymbol]: true })); // 해당 주식의 로딩 상태를 true로 설정
    fetch(`http://localhost:8080/stock-news`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ symbol: stockSymbol }),
    })
      .then((response) => response.json())
      .then((data) => {
        setNewsData((prev) => ({ ...prev, [stockSymbol]: data.data })); // 주식별 뉴스 데이터를 업데이트
        setNewsLoading((prev) => ({ ...prev, [stockSymbol]: false })); // 해당 주식의 로딩 상태를 false로 설정
      })
      .catch((error) => {
        console.error("Failed to fetch news:", error);
        setNewsLoading((prev) => ({ ...prev, [stockSymbol]: false })); // 로딩 상태 종료
      });
  };

  // 뉴스 섹션 열고 닫기
  const handleNewsToggle = (stockSymbol) => {
    setNewsVisible((prev) => ({ ...prev, [stockSymbol]: !prev[stockSymbol] }));
  };

  // 전체 뉴스 섹션 열고 닫기
  const handleNewsExpandedToggle = () => {
    setNewsExpanded(!newsExpanded);
  };

  if (isLoading) {
    return <div className="loading"><CircularProgress /></div>;
  }

  if (error) {
    return <Typography variant="h6" color="error" className="error-message">{error}</Typography>;
  }

  return (
    <Container className="container">
      <Typography variant="h4" align="center" gutterBottom>
        주식 데이터
      </Typography>

      {/* 주식 차트 */}
      <div style={{ marginBottom: '20px' }}>
        <StockChart stockData={savedStocks} />
      </div>

      {/* 주식 코드 입력 폼 */}
      <Paper className="paper" style={{ maxWidth: '500px', margin: '0 auto' }}> {/* 너비 조정 */}
        <Typography variant="h6" gutterBottom>
          새로운 주식 데이터 추가
        </Typography>
        <Grid container spacing={2}>
          <Grid item xs={9}>
            <TextField
              fullWidth
              label="주식 코드 (예: 005930)"
              value={symbol}
              onChange={(e) => setSymbol(e.target.value)}
              error={!!postError}
              helperText={postError}
              className="text-field"
            />
          </Grid>
          <Grid item xs={3}>
            <Button
              variant="contained"
              color="primary"
              fullWidth
              onClick={addStock}
              disabled={!symbol}
            >
              주식 추가
            </Button>
          </Grid>
        </Grid>
      </Paper>

      {/* 주식 리스트 */}
      <div style={{ marginTop: '20px' }}>
        {savedStocks.length > 0 ? (
          <div>
            <Typography variant="h6" style={{ marginTop: '20px' }}>저장된 주식 리스트</Typography>

            {/* Accordion을 사용한 전체 주식 리스트 */}
            <Accordion>
              <AccordionSummary
                expandIcon={<ExpandMoreIcon />}
                aria-controls="panel-content"
                id="panel-header"
              >
                <Typography>저장된 주식 전체 보기</Typography>
              </AccordionSummary>
              <AccordionDetails>
                <div>
                  {savedStocks.map((stock, index) => (
                    <Accordion key={index}>
                      <AccordionSummary
                        expandIcon={<ExpandMoreIcon />}
                        aria-controls={`panel${index}-content`}
                        id={`panel${index}-header`}
                      >
                        <Typography>
                          <strong>{stock.name}</strong> ({stock.symbol}) - {stock.price} 원
                        </Typography>
                      </AccordionSummary>
                      <AccordionDetails>
                        <Typography>
                          {/* 주식 상세 정보 출력 */}
                          <strong>기업명:</strong> {stock.name} <br />
                          <strong>주식 가격:</strong> {stock.price} 원 <br />
                          <strong>상장일:</strong> {stock.listingDate} <br />
                          {/* 다른 세부 정보를 추가할 수 있습니다. */}
                        </Typography>
                      </AccordionDetails>
                    </Accordion>
                  ))}
                </div>
              </AccordionDetails>
            </Accordion>
          </div>
        ) : (
          <Typography variant="body1">저장된 주식 데이터가 없습니다.</Typography>
        )}
      </div>

      {/* 뉴스 섹션 */}
      <div style={{ marginTop: '40px' }}>
        <Typography variant="h6" gutterBottom>
          주식 관련 뉴스
        </Typography>

        {/* 뉴스 전체 보기 Accordion */}
        <Accordion expanded={newsExpanded} onChange={handleNewsExpandedToggle}>
          <AccordionSummary
            expandIcon={<ExpandMoreIcon />}
            aria-controls="news-content"
            id="news-header"
          >
            <Typography>전체 뉴스 보기</Typography>
          </AccordionSummary>
          <AccordionDetails>
            <Grid container spacing={2}>
              {savedStocks.map((stock) => (
                <Grid item xs={12} sm={6} md={4} key={stock.symbol}>
                  <Paper className="paper" style={{ padding: '20px' }}>
                    <Typography variant="h6">{stock.name} ({stock.symbol}) 뉴스</Typography>
                    <Button
                      variant="contained"
                      color="primary"
                      onClick={() => {
                        if (!newsVisible[stock.symbol]) {
                          fetchNews(stock.symbol);
                        }
                        handleNewsToggle(stock.symbol); // 뉴스 보기/숨기기 토글
                      }}
                      disabled={newsLoading[stock.symbol]}
                    >
                      {newsLoading[stock.symbol] ? '뉴스 로딩 중...' : newsVisible[stock.symbol] ? '뉴스 숨기기' : '뉴스 보기'}
                    </Button>

                    {/* 뉴스 목록 표시 */}
                    {newsVisible[stock.symbol] && newsData[stock.symbol] && newsData[stock.symbol].length > 0 && (
                      <div style={{ marginTop: '20px' }}>
                        {newsData[stock.symbol].map((news, index) => (
                          <Accordion key={index}>
                            <AccordionSummary
                              expandIcon={<ExpandMoreIcon />}
                              aria-controls={`news-panel${index}-content`}
                              id={`news-panel${index}-header`}
                            >
                              <Typography>{news.title}</Typography>
                            </AccordionSummary>
                            <AccordionDetails>
                              <Typography>
                                {news.summary} <br />
                                <a href={news.url} target="_blank" rel="noopener noreferrer">Read more</a>
                              </Typography>
                            </AccordionDetails>
                          </Accordion>
                        ))}
                      </div>
                    )}
                  </Paper>
                </Grid>
              ))}
            </Grid>
          </AccordionDetails>
        </Accordion>
      </div>
    </Container>
  );
};

export default StockData;
