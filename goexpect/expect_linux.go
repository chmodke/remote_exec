// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package expect is a Go version of the classic TCL Expect.
package expect

import (
	"fmt"
	"os/exec"
	"remote_exec/goterm/term"
	"strings"
	"syscall"
	"time"
)

// SpawnWithArgs starts a new process and collects the output. The error
// channel returns the result of the command Spawned when it finishes.
// Arguments may contain spaces.
func SpawnWithArgs(command []string, timeout time.Duration, opts ...Option) (*GExpect, <-chan error, error) {
	pty, err := term.OpenPTY()
	if err != nil {
		return nil, nil, err
	}
	var t term.Termios
	t.Raw()
	t.Set(pty.Slave)

	if timeout < 1 {
		timeout = DefaultTimeout
	}
	// Get the command up and running
	cmd := exec.Command(command[0], command[1:]...)
	// This ties the commands Stdin,Stdout & Stderr to the virtual terminal we created
	cmd.Stdin, cmd.Stdout, cmd.Stderr = pty.Slave, pty.Slave, pty.Slave
	// New process needs to be the process leader and control of a tty
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid:  true,
		Setctty: true}
	e := &GExpect{
		rcv:         make(chan struct{}),
		snd:         make(chan string),
		cmd:         cmd,
		timeout:     timeout,
		chkDuration: checkDuration,
		pty:         pty,
		cls: func(e *GExpect) error {
			if e.cmd != nil {
				return e.cmd.Process.Kill()
			}
			return nil
		},
		chk: func(e *GExpect) bool {
			if e.cmd.Process == nil {
				return false
			}
			// Sending Signal 0 to a process returns nil if process can take a signal , something else if not.
			return e.cmd.Process.Signal(syscall.Signal(0)) == nil
		},
	}
	for _, o := range opts {
		o(e)
	}

	// Set the buffer size to the default if expect.BufferSize(...) is not utilized.
	if !e.bufferSizeIsSet {
		e.bufferSize = defaultBufferSize
	}

	res := make(chan error, 1)
	go e.runcmd(res)
	// Wait until command started
	return e, res, <-res
}

// Spawn starts a new process and collects the output. The error channel
// returns the result of the command Spawned when it finishes. Arguments may
// not contain spaces.
func Spawn(command string, timeout time.Duration, opts ...Option) (*GExpect, <-chan error, error) {
	return SpawnWithArgs(strings.Fields(command), timeout, opts...)
}

// String implements the stringer interface.
func (e *GExpect) String() string {
	res := fmt.Sprintf("%p: ", e)
	if e.pty != nil {
		_, name := e.pty.PTSName()
		res += fmt.Sprintf("pty: %s ", name)
	}
	switch {
	case e.cmd != nil:
		res += fmt.Sprintf("cmd: %s(%d) ", e.cmd.Path, e.cmd.Process.Pid)
	case e.ssh != nil:
		res += fmt.Sprint("ssh session ")
	}
	res += fmt.Sprintf("buf: %q", e.out.String())
	return res
}
