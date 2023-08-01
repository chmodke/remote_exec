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
	var (
		sshClient *ssh.Client
		expect    *goexpect.GExpect
		result    string
		err       error
	)

	if sshClient, err = SshClient(user, passwd, host, port); err != nil {
		log.Printf("ssh.Dial(%q) failed: %v\n", host, err)
		return
	}
	defer sshClient.Close()

	if expect, _, err = goexpect.SpawnSSH(sshClient, Timeout); err != nil {
		log.Printf("SpawnSSH(%q) failed: %v\n", host, err)
		return
	}
	defer expect.Close()

	expect.Send("su - root" + "\n")
	if result, _, err = expect.Expect(passwdPrompt, time.Duration(timeout)*time.Second); err != nil {
		log.Printf("%s: error: %v\n", "input passwd", err)
		return
	}
	expect.Send(rootPwd + "\n")
	if result, _, err = expect.Expect(rootPrompt, time.Duration(timeout)*time.Second); err != nil {
		log.Printf("%s: error: %v\n", "send passwd", err)
		return
	}
	expect.Send(command + "\n")
	if result, _, err = expect.Expect(rootPrompt, time.Duration(timeout)*time.Second); err != nil {
		log.Printf("%s: error: %v\n", "exec script", err)
		return
	} else {
		log.Printf("%s, result: \n%v\n", command, result)
	}
	expect.Send("exit\n")

	log.Println(term.Greenf("Done!"))

}
