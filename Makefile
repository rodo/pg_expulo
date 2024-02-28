.PHONY: test clean

bin/pg_expulo:
	go build -o $@

clean:
	rm -f bin/pg_expulo

test:
	go test
