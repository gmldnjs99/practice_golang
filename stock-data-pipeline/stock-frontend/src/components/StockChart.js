import React from 'react';
import { Line } from 'react-chartjs-2';
import { Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend, TimeScale } from 'chart.js';
import 'chartjs-adapter-date-fns'; // date adapter 추가

// 차트에 필요한 모듈을 등록
ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend, TimeScale);

const StockChart = ({ stockData }) => {
  // 랜덤 색상 생성 함수
  const getRandomColor = () => {
    const letters = '0123456789ABCDEF';
    let color = '#';
    for (let i = 0; i < 6; i++) {
      color += letters[Math.floor(Math.random() * 16)];
    }
    return color;
  };

  // 데이터가 올바르게 전달되고 있는지 콘솔로 확인
  console.log("Received stock data:", stockData);

  // 각 주식의 symbol별로 데이터셋을 구성
  const groupedData = stockData.reduce((acc, data) => {
    if (!acc[data.symbol]) {
      acc[data.symbol] = {
        label: data.name,  // 주식 이름을 label로 설정
        data: [],
        borderColor: getRandomColor(),
        backgroundColor: 'rgba(75,192,192,0.2)',
        fill: true,
      };
    }

    // 'x'는 Date 객체, 'y'는 price로 설정
    acc[data.symbol].data.push({ x: new Date(data.created_at), y: data.price }); // created_at을 x로 사용
    return acc;
  }, {});

  // datasets 배열로 변환
  const datasets = Object.values(groupedData);

  // 차트에 사용할 데이터 구성
  const chartData = {
    datasets: datasets, // datasets 배열에 주식 데이터가 들어있음
  };

  const options = {
    responsive: true,
    plugins: {
      title: {
        display: true,
        text: 'Stock Price Chart',
      },
      tooltip: {
        mode: 'index',
        intersect: false,
      },
    },
    scales: {
      x: {
        type: 'time',  // 시간 축을 사용
        time: {
          unit: 'day',  // 일 단위로 표시
          tooltipFormat: 'MM/dd/yyyy',  // 툴팁 형식 수정
        },
        title: {
          display: true,
          text: 'Date',
        },
      },
      y: {
        title: {
          display: true,
          text: 'Price (원)',
        },
      },
    },
  };

  return (
    <div>
      <h3>Stock Price Chart</h3>
      <Line data={chartData} options={options} />
    </div>
  );
};

export default StockChart;
