package parser

import (
	"bufio"
	"errors"
	"net"
	"strconv"
)

type Command struct {
	Args []string
}

type Parser struct {
	r *bufio.Reader
}

func NewParser(conn net.Conn) *Parser {
	return &Parser{r: bufio.NewReader(conn)}
}

func (p *Parser) ReadCommand() (Command, error) {
	prefix, err := p.r.ReadByte()
	if err != nil {
		return Command{}, err
	}

	if prefix == '*' {
		return p.readRESP()
	}

	line, err := p.readLine(prefix)
	if err != nil {
		return Command{}, err
	}
	return parseInline(line), nil
}

func (p *Parser) readLine(first byte) ([]byte, error) {
	buf, err := p.r.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	line := append([]byte{first}, buf...)
	return line, nil
}

func (p *Parser) readRESP() (Command, error) {
	line, err := p.r.ReadBytes('\r')
	if err != nil {
		return Command{}, err
	}
	p.r.ReadByte() // consume \n

	count, err := strconv.Atoi(string(line[:len(line)-1]))
	if err != nil || count <= 0 {
		return Command{}, errors.New("invalid array length")
	}

	args := make([]string, 0, count)

	for i := 0; i < count; i++ {
		typ, err := p.r.ReadByte()
		if err != nil {
			return Command{}, err
		}
		if typ != '$' {
			return Command{}, errors.New("expected bulk string")
		}
		lenLine, err := p.r.ReadBytes('\r')
		if err != nil {
			return Command{}, err
		}
		p.r.ReadByte()
		strLen, _ := strconv.Atoi(string(lenLine[:len(lenLine)-1]))

		data := make([]byte, strLen+2)
		p.r.Read(data)
		args = append(args, string(data[:strLen]))
	}

	return Command{Args: args}, nil
}

func parseInline(line []byte) Command {
	parts := []string{}
	curr := []byte{}
	for _, b := range line {
		if b == ' ' || b == '\n' || b == '\r' {
			if len(curr) > 0 {
				parts = append(parts, string(curr))
				curr = []byte{}
			}
		} else {
			curr = append(curr, b)
		}
	}
	if len(curr) > 0 {
		parts = append(parts, string(curr))
	}
	return Command{Args: parts}
}
