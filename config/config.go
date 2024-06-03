package Config

type config struct {
	HomeDir        string `json:"home_dir"`
	SafelineURL    string `json:"safeline_url"`
	API_CSFR_Path  string `json:"api_csfr_path"`
	API_Login_Path string `json:"api_login_path"`
	API_TFA_Path   string `json:"api_tfa_path"`
	API_Cert_Path  string `json:"api_cert_path"`
	LogLevel       string `json:"log_level"`
	LogFormat      string `json:"log_format"`
}

var SafelineConfig = &config{
	HomeDir:        ".safeline",
	SafelineURL:    "https://localhost:9443",
	API_CSFR_Path:  "/api/open/auth/csrf",
	API_Login_Path: "/api/open/auth/login",
	API_TFA_Path:   "/api/open/auth/tfa",
	API_Cert_Path:  "/api/open/cert",
	LogLevel:       "info",
	LogFormat:      "plain",
}
