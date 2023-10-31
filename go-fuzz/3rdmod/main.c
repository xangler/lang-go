#include <stdio.h>
#include "cfuzz.h"

int main(void){
    char* a = "mtxst";
    char* b = "this is a test";
    SCScoreAlignment* alignment = SCPartialRatioAlignment(a,b);
    printf("score %f\n", alignment->score);
    printf("src_start %d\n", alignment->src_start);
    printf("src_end %d\n", alignment->src_end);
    printf("dest_start %d\n", alignment->dest_start);
    printf("dest_end %d\n", alignment->dest_end);
    return 0;
}
