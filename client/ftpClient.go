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

var commands = map[string]func(string, string){
	"lcd":  localCD,
	"lls":  localList,
	"lpwd": localPwd,
	"cd":   remoteCD,
	"ls":   remoteList,
	"pwd":  remotePwd,
	"get":  getFile,
	"put":  putFile,
}

var conn net.Conn
var netReader *bufio.Reader

var serverErrorMsg = "Communication error with server. Server may be down.\nExiting..."

func main() {
	flag.Parse()
	connString := fmt.Sprintf("%s:%d", *host, *port)
	var err error
	conn, err = net.Dial("tcp", connString)
	if err != nil {
		log.Fatalf("Unable to establish connection to %s:%d\n%s\nExiting...", *host, *port, err)
	}
	// conn = cnxn
	defer conn.Close()
	netReader = bufio.NewReader(conn)

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

		function, ok := commands[command]
		if !ok {
			fmt.Printf("Unrecognized command: %s\n", command)
			continue
		}
		function(command, args)
	}
}

// localCD changes the working directory on the client
func localCD(cmd string, dir string) {
	if len(dir) > 0 {
		err := os.Chdir(dir)
		if err != nil {
			fmt.Printf("%s does not exist\n", dir)
		}
	}
}

// localList lists the client working directory contents
// Optionally may list contents of a supplied directory
func localList(cmd string, dir string) {
	listDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Unable to list directory contents")
		return
	}
	if len(dir) > 0 {
		listDir = dir // How does GO handle spaces in the dir argument?
	}
	dirContents, err := os.ReadDir(listDir)
	if err != nil {
		fmt.Println("Unable to list directory contents")
		return
	}
	for _, entry := range dirContents {
		fmt.Println(entry)
	}
}

// localPwd displays the client working directory
func localPwd(cmd string, _ string) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Unable to determine working directory")
		return
	}
	fmt.Println(cwd)
}

// remoteCD changes the working directory on the server
func remoteCD(cmd string, dir string) {
	err := sendServerCommand(conn, cmd, dir)
	if err != nil {
		log.Fatalln(serverErrorMsg)
	}
}

// remoteList lists the server working directory contents
// Optionally may list contents of a supplied directory
func remoteList(cmd string, dir string) {
	err := sendServerCommand(conn, cmd, dir)
	if err != nil {
		log.Fatalln(serverErrorMsg)
	}
	dirContents, err := receiveServerResponse(netReader)
	if err != nil {
		log.Fatalln("Unable to display remote working directory")
	}
	for _, entry := range strings.Split(*dirContents, " ; ") {
		fmt.Println(entry)
	}
}

// remotePwd displays the server working directory
func remotePwd(cmd string, args string) {
	err := sendServerCommand(conn, cmd, args)
	if err != nil {
		log.Fatalln(serverErrorMsg)
	}
	cwd, err := receiveServerResponse(netReader)
	if err != nil {
		log.Fatalln("Unable to display remote working directory")
	}
	fmt.Println(*cwd)
}

// getFile copies a file from the server to the client
func getFile(cmd string, filename string) {
}

// putFile sends a file to the server from the client
func putFile(cmd string, filename string) {
}

// sendServerCommand sends the command and arguments across the connection to the server
func sendServerCommand(conn net.Conn, command string, args string) error {
	_, err := io.WriteString(conn, command+"\n")
	if err != nil {
		return err
	}
	// if len(args) > 0 {
	_, err = io.WriteString(conn, args+"\n")
	if err != nil {
		return err
	}
	// }
	return nil
}

// receiveServerResponse reads the response from the server and strips the newline
func receiveServerResponse(reader *bufio.Reader) (*string, error) {
	response, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	response = strings.Replace(response, "\n", "", -1)
	return &response, nil
}
