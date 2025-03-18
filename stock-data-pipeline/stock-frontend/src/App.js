import React, { useState, useEffect } from "react";
import StockChart from './components/StockChart'; // StockChart 컴포넌트 import

const StockData = () => {
  const [savedStocks, setSavedStocks] = useState([]); // 주식 데이터 상태
  const [isLoading, setIsLoading] = useState(true); // 로딩 상태
  const [error, setError] = useState(null); // 에러 상태
  const [symbol, setSymbol] = useState(""); // 사용자로부터 입력받은 주식 코드 상태
  const [postError, setPostError] = useState(null); // POST 요청 에러 상태

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

  if (isLoading) {
    return <p>주식 데이터를 로드 중...</p>;
  }

  if (error) {
    return <p>{error}</p>;
  }

  return (
    <div>
      <h1>주식 데이터</h1>

      {/* 주식 코드 입력 폼 */}
      <div>
        <input
          type="text"
          value={symbol}
          onChange={(e) => setSymbol(e.target.value)}
          placeholder="주식 코드를 입력하세요"
        />
        <button onClick={addStock}>주식 데이터 추가</button>
        {postError && <p style={{ color: "red" }}>{postError}</p>}
      </div>

      {savedStocks.length > 0 ? (
        <div>
          {/* 주식 차트 시각화 추가 */}
          <StockChart stockData={savedStocks} />
          
          {/* 저장된 주식 리스트 */}
          <ul>
            {savedStocks.map((stock, index) => (
              <li key={index}>
                <strong>{stock.name}</strong> ({stock.symbol}) - {stock.price} 원
              </li>
            ))}
          </ul>
        </div>
      ) : (
        <p>저장된 주식 데이터가 없습니다.</p>
      )}
    </div>
  );
};

export default StockData;
