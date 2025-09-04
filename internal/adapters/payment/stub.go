package payment

type AuthResult struct {
	ID          uint
	Status      string // authorized
	ExternalRef string
}

func Authorize(amountCents int) (string, string) { // idStr, status
	return "stub-auth-123", "authorized"
}
