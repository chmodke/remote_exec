readonly BUILD_DIR='build'
readonly BIN_FILE='remote_linux_amd64'

readonly CUR_DIR=$(dirname "$0")

ssh kehao@chmodke.org rm -f /home/kehao/remote_exec/${BIN_FILE}
scp "${CUR_DIR}/../${BUILD_DIR}/${BIN_FILE}" kehao@chmodke.org:/home/kehao/remote_exec/
ssh kehao@chmodke.org chmod +x /home/kehao/remote_exec/${BIN_FILE}
