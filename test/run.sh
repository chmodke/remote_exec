docker build -t centos:sshd .

docker network rm test-network
docker network create --subnet=172.18.0.0/16 test-network

docker rm -f test-ssh1
docker rm -f test-ssh2
docker rm -f test-ssh3

docker run -itd --network test-network --ip 172.18.0.3 --name test-ssh1 centos:sshd
docker run -itd --network test-network --ip 172.18.0.4 --name test-ssh2 centos:sshd
docker run -itd --network test-network --ip 172.18.0.5 -p 23:22 --name test-ssh3 centos:sshd

docker exec -u kehao test-ssh1 bash -c 'ls -l ~'
docker exec -u kehao test-ssh2 bash -c 'ls -l ~'
docker exec -u kehao test-ssh3 bash -c 'ls -l ~'

docker exec -u kehao test-ssh1 bash -c 'rm ~/*'
docker exec -u kehao test-ssh2 bash -c 'rm ~/*'
docker exec -u kehao test-ssh3 bash -c 'rm ~/*'
