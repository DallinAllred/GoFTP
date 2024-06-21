package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

var port = flag.Int("p", 2020, "Port number")

func main() {
	flag.Parse()
	fmt.Printf("Server running on port %d\n", *port)
	connString := fmt.Sprintf("localhost:%d", *port)
	listener, err := net.Listen("tcp", connString)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(c net.Conn) {
	defer c.Close()
	reader := bufio.NewReader(c)
	for {
		command, err := reader.ReadString('\n')
		if err != nil {
			_, err = io.WriteString(c, err.Error()+"\n")
		}
		command = strings.Replace(command, "\n", "", -1)
		fmt.Printf("Command: %s\n", command)
		switch command {
		case "pwd":
			var response string
			cwd, err := os.Getwd()
			if err != nil {
				response = err.Error()
			} else {
				response = cwd
			}
			_, err = io.WriteString(c, response+"\n")
		case "cd":
			args, _ := reader.ReadString('\n')
			args = strings.Replace(args, "\n", "", -1)
			if len(args) > 0 {
				err := os.Chdir(args)
				if err != nil {
					_, err = io.WriteString(c, err.Error()+"\n")
				}
			}
		case "ls":
			args, _ := reader.ReadString('\n')
			args = strings.Replace(args, "\n", "", -1)
			dir, err := os.Getwd()
			if err != nil {
				_, err = io.WriteString(c, err.Error()+"\n")
			}
			if len(args) > 0 {
				dir = args
			}
			dirContents, err := os.ReadDir(dir)
			strContents := []string{}
			for _, entry := range dirContents {
				strContents = append(strContents, fmt.Sprint(entry))
			}
			response := strings.Join(strContents, ";")
			_, err = io.WriteString(c, response+"\n")
		}
		// _, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		// if err != nil {
		// 	return
		// }
		// time.Sleep(1 * time.Second)
	}
}
