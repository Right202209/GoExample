package main

import "fmt"

func sum(nums ...int) {
	fmt.Println(nums, " ")
	tatal := 0
	for _, num := range nums {
		tatal += num
	}
	fmt.Println(tatal)
}
func main() {
	sum(1, 2)
	sum(1, 2, 3)

	nums := []int{1, 2, 3, 4}
	sum(nums...)
}
