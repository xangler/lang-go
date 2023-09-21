package main

import "C"
import "fmt"

func main() {}

//export GoSayHello
func GoSayHello(s *C.char) {
	fmt.Printf("%v", C.GoString(s))
}
