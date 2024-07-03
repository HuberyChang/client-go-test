package main

func incr(a int) (b int) {
	defer func() {
		a++
		b++
	}()
	a++
	b = a
	return b
}

func incr2(a int) int {
	var b int
	defer func() {
		a++
		b++
	}()
	a++
	b = a
	return b
}

//func main() {
//	x := "hello"
//	fmt.Println(**x)
//}
