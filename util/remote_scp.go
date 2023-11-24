package util

import (
	"github.com/google/goterm/term"
	"github.com/pkg/sftp"
	"github.com/spf13/cast"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"path"
)

func RemotePut(host *Host, localPath, remoteDir string) bool {
	log.Printf("[%s:%v] upload %s to %s.\n", host.Host, host.Port, localPath, remoteDir)
	var (
		sshClient *ssh.Client
		ftpClient *sftp.Client
		err       error
	)
	if sshClient, err = SshClient(host); err != nil {
		log.Println(term.Redf("[%s:%v] ssh client connect error: %v", host.Host, host.Port, err))
		return false
	}
	defer sshClient.Close()

	if ftpClient, err = SftpConnect(sshClient); err != nil {
		log.Println(term.Redf("[%s:%v] sftp client connect error: %v", host.Host, host.Port, err))
		return false
	}
	defer ftpClient.Close()

	localFile, err := os.Open(localPath)
	if err != nil {
		log.Println(term.Redf("[%s:%v] open %s error: %v", host.Host, host.Port, localPath, err))
		return false
	}
	defer localFile.Close()

	ftpClient.MkdirAll(remoteDir)
	fileName := path.Base(localPath)

	remoteFile, err := ftpClient.Create(path.Join(remoteDir, fileName))
	if err != nil {
		log.Println(term.Redf("[%s:%v] sftp create %s error: %v", host.Host, host.Port, remoteDir, err))
		return false
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
	log.Println(term.Bluef("[%s:%v] upload %s finished!", host.Host, host.Port, localPath))
	return true
}

func RemoteGet(host *Host, remotePath, localDir string) bool {
	log.Printf("[%s:%v] download %s to %s.\n", host.Host, host.Port, remotePath, localDir)
	var (
		sshClient *ssh.Client
		ftpClient *sftp.Client
		err       error
	)
	if sshClient, err = SshClient(host); err != nil {
		log.Println(term.Redf("[%s:%v] ssh client connect error: %v", host.Host, host.Port, err))
		return false
	}
	defer sshClient.Close()

	if ftpClient, err = SftpConnect(sshClient); err != nil {
		log.Println(term.Redf("[%s:%v] sftp client connect error: %v", host.Host, host.Port, err))
		return false
	}
	defer ftpClient.Close()

	CreateDir(localDir)
	fileName := path.Base(remotePath)

	localFile, err := os.Create(path.Join(localDir, host.Host+"_"+cast.ToString(host.Port)+"_"+fileName))
	if err != nil {
		log.Println(term.Redf("[%s:%v] create %s error: %v", host.Host, host.Port, localDir, err))
		return false
	}
	defer localFile.Close()

	remoteFile, err := ftpClient.Open(remotePath)
	if err != nil {
		log.Println(term.Redf("[%s:%v] sftp open %s error: %v", host.Host, host.Port, remotePath, err))
		return false
	}
	defer remoteFile.Close()

	remoteFile.WriteTo(localFile)

	log.Println(term.Bluef("[%s:%v] download %s finished!", host.Host, host.Port, remotePath))
	return true
}
