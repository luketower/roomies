package main

import (
	"fmt"
	bp "github.com/luketower/roomies/billprocessor"
	br "github.com/luketower/roomies/billsreader"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) > 1 {
		if bp.Valid(args) {
			fmt.Println(bp.BillReport(args))
		} else {
			fmt.Println(bp.ErrorMsg(args))
		}
	} else {
		br.Run(args)
	}
}
