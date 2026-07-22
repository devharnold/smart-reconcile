package alerts

/*
- Sends scheduled alerts after the scheduler has completed fetching, matching and reconciling transactions
-
*/

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func finishedNormalizing() {
	_ = godotenv.Load("./root/.env")
	resendKey := os.Getenv("RESEND_KEY")
	if resendKey == "" {
		log.Fatal("Resend API Key is required")
	}
}
