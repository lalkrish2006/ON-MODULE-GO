package services

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
)

// SendEmail sends an email using Gmail SMTP
func SendEmail(to, subject, body string) error {
	from := "anishrkumar2k5@gmail.com"
	password := "veorxoeneclhtcgf" // App Password
	host := "smtp.gmail.com"
	port := "587" // Standard Port for STARTTLS

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	auth := smtp.PlainAuth("", from, password, host)

	// Custom TLS configuration to skip verification if needed (matching PHP 'verify_peer' => false)
	// In production, InsecureSkipVerify: true is bad, but matching local XAMPP setup parity.
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Dial the connection
	conn, err := tls.Dial("tcp", host+":465", tlsConfig)
	if err == nil {
		// If 465 works with direct SSL
		defer conn.Close()
		client, err := smtp.NewClient(conn, host)
		if err != nil {
			return err
		}
		defer client.Quit()

		if err = client.Auth(auth); err != nil {
			return err
		}
		if err = client.Mail(from); err != nil {
			return err
		}
		if err = client.Rcpt(to); err != nil {
			return err
		}
		w, err := client.Data()
		if err != nil {
			return err
		}
		_, err = w.Write([]byte(message))
		if err != nil {
			return err
		}
		err = w.Close()
		if err != nil {
			return err
		}
		log.Printf("Email sent successfully to %s via Port 465\n", to)
		return nil
	}

	// Fallback to 587 STARTTLS if 465 fails or logic prefers standard SendMail
	err = smtp.SendMail(host+":"+port, auth, from, []string{to}, []byte(message))
	if err != nil {
		log.Println("SMTP Error:", err)
		return err
	}
	
	log.Printf("Email sent successfully to %s via Port 587\n", to)
	return nil
}

// SendODNotification sends notification to mentors
func SendODNotification(mentorEmail, studentName, odType string) {
	subject := "OD Notification"
	// Matching PHP Body: 'You got a Notification Check the OD Module by Clicking this link ...'
	body := fmt.Sprintf(`You got a Notification Check the OD Module by Clicking this link <a href="http://localhost:8082/login">http://localhost:8082/login</a><br><br>
	Details:<br>
	Student: %s<br>
	Type: %s`, studentName, odType)

	go func() {
		err := SendEmail(mentorEmail, subject, body)
		if err != nil {
			log.Println("Failed to send email:", err)
		}
	}()
}
