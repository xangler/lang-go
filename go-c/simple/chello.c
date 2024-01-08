#include <stdio.h>
#include <stdlib.h>
#include "chello.h"
#include "libgohello.h"

void GSayHello(char*s){
   GoSayHello(s); 
}

void CSayHello(const char* s){
    printf("%s", s);
}

int CDiv(int a, int b){
    return a/b;
}

void SetCObj(void ** obj){
    struct CObj * nobj = (struct CObj *)malloc(sizeof(struct CObj));
    nobj->mAge = 5;
    nobj->mName = "demo";
    *obj = nobj;
}

void FreeCObj(void* obj){
    free((struct CObj *) obj);
}

void CPrintCObj(void* obj){
    struct CObj * nobj = (struct CObj *) obj;
    printf("cprint age: %d, name:%s\n", nobj->mAge, nobj->mName);
}

int CAddMth(int a, int b){
    int c = a + b;
    printf("%d + %d = %d\n", a, b,c);
    return c;
}