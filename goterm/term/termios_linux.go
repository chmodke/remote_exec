// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package term implements a subset of the C termios library to interface with Terminals.

This package allows the caller to get and set most Terminal capabilites
and sizes as well as create PTYs to enable writing things like script,
screen, tmux, and expect.

The Termios type is used for setting/getting Terminal capabilities while
the PTY type is used for handling virtual terminals.

Currently this part of this lib is Linux specific.

Also implements a simple version of readline in pure Go and some Stringers
for terminal colors and attributes.
*/
package term

import (
	"errors"
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

// Set Sets terminal t attributes on file.
func (t *Termios) Set(file *os.File) error {
	fd := file.Fd()
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(TCSETS), uintptr(unsafe.Pointer(t)))
	if errno != 0 {
		return errno
	}
	return nil
}

// Attr Gets (terminal related) attributes from file.
func Attr(file *os.File) (Termios, error) {
	var t Termios
	fd := file.Fd()
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(TCGETS), uintptr(unsafe.Pointer(&t)))
	if errno != 0 {
		return t, errno
	}
	t.Ispeed &= CBAUD | CBAUDEX
	t.Ospeed &= CBAUD | CBAUDEX
	return t, nil
}

// Isatty returns true if file is a tty.
func Isatty(file *os.File) bool {
	_, err := Attr(file)
	return err == nil
}

// GetPass reads password from a TTY with no echo.
func GetPass(prompt string, f *os.File, pbuf []byte) ([]byte, error) {
	t, err := Attr(f)
	if err != nil {
		return nil, err
	}
	defer t.Set(f)
	noecho := t
	noecho.Lflag = noecho.Lflag &^ ECHO
	if err := noecho.Set(f); err != nil {
		return nil, err
	}
	b := make([]byte, 1, 1)
	i := 0
	if _, err := f.Write([]byte(prompt)); err != nil {
		return nil, err
	}
	for ; i < len(pbuf); i++ {
		if _, err := f.Read(b); err != nil {
			b[0] = 0
			clearbuf(pbuf[:i+1])
		}
		if b[0] == '\n' || b[0] == '\r' {
			return pbuf[:i], nil
		}
		pbuf[i] = b[0]
		b[0] = 0
	}
	clearbuf(pbuf[:i+1])
	return nil, errors.New("ran out of bufferspace")
}

// PTSNumber return the pty number.
func (p *PTY) PTSNumber() (uint, error) {
	var ptyno uint
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(p.Master.Fd()), uintptr(TIOCGPTN), uintptr(unsafe.Pointer(&ptyno)))
	if errno != 0 {
		return 0, errno
	}
	return ptyno, nil
}

// Winsz Fetches the current terminal windowsize.
// example handling changing window sizes with PTYs:
//
// import "os"
// import "os/signal"
//
// var sig = make(chan os.Signal,2) 		// Channel to listen for UNIX SIGNALS on
// signal.Notify(sig, syscall.SIGWINCH) // That'd be the window changing
//
//	for {
//		<-sig
//		term.Winsz(os.Stdin)			// We got signaled our terminal changed size so we read in the new value
//	 term.Setwinsz(pty.Slave) // Copy it to our virtual Terminal
//	}
func (t *Termios) Winsz(file *os.File) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(file.Fd()), uintptr(TIOCGWINSZ), uintptr(unsafe.Pointer(&t.Wz)))
	if errno != 0 {
		return errno
	}
	return nil
}

// Setwinsz Sets the terminal window size.
func (t *Termios) Setwinsz(file *os.File) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(file.Fd()), uintptr(TIOCSWINSZ), uintptr(unsafe.Pointer(&t.Wz)))
	if errno != 0 {
		return errno
	}
	return nil
}

// OpenPTY Creates a new Master/Slave PTY pair.
func OpenPTY() (*PTY, error) {
	// Opening ptmx gives you the FD of a brand new PTY
	master, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	// unlock pty slave
	var unlock int // 0 => Unlock
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(master.Fd()), uintptr(TIOCSPTLCK), uintptr(unsafe.Pointer(&unlock))); errno != 0 {
		master.Close()
		return nil, errno
	}

	// get path of pts slave
	pty := &PTY{Master: master}
	slaveStr, err := pty.PTSName()
	if err != nil {
		master.Close()
		return nil, err
	}

	// open pty slave
	pty.Slave, err = os.OpenFile(slaveStr, os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		master.Close()
		return nil, err
	}

	return pty, nil
}

// PTSName return the name of the pty.
func (p *PTY) PTSName() (string, error) {
	n, err := p.PTSNumber()
	if err != nil {
		return "", err
	}
	return "/dev/pts/" + strconv.Itoa(int(n)), nil
}
