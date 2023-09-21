#include <stdlib.h>
#include <stdio.h>
#include "compose/cbuffer.h"

int main() {
    SCBuffer* obj = NewSCBuffer("hello cgo");
    int m = SCBufferSize(obj);
    printf("buffer size %d\n",m);
    DeleteSCBuffer(obj);    
    return 0;
}