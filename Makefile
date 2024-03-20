.PHONY: test clean


bin/pg_expulo: *.go
	go build -o $@

clean:
	rm -f bin/pg_expulo

test:
	go test

c.out: *.go
	go test -coverprofile=c.out

coverage.html: c.out
	go tool cover -html=c.out -o coverage.html
