package main

//#cgo CPPFLAGS: -I../../3rdmod
//#cgo LDFLAGS: -L../../3rdmod -lcfuzz -lstdc++ -lm
//
//#include "cfuzz.h"
import "C"

import (
	"fmt"
)

func main() {
	fmt.Printf("fuzz cgo test start..\n")
	n := (*C.struct_CScoreAlignment)(C.SCPartialRatioAlignment(C.CString("ago"), C.CString("hello cgo")))
	defer C.DeleteSCScoreAlignment(n)
	fmt.Printf("%+v\n", n.score)
	fmt.Printf("fuzz cgo test end..\n")
}
