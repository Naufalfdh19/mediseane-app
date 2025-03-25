package setup

import (
	"os"

	"github.com/resendlabs/resend-go"
)

func SetupResend() *resend.Client {
	apiKey := os.Getenv("RESEND_API_KEY")
	return resend.NewClient(apiKey)
}