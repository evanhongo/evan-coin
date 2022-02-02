package cmd

import (
	"fmt"
	"time"

	"github.com/evanhongo/blockchain-demo/pkg/socket"

	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type Message struct {
	Id   int    `json:"id"`
	Text string `json:"text"`
}

var rootCmd = &cobra.Command{
	Use:    "evan-coin",
	Short:  "Evan-coin client",
	Long:   `Welcome to Evan Coin~`,
	PreRun: func(cmd *cobra.Command, args []string) {},
	Run: func(cmd *cobra.Command, args []string) {
		c := socket.InitSocketClient()
		defer c.Close()
		for {
			var address string
			fmt.Println("Please provide your wallet address:")
			fmt.Scanln(&address)

			result, err := c.Ack("get-balance", address, time.Second*5)
			if err != nil {
				logger.Fatal(err)
			} else {
				logger.Println("result", result)
			}
		}
	},
}

func init() {
	// rootCmd.AddCommand(getBalanceCmd)
	// getBalanceCmd.Flags().StringVarP(&addr, "addr", "a", "", "address")
	// getBalanceCmd.MarkFlagRequired("addr")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Println(err)
	}
}
