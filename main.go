package main

import (
	"flag"
	"fmt"
	"github.com/anacrolix/utp"
	"io"
	"os"
	"os/signal"

	rand "github.com/dustin/randbo"
)

func beNetcat(con io.ReadWriteCloser, in io.Reader) {
	defer con.Close()

	go func() {
		_, err := io.Copy(os.Stdout, con)
		if err != io.EOF {
			fmt.Fprintf(os.Stderr, "Read error: %s\n", err)
			os.Exit(1)
		}
	}()
	_, err := io.Copy(con, in)
	if err != io.EOF {
		fmt.Fprintf(os.Stderr, "Write error: %s\n", err)
		os.Exit(1)
	}
}

func main() {
	list := flag.Bool("l", false, "listen on the given address")
	spew := flag.Bool("spew", false, "spew random data on the connection")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-l] <host> <port>\n", os.Args[0])
	}

	flag.Parse()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	if len(flag.Args()) < 2 {
		flag.Usage()
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

		defer sock.Close()

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

	var in io.Reader = os.Stdin
	if *spew {
		in = rand.New()
	}

	go func() {
		<-c
		con.Close()
	}()
	beNetcat(con, in)
}
