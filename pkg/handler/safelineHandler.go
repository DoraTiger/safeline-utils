package handler

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	cfg "github.com/DoraTiger/safeline-utils/config"
	"github.com/DoraTiger/safeline-utils/pkg/data"
	"github.com/DoraTiger/safeline-utils/pkg/utils"
	"github.com/sirupsen/logrus"
)

type CSRFResponse struct {
	Data struct {
		CSRFToken string `json:"csrf_token"`
	} `json:"data"`
}

type LoginResponse struct {
	Data struct {
		JWT string `json:"jwt"`
	} `json:"data"`
}

type TFAResponse struct {
	Data struct {
		JWT string `json:"jwt"`
	} `json:"data"`
	Err *string `json:"err"`
}

type CertListResponse struct {
	Data struct {
		Nodes []interface{} `json:"nodes"`
	} `json:"data"`
}

type CertResponse struct {
	Err *string `json:"err"`
}

type SafelineHandler struct {
	account data.Account
	JWT     string
	client  *http.Client
	logger  *logrus.Logger
}

func NewSafelineHandler() *SafelineHandler {
	return &SafelineHandler{
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

func (h *SafelineHandler) SetLogger(logger *logrus.Logger) {
	h.logger = logger
}

func (h *SafelineHandler) getCSRFToken() (string, error) {
	// Get CSRFResponse from Safeline
	resp, err := h.client.Get(fmt.Sprintf("%s%s", h.account.GetSafelineURL(), cfg.SafelineConfig.API_CSFR_Path))
	if err != nil {
		file, line := utils.GetErrorLocation()
		h.logger.Debug(fmt.Sprintf("Error in file %s, line %d: %v", file, line, err))
		return "unable to get CSRFResponse", err
	}
	defer resp.Body.Close()

	// Parse CSRF token
	var csrfResp CSRFResponse
	if err := json.NewDecoder(resp.Body).Decode(&csrfResp); err != nil {
		file, line := utils.GetErrorLocation()
		h.logger.Debug(fmt.Sprintf("Error in file %s, line %d: %v", file, line, err))
		return "unable to Parse CSRFResponse", err
	}

	return csrfResp.Data.CSRFToken, nil
}

func (h *SafelineHandler) validateTFA(code string) (string, error) {
	// Get CSRF token
	csrf_token, err := h.getCSRFToken()
	if err != nil {
		return csrf_token, err
	}

	headers := map[string]string{
		"authorization": "Bearer " + h.JWT,
	}

	payload := map[string]interface{}{
		"code":       code,
		"timestamp":  time.Now().UnixNano() / int64(time.Millisecond),
		"csrf_token": csrf_token,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s%s", h.account.GetSafelineURL(), cfg.SafelineConfig.API_TFA_Path), bytes.NewBuffer(body))
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, err := h.client.Do(req)
	if err != nil {
		file, line := utils.GetErrorLocation()
		h.logger.Debug(fmt.Sprintf("Error in file %s, line %d: %v", file, line, err))
		return "unable to get TFAResponse", err
	}
	defer resp.Body.Close()

	var tfaResp TFAResponse
	if err := json.NewDecoder(resp.Body).Decode(&tfaResp); err != nil {
		file, line := utils.GetErrorLocation()
		h.logger.Debug(fmt.Sprintf("Error in file %s, line %d: %v", file, line, err))
		return "unable to Parse TFAResponse", err
	}

	if tfaResp.Err == nil {
		h.JWT = tfaResp.Data.JWT
		return "validateTFA success", nil
	} else {
		h.JWT = ""
		return *tfaResp.Err, fmt.Errorf("validateTFA failed")
	}
}

func (h *SafelineHandler) Login(account *data.Account) (string, error) {
	h.account = *account
	csrf_token, err := h.getCSRFToken()
	if err != nil {
		return csrf_token, err
	}

	payload := map[string]string{
		"username":   h.account.GetUsername(),
		"password":   h.account.GetPassword(),
		"csrf_token": csrf_token,
	}

	body, _ := json.Marshal(payload)
	resp, err := h.client.Post(fmt.Sprintf("%s%s", h.account.GetSafelineURL(), cfg.SafelineConfig.API_Login_Path), "application/json", bytes.NewBuffer(body))
	if err != nil {
		file, line := utils.GetErrorLocation()
		h.logger.Debug(fmt.Sprintf("Error in file %s, line %d: %v", file, line, err))
		return "unable to get loginResp", err
	}
	defer resp.Body.Close()

	var loginResp LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		file, line := utils.GetErrorLocation()
		h.logger.Debug(fmt.Sprintf("Error in file %s, line %d: %v", file, line, err))
		return "unable to Parse loginResp", err
	}

	h.JWT = loginResp.Data.JWT

	if h.account.IsEnableOTP() {
		code, err := h.account.GetTOTPcode()
		if err != nil {
			return code, err
		}
		msg, err := h.validateTFA(code)
		if err != nil {
			return msg, err
		}
	}
	return "login success", nil
}

func (h *SafelineHandler) GetCertList() (string, error) {
	// Get CertListResponse from Safeline
	headers := map[string]string{
		"authorization": "Bearer " + h.JWT,
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s%s", h.account.GetSafelineURL(), cfg.SafelineConfig.API_Cert_Path), nil)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, err := h.client.Do(req)
	if err != nil {
		file, line := utils.GetErrorLocation()
		h.logger.Debug(fmt.Sprintf("Error in file %s, line %d: %v", file, line, err))
		return "unable to get certListResp", err
	}
	defer resp.Body.Close()

	// Parse CertListResponse
	var certListResp CertListResponse
	if err := json.NewDecoder(resp.Body).Decode(&certListResp); err != nil {
		file, line := utils.GetErrorLocation()
		h.logger.Debug(fmt.Sprintf("Error in file %s, line %d: %v", file, line, err))
		return "unable to Parse certListResp", err
	}

	// Get all nodes from CertListResponse
	nodes := ""
	for _, cert := range certListResp.Data.Nodes {
		certBytes, err := json.MarshalIndent(cert, "", "  ")
		if err != nil {
			file, line := utils.GetErrorLocation()
			h.logger.Debug(fmt.Sprintf("Error in file %s, line %d: %v", file, line, err))
			return "unable to Marshall cert", err
		}
		nodes += string(certBytes)
		nodes += ",\n"
	}

	return nodes, nil
}

func (h *SafelineHandler) UpdateCert(certPath *data.CertPath) (string, error) {
	// Get UpdateCertResponse from Safeline
	headers := map[string]string{
		"authorization": "Bearer " + h.JWT,
	}

	payload := map[string]interface{}{
		"id": certPath.GetID(),
		"manual": map[string]string{
			"crt": certPath.GetCertStr(),
			"key": certPath.GetKeyStr(),
		},
		"type": 2,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s%s/%s", h.account.GetSafelineURL(), cfg.SafelineConfig.API_Cert_Path, certPath.GetID()), bytes.NewBuffer(body))
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, err := h.client.Do(req)
	if err != nil {
		file, line := utils.GetErrorLocation()
		h.logger.Debug(fmt.Sprintf("Error in file %s, line %d: %v", file, line, err))
		return "unable to get updateCertResp", err
	}
	defer resp.Body.Close()

	// Parse updateCertResp
	var updateCertResp CertResponse
	if err := json.NewDecoder(resp.Body).Decode(&updateCertResp); err != nil {
		file, line := utils.GetErrorLocation()
		h.logger.Debug(fmt.Sprintf("Error in file %s, line %d: %v", file, line, err))
		return "unable to Parse updateCertResp", err
	}

	if updateCertResp.Err == nil {
		return "update cert success", nil
	} else {
		return *updateCertResp.Err, fmt.Errorf("update cert failed")
	}
}

func (h *SafelineHandler) DeleteCert(certPath *data.CertPath) (string, error) {
	// Get UpdateCertResponse from Safeline
	headers := map[string]string{
		"authorization": "Bearer " + h.JWT,
	}

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s%s/%s", h.account.GetSafelineURL(), cfg.SafelineConfig.API_Cert_Path, certPath.GetID()), nil)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, err := h.client.Do(req)
	if err != nil {
		file, line := utils.GetErrorLocation()
		h.logger.Debug(fmt.Sprintf("Error in file %s, line %d: %v", file, line, err))
		return "unable to get deleteCertResp", err
	}
	defer resp.Body.Close()

	// Parse updateCertResp
	var deleteCertResp CertResponse
	if err := json.NewDecoder(resp.Body).Decode(&deleteCertResp); err != nil {
		file, line := utils.GetErrorLocation()
		h.logger.Debug(fmt.Sprintf("Error in file %s, line %d: %v", file, line, err))
		return "unable to Parse deleteCertResp", err
	}

	if deleteCertResp.Err == nil {
		return "delete cert success", nil
	} else {
		return *deleteCertResp.Err, fmt.Errorf("delete cert failed")
	}
}

func (h *SafelineHandler) CreateCert(certPath *data.CertPath) (string, error) {
	// Get CreateCertResponse from Safeline
	headers := map[string]string{
		"authorization": "Bearer " + h.JWT,
	}

	payload := map[string]interface{}{
		"manual": map[string]string{
			"crt": certPath.GetCertStr(),
			"key": certPath.GetKeyStr(),
		},
		"type": 2,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s%s", h.account.GetSafelineURL(), cfg.SafelineConfig.API_Cert_Path), bytes.NewBuffer(body))
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, err := h.client.Do(req)
	if err != nil {
		file, line := utils.GetErrorLocation()
		h.logger.Debug(fmt.Sprintf("Error in file %s, line %d: %v", file, line, err))
		return "unable to get createCertResp", err
	}
	defer resp.Body.Close()

	// Parse createCertResp
	var createCertResp CertResponse
	if err := json.NewDecoder(resp.Body).Decode(&createCertResp); err != nil {
		file, line := utils.GetErrorLocation()
		h.logger.Debug(fmt.Sprintf("Error in file %s, line %d: %v", file, line, err))
		return "unable to Parse createCertResp", err
	}

	if createCertResp.Err == nil {
		return "create cert success", nil
	} else {
		return *createCertResp.Err, fmt.Errorf("create cert failed")
	}
}
