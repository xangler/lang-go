#include "cbuffer.h"
#include "zbuffer.h"

struct CBuffer: public ZBuffer{
    CBuffer(char * data):ZBuffer(data){}
    ~CBuffer(){}
};

SCBuffer* NewSCBuffer(char * data){
    auto p = new SCBuffer(data);
    return p;
}

void DeleteSCBuffer(SCBuffer* p){
    delete p;
}

char* SCBufferData(SCBuffer* p){
    return p->Data();
}

int SCBufferSize(SCBuffer* p){
    return p->Size();
}
