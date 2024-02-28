.PHONY: test clean

clean:
	rm -f bin/pg_expulo

test:
	go test

bin/pg_expulo:
	go build -o $@
