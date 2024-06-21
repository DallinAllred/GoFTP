package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"runtime"
	"strings"
)

var host = flag.String("h", "localhost", "FTP Server")
var port = flag.Int("p", 2020, "Port number")

func main() {
	flag.Parse()
	connString := fmt.Sprintf("%s:%d", *host, *port)
	conn, err := net.Dial("tcp", connString)
	if err != nil {
		log.Fatalf("Unable to establish connection to %s:%d\n%s\n", *host, *port, err)
	}
	defer conn.Close()
	netReader := bufio.NewReader(conn)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("GoFTP connection established")
	fmt.Println("----------------------------")
	for {
		fmt.Printf("GoFTP(%s)> ", *host)
		text, _ := reader.ReadString('\n')
		if runtime.GOOS == "windows" {
			text = strings.Replace(text, "\r\n", "", -1)
		} else {
			text = strings.Replace(text, "\n", "", -1)
		}
		if text == "exit" {
			break
		}
		input := strings.Split(text, " ")
		command := input[0]
		args := strings.Join(input[1:], " ")
		switch command {
		// Local commands
		case "lcd":
			if len(input) > 1 {
				err := os.Chdir(input[1])
				if err != nil {
					fmt.Printf("%s does not exist\n", input[1])
				}
			}
		case "lls":
			dir, err := os.Getwd()
			if err != nil {
				fmt.Println("Unable to list directory contents")
				continue
			}
			if len(input) > 1 {
				dir = input[1] // How does GO handle spaces in the dir argument?
			}
			dirContents, err := os.ReadDir(dir)
			if err != nil {
				fmt.Println("Unable to list directory contents")
				continue
			}
			for _, entry := range dirContents {
				fmt.Println(entry)
			}
		case "lpwd":
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Println("Unable to determine working directory")
				continue
			}
			fmt.Println(cwd)

		// Remote commands
		case "cd":
			err := sendServerCommand(conn, command, args)
			if err != nil {
				fmt.Println("Communication error with server")
			}
		case "ls":
			err := sendServerCommand(conn, command, args)
			if err != nil {
				fmt.Println("Communication error with server")
			}
			dirContents, err := receiveServerResponse(netReader)
			if err != nil {
				fmt.Println("Unable to display remote working directory")
			}
			for _, entry := range strings.Split(*dirContents, ";") {
				fmt.Println(entry)
			}
		case "pwd":
			err := sendServerCommand(conn, command, args)
			if err != nil {
				fmt.Println("Communication error with server")
			}
			cwd, err := receiveServerResponse(netReader)
			if err != nil {
				fmt.Println("Unable to display remote working directory")
			}
			fmt.Println(*cwd)

		// File transfer commands
		case "get":
		case "put":
		}
	}
}

func sendServerCommand(conn net.Conn, command string, args string) error {
	_, err := io.WriteString(conn, command+"\n")
	if err != nil {
		return err
	}
	_, err = io.WriteString(conn, args+"\n")
	if err != nil {
		return err
	}
	return nil
}

func receiveServerResponse(reader *bufio.Reader) (*string, error) {
	response, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	response = strings.Replace(response, "\n", "", -1)
	return &response, nil
}

// mustCopy(os.Stdout, conn)

// func mustCopy(dst io.Writer, src io.Reader) {
// 	if _, err := io.Copy(dst, src); err != nil {
// 		log.Fatal(err)
// 	}
// }
