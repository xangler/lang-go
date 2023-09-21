#ifndef Z_BUFFER_H
#define Z_BUFFER_H
#include<string> 

class ZBuffer{
    public:
    std::string mps;

    ZBuffer(char* data);
    ~ZBuffer();
    int Size();
    char* Data();
};

#endif // !Z_BUFFER_H
