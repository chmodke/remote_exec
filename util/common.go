package util

import (
	"bufio"
	"fmt"
	"github.com/pkg/sftp"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Host struct {
	Port    int `default:"22"`
	Host    string
	User    string
	Passwd  string
	RootPwd string
}

const (
	Timeout = 10 * time.Second
)

func SshClient(host *Host) (*ssh.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(host.Passwd))

	// get host public key
	//hostKey := getHostKey(host.Host)

	clientConfig = &ssh.ClientConfig{
		User:    host.User,
		Auth:    auth,
		Timeout: 30 * time.Second,
		// allow any host key to be used (non-prod)
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),

		// verify host public key
		//HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	// connect to ssh
	if strings.Contains(host.Host, ":") {
		addr = fmt.Sprintf("[%s]:%d", host.Host, host.Port)
	} else {
		addr = fmt.Sprintf("%s:%d", host.Host, host.Port)
	}
	return ssh.Dial("tcp", addr, clientConfig)
}

func SftpConnect(sshClient *ssh.Client) (*sftp.Client, error) {
	return sftp.NewClient(sshClient)
}

func getHostKey(host string) ssh.PublicKey {
	// parse OpenSSH known_hosts file
	// ssh or use ssh-keyscan to pull key
	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], host) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				log.Fatalf("error parsing %q: %v", fields[2], err)
			}
			break
		}
	}

	if hostKey == nil {
		log.Fatalf("no hostkey found for %s", host)
	}

	return hostKey
}

func LoadCfg(configName string) (*viper.Viper, error) {
	var v = viper.New()
	v.SetConfigName(configName)
	v.SetConfigType("yaml")
	v.AddConfigPath("./")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file %s.%s not found", configName, "yaml")
		} else {
			return nil, err
		}
	}

	return v, nil
}

// CreateDir 创建文件夹
func CreateDir(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		// 创建文件夹
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Printf("mkdir failed![%v]\n", err)
		} else {
			return true, nil
		}
	}
	return false, err
}
