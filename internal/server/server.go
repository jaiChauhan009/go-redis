package server

import (
	"log"
	"net"
	"redis-clone/internal/parser"
	"redis-clone/internal/command"
)

func Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	log.Println("Listening on", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Accept error:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	p := parser.NewParser(conn)
	for {
		cmd, err := p.ReadCommand()
		if err != nil {
			conn.Write([]byte("-ERR " + err.Error() + "\r\n"))
			return
		}
		if !command.Execute(cmd, conn) {
			return
		}
	}
}
