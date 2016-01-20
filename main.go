package main

import (
	"flag"
	"fmt"
	"github.com/anacrolix/utp"
	"io"
	"os"
)

func beNetcat(con io.ReadWriteCloser) {
	defer con.Close()

	go io.Copy(os.Stdout, con)
	io.Copy(con, os.Stdin)
}

func main() {
	list := flag.Bool("l", false, "listen on the given address")

	flag.Parse()

	if len(flag.Args()) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-l] <host> <port>\n", os.Args[0])
		os.Exit(1)
	}

	addr := fmt.Sprintf("%s:%s", flag.Arg(0), flag.Arg(1))

	var con io.ReadWriteCloser
	if *list {
		sock, err := utp.NewSocket("udp", addr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "create socket failed: %s\n", err)
			os.Exit(1)
		}

		utpcon, err := sock.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "accept failed: %s\n", err)
			os.Exit(1)
		}

		con = utpcon
	} else {
		utpcon, err := utp.Dial(addr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "dial failed: %s\n", err)
			os.Exit(1)
		}

		con = utpcon
	}

	beNetcat(con)
}
