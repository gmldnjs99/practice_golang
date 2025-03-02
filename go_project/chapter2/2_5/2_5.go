package main

import "fmt"

func main() {
	a := 1              // int
	var b float64 = 3.5 // float64

	var c int = int(b)  // float64에서 int로 변환
	d := float64(a * c) // int에서 float64로 변환

	var e int64 = 7
	f := int64(d) * e // float64에서 int64로 전환

	var g int = int(b * 3) // float64에서 int로 변환
	var h int = int(b) * 3 // float64에서 int로 변환, g와 값이 다름
	fmt.Println(g, h, f)
}
