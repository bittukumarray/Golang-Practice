package testfact


func Fact(n int) int{
	if n<0{
		return -1
	}
	if n==0{
		return 1
	}
	return n*Fact(n-1)
}