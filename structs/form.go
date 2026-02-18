package structs

type User struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type InternalPaymentRequest struct {
	FromAccount string  `json:"from_account"`
	ToAccount   string  `json:"to_account"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
}

type RecipientDetails struct {
	RecipientNumber string `json:"recipientNumber"`
	RecipientName   string `json:"recipientName"`
}

type ExternalPaymentRequest struct {
	Account         string           `json:"from_account"`
	Amount          float64          `json:"amount"`
	Currency        string           `json:"currency"`
	TransactionType string           `json:"transaction_type"`
	Recipient       RecipientDetails `json:"recipient"`
}

type ExternalPaymentResponse struct {
	PaymentID      string           `json:"payment_id"`
	Recipient      RecipientDetails `json:"recipient"`
	Status         string           `json:"status"`
	ProviderStatus string           `json:"provider_status"`
}

type TopUp struct {
	Account  string  `json:"account"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type TopUpResponse struct {
	PaymentID string `json:"payment_id"`
	Status    string `json:"status"`
}
