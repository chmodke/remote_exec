package util

import (
	"github.com/google/goterm/term"
	"github.com/pkg/sftp"
	"github.com/spf13/cast"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"path"
	"strings"
)

func RemotePut(host *Host, files []string) bool {
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

	for _, file := range files {
		params := strings.Split(file, "#")
		var (
			localPath string
			remoteDir string
		)
		localPath = params[0]
		remoteDir = path.Dir(params[0])
		if len(params) == 2 {
			remoteDir = params[1]
		}
		put(host, ftpClient, localPath, remoteDir)
	}

	return true
}

func put(host *Host, ftpClient *sftp.Client, localPath, remoteDir string) bool {
	log.Printf("[%s:%v] upload %s to %s.\n", host.Host, host.Port, localPath, remoteDir)
	localFile, err := os.Open(localPath)
	if err != nil {
		log.Println(term.Redf("[%s:%v] open %s error: %v", host.Host, host.Port, localPath, err))
		return false
	}
	defer localFile.Close()

	if exists, _ := DirExists(ftpClient, remoteDir); !exists {
		ftpClient.MkdirAll(remoteDir)
	}

	fileName := path.Base(localPath)
	remotePath := path.Join(remoteDir, fileName)
	if exists, _ := FileExists(ftpClient, remotePath); exists {
		ftpClient.Remove(remotePath)
	}

	remoteFile, err := ftpClient.Create(remotePath)
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

func RemoteGet(host *Host, files []string) bool {
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

	for _, file := range files {
		params := strings.Split(file, "#")
		var (
			remotePath string
			localDir   string
		)
		remotePath = params[0]
		localDir = path.Dir(params[0])
		if len(params) == 2 {
			localDir = params[1]
		}
		get(host, ftpClient, remotePath, localDir)
	}
	return true
}

func get(host *Host, ftpClient *sftp.Client, remotePath, localDir string) bool {
	log.Printf("[%s:%v] download file from %s to %s.\n", host.Host, host.Port, remotePath, localDir)
	if exists, _ := DirExists(nil, localDir); !exists {
		CreateDir(localDir)
	}

	fileName := path.Base(remotePath)
	localPath := path.Join(localDir, host.Host+"_"+cast.ToString(host.Port)+"_"+fileName)
	if exists, _ := FileExists(nil, localPath); exists {
		os.Remove(localPath)
	}
	localFile, err := os.Create(localPath)
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
