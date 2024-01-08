package rapid_fuzz

//#cgo CPPFLAGS: -I../../3rdmod/RapidFuzz
//#cgo LDFLAGS: -L../../3rdmod/RapidFuzz -lcfuzz -lstdc++ -lm
//
//#include <stdlib.h>
//#include "cfuzz.h"
import "C"
import "unsafe"

type GScoreAlignment struct {
	Score    float64
	SrcStart int
	SrcEnd   int
	DstStart int
	DstEnd   int
	Crash bool
}

func PartialRatioAlignment(src, dst string) *GScoreAlignment {
	csrc := C.CString(src)
	cdst := C.CString(dst)
	defer C.free(unsafe.Pointer(csrc))
	defer C.free(unsafe.Pointer(cdst))
	sa := (*C.struct_CScoreAlignment)(C.SCPartialRatioAlignment(csrc, cdst))
	defer C.DeleteSCScoreAlignment(sa)
	return &GScoreAlignment{
		Score:    float64(sa.score),
		SrcStart: int(sa.src_start),
		SrcEnd:   int(sa.src_end),
		DstStart: int(sa.dest_start),
		DstEnd:   int(sa.dest_end),
		Crash:    bool(sa.crash),
	}
}
