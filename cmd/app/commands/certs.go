package commands

import (
	"fmt"

	cfg "github.com/DoraTiger/safeline-utils/config"
	"github.com/DoraTiger/safeline-utils/pkg/data"
	"github.com/DoraTiger/safeline-utils/pkg/handler"
	"github.com/spf13/cobra"
)

var (
	safelineURL string
	username    string
	password    string
	optSecret   string
	certID      string
	certFile    string
	keyFile     string
)

func init() {
	registerFlagsCertCmd(CertCmd)
}

func registerFlagsCertCmd(CertCmd *cobra.Command) {

	// listCertsCmd
	CertCmd.AddCommand(listCertsCmd)
	listCertsCmd.Flags().StringVarP(&safelineURL, "url", "u", cfg.SafelineConfig.SafelineURL, "Server URL")
	listCertsCmd.Flags().StringVarP(&username, "username", "n", "", "Username")
	listCertsCmd.Flags().StringVarP(&password, "password", "p", "", "Password")
	listCertsCmd.Flags().StringVarP(&optSecret, "otpSecret", "o", "", "OTPSecret")
	listCertsCmd.MarkFlagRequired("username")
	listCertsCmd.MarkFlagRequired("password")

	// updateCertCmd
	CertCmd.AddCommand(updateCertCmd)

	updateCertCmd.Flags().StringVarP(&safelineURL, "url", "u", cfg.SafelineConfig.SafelineURL, "Server URL")
	updateCertCmd.Flags().StringVarP(&username, "username", "n", "", "Username")
	updateCertCmd.Flags().StringVarP(&password, "password", "p", "", "Password")
	updateCertCmd.Flags().StringVarP(&optSecret, "otpSecret", "o", "", "OTPSecret")
	updateCertCmd.MarkFlagRequired("username")
	updateCertCmd.MarkFlagRequired("password")

	updateCertCmd.Flags().StringVarP(&certID, "cert-id", "i", "", "Certificate ID")
	updateCertCmd.Flags().StringVarP(&certFile, "cert-file", "c", "", "Certificate file")
	updateCertCmd.Flags().StringVarP(&keyFile, "key-file", "k", "", "Key file")
	updateCertCmd.MarkFlagRequired("cert-id")
	updateCertCmd.MarkFlagRequired("cert-file")
	updateCertCmd.MarkFlagRequired("key-file")

	// deleteCertCmd
	CertCmd.AddCommand(deleteCertCmd)
	deleteCertCmd.Flags().StringVarP(&safelineURL, "url", "u", cfg.SafelineConfig.SafelineURL, "Server URL")
	deleteCertCmd.Flags().StringVarP(&username, "username", "n", "", "Username")
	deleteCertCmd.Flags().StringVarP(&password, "password", "p", "", "Password")
	deleteCertCmd.Flags().StringVarP(&optSecret, "otpSecret", "o", "", "OTPSecret")
	deleteCertCmd.MarkFlagRequired("username")
	deleteCertCmd.MarkFlagRequired("password")
	deleteCertCmd.Flags().StringVarP(&certID, "cert-id", "i", "", "Certificate ID")
	deleteCertCmd.MarkFlagRequired("cert-id")

	// createCertCmd
	CertCmd.AddCommand(createCertCmd)

	createCertCmd.Flags().StringVarP(&safelineURL, "url", "u", cfg.SafelineConfig.SafelineURL, "Server URL")
	createCertCmd.Flags().StringVarP(&username, "username", "n", "", "Username")
	createCertCmd.Flags().StringVarP(&password, "password", "p", "", "Password")
	createCertCmd.Flags().StringVarP(&optSecret, "otpSecret", "o", "", "OTPSecret")
	createCertCmd.MarkFlagRequired("username")
	createCertCmd.MarkFlagRequired("password")

	createCertCmd.Flags().StringVarP(&certFile, "cert-file", "c", "", "Certificate file")
	createCertCmd.Flags().StringVarP(&keyFile, "key-file", "k", "", "Key file")
	createCertCmd.MarkFlagRequired("cert-file")
	createCertCmd.MarkFlagRequired("key-file")

}

var CertCmd = &cobra.Command{
	Use:   "cert",
	Short: "Manage certificates",
}

var listCertsCmd = &cobra.Command{
	Use:   "list",
	Short: "List certificates",
	Run: func(cmd *cobra.Command, args []string) {

		account := data.NewAccount(safelineURL, username, password, optSecret)
		h := handler.NewSafelineHandler()
		h.SetLogger(logger)
		msg, err := h.Login(account)
		if err != nil {
			fmt.Printf("Error login: %v\n", msg)
			return
		}
		msg, err = h.GetCertList()
		if err != nil {
			fmt.Printf("Error GetCertList: %v\n", msg)
			return
		}
		fmt.Println(msg)
	},
}

var updateCertCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a certificate",
	Run: func(cmd *cobra.Command, args []string) {
		account := data.NewAccount(safelineURL, username, password, optSecret)
		certPath := data.NewCertPath(certID, certFile, keyFile)

		h := handler.NewSafelineHandler()
		h.SetLogger(logger)
		msg, err := h.Login(account)
		if err != nil {
			fmt.Printf("Error login: %v\n", msg)
			return
		}
		msg, err = h.UpdateCert(certPath)
		if err != nil {
			fmt.Printf("Error UpdateCert: %v\n", msg)
			return
		}
		fmt.Println(msg)
	},
}

var deleteCertCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete a certificate",
	Run: func(cmd *cobra.Command, args []string) {
		account := data.NewAccount(safelineURL, username, password, optSecret)
		certPath := data.NewCertPath(certID, "", "")

		h := handler.NewSafelineHandler()
		h.SetLogger(logger)
		msg, err := h.Login(account)
		if err != nil {
			fmt.Printf("Error login: %v\n", msg)
			return
		}
		msg, err = h.DeleteCert(certPath)
		if err != nil {
			fmt.Printf("Error DeleteCert: %v\n", msg)
			return
		}
		fmt.Println(msg)
	},
}

var createCertCmd = &cobra.Command{
	Use:   "create",
	Short: "Update a certificate",
	Run: func(cmd *cobra.Command, args []string) {
		account := data.NewAccount(safelineURL, username, password, optSecret)
		certPath := data.NewCertPath("", certFile, keyFile)

		h := handler.NewSafelineHandler()
		h.SetLogger(logger)
		msg, err := h.Login(account)
		if err != nil {
			fmt.Printf("Error login: %v\n", msg)
			return
		}
		msg, err = h.CreateCert(certPath)
		if err != nil {
			fmt.Printf("Error CreateCert: %v\n", msg)
			return
		}
		fmt.Println(msg)
	},
}
