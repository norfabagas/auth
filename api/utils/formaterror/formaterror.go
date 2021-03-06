package formaterror

import (
	"errors"
	"log"
	"strings"
)

func FormatError(err string) error {
	defer log.Println(err)

	if strings.Contains(err, "nam") {
		return errors.New("Name Already Taken")
	}
	if strings.Contains(err, "email") {
		return errors.New("Email Already Taken")
	}
	if strings.Contains(err, "hashedPassword") {
		return errors.New("Incorrect Email or Password")
	}
	return errors.New("Incorrect Details")
}
