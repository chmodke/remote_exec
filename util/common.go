package util

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/pkg/sftp"
	"github.com/spf13/cast"
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

// ParseHosts parse host info from config.yaml
func ParseHosts() ([]*Host, error) {
	var (
		config   *viper.Viper
		err      error
		hosts    []*Host
		port     = 22
		user     string
		passwd   string
		rootPwd  string
		hostList []string
	)

	if config, err = LoadCfg("config"); err != nil {
		return hosts, errors.New("load config.yaml failed")
	}

	if config.IsSet("port") {
		port = config.GetInt("port")
	}
	user = config.GetString("user")
	passwd = config.GetString("passwd")
	rootPwd = config.GetString("rootPwd")
	hostList = config.GetStringSlice("hosts")

	for _, host := range hostList {
		h := &Host{Port: port, Host: host, User: user, Passwd: passwd, RootPwd: rootPwd}
		hosts = append(hosts, h)
	}

	if config.IsSet("spc_hosts") {
		spcHosts := config.GetStringSlice("spc_hosts")
		for _, spcHost := range spcHosts {
			params := strings.Split(spcHost, " ")
			h := &Host{User: user, Host: params[0], Port: cast.ToInt(params[1]), Passwd: params[2], RootPwd: params[3]}
			hosts = append(hosts, h)
		}
	}
	return hosts, nil
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
