#include "zbuffer.h"

ZBuffer::ZBuffer(char* data){
    mps = data;
}

ZBuffer::~ZBuffer(){
}

int ZBuffer::Size(){
    return mps.size();
}

char* ZBuffer::Data(){
    return (char*)mps.data();
}