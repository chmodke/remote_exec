docker build -t centos:sshd .

docker network rm test-network
docker network create --subnet=172.18.0.0/16 test-network

docker rm -f test-ssh1 test-ssh2 test-ssh3 test-ssh4 test-ssh5 test-ssh6 test-ssh7 test-ssh8

docker run -itd --network test-network --ip 172.18.0.3 --name test-ssh1 centos:sshd
docker run -itd --network test-network --ip 172.18.0.4 --name test-ssh2 centos:sshd
docker run -itd --network test-network --ip 172.18.0.5 -p 23:22 --name test-ssh3 centos:sshd
docker run -itd --network test-network --ip 172.18.0.6 --name test-ssh4 centos:sshd
docker run -itd --network test-network --ip 172.18.0.7 --name test-ssh5 centos:sshd
docker run -itd --network test-network --ip 172.18.0.8 --name test-ssh6 centos:sshd
docker run -itd --network test-network --ip 172.18.0.9 --name test-ssh7 centos:sshd
docker run -itd --network test-network --ip 172.18.0.10 --name test-ssh8 centos:sshd

docker exec -u kehao test-ssh1 bash -c 'ls -l ~'
docker exec -u kehao test-ssh2 bash -c 'ls -l ~'
docker exec -u kehao test-ssh3 bash -c 'ls -l ~'

docker exec -u kehao test-ssh1 bash -c 'rm -rf ~/*'
docker exec -u kehao test-ssh2 bash -c 'rm -rf ~/*'
docker exec -u kehao test-ssh3 bash -c 'rm -rf ~/*'
docker exec -u kehao test-ssh4 bash -c 'rm -rf ~/*'
docker exec -u kehao test-ssh5 bash -c 'rm -rf ~/*'
docker exec -u kehao test-ssh6 bash -c 'rm -rf ~/*'
docker exec -u kehao test-ssh7 bash -c 'rm -rf ~/*'
docker exec -u kehao test-ssh8 bash -c 'rm -rf ~/*'

docker start test-ssh1 test-ssh2 test-ssh3 test-ssh4 test-ssh5 test-ssh6 test-ssh7 test-ssh8
docker stop test-ssh1 test-ssh2 test-ssh3 test-ssh4 test-ssh5 test-ssh6 test-ssh7 test-ssh8
