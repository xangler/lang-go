cmake_minimum_required (VERSION 3.12...3.24)

set(CMAKE_INTERPROCEDURAL_OPTIMIZATION TRUE)
set(CMAKE_POSITION_INDEPENDENT_CODE ON)
set(SKBUILD_LINK_LIBRARIES_KEYWORD PRIVATE)
set(THREADS_PREFER_PTHREAD_FLAG ON)

project (Fuzz)

add_compile_options(-O3 -DNDEBUG -Wall -Wextra -pedantic -Wno-psabi -Wno-overlength-strings -mavx2)

aux_source_directory(. DIR_SRCS)
add_subdirectory(../../3rdparty/RapidFuzz ../../3rdparty/RapidFuzz)

add_executable(Mains.out mains.cpp)
target_link_libraries(Mains.out rapidfuzz::rapidfuzz)

add_library(cfuzz STATIC cfuzz.cpp)
target_include_directories(cfuzz PUBLIC .)
target_link_libraries(cfuzz rapidfuzz::rapidfuzz)

add_executable(Mainc.out mainc.c)
target_link_libraries(Mainc.out cfuzz)

