package lableuchtung

import (
	"errors"
	"fmt"
	"io"
	"time"
)

type Package []byte

func (p Package) checkSum() byte {
	return p[0] ^ p[1] ^ p[2] ^ p[3]
}

var (
	EnableAutomode = Package{251, 0, 0, 0}
	EnableBeacon   = Package{252, 0, 0, 0}
	DisableBeacon  = Package{253, 0, 0, 0}
)

type LabLeucht struct {
	io.ReadWriteCloser
	ResponseTimeout time.Duration
}

func (l *LabLeucht) SendPackage(pkg Package) error {

	//fmt.Printf("Sending package [%3d %3d %3d %3d]\n", pkg[0], pkg[1], pkg[2], pkg[3])

	var timeoutMult time.Duration

	if pkg[0] < 251 {
		timeoutMult = time.Duration(pkg[0]) * 100
	}

	//fmt.Println("timeout is", timeoutMult*time.Millisecond+l.ResponseTimeout)

	timeout := time.After(timeoutMult*time.Millisecond + l.ResponseTimeout)

	_, err := l.Write([]byte(pkg))
	if err != nil {
		return err
	}

	rspCh := make(chan byte)
	errCh := make(chan error)

	go func() {
		buf := make([]byte, 1)
		for {
			n, err := l.Read(buf)
			if err != nil && err != io.EOF {
				errCh <- err
				return
			}
			if n != 0 {
				rspCh <- buf[0]
				return
			}
		}
	}()

	select {
	case err := <-errCh:
		return err

	case r := <-rspCh:
		if r != pkg.checkSum() {
			return errors.New("bad response: got " + fmt.Sprint(r) + ", want " + fmt.Sprint(pkg.checkSum()))
		}

	case <-timeout:
		return errors.New("response timed out")
	}

	return nil
}
