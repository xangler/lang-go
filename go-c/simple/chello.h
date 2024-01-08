#ifndef C_HELLO_H
#define C_HELLO_H

struct CObj{
    int mAge;
    char* mName;
};

void CPrintCObj(void* obj);
void SetCObj(void ** obj);
void FreeCObj(void* obj);

void GSayHello(char *s);
int CDiv(int a, int b);
void CSayHello(const char* s);
int CAddMth(int a, int b);

#endif // !C_HELLO_H
