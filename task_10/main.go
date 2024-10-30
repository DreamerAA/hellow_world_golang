package main

import (
	"flag"
	"fmt"
)

func main() {
	length := flag.Int("length", 8, "length of password")
	useNumbers := flag.Bool("numbers", true, "use numbers")
	useSymbols := flag.Bool("symbols", false, "use symbols")
	useUpper := flag.Bool("uppercase", false, "use uppercase")
	useLower := flag.Bool("lowercase", true, "use lowercase")

	flag.Parse()
	fmt.Println("Start generating")
	password, err := GeneratePassword(*length, *useNumbers, *useSymbols, *useUpper, *useLower)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Created password:", password)
}
