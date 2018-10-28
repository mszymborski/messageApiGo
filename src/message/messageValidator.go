package message

import "net/mail"

func validateMailAddress(address string) error {
	if _, err := mail.ParseAddress(address); err != nil{
		return err
	}
	return nil
}

func validateMessageModel(message *POSTMessageModel) error {

	if err := validateMailAddress(message.EmailAddress); err != nil {
		return err
	}

	return nil
}