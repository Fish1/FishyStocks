install:
	go build -o bin/fishystocks

run: install
	./bin/fishystocks
