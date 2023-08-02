package util

import (
	goexpect "github.com/google/goexpect"
	"github.com/google/goterm/term"
	"golang.org/x/crypto/ssh"
	"log"
	"regexp"
	"time"
)

func RemoteExec(port int, host, user, passwd, rootPwd, command string, rootPrompt, passwdPrompt *regexp.Regexp, timeout int64) {
	log.Printf("[%s] execute (%v).\n", host, command)
	var (
		sshClient *ssh.Client
		expect    *goexpect.GExpect
		result    string
		err       error
	)

	if sshClient, err = SshClient(user, passwd, host, port); err != nil {
		log.Println(term.Redf("[%s] ssh client connect error: %v", host, err))
		return
	}
	defer sshClient.Close()

	if expect, _, err = goexpect.SpawnSSH(sshClient, Timeout); err != nil {
		log.Println(term.Redf("[%s] expect ssh error: %v", host, err))
		return
	}
	defer expect.Close()

	expect.Send("su - root" + "\n")
	if result, _, err = expect.Expect(passwdPrompt, time.Duration(timeout)*time.Second); err != nil {
		log.Println(term.Redf("[%s] change user error: %v", host, err))
		return
	}
	expect.Send(rootPwd + "\n")
	if result, _, err = expect.Expect(rootPrompt, time.Duration(timeout)*time.Second); err != nil {
		log.Println(term.Redf("[%s] send passwd error: %v", host, err))
		return
	}
	expect.Send(command + "\n")
	if result, _, err = expect.Expect(rootPrompt, time.Duration(timeout)*time.Second); err != nil {
		log.Println(term.Redf("[%s] execute command error: %v", host, err))
		return
	} else {
		log.Printf("[%s] execute (%v) result: \n%v\n", host, command, result)
	}
	//expect.Send("exit\n")

	log.Println(term.Bluef("[%s] execute command finished!", host))

}
