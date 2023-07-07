
vet:
	go vet -vettool=$(which statictest-darwin-amd64) ./...


server:
	@echo "building server"
	go build -o ./cmd/gophermart/gophermart ./cmd/gophermart/*.go

clean:
	rm -f ./cmd/gophermart/gophermart

lnt:
	# golangci-lint run -v --enable-all --disable gochecknoglobals --disable paralleltest --disable exhaustivestruct --disable depguard --disable wsl
	golangci-lint run -v

fmt:
	# to install it:
	# go install mvdan.cc/gofumpt@latest
	gofumpt -l -w .

gci:
	# to install it:
	# go install github.com/daixiang0/gci@latest
	gci write --skip-generated -s default .

gofmt:
	gofmt -s -w .

fix: gofmt gci fmt

cover:
	rm -f ./cover.html cover.out coverage.txt
	go test -coverprofile cover.out  ./... ./internal/... -coverpkg=./...
	go tool cover -html=cover.out -o cover.html

