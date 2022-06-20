/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"crypto/tls"
	"os"

	"github.com/jlaffaye/ftp"
	"github.com/spf13/cobra"
)

var sendCommandOptions struct {
	LocalFilePath  string `validate:"required"`
	ServerFilePath string `validate:"required"`
	ServerAddress  string `validate:"required"`
	User           string `validate:"required"`
	Password       string `validate:"required"`
}

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "send file to ftpserver",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := ftp.Dial(sendCommandOptions.ServerAddress, ftp.DialWithExplicitTLS(&tls.Config{
			InsecureSkipVerify: true,
		}), ftp.DialWithDebugOutput(os.Stdout), ftp.DialWithDisabledEPSV(true))
		if err != nil {
			panic(err)
		}
		defer client.Quit()

		data, err := os.Open(sendCommandOptions.LocalFilePath)
		if err != nil {
			panic(err)
		}

		err = client.Login(sendCommandOptions.User, sendCommandOptions.Password)
		if err != nil {
			panic(err)
		}
		defer client.Logout()

		err = client.Stor(sendCommandOptions.ServerFilePath, data)
		if err != nil {
			panic(err)
		}
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateParams(sendCommandOptions)
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
	sendCmd.PersistentFlags().StringVarP(&sendCommandOptions.LocalFilePath, "file", "f", "", "send file path")
	sendCmd.PersistentFlags().StringVarP(&sendCommandOptions.ServerFilePath, "spath", "s", "", "server file path")
	sendCmd.PersistentFlags().StringVarP(&sendCommandOptions.ServerAddress, "addr", "a", "", "ftp server address")
	sendCmd.PersistentFlags().StringVarP(&sendCommandOptions.User, "user", "u", "", "ftp user")
	sendCmd.PersistentFlags().StringVarP(&sendCommandOptions.Password, "pass", "p", "", "ftp password")
}
