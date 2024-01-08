package main

import (
	"fmt"

	rapid_fuzz "github.com/learn-go/go-fuzz/3rdpkg/rapid-fuzz"
)

func main() {
	fmt.Printf("fuzz cgo test start..\n")
	sa := rapid_fuzz.PartialRatioAlignment("\x6d\x4b\x8b\xd5", "hello cgo")
	fmt.Printf("%+v\n", sa)
	fmt.Printf("fuzz cgo test end..\n")
}
