package util

import (
	"github.com/google/goterm/term"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"path"
)

func RemotePut(host *Host, localPath, remoteDir string) {
	log.Printf("[%s] upload %s to %s.\n", host.Host, localPath, remoteDir)
	var (
		sshClient *ssh.Client
		ftpClient *sftp.Client
		err       error
	)
	if sshClient, err = SshClient(host); err != nil {
		log.Println(term.Redf("[%s] ssh client connect error: %v", host.Host, err))
		return
	}
	defer sshClient.Close()

	if ftpClient, err = SftpConnect(sshClient); err != nil {
		log.Println(term.Redf("[%s] sftp client connect error: %v", host.Host, err))
		return
	}
	defer ftpClient.Close()

	localFile, err := os.Open(localPath)
	if err != nil {
		log.Println(term.Redf("[%s] open %s error: %v", host.Host, localPath, err))
		return

	}
	defer localFile.Close()

	ftpClient.MkdirAll(remoteDir)
	fileName := path.Base(localPath)

	remoteFile, err := ftpClient.Create(path.Join(remoteDir, fileName))
	if err != nil {
		log.Println(term.Redf("[%s] sftp create %s error: %v", host.Host, remoteDir, err))
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
	log.Println(term.Bluef("[%s] upload %s finished!", host.Host, localPath))
}

func RemoteGet(host *Host, remotePath, localDir string) {
	log.Printf("[%s] download %s to %s.\n", host.Host, remotePath, localDir)
	var (
		sshClient *ssh.Client
		ftpClient *sftp.Client
		err       error
	)
	if sshClient, err = SshClient(host); err != nil {
		log.Println(term.Redf("[%s] ssh client connect error: %v", host.Host, err))
		return
	}
	defer sshClient.Close()

	if ftpClient, err = SftpConnect(sshClient); err != nil {
		log.Println(term.Redf("[%s] sftp client connect error: %v", host.Host, err))
		return
	}
	defer ftpClient.Close()

	CreateDir(localDir)
	fileName := path.Base(remotePath)

	localFile, err := os.Create(path.Join(localDir, host.Host+"_"+fileName))
	if err != nil {
		log.Println(term.Redf("[%s] create %s error: %v", host.Host, localDir, err))
		return

	}
	defer localFile.Close()

	remoteFile, err := ftpClient.Open(remotePath)
	if err != nil {
		log.Println(term.Redf("[%s] sftp open %s error: %v", host.Host, remotePath, err))
		return

	}
	defer remoteFile.Close()

	remoteFile.WriteTo(localFile)

	log.Println(term.Bluef("[%s] download %s finished!", host.Host, remotePath))
}
