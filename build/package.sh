export GO111MODULE=on
export CGO_ENABLED=0
go mod tidy

rm -rf ../dest
mkdir ../dest

export GOOS=linux
export GOARCH=amd64
go build -ldflags "-s -w" -o ../dest/remote_x86 ../*.go
upx ../dest/remote_x86

export GOOS=linux
export GOARCH=arm64
go build -ldflags "-s -w" -o ../dest/remote_arm ../*.go
upx ../dest/remote_arm

cp ../config.yaml ../dest/
cp ../command.yaml ../dest/
cp ../README.md ../dest/

cd ../dest

tar -czf remote_exec.tar.gz *
