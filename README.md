# remote_exec

linux服务器批量操作工具。

提供三个功能：

1. 批量上传文件
2. 以root用户批量执行命令
3. 批量下载文件

通过组合这三个功能能完成大多数Linux服务器上的操作。

## 配置方法

### 服务器配置

config.yaml文件

```yaml
port: 22
root-prompt: '#'
passwd-prompt: '.*assword.*'
timeout: 10
user: 'kehao'
passwd: '123456'
root-passwd: '123456'
hosts:
  - '172.18.0.3'
  - '172.18.0.4'
  - '172.18.0.6'
  - '172.18.0.7'
  - '172.18.0.8'
  - '172.18.0.9'
  - '172.18.0.10'
spc-hosts:
  - '127.0.0.1 23 123456 123456'
```

- port: ssh端口；
- rootPrompt：root用户的提示信息，golang正则表达式；
- passwdPrompt：root用户输入password的提示信息，golang正则表达式；
- timeout：等待登录的超时时长，单位：秒；
- user：ssh登录用户；
- passwd：ssh登录密码；
- rootPwd：root用户密码；
- hosts：服务器IP地址列表；
- spc_hosts：特殊密码的服务器列表，格式：`IP PORT 用户密码 root用户密码`。

### 命令配置文件

command.yaml文件

```yaml
put:
  - /home/kehao/test.txt#/home/kehao
  - /home/kehao/test.sh#/home/kehao
exec:
  - /bin/sh /home/kehao/test.sh
get:
  - /home/kehao/result.txt#/home/kehao/result
```

- put：上传文件列表，格式：`本地文件路径#远程目录名称`，远程目录名称可缺省，默认与本地文件路径的目录名相同，ssh登录用户需要具有远程目录的写权限；
- exec：需要执行的命令列表，命令将以`root`用户执行；
- get：下载的文件列表，格式：`远程文件路径#本地目录名称`，本地目录名称可缺省，默认与远程文件路径的目录名相同，ssh登录用户需要具有远程文件的读权限。

## 使用方法

### 命令

#### 批量上传

```shell
remote put
```

#### 批量执行

```shell
remote exec
```

#### 批量下载

```shell
remote get
```

### 参数

```shell
remote --help
remote execute tool

Usage:
  remote [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  exec        execute command on remote
  get         get file from remote
  help        Help about any command
  put         put file to remote

Flags:
  -c, --command string   Specify commands configuration (default "command.yaml")
  -f, --config string    Specify servers configuration (default "config.yaml")
  -h, --help             help for remote
  -t, --thread int       maximum number of concurrent (0 < t <= 16) (default 1)
  -v, --version          version for remote

Use "remote [command] --help" for more information about a command.
```

- 服务器配置文件（f），默认值config.yaml
- 命令配置文件（c），默认值command.yaml
- 最大并发数（t），0 < t <= 16，默认值1，开启并发执行能提高执行效率，但是输出结果会乱序。

## 使用样例

测试脚本：

test.sh

```shell
hostname -s > /home/kehao/result.txt
hostname -i >> /home/kehao/result.txt
id >> /home/kehao/result.txt
ls -l /home/kehao/test.txt  >> /home/kehao/result.txt
chown kehao: /home/kehao/result.txt
```

运行日志：

```powershell
kehao@kehaopcs ~/remote_exec $ ./remote put
2023/08/05 11:39:40 start put file...
2023/08/05 11:39:40 [172.18.0.3:22] put /home/kehao/test.txt to /home/kehao.
2023/08/05 11:39:40 [172.18.0.3:22] put /home/kehao/test.txt finished!
2023/08/05 11:39:40 [172.18.0.3:22] put /home/kehao/test.sh to /home/kehao.
2023/08/05 11:39:41 [172.18.0.3:22] put /home/kehao/test.sh finished!
2023/08/05 11:39:41 [172.18.0.4:22] put /home/kehao/test.txt to /home/kehao.
2023/08/05 11:39:41 [172.18.0.4:22] put /home/kehao/test.txt finished!
2023/08/05 11:39:41 [172.18.0.4:22] put /home/kehao/test.sh to /home/kehao.
2023/08/05 11:39:41 [172.18.0.4:22] put /home/kehao/test.sh finished!
2023/08/05 11:39:41 [127.0.0.1:22] put /home/kehao/test.txt to /home/kehao.
2023/08/05 11:39:42 [127.0.0.1:22] put /home/kehao/test.txt finished!
2023/08/05 11:39:42 [127.0.0.1:22] put /home/kehao/test.sh to /home/kehao.
2023/08/05 11:39:42 [127.0.0.1:22] put /home/kehao/test.sh finished!
2023/08/05 11:39:42 put file finished.

kehao@kehaopcs ~/remote_exec $ ./remote exec
2023/08/05 11:39:45 start execute command...
2023/08/05 11:39:45 [172.18.0.3:22] execute (/bin/sh /home/kehao/test.sh).
2023/08/05 11:39:46 [172.18.0.3:22] execute (/bin/sh /home/kehao/test.sh) result: 
I am 172.18.0.3 :)
[root@5421d92524c8 ~]# 
2023/08/05 11:39:46 [172.18.0.3:22] execute command finished!
2023/08/05 11:39:46 [172.18.0.4:22] execute (/bin/sh /home/kehao/test.sh).
2023/08/05 11:39:48 [172.18.0.4:22] execute (/bin/sh /home/kehao/test.sh) result: 
I am 172.18.0.4 :)
[root@9c0bbfbc1446 ~]# 
2023/08/05 11:39:48 [172.18.0.4:22] execute command finished!
2023/08/05 11:39:48 [127.0.0.1:22] execute (/bin/sh /home/kehao/test.sh).
2023/08/05 11:39:49 [127.0.0.1:22] execute (/bin/sh /home/kehao/test.sh) result: 
I am 172.18.0.5 :)
[root@93d37bd5e970 ~]# 
2023/08/05 11:39:49 [127.0.0.1:22] execute command finished!
2023/08/05 11:39:49 execute command finished.

kehao@kehaopcs ~/remote_exec $ ./remote get
2023/08/05 11:40:09 start get file...
2023/08/05 11:40:09 [172.18.0.3:22] get file from /home/kehao/result.txt to /home/kehao/result.
2023/08/05 11:40:09 [172.18.0.3:22] get /home/kehao/result.txt finished!
2023/08/05 11:40:09 [172.18.0.4:22] get file from /home/kehao/result.txt to /home/kehao/result.
2023/08/05 11:40:10 [172.18.0.4:22] get /home/kehao/result.txt finished!
2023/08/05 11:40:10 [127.0.0.1:22] get file from /home/kehao/result.txt to /home/kehao/result.
2023/08/05 11:40:10 [127.0.0.1:22] get /home/kehao/result.txt finished!
2023/08/05 11:40:10 get file finished.
```

操作结果：

远端

```powershell
[root@kehaopcs ~]# docker exec -u kehao test-ssh1 bash -c 'ls -l ~'
total 12
-rw-r--r--. 1 kehao kehao 126 Aug  5 11:39 result.txt
-rw-rw-r--. 1 kehao kehao 194 Aug  5 11:39 test.sh
-rw-rw-r--. 1 kehao kehao  10 Aug  5 11:39 test.txt

[root@kehaopcs ~]# docker exec -u kehao test-ssh2 bash -c 'ls -l ~'
total 12
-rw-r--r--. 1 kehao kehao 126 Aug  5 11:39 result.txt
-rw-rw-r--. 1 kehao kehao 194 Aug  5 11:39 test.sh
-rw-rw-r--. 1 kehao kehao  10 Aug  5 11:39 test.txt

[root@kehaopcs ~]# docker exec -u kehao test-ssh3 bash -c 'ls -l ~'
total 12
-rw-r--r--. 1 kehao kehao 126 Aug  5 11:39 result.txt
-rw-rw-r--. 1 kehao kehao 194 Aug  5 11:39 test.sh
-rw-rw-r--. 1 kehao kehao  10 Aug  5 11:39 test.txt
```

本地：

```powershell
kehao@kehaopcs ~/remote_exec $ ls -l /home/kehao/result/
total 12
-rw-r--r--. 1 kehao manager 126 Aug  5 11:40 172.18.0.3_22_result.txt
-rw-r--r--. 1 kehao manager 126 Aug  5 11:40 172.18.0.4_22_result.txt
-rw-r--r--. 1 kehao manager 126 Aug  5 11:40 127.0.0.1_23_result.txt

kehao@kehaopcs ~/remote_exec $ cat /home/kehao/result/172.18.0.3_22_result.txt 
5421d92524c8
172.18.0.3
uid=0(root) gid=0(root) groups=0(root)
-rw-rw-r--. 1 kehao kehao 10 Aug  5 11:39 /home/kehao/test.txt

kehao@kehaopcs ~/remote_exec $ cat /home/kehao/result/172.18.0.4_22_result.txt 
9c0bbfbc1446
172.18.0.4
uid=0(root) gid=0(root) groups=0(root)
-rw-rw-r--. 1 kehao kehao 10 Aug  5 11:39 /home/kehao/test.txt

kehao@kehaopcs ~/remote_exec $ cat /home/kehao/result/127.0.0.1_23_result.txt 
93d37bd5e970
172.18.0.5
uid=0(root) gid=0(root) groups=0(root)
-rw-rw-r--. 1 kehao kehao 10 Aug  5 11:39 /home/kehao/test.txt
```