package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

type TerminalSize struct {
	Cols int `json:"cols"`
	Rows int `json:"rows"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func sshHandler(w http.ResponseWriter, r *http.Request) {
	user := "root"
	pass := "plokij2wsx3edc2$"
	host := "172.16.0.69"
	port := "22"

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

	var termSize TerminalSize
	if err := wsConn.ReadJSON(&termSize); err != nil {
		log.Println("Read terminal size error:", err)
		return
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 8192,
		ssh.TTY_OP_OSPEED: 8192,
		ssh.IEXTEN:        0,
	}

	if err := session.RequestPty("xterm", termSize.Cols, termSize.Rows, modes); err != nil {
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
		if messageType == websocket.TextMessage {
			var newSize TerminalSize
			if err := json.Unmarshal(p, &newSize); err == nil {
				// Resize terminal
				//log.Printf("Resizing terminal to cols: %d, rows: %d", newSize.Cols, newSize.Rows)
				session.WindowChange(newSize.Rows, newSize.Cols)
			} else {
				// Write to SSH stdin
				_, err = sshIn.Write(p)
				if err != nil {
					log.Println("Write to SSH stdin error:", err)
					return
				}
			}
		}
	}
}

func main() {
	http.HandleFunc("/ssh", sshHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/xterm.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "xterm.css")
	})
	http.HandleFunc("/xterm.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "xterm.js")
	})

	http.HandleFunc("/xterm-addon-fit.min.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "xterm-addon-fit.min.js")
	})

	fmt.Println("Starting server on :8280")
	log.Fatal(http.ListenAndServe(":8280", nil))
}
