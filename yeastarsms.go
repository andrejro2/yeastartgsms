package yeastarsms

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"
	"time"
)

type Connection struct {
	Host     string
	Username string
	Secret   string
}

var conn *Connection

// SetConnection initializes the Connection struct with details
func SetConnection(host string, username string, password string) {
	conn = &Connection{
		Host:     host,
		Username: username,
		Secret:   password,
	}
}

// GetConnection returns the initialized Connection struct
func GetConnection() *Connection {
	return conn
}

func ConnectToService() (*Connection, error) {
	if conn != nil {
		return conn, nil
	}
	return nil, errors.New("error: Connection variable is empty or does not exist")
}

// SendSMS connects to the GSM gateway, logs in, and sends an SMS.
func SendSMS(port int, dest, message string, id string) error {
	ConnectToService()
	connect, err := net.Dial("tcp", conn.Host)
	if err != nil {
		return fmt.Errorf("failed to connect to gateway: %w", err)
	}
	defer connect.Close()

	reader := bufio.NewReader(connect)
	// Step 0: Wait for initial greeting
	initialResp, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read initial response: %w", err)
	}

	if !strings.Contains(initialResp, "Asterisk Call Manager/1.1") {
		return fmt.Errorf("unexpected initial response: %s", initialResp)
	}
	// Step 1: Send login command
	loginCmd := fmt.Sprintf("Action: Login\r\nUsername: %s\r\nSecret: %s\r\n\r\n", conn.Username, conn.Secret)
	_, err = connect.Write([]byte(loginCmd))
	if err != nil {
		return fmt.Errorf("failed to send login command: %w", err)
	}

	// Step 2: Wait for login response (2 lines expected)
	loginResp1, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read first line of login response: %w", err)
	}
	if !strings.Contains(loginResp1, "Success") {
		return fmt.Errorf("login failed, authentication not accepted: %s", loginResp1)
	}
	loginResp2, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read second line of login response: %w", err)
	}

	// Verify that the second line contains "Authentication accepted"
	if !strings.Contains(loginResp2, "Authentication accepted") {
		return fmt.Errorf("login failed, authentication not accepted: %s", loginResp2)
	}
	// Step 3: Prepare the SMS command
	encodedMessage := url.QueryEscape(message)
	smsCmd := fmt.Sprintf(
		"Action: smscommand\r\ncommand: gsm send sms %d %s \"%s\" %s\r\n\r\n",
		port+1, dest, encodedMessage, id,
	)

	// Step 4: Send SMS command
	_, err = connect.Write([]byte(smsCmd))
	if err != nil {
		return fmt.Errorf("failed to send SMS command: %w", err)
	}

	// Step 5: Wait for SMS response (optional but recommended)
	connect.SetReadDeadline(time.Now().Add(5 * time.Second)) // Timeout after 5 seconds
	sendresult := false
	readresult := ""
	for i := 0; i < 3; i++ {
		smsResp, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read SMS response: %w", err)
		}
		if strings.Contains(smsResp, "Response: Follows") {
			log.Println("SMS sent successfully!")
			sendresult = true
		}
		readresult += smsResp
	}
	if !sendresult {
		return fmt.Errorf("SMS command failed: %s", readresult)
	}

	return nil
}
