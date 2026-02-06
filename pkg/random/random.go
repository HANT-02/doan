package random

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateSixDigitOtp() string {
	rand.Seed(time.Now().UnixNano())
	otp := rand.Intn(900000) + 100000 // ensures a 6-digit number
	return fmt.Sprintf("%06d", otp)
}
