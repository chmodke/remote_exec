package util

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/pkg/sftp"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"os"
	"path/filepath"
	"remote_exec/goterm/term"
	"strings"
	"sync"
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

func LoadCfg(configPath string, defaultConfig string) (*viper.Viper, error) {
	var (
		v       = viper.New()
		cfgPath string
	)
	if r, _ := FileExists(nil, configPath); r {
		cfgPath = configPath
	} else {
		log.Println(term.Yellowf("config file %s not found, will use %s.", configPath, defaultConfig))
		cfgPath = defaultConfig
	}

	v.SetConfigFile(cfgPath)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file %s not found", cfgPath)
		} else {
			return nil, err
		}
	}

	return v, nil
}

// ParseHosts parse host info from config.yaml
func ParseHosts(configPath string, cmd *cobra.Command) ([]*Host, error) {
	var (
		config     *viper.Viper
		err        error
		hosts      []*Host
		port       = 22
		user       string
		passwd     string
		rootPwd    string
		hostList   []string
		netMask, _ = cmd.Flags().GetString(ConstNetMask)
	)

	if config, err = LoadCfg(configPath, DefaultConfig); err != nil {
		return hosts, errors.New(err.Error())
	}

	if config.IsSet("port") {
		port = config.GetInt("port")
	}
	user = config.GetString("user")
	passwd = config.GetString("passwd")
	rootPwd = config.GetString("root-passwd")
	hostList = config.GetStringSlice("hosts")

	for _, host := range hostList {
		if len(netMask) == 0 || IpAllow(host, netMask) {
			h := &Host{Port: port, Host: host, User: user, Passwd: passwd, RootPwd: rootPwd}
			hosts = append(hosts, h)
		}
	}

	if config.IsSet("spc-hosts") {
		spcHosts := config.GetStringSlice("spc-hosts")
		for _, spcHost := range spcHosts {
			params := strings.Split(spcHost, " ")
			if len(netMask) == 0 || IpAllow(params[0], netMask) {
				h := &Host{User: user, Host: params[0], Port: cast.ToInt(params[1]), Passwd: params[2], RootPwd: params[3]}
				hosts = append(hosts, h)
			}
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

// DirExists 文件夹是否存在
func DirExists(ftpClient *sftp.Client, path string) (bool, error) {
	var (
		stat os.FileInfo
		err  error
	)
	if ftpClient != nil {
		stat, err = ftpClient.Stat(path)
	} else {
		stat, err = os.Stat(path)
	}
	if err == nil {
		return stat.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// FileExists 文件是否存在
func FileExists(ftpClient *sftp.Client, path string) (bool, error) {
	var (
		stat os.FileInfo
		err  error
	)
	if ftpClient != nil {
		stat, err = ftpClient.Stat(path)
	} else {
		stat, err = os.Stat(path)
	}
	if err == nil {
		return !stat.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Process is concurrent controller
func Process(taskLimit int, hosts []*Host, exec func(*Host)) {
	var i, j int
	var jobGroup sync.WaitGroup
	hostChan := make(chan *Host, len(hosts))

	for j = 0; j < len(hosts); j++ {
		hostChan <- hosts[j]
	}
	close(hostChan)

	jobGroup.Add(len(hosts))

	for i = 0; i < Max(Min(taskLimit, 16), 1); i++ {
		go func(taskChan chan *Host) {
			for {
				if host, ok := <-taskChan; ok {
					log.Printf("progress [%v/%v]...\n", len(taskChan), cap(taskChan))
					exec(host)
					jobGroup.Done()
				} else {
					break
				}
			}
		}(hostChan)
	}
	jobGroup.Wait()
}

func Min(x, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}

func Max(x, y int) int {
	if x > y {
		return x
	} else {
		return y
	}
}

// IpAllow judge the ip address matches the mask
// ip 192.168.1.1
// netmask "192.168.0.0/24" "192.168.1.1" "192.168.1.1,192.168.1.2"
func IpAllow(ip string, netmask string) bool {
	if ip == netmask {
		return true
	}
	for _, sub := range strings.Split(netmask, ",") {
		if sub == ip {
			return true
		}
	}
	if !strings.Contains(netmask, "/") {
		return false
	}
	ips := net.ParseIP(ip)
	if _, ipNet, err := net.ParseCIDR(netmask); err != nil {
		return true
	} else {
		return ipNet.Contains(ips)
	}
}
