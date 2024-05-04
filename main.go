package main

import (
	"fmt"
)

func intSeq() func() int {
	i := 0
	return func() int {
		i++
		return i
	}
}

func main() {
	//breaking jdsl
	fmt.Println("I'm thinking of using the ::19990")
}
