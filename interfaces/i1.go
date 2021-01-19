package main

import (
	"fmt"
)

type Abser interface{
	Abs() int
	Dub() int
}

type Vr struct{
	x,y int
}

type myInt int

func (V Vr) Abs() int{
	return V.x+V.y
}

func (V Vr) Dub() int{
	return V.x -V.y
}

func (i myInt) Abs() int{
	return int(2*i)
}

func (i myInt) Dub() int{

	return int(i)
}



func main(){
	i:=myInt(3)
	v:=Vr{10,20}

	var a Abser
	
	a=i
	fmt.Println(a.Abs())

	a=v
	fmt.Println(a)

	a=&v

	fmt.Println(a)

	a=i

	fmt.Println(a.Dub())
}

