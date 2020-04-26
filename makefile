install:
	go build -o bin/fishstock

run: install
	./bin/fishstock
