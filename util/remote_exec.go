package util

import (
	goexpect "github.com/google/goexpect"
	"github.com/google/goterm/term"
	"golang.org/x/crypto/ssh"
	"log"
	"regexp"
	"time"
)

func RemoteExec(host *Host, commands []string, rootPrompt, passwdPrompt *regexp.Regexp, timeout time.Duration) bool {
	var (
		sshClient *ssh.Client
		expect    *goexpect.GExpect
		result    string
		err       error
	)

	if sshClient, err = SshClient(host); err != nil {
		log.Println(term.Redf("[%s:%v] ssh client connect error: %v", host.Host, host.Port, err))
		return false
	}
	defer sshClient.Close()

	if expect, _, err = goexpect.SpawnSSH(sshClient, Timeout); err != nil {
		log.Println(term.Redf("[%s:%v] expect ssh error: %v", host.Host, host.Port, err))
		return false
	}
	defer expect.Close()

	expect.Send("su - root" + "\n")
	if result, _, err = expect.Expect(passwdPrompt, timeout); err != nil {
		log.Println(term.Yellowf("[%s:%v] change user, expected: (%v), actual: (%v)", host.Host, host.Port, passwdPrompt, result))
		log.Println(term.Redf("[%s:%v] change user error: %v", host.Host, host.Port, err))
		return false
	}
	expect.Send(host.RootPwd + "\n")
	if result, _, err = expect.Expect(rootPrompt, timeout); err != nil {
		log.Println(term.Yellowf("[%s:%v] send passwd, expected: (%v), actual: (%v)", host.Host, host.Port, rootPrompt, result))
		log.Println(term.Redf("[%s:%v] send passwd error: %v", host.Host, host.Port, err))
		return false
	}
	for _, command := range commands {
		log.Printf("[%s:%v] execute (%v).\n", host.Host, host.Port, command)
		expect.Send(command + "\n")
		if result, _, err = expect.Expect(rootPrompt, timeout); err != nil {
			log.Println(term.Yellowf("[%s:%v] execute command, expected: (%v), actual: (%v)", host.Host, host.Port, rootPrompt, result))
			log.Println(term.Redf("[%s:%v] execute command error: %v", host.Host, host.Port, err))
			return false
		} else {
			log.Printf("[%s:%v] execute (%v) result: \n%v\n", host.Host, host.Port, command, result)
		}
	}
	//expect.Send("exit\n")

	log.Println(term.Bluef("[%s:%v] execute command finished!", host.Host, host.Port))
	return true
}
