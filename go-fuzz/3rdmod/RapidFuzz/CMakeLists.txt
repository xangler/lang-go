cmake_minimum_required (VERSION 3.0)

project (Fuzz)

aux_source_directory(. DIR_SRCS)
add_subdirectory(../../3rdparty/RapidFuzz ../../3rdparty/RapidFuzz)

add_library(cfuzz STATIC cfuzz.cpp)
target_include_directories(cfuzz PUBLIC .)
target_link_libraries(cfuzz rapidfuzz::rapidfuzz)

add_executable(Demo.out main.c)
target_link_libraries(Demo.out cfuzz)

