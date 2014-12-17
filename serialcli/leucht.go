package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/openlab-aux/golableuchtung"

	serial "github.com/huin/goserial"
)

func main() {

	if len(os.Args) != 5 {
		log.Fatal("invalid number of arguments")
	}

	leucht := lableuchtung.LabLeucht{}

	c := &serial.Config{
		Name: "/dev/ttyACM0",
		Baud: 115200,
	}

	var err error

	leucht.ReadWriteCloser, err = serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	leucht.ResponseTimeout = 10 * time.Millisecond
	ioutil.ReadAll(leucht) // consume fnord

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
