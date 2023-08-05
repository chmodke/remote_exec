FROM chmodke.org:8082/centos:centos7.9.2009

ENV DEBIAN_FRONTEND noninteractive
ENV TZ Asia/Shanghai

ARG V_USER="kehao"

USER root

RUN ln -fs /usr/share/zoneinfo/${TZ} /etc/localtime && echo ${TZ} > /etc/timezone \
    && curl http://mirrors.aliyun.com/repo/Centos-7.repo -o /etc/yum.repos.d/CentOS-Base.repo \
    && yum clean all && yum makecache \
    && yum install -y openssh-server  \
    && ssh-keygen -t rsa -f /etc/ssh/ssh_host_rsa_key \
    && ssh-keygen -t rsa -f /etc/ssh/ssh_host_ecdsa_key \
    && ssh-keygen -t rsa -f /etc/ssh/ssh_host_ed25519_key \
    && mkdir -p /var/run/sshd \
    && mkdir -p /root/.ssh \
    && sed -ri 's/session    required     pam_loginuid.so/#session    required     pam_loginuid.so/g' /etc/pam.d/sshd \
    && useradd -d /home/${V_USER}  -s /bin/bash -m ${V_USER} \
    && echo "${V_USER}:123456" | chpasswd \
    && echo 'root:123456' | chpasswd

EXPOSE 22

ENTRYPOINT ["/usr/sbin/sshd"]
CMD ["-D"]