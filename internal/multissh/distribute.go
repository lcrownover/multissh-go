package multissh

import (
	"bufio"
	"fmt"
	"io"
	// "os"
	"sync"

	"golang.org/x/crypto/ssh"
)

type message struct {
	hostname string
	line     string
}

// Distribute : for each of the nodes, spin up a goroutine that
// starts a connection, then runs the given command in that connection
func Distribute(config *ssh.ClientConfig, p CommandParameters) {
	var wg sync.WaitGroup
	command := p.Command

	for _, node := range p.Nodes {
		wg.Add(1)
		go func(node string) {
			defer wg.Done()

			c, err := Connect(node, config)
            if err != nil {
                return
            }
			defer c.connection.Close()
			RunCommand(command, c)
		}(node)
	}

	wg.Wait()
}

// RunCommand: run the given command in the ssh connection
func RunCommand(command string, c connection) {
	// Try to establish new session
    conn := c.connection
	sess, err := conn.NewSession()
	if err != nil {
		logger("distribute").Error(fmt.Sprintf("Failed to establish session: %v", err))
        return
	}
	defer sess.Close()

	sessStdout, err := sess.StdoutPipe()
	if err != nil {
		logger("distribute").Error(fmt.Sprintf("Failed to connect stdout: %v", err))
        return
	}
	sessStderr, err := sess.StderrPipe()
	if err != nil {
		logger("distribute").Error(fmt.Sprintf("Failed to connect stderr: %v", err))
        return
	}

    m := message{
        hostname: c.host,
    }

    soc := make(chan(message), 1000)
    sec := make(chan(message), 1000)

    // Set up the receivers
    go messageReceiver(soc)
    go messageReceiver(sec)

    // Senders
	go scanSessionMessageToChannel(sessStdout, soc, m)
	go scanSessionMessageToChannel(sessStderr, sec, m)

	sess.Run(command)

}


func messageReceiver(c chan message) {
    for {
        select {
            case m := <-c:
                fmt.Printf("[%s] %s\n", m.hostname, m.line)
        }
    }
}

func bufferedMessageReceiver(c chan message) {
    var msgs []string
    for {
        select {
            case m := <-c:
                msgs = append(msgs, fmt.Sprintf("[%s] %s\n", m.hostname, m.line))
        }
    }
    fmt.Println(msgs)
    for _,line := range msgs {
        fmt.Printf(line)
    }
}

func scanSessionMessageToChannel(s io.Reader, c chan message, m message) {
    scanner := bufio.NewScanner(s)
    for {
        if tkn := scanner.Scan(); tkn {
            rcv := scanner.Bytes()
            raw := make([]byte, len(rcv))
            copy(raw, rcv)
            m.line = string(raw)
            c <- m
        }
    }
}
