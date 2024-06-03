package data

import (
	"encoding/json"
	"time"

	"github.com/pquerna/otp/totp"
)

type Account struct {
	SafelineURL string `json:"safeline_url"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	OTPSecret   string `json:"otp_secret"`
}

func NewAccount(safelineURL, username, password, otpSecret string) *Account {
	return &Account{
		SafelineURL: safelineURL,
		Username:    username,
		Password:    password,
		OTPSecret:   otpSecret,
	}
}

func NewAccountFromJSON(data []byte) (*Account, error) {
	var account Account
	err := json.Unmarshal(data, &account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}
func NewAccountWithoutOTP(safelineURL, username, password string) *Account {
	return NewAccount(safelineURL, username, password, "")
}

func (a *Account) GetSafelineURL() string {
	return a.SafelineURL
}

func (a *Account) GetUsername() string {
	return a.Username
}

func (a *Account) GetPassword() string {
	return a.Password
}

func (a *Account) GetOTPSecret() string {
	return a.OTPSecret
}

func (a *Account) GetTOTPcode() (string, error) {
	// Generate TOTP code
	code, err := totp.GenerateCode(a.GetOTPSecret(), time.Now())
	if err != nil {
		return "unable to Generate TOTPCode", err
	}
	return code, nil
}

func (a *Account) IsEnableOTP() bool {
	return a.OTPSecret != ""
}

func (a *Account) ToJSON() ([]byte, error) {
	return json.Marshal(a)
}
