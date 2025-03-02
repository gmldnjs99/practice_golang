package main

import "fmt"

func main() {
	var a int = 1 // 기본 형태
	var b int     // 초기값 생략, 초기값은 타입별 기본값으로 대체
	var c = 3     // 타입 생략, 변수 타입은 우변 값의 타입이 됨
	d := 4        // 선언 대입문 :=을 사용해서 var 키워드와 타입 생략, 선언과 대입을 한꺼번에 하는 구문

	/*
		var a := 3.1234     - a는 float64 타입으로 자동 지정됩니다.
		b := 365            - b는 int 타입으로 자동 지정됩니다.
		c := "Hello golang" - c는 string 타입으로 자동 지정됩니다.
	*/

	fmt.Println(a, b, c, d)
}
