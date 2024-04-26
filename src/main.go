package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func sshHandler(w http.ResponseWriter, r *http.Request) {
	user := os.Getenv("SSH_USER")
	pass := os.Getenv("SSH_PASS")
	host := os.Getenv("SSH_HOST")
	port := os.Getenv("SSH_PORT")

	hostport := fmt.Sprintf("%s:%s", host, port)

	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer wsConn.Close()

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshConn, err := ssh.Dial("tcp", hostport, config)
	if err != nil {
		log.Println("SSH dial error:", err)
		return
	}
	defer sshConn.Close()

	session, err := sshConn.NewSession()
	if err != nil {
		log.Println("SSH session error:", err)
		return
	}
	defer session.Close()

	sshOut, err := session.StdoutPipe()
	if err != nil {
		log.Println("STDOUT pipe error:", err)
		return
	}

	sshIn, err := session.StdinPipe()
	if err != nil {
		log.Println("STDIN pipe error:", err)
		return
	}

	if err := session.RequestPty("xterm", 80, 40, ssh.TerminalModes{}); err != nil {
		log.Println("Request PTY error:", err)
		return
	}

	if err := session.Shell(); err != nil {
		log.Println("Start shell error:", err)
		return
	}

	go func() {
		defer session.Close()
		buf := make([]byte, 1024)
		for {
			n, err := sshOut.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Println("Read from SSH stdout error:", err)
				}
				return
			}
			if n > 0 {
				err = wsConn.WriteMessage(websocket.BinaryMessage, buf[:n])
				if err != nil {
					log.Println("Write to WebSocket error:", err)
					return
				}
			}
		}
	}()

	for {
		messageType, p, err := wsConn.ReadMessage()
		if err != nil {
			if err != io.EOF {
				log.Println("Read from WebSocket error:", err)
			}
			return
		}
		if messageType == websocket.BinaryMessage || messageType == websocket.TextMessage {
			_, err = sshIn.Write(p)
			if err != nil {
				log.Println("Write to SSH stdin error:", err)
				return
			}
		}
	}
}

func checkForVariables() error {
	if os.Getenv("SSH_USER") == "" {
		return fmt.Errorf("SSH_USER is not set")
	}
	if os.Getenv("SSH_PASS") == "" {
		return fmt.Errorf("SSH_PASS is not set")
	}
	if os.Getenv("SSH_HOST") == "" {
		return fmt.Errorf("SSH_HOST is not set")
	}
	if os.Getenv("SSH_PORT") == "" {
		return fmt.Errorf("SSH_PORT is not set")
	}
	return nil
}

func printUsage() {
	fmt.Println("Usage (set proper environmental variables): ")
	fmt.Println("SSH_USER - username for ssh connection")
	fmt.Println("SSH_PASS - password for ssh connection")
	fmt.Println("SSH_HOST - host for ssh connection")
	fmt.Println("SSH_PORT - port for ssh connection")
	fmt.Println("MOUNT_HTML - mount html files (default true, set to false to disable)")

}
func main() {
	errOf := checkForVariables()
	if errOf != nil {

		printUsage()
		log.Fatal(errOf)
	}
	shouldMountHTML := os.Getenv("MOUNT_HTML") == "true" || os.Getenv("MOUNT_HTML") == ""
	http.HandleFunc("/ssh", sshHandler)

	if shouldMountHTML {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "index.html")
		})
		http.HandleFunc("/xterm.css", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "xterm.css")
		})
		http.HandleFunc("/xterm.js", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "xterm.js")
		})
	}

	fmt.Println("Starting server on :8280")
	log.Fatal(http.ListenAndServe(":8280", nil))

}
