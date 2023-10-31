build:
	git submodule init
	git submodule update
	cd 3rdmod/RapidFuzz && cmake . && make
	go build -o fuzz.out cmd/fuzz/main.go

.PHONY:
	build