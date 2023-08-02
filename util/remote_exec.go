package util

import (
	goexpect "github.com/google/goexpect"
	"github.com/google/goterm/term"
	"golang.org/x/crypto/ssh"
	"log"
	"regexp"
	"time"
)

func RemoteExec(host *Host, command string, rootPrompt, passwdPrompt *regexp.Regexp, timeout time.Duration) {
	log.Printf("[%s] execute (%v).\n", host.Host, command)
	var (
		sshClient *ssh.Client
		expect    *goexpect.GExpect
		result    string
		err       error
	)

	if sshClient, err = SshClient(host); err != nil {
		log.Println(term.Redf("[%s] ssh client connect error: %v", host.Host, err))
		return
	}
	defer sshClient.Close()

	if expect, _, err = goexpect.SpawnSSH(sshClient, Timeout); err != nil {
		log.Println(term.Redf("[%s] expect ssh error: %v", host.Host, err))
		return
	}
	defer expect.Close()

	expect.Send("su - root" + "\n")
	if result, _, err = expect.Expect(passwdPrompt, timeout); err != nil {
		log.Println(term.Redf("[%s] change user error: %v", host.Host, err))
		return
	}
	expect.Send(host.RootPwd + "\n")
	if result, _, err = expect.Expect(rootPrompt, timeout); err != nil {
		log.Println(term.Redf("[%s] send passwd error: %v", host.Host, err))
		return
	}
	expect.Send(command + "\n")
	if result, _, err = expect.Expect(rootPrompt, timeout); err != nil {
		log.Println(term.Redf("[%s] execute command error: %v", host.Host, err))
		return
	} else {
		log.Printf("[%s] execute (%v) result: \n%v\n", host.Host, command, result)
	}
	//expect.Send("exit\n")

	log.Println(term.Bluef("[%s] execute command finished!", host.Host))

}
