docker build -t centos:sshd .

docker run -itd --name test-ssh1 centos:sshd
docker run -itd --name test-ssh2 centos:sshd
docker run -itd --name test-ssh3 centos:sshd

docker exec -u kehao test-ssh1 bash -c 'ls -l ~'
docker exec -u kehao test-ssh2 bash -c 'ls -l ~'
docker exec -u kehao test-ssh3 bash -c 'ls -l ~'

docker exec -u kehao test-ssh1 bash -c 'rm ~/*'
docker exec -u kehao test-ssh2 bash -c 'rm ~/*'
docker exec -u kehao test-ssh3 bash -c 'rm ~/*'
