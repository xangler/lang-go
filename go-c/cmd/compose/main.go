package main

//#cgo CPPFLAGS: -I../../compose
//#cgo LDFLAGS: -L../../compose -lcbuffer -lstdc++
//
//#include "cbuffer.h"
import "C"

import (
	"fmt"
)

func main() {
	fmt.Printf("compose cgo test start..\n")
	n := (*C.struct_CBuffer)(C.NewSCBuffer(C.CString("hello cgo")))
	v := C.SCBufferSize(n)
	d := C.SCBufferData(n)
	fmt.Printf("C buffer size:%d, content:%s\n", v, C.GoString(d))
	C.DeleteSCBuffer(n)
	fmt.Printf("compose cgo test end..\n")
}
