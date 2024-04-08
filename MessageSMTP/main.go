package main

import (
	"errors"
	"fmt"
	"log"
	"net/smtp"
)

/*type plainAuth struct {
	identity, username, password string
	host                         string
}

func PlainAuth(identity, username, password, host string) smtp.Auth {
	return &plainAuth{identity, username, password, host}
}

func isLocalhost(name string) bool {
	return name == "localhost" || name == "127.0.0.1" || name == "::1"
}

func (a *plainAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	//if !server.TLS && !isLocalhost(server.Name) {
	//	return "", nil, errors.New("unencrypted connection")
	//}

	if server.Name != a.host {
		return "", nil, errors.New("wrong host name")
	}

	resp := []byte(a.identity + "\x00" + a.username + "\x00" + a.password)
	return "PLAIN", resp, nil
}
func (a *plainAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		// We've already sent everything.
		return nil, errors.New("unexpected server challenge")
	}
	return nil, nil
}*/

type loginAuth struct {
	username, password string
}

// LoginAuth is used for smtp login auth
func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unknown from server")
		}
	}
	return nil, nil
}

func main() {
	from := "valid"
	password := "password"

	to := []string{
		"fedasov03@mail.ru",
		//"fedasovsergey00@gmail.com",
	}

	host := "mail.mailhub.su" //host := "smtp.mail.ru"
	port := "25"              //port := "587"
	address := host + ":" + port

	subject := "Subject: Hello\n"
	body := "Hello"
	message := []byte(subject + body)

	//auth := PlainAuth("", from, password, host)
	auth := LoginAuth(from, password)

	err := smtp.SendMail(address, auth, "test@example.org", to, message)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("OK")
}

/*package main

import (
	"fmt"
	"net/smtp"
)

func main() {

	from := "fedasov03@mail.ru"
	password := "a8ibbs7qDWTxqE4qhQcQ"
	//from := "test@mailhub.su"
	//password := "password"

	to := []string{
		"fedasovsergey00@gmail.com",
	}

	host := "smtp.mail.ru" //host := "smtp.mail.ru"
	port := "25"           //port := "587"
	address := host + ":" + port

	subject := "Subject: Hello\n"
	body := "Hello"
	message := []byte(subject + body)

	auth := smtp.PlainAuth("", from, password, host)

	err := smtp.SendMail(address, auth, from, to, message)
	if err != nil {
		panic(err)
	}

	fmt.Println("Message sent successfully")
}*/

/*package main

import (
	"fmt"
	"log"
	"net/smtp"
)

type unencryptedAuth struct {
	smtp.Auth
}

func (a unencryptedAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	s := *server
	s.TLS = true
	return a.Auth.Start(&s)
}

func main() {
	from := "test@mailhub.su"
	password := "password"

	to := []string{
		"fedasovsergey00@gmail.com",
	}

	host := "mail.mailhub.su" //host := "smtp.mail.ru"
	port := "25"              //port := "587"
	address := host + ":" + port

	subject := "Subject: Hello\n"
	body := "Hello"
	message := []byte(subject + body)

	//auth := smtp.PlainAuth("", from, password, host)

	auth := unencryptedAuth{smtp.PlainAuth("", from, password, host)}

	err := smtp.SendMail(address, auth, "test@example.org", to, message)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("OK")
}*/

/*package main

import (
	"github.com/wneessen/go-mail"
	"log"
)

func main() {
	m := mail.NewMsg()
	if err := m.From("sergey@mailhub.su"); err != nil {
		log.Fatalf("failed to set From address: %s", err)
	}
	if err := m.To("fedasov03@mail.ru"); err != nil {
		log.Fatalf("failed to set To address: %s", err)
	}

	m.Subject("This is my first mail with go-mail!")
	m.SetBodyString(mail.TypeTextPlain, "Do you like this mail? I certainly do!")

	c, err := mail.NewClient("smtp.mailhub.su", mail.WithPort(25), mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername("Sergey"), mail.WithPassword("1234"))
	if err != nil {
		log.Fatalf("failed to create mail client: %s", err)
	}

	if err := c.DialAndSend(m); err != nil {
		log.Fatalf("failed to send mail: %s", err)
	}
}*/
