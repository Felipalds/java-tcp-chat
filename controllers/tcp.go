package controllers

import (
	"bufio"
	"crypto/rsa"
	"fmt"
	"go-tcp-chat/encrypt"
	"go-tcp-chat/models"
	"go-tcp-chat/utils"
	"net"
	"strings"
)

func HandleClient(conn net.Conn, a *int) {
	*a += 1
	b := *a
	var user models.User
	var pk *rsa.PrivateKey
	var aesKey []byte
	var auth bool
	auth = false

	fmt.Println("Handle connection %d", b)
	defer conn.Close()
	//buffer := make([]byte, 2048)

	for {
		//_, err := conn.Read(buffer)
		status, err := bufio.NewReader(conn).ReadString('\n')
		fmt.Println(status)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		buffStr := status
		fmt.Println(buffStr)
		buffStr = strings.Trim(buffStr, "\r\n")
		buffStr = strings.ReplaceAll(buffStr, "\x00", "")
		if auth {
			buffStr, err = encrypt.Decrypt(buffStr, aesKey)
			if err != nil {
				fmt.Println("Error:", err)
				panic(err)
			}
		}

		buffParts := strings.Split(buffStr, " ")
		if !utils.IsRequestValid(buffParts) {
			fmt.Fprintf(conn, "ERRO: faltou passar argumentos")
		}

		buffParts[len(buffParts)-1] = strings.ReplaceAll(buffParts[len(buffParts)-1], "\x00", "")
		buffParts[len(buffParts)-1] = strings.ReplaceAll(buffParts[len(buffParts)-1], "\n", "")
		msg, err := HandleRequest(&conn, buffParts, &user, &pk, &aesKey, &auth)
		var encryptErr error

		if auth {
			msg, encryptErr = encrypt.Encrypt(msg, aesKey)
			if encryptErr != nil {
				fmt.Println("Error encrypting msg to client:", err)
			}
		}

		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(conn, "ERRO:%s\n", err.Error())
		} else {
			msg += "\n"
			fmt.Fprintf(conn, msg)
		}
		//buffer = make([]byte, 2048)
	}
}
