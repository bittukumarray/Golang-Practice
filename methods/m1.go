package main

import (
	"fmt"
)

type myInt int

type Vr struct{
	x,y int
}

func (I *myInt) DoubleNo() int{
	return int(2*(*I))
}

func (V *Vr) DoubleNo() int{
	V.x=0
	V.y=0
	return V.x*2+V.y*2
}


// func DoubleNoF(I myInt) int {
// 	return int(2*I)
// }

// func DoubleNoF(V Vr) int{

// 	return V.x*2+V.y*2
// }

func main(){

	i := myInt(4)

	fmt.Println(i.DoubleNo())
	//fmt.Println(DoubleNoF(i))

	v :=Vr{10,20}
	v2:=&v
	v2.DoubleNo()
	// v.DoubleNo()
	// fmt.Println(DoubleNoF(v))
	fmt.Println(v)

}