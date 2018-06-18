.PHONY: go c run

go:
	go build -o ogame.so -buildmode=c-shared main.go

c:
	gcc -o main main.c ./ogame.so

run:
	./main
