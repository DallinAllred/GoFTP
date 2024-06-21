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
		case "lpwd":
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Println("Unable to determine working directory")
				continue
			}
			fmt.Println(cwd)
		case "lcd":
			if len(input) > 1 {
				err := os.Chdir(input[1])
				if err != nil {
					fmt.Printf("%s does not exist\n", input[1])
				}
			}
		case "lls":
			var dir string
			if len(input) == 1 {
				cwd, err := os.Getwd()
				if err != nil {
					fmt.Println("Unable to list directory contents")
					continue
				}
				dir = cwd
			} else {
				dir = input[1] // How does GO handle spaces in the dir argument?
			}
			dirContents, err := os.ReadDir(dir)
			if err != nil {
				fmt.Println("Unable to list directory contents")
				continue
			}
			for _, item := range dirContents {
				fmt.Println(item)
			}
		case "pwd":
			_, err := io.WriteString(conn, command+"\n")
			if err != nil {
				fmt.Println("Communication error with server")
			}
			err = handleResponse(netReader)
			if err != nil {
				fmt.Println("Unable to display remote working directory")
			}
		case "cd":
			_, err := io.WriteString(conn, command+"\n")
			if err != nil {
				fmt.Println("Communication error with server")
			}
			_, err = io.WriteString(conn, args+"\n")
			if err != nil {
				fmt.Println("Communication error with server")
			}
		case "ls":
		case "get":
		case "put":
		}
	}
}

func handleResponse(reader *bufio.Reader) error {
	response, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	fmt.Print(response)
	// for _, item := range response {
	// 	fmt.Print(item)
	// }
	return nil
}

// mustCopy(os.Stdout, conn)

// func mustCopy(dst io.Writer, src io.Reader) {
// 	if _, err := io.Copy(dst, src); err != nil {
// 		log.Fatal(err)
// 	}
// }
