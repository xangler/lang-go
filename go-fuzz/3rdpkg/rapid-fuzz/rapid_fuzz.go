package rapid_fuzz

//#cgo CPPFLAGS: -I../../3rdmod/RapidFuzz
//#cgo LDFLAGS: -L../../3rdmod/RapidFuzz -lcfuzz -lstdc++ -lm
//
//#include "cfuzz.h"
import "C"

type GScoreAlignment struct {
	Score    float64
	SrcStart uint
	SrcEnd   uint
	DstStart uint
	DstEnd   uint
}

func PartialRatioAlignment(src, dst string) *GScoreAlignment {
	sa := (*C.struct_CScoreAlignment)(C.SCPartialRatioAlignment(C.CString(src), C.CString(dst)))
	defer C.DeleteSCScoreAlignment(sa)
	return &GScoreAlignment{
		Score:    float64(sa.score),
		SrcStart: uint(sa.src_start),
		SrcEnd:   uint(sa.src_end),
		DstStart: uint(sa.dest_start),
		DstEnd:   uint(sa.dest_end),
	}
}
