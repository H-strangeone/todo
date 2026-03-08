.PHONY: build install clean

build:
	go build -o bin/todo cmd/todo-tui/main.go

install:
	go install cmd/todo-tui/main.go

clean:
	rm -rf bin/