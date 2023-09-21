#ifndef C_BUFFER_H
#define C_BUFFER_H

#ifdef __cplusplus
extern "C"{
#endif
    typedef struct CBuffer SCBuffer;
    SCBuffer* NewSCBuffer(char* data);
    void DeleteSCBuffer(SCBuffer* p);

    char* SCBufferData(SCBuffer* p);
    int SCBufferSize(SCBuffer* p);
#ifdef __cplusplus
};
#endif
#endif // !C_BUFFER_H
