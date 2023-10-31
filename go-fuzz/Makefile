build:
	cd 3rdmod && cmake . && make
	go build -o fuzz.out cmd/fuzz/main.go

.PHONY:
	build