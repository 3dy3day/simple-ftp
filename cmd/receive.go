/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"crypto/tls"
	"io/ioutil"
	"os"

	"github.com/jlaffaye/ftp"
	"github.com/spf13/cobra"
)

var receiveCommandOptions struct {
	LocalFilePath  string `validate:"required"`
	ServerFilePath string `validate:"required"`
	ServerAddress  string `validate:"required"`
	User           string `validate:"required"`
	Password       string `validate:"required"`
}

// receiveCmd represents the receive command
var receiveCmd = &cobra.Command{
	Use:   "receive",
	Short: "receive file from ftpserver",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := ftp.Dial(receiveCommandOptions.ServerAddress, ftp.DialWithExplicitTLS(&tls.Config{
			InsecureSkipVerify: true,
		}), ftp.DialWithDebugOutput(os.Stdout), ftp.DialWithDisabledEPSV(true))
		if err != nil {
			panic(err)
		}
		defer client.Quit()

		err = client.Login(receiveCommandOptions.User, receiveCommandOptions.Password)
		if err != nil {
			panic(err)
		}
		defer client.Logout()

		res, err := client.Retr(receiveCommandOptions.ServerFilePath)
		if err != nil {
			panic(err)
		}
		defer res.Close()
		data, err := ioutil.ReadAll(res)
		if err != nil {
			panic(err)
		}
		ioutil.WriteFile(receiveCommandOptions.LocalFilePath, data, 0777)
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateParams(receiveCommandOptions)
	},
}

func init() {
	rootCmd.AddCommand(receiveCmd)
	receiveCmd.PersistentFlags().StringVarP(&receiveCommandOptions.LocalFilePath, "file", "f", "", "receive file path")
	receiveCmd.PersistentFlags().StringVarP(&receiveCommandOptions.ServerFilePath, "spath", "s", "", "server file path")
	receiveCmd.PersistentFlags().StringVarP(&receiveCommandOptions.ServerAddress, "addr", "a", "", "ftp server address")
	receiveCmd.PersistentFlags().StringVarP(&receiveCommandOptions.User, "user", "u", "", "ftp user")
	receiveCmd.PersistentFlags().StringVarP(&receiveCommandOptions.Password, "pass", "p", "", "ftp password")
}
