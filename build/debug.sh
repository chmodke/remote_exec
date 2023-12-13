export GO111MODULE=on
export CGO_ENABLED=0
go mod tidy

rm -rf ../dest
mkdir ../dest

export GOOS=linux
export GOARCH=amd64
go build -ldflags "-s -w" -o ../dest/remote_x86 ../*.go
upx ../dest/remote_x86

#export GOOS=linux
#export GOARCH=arm64
#go build -ldflags "-s -w" -o ../dest/remote_arm ../*.go
#upx ../dest/remote_arm

cp ../config.yaml ../dest/
cp ../command.yaml ../dest/
cp ../README.md ../dest/

ssh kehao@chmodke.org rm -f /home/kehao/remote_exec/remote_x86
scp ../dest/remote_x86 kehao@chmodke.org:/home/kehao/remote_exec/
ssh kehao@chmodke.org chmod +x /home/kehao/remote_exec/remote_x86
