package sendemail

import (
	"fmt"
	"montelukast/modules/pharmacist/entity"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	"os"
	"strings"

	"gopkg.in/gomail.v2"
)

func SendEmail(to []string, contentType, body string, mailer *gomail.Message) error {
	mailer.SetHeader("From", os.Getenv("EMAIL_SENDER"))
	mailer.SetHeader("To", strings.Join(to, ","))
	mailer.SetBody(contentType, body)

	dialer := gomail.NewDialer(
		appconstant.CONFIG_SMTP_HOST,
		appconstant.CONFIG_SMTP_PORT,
		os.Getenv("EMAIL_SENDER"),
		os.Getenv("EMAIL_AUTH"),
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		return apperror.ErrSendEmail
	}
	return nil
}

func SendMailForVerification(to []string, token string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("Subject", "[mediSEAne] Account Verification")

	contentType := "text/html"
	body := fmt.Sprintf(`Hello, <b>verify your email at mediSEAne here!</b>
		<button><a href="%s%s">Verify Email</a></button>`, os.Getenv("WEB_URL"), appconstant.VERIFY_EMAIL_PATH+token)
	err := SendEmail(to, contentType, body, mailer)
	if err != nil {
		return err
	}

	return nil
}

func SendMailForResetPassword(to []string, token string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("Subject", "[mediSEAne] Account Reset Password")

	contentType := "text/html"
	body := fmt.Sprintf(`Hello, <b>reset your password at mediSEAne here!</b>
		<button><a href="%s%s">Reset Password</a></button>`, os.Getenv("WEB_URL"), appconstant.RESET_PASSWORD_PATH+token)
	err := SendEmail(to, contentType, body, mailer)
	if err != nil {
		return err
	}

	return nil
}

func SendMailForPharmacistAccountDetail(pharmacist entity.Pharmacist) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("Subject", "[mediSEAne] Account Pharmacist Details")

	contentType := "text/html"
	body := fmt.Sprintf(`Hello, <b>below is your account details!</b>
		<div>
			<p>Name: %s</p>
			<p>Password: %s</p>
			<p>SIPA Number: %s</p>
			<p>Phone Number: %s</p>
			<p>Year Of Experience: %d</p>
		</div>
	`, pharmacist.Name, pharmacist.Password, pharmacist.SipaNumber, pharmacist.PhoneNumber, pharmacist.YearOfExperience)
	err := SendEmail([]string{pharmacist.Email}, contentType, body, mailer)
	if err != nil {
		return err
	}

	return nil
}
