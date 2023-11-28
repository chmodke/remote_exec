package util

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func InitLog() func() {
	date := time.Now().Format("20060102")
	log.SetFlags(log.Ldate | log.Ltime)
	f, err := os.OpenFile(fmt.Sprintf("remote_%s.log", date), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0640)
	if err != nil {
		return nil
	}

	multiWriter := io.MultiWriter(os.Stdout, f)
	log.SetOutput(multiWriter)
	return func() {
		f.Close()
	}
}
