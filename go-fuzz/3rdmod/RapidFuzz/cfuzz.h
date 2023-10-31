#ifndef C_SCORE_ALIGNMENT_H
#define C_SCORE_ALIGNMENT_H

#ifdef __cplusplus
extern "C"{
#endif
    struct CScoreAlignment {
        double score;
        size_t src_start;
        size_t src_end;
        size_t dest_start;
        size_t dest_end;
    };
    typedef struct CScoreAlignment SCScoreAlignment;
    SCScoreAlignment* SCPartialRatioAlignment(char* src, char* dst);
    void DeleteSCScoreAlignment(SCScoreAlignment* p);

#ifdef __cplusplus
};
#endif
#endif // !C_SCORE_ALIGNMENT_H
