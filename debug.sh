export GO111MODULE=on
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
go mod tidy
go build -o remote *.go