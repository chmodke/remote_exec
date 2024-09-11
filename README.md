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
- root-prompt：root用户的提示信息，golang正则表达式；
- passwd-prompt：root用户输入password的提示信息，golang正则表达式；
- timeout：等待登录的超时时长，单位：秒；
- user：ssh登录用户；
- passwd：ssh登录密码；
- root-passwd：root用户密码；
- hosts：服务器IP地址列表；
- spc_hosts：特殊密码的服务器列表，格式：`IP PORT 用户密码 root用户密码`。

### 命令配置文件

command.yaml文件

```yaml
put:
  - /home/kehao/test1.txt#/home/kehao
  - /home/kehao/test2.txt#/home/kehao
  - /home/kehao/test1.sh#/home/kehao
  - /home/kehao/test2.sh#/home/kehao
exec:
  - /bin/sh /home/kehao/test1.sh
  - /bin/sh /home/kehao/test2.sh
get:
  - /home/kehao/result1.txt#/home/kehao/result
  - /home/kehao/result2.txt#/home/kehao/result

arthas:
  put:
    - /home/kehao/arthas.sh#/home/kehao
    - /home/kehao/jdk8.sh#/home/kehao
  exec:
    - /bin/sh /home/kehao/arthas.sh
    - /bin/sh /home/kehao/jdk8.sh
  get:
    - /home/kehao/arthas-result.txt#/home/kehao/result
    - /home/kehao/jdk8-result.txt#/home/kehao/result
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
  -m, --netmask string   ip filter, e.g. 192.168.1.1 192.168.1.1,192.168.1.2 192.168.0.0/24
  -t, --thread int       maximum number of concurrent (0 < t <= 16) (default 1)
  -v, --version          version for remote

Use "remote [command] --help" for more information about a command.
```

- 服务器配置文件（f），默认值config.yaml
- 命令配置文件（c），默认值command.yaml
- 最大并发数（t），0 < t <= 16，默认值1，开启并发执行能提高执行效率，但是输出结果会乱序
- IP地址过滤（m），默认为空，支持192.168.1.1、192.168.1.1,192.168.1.2、192.168.0.0/24三种写法，匹配的主机才会执行

## 使用样例

测试脚本：

test1.sh

```shell
echo "test1.sh" >/home/kehao/result1.txt
hostname -s >>/home/kehao/result1.txt
hostname -i >>/home/kehao/result1.txt
id >>/home/kehao/result1.txt
ls -l /home/kehao/test1.txt >>/home/kehao/result1.txt
cat /home/kehao/test1.txt >>/home/kehao/result1.txt
chown kehao: /home/kehao/result1.txt
echo "I am $(hostname -i) :)"
```

test2.sh

```shell
echo "test2.sh" >/home/kehao/result2.txt
hostname -s >>/home/kehao/result2.txt
hostname -i >>/home/kehao/result2.txt
id >>/home/kehao/result2.txt
ls -l /home/kehao/test2.txt >>/home/kehao/result2.txt
cat /home/kehao/test2.txt >>/home/kehao/result2.txt
chown kehao: /home/kehao/result2.txt
echo "I am $(hostname -i) :)"
```

运行日志：

```powershell
kehao@kehaopcs ~/remote_exec $ ./remote_x86 put
2024/02/04 16:59:32 start put file...
2024/02/04 16:59:32 progress [7/8]...
2024/02/04 16:59:32 [172.18.0.3:22] upload /home/kehao/test1.txt to /home/kehao.
2024/02/04 16:59:32 [172.18.0.3:22] upload /home/kehao/test1.txt finished!
2024/02/04 16:59:32 [172.18.0.3:22] upload /home/kehao/test2.txt to /home/kehao.
2024/02/04 16:59:32 [172.18.0.3:22] upload /home/kehao/test2.txt finished!
2024/02/04 16:59:32 [172.18.0.3:22] upload /home/kehao/test1.sh to /home/kehao.
2024/02/04 16:59:32 [172.18.0.3:22] upload /home/kehao/test1.sh finished!
2024/02/04 16:59:32 [172.18.0.3:22] upload /home/kehao/test2.sh to /home/kehao.
2024/02/04 16:59:32 [172.18.0.3:22] upload /home/kehao/test2.sh finished!
2024/02/04 16:59:32 progress [6/8]...
2024/02/04 16:59:32 [172.18.0.4:22] upload /home/kehao/test1.txt to /home/kehao.
2024/02/04 16:59:32 [172.18.0.4:22] upload /home/kehao/test1.txt finished!
2024/02/04 16:59:32 [172.18.0.4:22] upload /home/kehao/test2.txt to /home/kehao.
2024/02/04 16:59:32 [172.18.0.4:22] upload /home/kehao/test2.txt finished!
2024/02/04 16:59:32 [172.18.0.4:22] upload /home/kehao/test1.sh to /home/kehao.
2024/02/04 16:59:32 [172.18.0.4:22] upload /home/kehao/test1.sh finished!
2024/02/04 16:59:32 [172.18.0.4:22] upload /home/kehao/test2.sh to /home/kehao.
2024/02/04 16:59:32 [172.18.0.4:22] upload /home/kehao/test2.sh finished!
2024/02/04 16:59:32 progress [5/8]...
2024/02/04 16:59:32 [172.18.0.6:22] upload /home/kehao/test1.txt to /home/kehao.
2024/02/04 16:59:32 [172.18.0.6:22] upload /home/kehao/test1.txt finished!
2024/02/04 16:59:32 [172.18.0.6:22] upload /home/kehao/test2.txt to /home/kehao.
2024/02/04 16:59:32 [172.18.0.6:22] upload /home/kehao/test2.txt finished!
2024/02/04 16:59:32 [172.18.0.6:22] upload /home/kehao/test1.sh to /home/kehao.
2024/02/04 16:59:32 [172.18.0.6:22] upload /home/kehao/test1.sh finished!
2024/02/04 16:59:32 [172.18.0.6:22] upload /home/kehao/test2.sh to /home/kehao.
2024/02/04 16:59:32 [172.18.0.6:22] upload /home/kehao/test2.sh finished!
2024/02/04 16:59:32 progress [4/8]...
2024/02/04 16:59:33 [172.18.0.7:22] upload /home/kehao/test1.txt to /home/kehao.
2024/02/04 16:59:33 [172.18.0.7:22] upload /home/kehao/test1.txt finished!
2024/02/04 16:59:33 [172.18.0.7:22] upload /home/kehao/test2.txt to /home/kehao.
2024/02/04 16:59:33 [172.18.0.7:22] upload /home/kehao/test2.txt finished!
2024/02/04 16:59:33 [172.18.0.7:22] upload /home/kehao/test1.sh to /home/kehao.
2024/02/04 16:59:33 [172.18.0.7:22] upload /home/kehao/test1.sh finished!
2024/02/04 16:59:33 [172.18.0.7:22] upload /home/kehao/test2.sh to /home/kehao.
2024/02/04 16:59:33 [172.18.0.7:22] upload /home/kehao/test2.sh finished!
2024/02/04 16:59:33 progress [3/8]...
2024/02/04 16:59:33 [172.18.0.8:22] upload /home/kehao/test1.txt to /home/kehao.
2024/02/04 16:59:33 [172.18.0.8:22] upload /home/kehao/test1.txt finished!
2024/02/04 16:59:33 [172.18.0.8:22] upload /home/kehao/test2.txt to /home/kehao.
2024/02/04 16:59:33 [172.18.0.8:22] upload /home/kehao/test2.txt finished!
2024/02/04 16:59:33 [172.18.0.8:22] upload /home/kehao/test1.sh to /home/kehao.
2024/02/04 16:59:33 [172.18.0.8:22] upload /home/kehao/test1.sh finished!
2024/02/04 16:59:33 [172.18.0.8:22] upload /home/kehao/test2.sh to /home/kehao.
2024/02/04 16:59:33 [172.18.0.8:22] upload /home/kehao/test2.sh finished!
2024/02/04 16:59:33 progress [2/8]...
2024/02/04 16:59:33 [172.18.0.9:22] upload /home/kehao/test1.txt to /home/kehao.
2024/02/04 16:59:33 [172.18.0.9:22] upload /home/kehao/test1.txt finished!
2024/02/04 16:59:33 [172.18.0.9:22] upload /home/kehao/test2.txt to /home/kehao.
2024/02/04 16:59:33 [172.18.0.9:22] upload /home/kehao/test2.txt finished!
2024/02/04 16:59:33 [172.18.0.9:22] upload /home/kehao/test1.sh to /home/kehao.
2024/02/04 16:59:33 [172.18.0.9:22] upload /home/kehao/test1.sh finished!
2024/02/04 16:59:33 [172.18.0.9:22] upload /home/kehao/test2.sh to /home/kehao.
2024/02/04 16:59:33 [172.18.0.9:22] upload /home/kehao/test2.sh finished!
2024/02/04 16:59:33 progress [1/8]...
2024/02/04 16:59:33 [172.18.0.10:22] upload /home/kehao/test1.txt to /home/kehao.
2024/02/04 16:59:33 [172.18.0.10:22] upload /home/kehao/test1.txt finished!
2024/02/04 16:59:33 [172.18.0.10:22] upload /home/kehao/test2.txt to /home/kehao.
2024/02/04 16:59:33 [172.18.0.10:22] upload /home/kehao/test2.txt finished!
2024/02/04 16:59:33 [172.18.0.10:22] upload /home/kehao/test1.sh to /home/kehao.
2024/02/04 16:59:33 [172.18.0.10:22] upload /home/kehao/test1.sh finished!
2024/02/04 16:59:33 [172.18.0.10:22] upload /home/kehao/test2.sh to /home/kehao.
2024/02/04 16:59:33 [172.18.0.10:22] upload /home/kehao/test2.sh finished!
2024/02/04 16:59:33 progress [0/8]...
2024/02/04 16:59:34 [127.0.0.1:23] upload /home/kehao/test1.txt to /home/kehao.
2024/02/04 16:59:34 [127.0.0.1:23] upload /home/kehao/test1.txt finished!
2024/02/04 16:59:34 [127.0.0.1:23] upload /home/kehao/test2.txt to /home/kehao.
2024/02/04 16:59:34 [127.0.0.1:23] upload /home/kehao/test2.txt finished!
2024/02/04 16:59:34 [127.0.0.1:23] upload /home/kehao/test1.sh to /home/kehao.
2024/02/04 16:59:34 [127.0.0.1:23] upload /home/kehao/test1.sh finished!
2024/02/04 16:59:34 [127.0.0.1:23] upload /home/kehao/test2.sh to /home/kehao.
2024/02/04 16:59:34 [127.0.0.1:23] upload /home/kehao/test2.sh finished!
2024/02/04 16:59:34 put file finished.
kehao@kehaopcs ~/remote_exec $ 
kehao@kehaopcs ~/remote_exec $ ./remote_x86 exec
2024/02/04 16:59:37 start execute command...
2024/02/04 16:59:37 progress [7/8]...
2024/02/04 16:59:38 [172.18.0.3:22] execute (/bin/sh /home/kehao/test1.sh).
2024/02/04 16:59:38 [172.18.0.3:22] execute (/bin/sh /home/kehao/test1.sh) result: 
I am 172.18.0.3 :)
[root@5f0c9135f133 ~]# 
2024/02/04 16:59:38 [172.18.0.3:22] execute (/bin/sh /home/kehao/test2.sh).
2024/02/04 16:59:38 [172.18.0.3:22] execute (/bin/sh /home/kehao/test2.sh) result: 
I am 172.18.0.3 :)
[root@5f0c9135f133 ~]# 
2024/02/04 16:59:38 [172.18.0.3:22] execute command finished!
2024/02/04 16:59:38 progress [6/8]...
2024/02/04 16:59:38 [172.18.0.4:22] execute (/bin/sh /home/kehao/test1.sh).
2024/02/04 16:59:38 [172.18.0.4:22] execute (/bin/sh /home/kehao/test1.sh) result: 
I am 172.18.0.4 :)
[root@05d92c6bf1a1 ~]# 
2024/02/04 16:59:38 [172.18.0.4:22] execute (/bin/sh /home/kehao/test2.sh).
2024/02/04 16:59:38 [172.18.0.4:22] execute (/bin/sh /home/kehao/test2.sh) result: 
I am 172.18.0.4 :)
[root@05d92c6bf1a1 ~]# 
2024/02/04 16:59:38 [172.18.0.4:22] execute command finished!
2024/02/04 16:59:38 progress [5/8]...
2024/02/04 16:59:39 [172.18.0.6:22] execute (/bin/sh /home/kehao/test1.sh).
2024/02/04 16:59:39 [172.18.0.6:22] execute (/bin/sh /home/kehao/test1.sh) result: 
I am 172.18.0.6 :)
[root@48aa865c7d59 ~]# 
2024/02/04 16:59:39 [172.18.0.6:22] execute (/bin/sh /home/kehao/test2.sh).
2024/02/04 16:59:39 [172.18.0.6:22] execute (/bin/sh /home/kehao/test2.sh) result: 
I am 172.18.0.6 :)
[root@48aa865c7d59 ~]# 
2024/02/04 16:59:39 [172.18.0.6:22] execute command finished!
2024/02/04 16:59:39 progress [4/8]...
2024/02/04 16:59:39 [172.18.0.7:22] execute (/bin/sh /home/kehao/test1.sh).
2024/02/04 16:59:39 [172.18.0.7:22] execute (/bin/sh /home/kehao/test1.sh) result: 
I am 172.18.0.7 :)
[root@26d26a913b2f ~]# 
2024/02/04 16:59:39 [172.18.0.7:22] execute (/bin/sh /home/kehao/test2.sh).
2024/02/04 16:59:39 [172.18.0.7:22] execute (/bin/sh /home/kehao/test2.sh) result: 
I am 172.18.0.7 :)
[root@26d26a913b2f ~]# 
2024/02/04 16:59:39 [172.18.0.7:22] execute command finished!
2024/02/04 16:59:39 progress [3/8]...
2024/02/04 16:59:39 [172.18.0.8:22] execute (/bin/sh /home/kehao/test1.sh).
2024/02/04 16:59:39 [172.18.0.8:22] execute (/bin/sh /home/kehao/test1.sh) result: 
I am 172.18.0.8 :)
[root@86f1b1085a0d ~]# 
2024/02/04 16:59:39 [172.18.0.8:22] execute (/bin/sh /home/kehao/test2.sh).
2024/02/04 16:59:40 [172.18.0.8:22] execute (/bin/sh /home/kehao/test2.sh) result: 
I am 172.18.0.8 :)
[root@86f1b1085a0d ~]# 
2024/02/04 16:59:40 [172.18.0.8:22] execute command finished!
2024/02/04 16:59:40 progress [2/8]...
2024/02/04 16:59:40 [172.18.0.9:22] execute (/bin/sh /home/kehao/test1.sh).
2024/02/04 16:59:40 [172.18.0.9:22] execute (/bin/sh /home/kehao/test1.sh) result: 
I am 172.18.0.9 :)
[root@4c2db381a3d8 ~]# 
2024/02/04 16:59:40 [172.18.0.9:22] execute (/bin/sh /home/kehao/test2.sh).
2024/02/04 16:59:40 [172.18.0.9:22] execute (/bin/sh /home/kehao/test2.sh) result: 
I am 172.18.0.9 :)
[root@4c2db381a3d8 ~]# 
2024/02/04 16:59:40 [172.18.0.9:22] execute command finished!
2024/02/04 16:59:40 progress [1/8]...
2024/02/04 16:59:42 [172.18.0.10:22] execute (/bin/sh /home/kehao/test1.sh).
2024/02/04 16:59:42 [172.18.0.10:22] execute (/bin/sh /home/kehao/test1.sh) result: 
I am 172.18.0.10 :)
[root@abf48a7fca3b ~]# 
2024/02/04 16:59:42 [172.18.0.10:22] execute (/bin/sh /home/kehao/test2.sh).
2024/02/04 16:59:44 [172.18.0.10:22] execute (/bin/sh /home/kehao/test2.sh) result: 
I am 172.18.0.10 :)
[root@abf48a7fca3b ~]# 
2024/02/04 16:59:44 [172.18.0.10:22] execute command finished!
2024/02/04 16:59:44 progress [0/8]...
2024/02/04 16:59:45 [127.0.0.1:23] execute (/bin/sh /home/kehao/test1.sh).
2024/02/04 16:59:45 [127.0.0.1:23] execute (/bin/sh /home/kehao/test1.sh) result: 
I am 172.18.0.5 :)
[root@1244ba6621fc ~]# 
2024/02/04 16:59:45 [127.0.0.1:23] execute (/bin/sh /home/kehao/test2.sh).
2024/02/04 16:59:45 [127.0.0.1:23] execute (/bin/sh /home/kehao/test2.sh) result: 
I am 172.18.0.5 :)
[root@1244ba6621fc ~]# 
2024/02/04 16:59:45 [127.0.0.1:23] execute command finished!
2024/02/04 16:59:45 execute command finished.
kehao@kehaopcs ~/remote_exec $ 
kehao@kehaopcs ~/remote_exec $ ./remote_x86 get
2024/02/04 16:59:51 start get file...
2024/02/04 16:59:51 progress [7/8]...
2024/02/04 16:59:51 [172.18.0.3:22] download file from /home/kehao/result1.txt to /home/kehao/result.
2024/02/04 16:59:51 [172.18.0.3:22] download /home/kehao/result1.txt finished!
2024/02/04 16:59:51 [172.18.0.3:22] download file from /home/kehao/result2.txt to /home/kehao/result.
2024/02/04 16:59:51 [172.18.0.3:22] download /home/kehao/result2.txt finished!
2024/02/04 16:59:51 progress [6/8]...
2024/02/04 16:59:51 [172.18.0.4:22] download file from /home/kehao/result1.txt to /home/kehao/result.
2024/02/04 16:59:51 [172.18.0.4:22] download /home/kehao/result1.txt finished!
2024/02/04 16:59:51 [172.18.0.4:22] download file from /home/kehao/result2.txt to /home/kehao/result.
2024/02/04 16:59:51 [172.18.0.4:22] download /home/kehao/result2.txt finished!
2024/02/04 16:59:51 progress [5/8]...
2024/02/04 16:59:52 [172.18.0.6:22] download file from /home/kehao/result1.txt to /home/kehao/result.
2024/02/04 16:59:52 [172.18.0.6:22] download /home/kehao/result1.txt finished!
2024/02/04 16:59:52 [172.18.0.6:22] download file from /home/kehao/result2.txt to /home/kehao/result.
2024/02/04 16:59:52 [172.18.0.6:22] download /home/kehao/result2.txt finished!
2024/02/04 16:59:52 progress [4/8]...
2024/02/04 16:59:52 [172.18.0.7:22] download file from /home/kehao/result1.txt to /home/kehao/result.
2024/02/04 16:59:52 [172.18.0.7:22] download /home/kehao/result1.txt finished!
2024/02/04 16:59:52 [172.18.0.7:22] download file from /home/kehao/result2.txt to /home/kehao/result.
2024/02/04 16:59:52 [172.18.0.7:22] download /home/kehao/result2.txt finished!
2024/02/04 16:59:52 progress [3/8]...
2024/02/04 16:59:52 [172.18.0.8:22] download file from /home/kehao/result1.txt to /home/kehao/result.
2024/02/04 16:59:52 [172.18.0.8:22] download /home/kehao/result1.txt finished!
2024/02/04 16:59:52 [172.18.0.8:22] download file from /home/kehao/result2.txt to /home/kehao/result.
2024/02/04 16:59:52 [172.18.0.8:22] download /home/kehao/result2.txt finished!
2024/02/04 16:59:52 progress [2/8]...
2024/02/04 16:59:52 [172.18.0.9:22] download file from /home/kehao/result1.txt to /home/kehao/result.
2024/02/04 16:59:52 [172.18.0.9:22] download /home/kehao/result1.txt finished!
2024/02/04 16:59:52 [172.18.0.9:22] download file from /home/kehao/result2.txt to /home/kehao/result.
2024/02/04 16:59:52 [172.18.0.9:22] download /home/kehao/result2.txt finished!
2024/02/04 16:59:52 progress [1/8]...
2024/02/04 16:59:53 [172.18.0.10:22] download file from /home/kehao/result1.txt to /home/kehao/result.
2024/02/04 16:59:53 [172.18.0.10:22] download /home/kehao/result1.txt finished!
2024/02/04 16:59:53 [172.18.0.10:22] download file from /home/kehao/result2.txt to /home/kehao/result.
2024/02/04 16:59:53 [172.18.0.10:22] download /home/kehao/result2.txt finished!
2024/02/04 16:59:53 progress [0/8]...
2024/02/04 16:59:53 [127.0.0.1:23] download file from /home/kehao/result1.txt to /home/kehao/result.
2024/02/04 16:59:53 [127.0.0.1:23] download /home/kehao/result1.txt finished!
2024/02/04 16:59:53 [127.0.0.1:23] download file from /home/kehao/result2.txt to /home/kehao/result.
2024/02/04 16:59:53 [127.0.0.1:23] download /home/kehao/result2.txt finished!
2024/02/04 16:59:53 get file finished.
```

操作结果：

远端

```powershell
kehao@kehaopcs ~/remote_exec $ docker exec -u kehao test-ssh1 bash -c 'ls -l ~'
total 24
-rw-r--r--. 1 kehao kehao 146 Feb  4 16:59 result1.txt
-rw-r--r--. 1 kehao kehao 146 Feb  4 16:59 result2.txt
-rw-rw-r--. 1 kehao kehao 319 Feb  4 16:59 test1.sh
-rw-rw-r--. 1 kehao kehao  10 Feb  4 16:59 test1.txt
-rw-rw-r--. 1 kehao kehao 319 Feb  4 16:59 test2.sh
-rw-rw-r--. 1 kehao kehao  10 Feb  4 16:59 test2.txt
kehao@kehaopcs ~/remote_exec $ docker exec -u kehao test-ssh2 bash -c 'ls -l ~'
total 24
-rw-r--r--. 1 kehao kehao 146 Feb  4 16:59 result1.txt
-rw-r--r--. 1 kehao kehao 146 Feb  4 16:59 result2.txt
-rw-rw-r--. 1 kehao kehao 319 Feb  4 16:59 test1.sh
-rw-rw-r--. 1 kehao kehao  10 Feb  4 16:59 test1.txt
-rw-rw-r--. 1 kehao kehao 319 Feb  4 16:59 test2.sh
-rw-rw-r--. 1 kehao kehao  10 Feb  4 16:59 test2.txt
kehao@kehaopcs ~/remote_exec $ docker exec -u kehao test-ssh3 bash -c 'ls -l ~'
total 24
-rw-r--r--. 1 kehao kehao 146 Feb  4 16:59 result1.txt
-rw-r--r--. 1 kehao kehao 146 Feb  4 16:59 result2.txt
-rw-rw-r--. 1 kehao kehao 319 Feb  4 16:59 test1.sh
-rw-rw-r--. 1 kehao kehao  10 Feb  4 16:59 test1.txt
-rw-rw-r--. 1 kehao kehao 319 Feb  4 16:59 test2.sh
-rw-rw-r--. 1 kehao kehao  10 Feb  4 16:59 test2.txt
```

本地：

```powershell
kehao@kehaopcs ~/remote_exec $ ls -l /home/kehao/result/
total 64
-rw-r--r--. 1 kehao manager 146 Feb  4 16:59 127.0.0.1_23_result1.txt
-rw-r--r--. 1 kehao manager 146 Feb  4 16:59 127.0.0.1_23_result2.txt
-rw-r--r--. 1 kehao manager 147 Feb  4 16:59 172.18.0.10_22_result1.txt
-rw-r--r--. 1 kehao manager 147 Feb  4 16:59 172.18.0.10_22_result2.txt
-rw-r--r--. 1 kehao manager 146 Feb  4 16:59 172.18.0.3_22_result1.txt
-rw-r--r--. 1 kehao manager 146 Feb  4 16:59 172.18.0.3_22_result2.txt
-rw-r--r--. 1 kehao manager 146 Feb  4 16:59 172.18.0.4_22_result1.txt
-rw-r--r--. 1 kehao manager 146 Feb  4 16:59 172.18.0.4_22_result2.txt
-rw-r--r--. 1 kehao manager 146 Feb  4 16:59 172.18.0.6_22_result1.txt
-rw-r--r--. 1 kehao manager 146 Feb  4 16:59 172.18.0.6_22_result2.txt
-rw-r--r--. 1 kehao manager 146 Feb  4 16:59 172.18.0.7_22_result1.txt
-rw-r--r--. 1 kehao manager 146 Feb  4 16:59 172.18.0.7_22_result2.txt
-rw-r--r--. 1 kehao manager 146 Feb  4 16:59 172.18.0.8_22_result1.txt
-rw-r--r--. 1 kehao manager 146 Feb  4 16:59 172.18.0.8_22_result2.txt
-rw-r--r--. 1 kehao manager 146 Feb  4 16:59 172.18.0.9_22_result1.txt
-rw-r--r--. 1 kehao manager 146 Feb  4 16:59 172.18.0.9_22_result2.txt
kehao@kehaopcs ~/remote_exec $ cat /home/kehao/result/172.18.0.9_22_result2.txt
test2.sh
4c2db381a3d8
172.18.0.9
uid=0(root) gid=0(root) groups=0(root)
-rw-rw-r--. 1 kehao kehao 10 Feb  4 16:59 /home/kehao/test2.txt
0987654321
```
