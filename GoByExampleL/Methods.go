package main

import "fmt"

type rect1 struct {
	width, height int
}

func (r *rect1) area() int {
	return r.width * r.height
}
func (r *rect1) perim() int {
	return 2*r.width + 2*r.height
}
func main() {
	r := rect1{width: 10, height: 5}

	fmt.Println("area:", r.area())
	fmt.Println("perim:", r.perim())

	rp := &r
	fmt.Println("area:", rp.area())
	fmt.Println("perim:", rp.perim())

}
