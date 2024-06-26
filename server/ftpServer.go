package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var port = flag.Int("p", 2020, "Port number")

var fsRoots = `^(?i)\/|[A-Z]:`

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
		fmt.Printf("Connection established with %s\n", conn.RemoteAddr())
		if err != nil {
			log.Print(err)
			continue
		}
		sessionDir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn, sessionDir)
	}
}

func handleConn(c net.Conn, sessionDir string) {
	defer c.Close()
	reader := bufio.NewReader(c)
	for {
		command, err := reader.ReadString('\n')
		if err != nil {
			_, err = io.WriteString(c, err.Error()+"\n")
		}
		command = strings.Replace(command, "\n", "", -1)

		args, err := reader.ReadString('\n')
		if err != nil {
			_, err = io.WriteString(c, err.Error()+"\n")
		}
		args = strings.Replace(args, "\n", "", -1)

		if command == "" {
			break
		}
		fmt.Printf("%s -- Command: %s %s\n", c.RemoteAddr(), command, args)
		switch command {
		case "pwd":
			_, err = io.WriteString(c, sessionDir+"\n")

		case "cd":
			if len(args) > 0 {
				proposedPath := args
				root, err := regexp.Match(fsRoots, []byte(args))
				if err != nil {
					_, err = io.WriteString(c, sessionDir+"\n")
					return
				}
				if !root {
					pathSegments := []string{sessionDir, proposedPath}
					proposedPath = filepath.Join(pathSegments...)
					proposedPath, _ = filepath.Abs(proposedPath)
				}
				_, err = os.Stat(proposedPath)
				if err != nil {
					_, err = io.WriteString(c, err.Error()+"\n")
				} else {
					sessionDir = proposedPath
					_, err = io.WriteString(c, "\n")
				}
				fmt.Println(sessionDir)
			}

		case "ls":
			dir := sessionDir
			if err != nil {
				_, err = io.WriteString(c, err.Error()+"\n")
			}
			if len(args) > 0 {
				dir = args
			}
			dirContents, err := os.ReadDir(dir)
			if err != nil {
				_, err = io.WriteString(c, err.Error()+"\n")
			}
			strContents := []string{}
			for _, entry := range dirContents {
				strContents = append(strContents, fmt.Sprint(entry))
			}
			response := strings.Join(strContents, " ; ")
			_, err = io.WriteString(c, response+"\n")

		case "get":
			if len(args) == 0 {
				continue
			}
			file, err := os.Open(args)
			if err != nil {
				_, err = io.WriteString(c, err.Error()+"\n")
			}
			n, err := io.Copy(c, file)
			file.Close()
			fmt.Println(n, "bytes sent")
			if err != nil {
				_, err = io.WriteString(c, err.Error()+"\n")
			}
		case "put":
			if len(args) == 0 {
				continue
			}
			file, err := os.Create(args)
			if err != nil {
				_, err = io.WriteString(c, err.Error()+"\n")
			}
			var currentByte int64 = 0
			totalBytesReceived := 0
			bufferSize := 64
			buffer := make([]byte, bufferSize)
			for {
				n, err := c.Read(buffer)
				_, err = file.WriteAt(buffer[:n], currentByte)
				totalBytesReceived += n
				if err == io.EOF || n < bufferSize {
					break
				}
				currentByte += int64(n)
			}
			file.Close()
			fmt.Printf("%d bytes received by server\n", totalBytesReceived)
			_, err = io.WriteString(c, fmt.Sprintf("%d bytes received by server\n", totalBytesReceived))
		}
	}
	fmt.Printf("Connection with %s terminated\n", c.RemoteAddr())
}
