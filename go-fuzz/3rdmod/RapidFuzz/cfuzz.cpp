#include <stdlib.h>
#include <stdio.h>
#include <locale>
#include <codecvt>
#include "cfuzz.h"
#include "rapidfuzz/fuzz.hpp"
#include "rapidfuzz/details/types.hpp"

SCScoreAlignment* NewSCScoreAlignment(bool crash){
    SCScoreAlignment * p = (SCScoreAlignment *)malloc(sizeof(SCScoreAlignment));
    p->score = 0;
    p->src_start = 0;
    p->src_end = 0;
    p->dest_start = 0;
    p->dest_end = 0;
    p->crash = crash;
    return p;
}

SCScoreAlignment* NewSCScoreAlignment(rapidfuzz::ScoreAlignment<double>* sa){
    SCScoreAlignment * p = (SCScoreAlignment *)malloc(sizeof(SCScoreAlignment));
    p->score = sa->score;
    p->src_start = sa->src_start;
    p->src_end = sa->src_end;
    p->dest_start = sa->dest_start;
    p->dest_end = sa->dest_end;
    p->crash = false;
    return p;
}

SCScoreAlignment* SCPartialRatioAlignment(char* src, char* dst){
    try{
        std::wstring_convert<std::codecvt_utf8<wchar_t>> converter;
        auto sa = rapidfuzz::fuzz::partial_ratio_alignment(converter.from_bytes(src), converter.from_bytes(dst));
        SCScoreAlignment * p = NewSCScoreAlignment(&sa);
        return p;
    } catch(const std::exception &e){
        SCScoreAlignment * p = NewSCScoreAlignment(true);
        return p;
    }
}

void DeleteSCScoreAlignment(SCScoreAlignment* p){
    delete p;
}
