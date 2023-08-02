package util

import (
	"github.com/google/goterm/term"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"path"
)

func RemotePut(port int, host, user, passwd, localPath, remoteDir string) {
	log.Printf("[%s] upload %s to %s.\n", host, localPath, remoteDir)
	var (
		sshClient *ssh.Client
		ftpClient *sftp.Client
		err       error
	)
	if sshClient, err = SshClient(user, passwd, host, port); err != nil {
		log.Println(term.Redf("[%s] ssh client connect error: %v", host, err))
		return
	}
	defer sshClient.Close()

	if ftpClient, err = SftpConnect(sshClient); err != nil {
		log.Println(term.Redf("[%s] sftp client connect error: %v", host, err))
		return
	}
	defer ftpClient.Close()

	localFile, err := os.Open(localPath)
	if err != nil {
		log.Println(term.Redf("[%s] open %s error: %v", host, localPath, err))
		return

	}
	defer localFile.Close()

	ftpClient.MkdirAll(remoteDir)
	fileName := path.Base(localPath)

	remoteFile, err := ftpClient.Create(path.Join(remoteDir, fileName))
	if err != nil {
		log.Println(term.Redf("[%s] sftp create %s error: %v", host, remoteDir, err))
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
	log.Println(term.Bluef("[%s] copy %s finished!", host, localPath))
}

func RemoteGet(port int, host, user, passwd, remotePath, localDir string) {
	log.Printf("[%s] download %s to %s.\n", host, remotePath, localDir)
	var (
		sshClient *ssh.Client
		ftpClient *sftp.Client
		err       error
	)
	if sshClient, err = SshClient(user, passwd, host, port); err != nil {
		log.Println(term.Redf("[%s] ssh client connect error: %v", host, err))
		return
	}
	defer sshClient.Close()

	if ftpClient, err = SftpConnect(sshClient); err != nil {
		log.Println(term.Redf("[%s] sftp client connect error: %v", host, err))
		return
	}
	defer ftpClient.Close()

	CreateDir(localDir)
	fileName := path.Base(remotePath)

	localFile, err := os.Create(path.Join(localDir, host+"_"+fileName))
	if err != nil {
		log.Println(term.Redf("[%s] create %s error: %v", host, localDir, err))
		return

	}
	defer localFile.Close()

	remoteFile, err := ftpClient.Open(remotePath)
	if err != nil {
		log.Println(term.Redf("[%s] sftp open %s error: %v", host, remotePath, err))
		return

	}
	defer remoteFile.Close()

	remoteFile.WriteTo(localFile)

	log.Println(term.Bluef("[%s] copy %s finished!", host, remotePath))
}
