package main

//#cgo CFLAGS: -I../../simple
//#cgo LDFLAGS: -L../../simple -lchello
//
//#include "chello.h"
import "C"

import (
	"fmt"
	"unsafe"
)

//export GoSayHello
func GoSayHello(s *C.char) {
	fmt.Printf(C.GoString(s))
}

func main() {
	fmt.Printf("simple cgo test start..\n")
	C.CSayHello(C.CString("chello cgo\n"))
	C.GSayHello(C.CString("gohello cgo\n"))
	v := C.CAddMth(C.int(3), C.int(2))
	fmt.Printf("C call return:%v\n", v)
	var m C.struct_CObj
	m.mAge = C.int(5)
	fmt.Printf("C struct member mAge:%v\n", m.mAge)
	var n unsafe.Pointer
	C.SetCObj(&n)
	C.CPrintCObj(n)
	x := (*C.struct_CObj)(n)
	fmt.Printf("C struct member mAge:%v\n", x.mAge)
	fmt.Printf("C struct member mName:%v\n", C.GoString(x.mName))
	C.FreeCObj(n)
	fmt.Printf("simple cgo test end..\n")
}
