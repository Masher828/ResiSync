package models

type ResetPassword struct {
	Otp      string `json:"otp"`
	Contact  string `json:"contact,omitempty"`
	Method   string `json:"method"`
	Password string `json:"password"`
}
