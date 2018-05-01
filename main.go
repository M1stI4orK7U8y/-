package main

import "fmt"

//max and min
const (
	MaxUint = ^uint(0)
	MinUint = 0
	MaxInt  = int(MaxUint >> 1)
	MinInt  = -MaxInt - 1
)

func main() {
	fmt.Println("區塊鏈小作品展示")
	bc := NewBlockchain()
	defer bc.db.Close()

	cmd := CMD{bc}
	cmd.Run()
}
