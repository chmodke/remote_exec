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
rootPrompt: '#'
passwdPrompt: '.*assword.*'
timeout: 10
user: 'kehao'
passwd: '123456'
rootPwd: '123456'
hosts:
  - '172.17.0.3'
  - '172.17.0.4'
spc_hosts:
  - '172.17.0.5 123456 123456'
```

- port: ssh端口；
- rootPrompt：root用户的提示信息，golang正则表达式；
- passwdPrompt：root用户输入password的提示信息，golang正则表达式；
- timeout：等待登录的超时时长，单位：秒；
- user：ssh登录用户；
- passwd：ssh登录密码；
- rootPwd：root用户密码；
- hosts：服务器IP地址列表；
- spc_hosts：特殊密码的服务器列表，格式：`IP 用户密码 root用户密码`。

### 命令配置文件

command.yaml文件

```yaml
upload:
  - /home/kehao/test.txt#/home/kehao
  - /home/kehao/test.sh#/home/kehao
commands:
  - /bin/sh /home/kehao/test.sh
download:
  - /home/kehao/result.txt#/home/kehao/result
```

- upload：上传文件列表，格式：`本地文件路径#远程目录名称`，远程目录名称可缺省，默认与本地文件路径的目录名相同，ssh登录用户需要具有远程目录的写权限；
- commands：需要执行的命令列表，命令将以`root`用户执行；
- download：下载的文件列表，格式：`远程文件路径#本地目录名称`，本地目录名称可缺省，默认与远程文件路径的目录名相同，ssh登录用户需要具有远程文件的读权限。

## 使用方法

### 批量上传

```shell
remote upload
```

### 批量执行

```shell
remote execute
```

### 批量下载

```shell
remote download
```

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
kehao@kehaopcs ~/remote_exec $ ./remote upload
2023/08/05 11:39:40 start upload file...
2023/08/05 11:39:40 [172.17.0.3] upload /home/kehao/test.txt to /home/kehao.
2023/08/05 11:39:40 [172.17.0.3] upload /home/kehao/test.txt finished!
2023/08/05 11:39:40 [172.17.0.3] upload /home/kehao/test.sh to /home/kehao.
2023/08/05 11:39:41 [172.17.0.3] upload /home/kehao/test.sh finished!
2023/08/05 11:39:41 [172.17.0.4] upload /home/kehao/test.txt to /home/kehao.
2023/08/05 11:39:41 [172.17.0.4] upload /home/kehao/test.txt finished!
2023/08/05 11:39:41 [172.17.0.4] upload /home/kehao/test.sh to /home/kehao.
2023/08/05 11:39:41 [172.17.0.4] upload /home/kehao/test.sh finished!
2023/08/05 11:39:41 [172.17.0.5] upload /home/kehao/test.txt to /home/kehao.
2023/08/05 11:39:42 [172.17.0.5] upload /home/kehao/test.txt finished!
2023/08/05 11:39:42 [172.17.0.5] upload /home/kehao/test.sh to /home/kehao.
2023/08/05 11:39:42 [172.17.0.5] upload /home/kehao/test.sh finished!
2023/08/05 11:39:42 upload file finished.

kehao@kehaopcs ~/remote_exec $ ./remote execute
2023/08/05 11:39:45 start execute command...
2023/08/05 11:39:45 [172.17.0.3] execute (/bin/sh /home/kehao/test.sh).
2023/08/05 11:39:46 [172.17.0.3] execute (/bin/sh /home/kehao/test.sh) result: 
I am 172.17.0.3 :)
[root@5421d92524c8 ~]# 
2023/08/05 11:39:46 [172.17.0.3] execute command finished!
2023/08/05 11:39:46 [172.17.0.4] execute (/bin/sh /home/kehao/test.sh).
2023/08/05 11:39:48 [172.17.0.4] execute (/bin/sh /home/kehao/test.sh) result: 
I am 172.17.0.4 :)
[root@9c0bbfbc1446 ~]# 
2023/08/05 11:39:48 [172.17.0.4] execute command finished!
2023/08/05 11:39:48 [172.17.0.5] execute (/bin/sh /home/kehao/test.sh).
2023/08/05 11:39:49 [172.17.0.5] execute (/bin/sh /home/kehao/test.sh) result: 
I am 172.17.0.5 :)
[root@93d37bd5e970 ~]# 
2023/08/05 11:39:49 [172.17.0.5] execute command finished!
2023/08/05 11:39:49 execute command finished.

kehao@kehaopcs ~/remote_exec $ ./remote download
2023/08/05 11:40:09 start download file...
2023/08/05 11:40:09 [172.17.0.3] download /home/kehao/result.txt to /home/kehao/result.
2023/08/05 11:40:09 [172.17.0.3] download /home/kehao/result.txt finished!
2023/08/05 11:40:09 [172.17.0.4] download /home/kehao/result.txt to /home/kehao/result.
2023/08/05 11:40:10 [172.17.0.4] download /home/kehao/result.txt finished!
2023/08/05 11:40:10 [172.17.0.5] download /home/kehao/result.txt to /home/kehao/result.
2023/08/05 11:40:10 [172.17.0.5] download /home/kehao/result.txt finished!
2023/08/05 11:40:10 download file finished.
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
-rw-r--r--. 1 kehao manager 126 Aug  5 11:40 172.17.0.3_result.txt
-rw-r--r--. 1 kehao manager 126 Aug  5 11:40 172.17.0.4_result.txt
-rw-r--r--. 1 kehao manager 126 Aug  5 11:40 172.17.0.5_result.txt

kehao@kehaopcs ~/remote_exec $ cat /home/kehao/result/172.17.0.3_result.txt 
5421d92524c8
172.17.0.3
uid=0(root) gid=0(root) groups=0(root)
-rw-rw-r--. 1 kehao kehao 10 Aug  5 11:39 /home/kehao/test.txt

kehao@kehaopcs ~/remote_exec $ cat /home/kehao/result/172.17.0.4_result.txt 
9c0bbfbc1446
172.17.0.4
uid=0(root) gid=0(root) groups=0(root)
-rw-rw-r--. 1 kehao kehao 10 Aug  5 11:39 /home/kehao/test.txt

kehao@kehaopcs ~/remote_exec $ cat /home/kehao/result/172.17.0.5_result.txt 
93d37bd5e970
172.17.0.5
uid=0(root) gid=0(root) groups=0(root)
-rw-rw-r--. 1 kehao kehao 10 Aug  5 11:39 /home/kehao/test.txt
```