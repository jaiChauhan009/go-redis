package command

import (
	"fmt"
	"net"
	"redis-clone/internal/store"
	"strings"
	"time"
	"strconv"
	"redis-clone/internal/parser"
)

func Execute(cmd parser.Command, conn net.Conn) bool {
	if len(cmd.Args) == 0 {
		conn.Write([]byte("-ERR empty command\r\n"))
		return true
	}

	switch strings.ToUpper(cmd.Args[0]) {
	case "GET":
		handleGet(cmd, conn)
	case "SET":
		handleSet(cmd, conn)
	case "DEL":
		handleDel(cmd, conn)
	case "QUIT":
		conn.Write([]byte("+OK\r\n"))
		return false
	default:
		conn.Write([]byte(fmt.Sprintf("-ERR unknown command: %s\r\n", cmd.Args[0])))
	}
	return true
}

func handleGet(cmd parser.Command, conn net.Conn) {
	if len(cmd.Args) != 2 {
		conn.Write([]byte("-ERR wrong number of arguments for GET\r\n"))
		return
	}
	val, ok := store.Get(cmd.Args[1])
	if !ok {
		conn.Write([]byte("$-1\r\n"))
		return
	}
	conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(val), val)))
}

func handleSet(cmd parser.Command, conn net.Conn) {
	if len(cmd.Args) < 3 {
		conn.Write([]byte("-ERR wrong number of arguments for SET\r\n"))
		return
	}
	key := cmd.Args[1]
	value := cmd.Args[2]
	options := cmd.Args[3:]

	var ttl time.Duration = 0
	var nx, xx bool

	for i := 0; i < len(options); i++ {
		opt := strings.ToUpper(options[i])
		switch opt {
		case "EX":
			if i+1 < len(options) {
				sec, _ := strconv.Atoi(options[i+1])
				ttl = time.Duration(sec) * time.Second
				i++
			}
		case "PX":
			if i+1 < len(options) {
				ms, _ := strconv.Atoi(options[i+1])
				ttl = time.Duration(ms) * time.Millisecond
				i++
			}
		case "NX":
			nx = true
		case "XX":
			xx = true
		}
	}

	_, exists := store.Get(key)

	if (nx && exists) || (xx && !exists) {
		conn.Write([]byte("$-1\r\n"))
		return
	}

	store.Set(key, value, ttl)
	conn.Write([]byte("+OK\r\n"))
}

func handleDel(cmd parser.Command, conn net.Conn) {
	if len(cmd.Args) < 2 {
		conn.Write([]byte("-ERR wrong number of arguments for DEL\r\n"))
		return
	}
	count := 0
	for _, key := range cmd.Args[1:] {
		if store.Delete(key) {
			count++
		}
	}
	conn.Write([]byte(fmt.Sprintf(":%d\r\n", count)))
}
