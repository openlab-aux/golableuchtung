package main

import (
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/openlab-aux/golableuchtung"
)

func main() {

	if len(os.Args) != 5 {
		log.Fatal("invalid number of arguments")
	}

	leucht := lableuchtung.LabLeucht{}

	serverAddr, err := net.ResolveUDPAddr("udp", "10.11.7.3:1337")
	if err != nil {
		log.Fatal(err)
	}

	leucht.ReadWriteCloser, err = net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		log.Fatal(err)
	}

	defer leucht.Close()

	leucht.ResponseTimeout = 100 * time.Millisecond

	pkg := lableuchtung.Package{0, 0, 0, 0}

	for i := 0; i < 4; i++ {
		v, err := strconv.ParseUint(os.Args[i+1], 10, 8)
		if err != nil {
			log.Fatal(err)
		}
		pkg[i] = byte(v)
	}

	err = leucht.SendPackage(pkg)
	if err != nil {
		log.Fatal(err)
	}

}
