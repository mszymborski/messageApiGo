package message

import (
	"Utils"
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
)

func (message *DBMessageModel) buildMessage() string {

	messageOut := ""
	messageOut += fmt.Sprintf("From: %s\r\n", Utils.SMTP_USERNAME)
	messageOut += fmt.Sprintf("To: %s\r\n", message.EmailAddress)

	messageOut += fmt.Sprintf("Subject: %s\r\n", message.Title)
	messageOut += "\r\n" + message.Content

	return messageOut
}

func sendMails(model []DBMessageModel) error {

	for _, message := range model {
		client, err := prepareSMTPClient()

		if err != nil {
			return err
		}

		if err = client.Rcpt(message.EmailAddress); err != nil {
			log.Fatal(err)
			return err
		}


		w, err := client.Data()
		if err != nil {
			log.Fatal(err)
			return err
		}

		_, err = w.Write([]byte(message.buildMessage()))
		if err != nil {
			log.Fatal(err)
			return err
		}

		err = w.Close()
		if err != nil {
			log.Fatal(err)
			return err
		}

		client.Quit()

		if err = delete(message.Id); err != nil {
			log.Fatal(err)
			return err
		}
	}

	return nil
}

func prepareSMTPClient() (*smtp.Client, error) {

	serverAddress := Utils.SMTP_HOST + ":" + Utils.SMTP_PORT

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         Utils.SMTP_HOST,
	}

	conn, err := tls.Dial("tcp", serverAddress, tlsconfig)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	auth := smtp.PlainAuth("", Utils.SMTP_USERNAME , Utils.SMTP_PASSWORD, Utils.SMTP_HOST)

	client, err := smtp.NewClient(conn, Utils.SMTP_HOST)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if err = client.Auth(auth); err != nil {
		log.Fatal(err)
		return nil, err
	}

	if err = client.Mail(Utils.SMTP_USERNAME); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return client, nil
}