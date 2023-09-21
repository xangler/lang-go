build-simple:
	go build -buildmode=c-archive -o simple/libgohello.a simple/gohello.go
	gcc -c -o simple/chello.o simple/chello.c
	ar rcs simple/libchello.a simple/chello.o
	go build -o simple.out cmd/simple/main.go

compile-compose:
	g++ --std=c++11 -c -o compose/zbuffer.o compose/zbuffer.cpp
	g++ --std=c++11 -c -o compose/cbuffer.o compose/cbuffer.cpp

build-cppdemo: compile-compose
	gcc cmd/cppdemo/main.c -I ./ -o cppdemo.out compose/cbuffer.o compose/zbuffer.o -lstdc++

build-compose: compile-compose
	ar rcs compose/libcbuffer.a compose/cbuffer.o compose/zbuffer.o 
	go build -o compose.out cmd/compose/main.go

build: build-simple build-cppdemo build-compose
	echo "build all succ"

.PHONY:
	build