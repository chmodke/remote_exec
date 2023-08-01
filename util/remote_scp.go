package util

import (
	"github.com/google/goterm/term"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"path"
)

func RemotePut(port int, host, user, passwd, localPath, remotePath string) {
	var (
		sshClient *ssh.Client
		ftpClient *sftp.Client
		err       error
	)
	if sshClient, err = SshClient(user, passwd, host, port); err != nil {
		log.Printf("ssh.Dial(%q) failed: %v\n", host, err)
		return
	}
	defer sshClient.Close()

	if ftpClient, err = SftpConnect(sshClient); err != nil {
		log.Printf("SftpConnect(%q) failed: %v\n", host, err)
		return
	}
	defer ftpClient.Close()

	localFile, err := os.Open(localPath)
	if err != nil {
		log.Println("os.Open error : ", localPath)
		return

	}
	defer localFile.Close()

	err = ftpClient.Remove(remotePath)

	base := path.Dir(remotePath)
	ftpClient.MkdirAll(base)

	remoteFile, err := ftpClient.Create(remotePath)
	if err != nil {
		log.Println("sftpClient.Create error : ", remotePath)
		return

	}
	defer remoteFile.Close()

	var buf = make([]byte, 1024)
	for {
		var len = 0
		len, err = localFile.Read(buf)
		if err != nil {
			break
		} else if len == 0 {
			break
		} else {
			remoteFile.Write(buf[:len])
		}
	}
	log.Printf("copy %v to %v finished!\n", localPath, host)
	log.Println(term.Bluef("Write file to remote"))
}

func RemoteGet(port int, host, user, passwd, remotePath, localPath string) {
	var (
		sshClient *ssh.Client
		ftpClient *sftp.Client
		err       error
	)
	if sshClient, err = SshClient(user, passwd, host, port); err != nil {
		log.Printf("ssh.Dial(%q) failed: %v\n", host, err)
		return
	}
	defer sshClient.Close()

	if ftpClient, err = SftpConnect(sshClient); err != nil {
		log.Printf("SftpConnect(%q) failed: %v\n", host, err)
		return
	}
	defer ftpClient.Close()

	base := path.Dir(localPath)
	fileName := path.Base(localPath)

	CreateDir(base)

	localFile, err := os.Create(path.Join(base, host+"_"+fileName))
	if err != nil {
		log.Println("os.Create error : ", localPath)
		return

	}
	defer localFile.Close()

	remoteFile, err := ftpClient.Open(remotePath)
	if err != nil {
		log.Println("sftpClient.Create error : ", remotePath)
		return

	}
	defer remoteFile.Close()

	remoteFile.WriteTo(localFile)

	log.Printf("copy %v from %s finished!\n", remotePath, host)
	log.Println(term.Bluef("Get file from remote"))
}
