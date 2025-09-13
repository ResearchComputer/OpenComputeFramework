package cmd

import (
	"fmt"
	"ocf/internal/wallet"

	"github.com/spf13/cobra"
)

var walletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Wallet management commands",
}

var walletCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new wallet",
	Run: func(cmd *cobra.Command, args []string) {
		wm := wallet.NewWalletManager()
		if wm.WalletExists() {
			fmt.Printf("Wallet already exists at %s\n", wm.GetWalletPath())
			return
		}

		if err := wm.CreateWallet(); err != nil {
			fmt.Printf("Failed to create wallet: %v\n", err)
			return
		}
		fmt.Printf("Wallet created successfully at %s\n", wm.GetWalletPath())
		fmt.Printf("Public Key: %s\n", wm.GetPublicKey())
	},
}

var walletLoadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load an existing wallet",
	Run: func(cmd *cobra.Command, args []string) {
		wm := wallet.NewWalletManager()
		if !wm.WalletExists() {
			fmt.Printf("Wallet not found at %s\n", wm.GetWalletPath())
			return
		}

		if err := wm.LoadWallet(); err != nil {
			fmt.Printf("Failed to load wallet: %v\n", err)
			return
		}
		fmt.Printf("Wallet loaded successfully\n")
		fmt.Printf("Public Key: %s\n", wm.GetPublicKey())
	},
}

var walletInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show wallet information",
	Run: func(cmd *cobra.Command, args []string) {
		wm := wallet.NewWalletManager()
		if !wm.WalletExists() {
			fmt.Printf("Wallet not found at %s\n", wm.GetWalletPath())
			return
		}

		if err := wm.LoadWallet(); err != nil {
			fmt.Printf("Failed to load wallet: %v\n", err)
			return
		}
		fmt.Printf("Wallet Path: %s\n", wm.GetWalletPath())
		fmt.Printf("Public Key: %s\n", wm.GetPublicKey())
	},
}

func init() {
	walletCmd.AddCommand(walletCreateCmd)
	walletCmd.AddCommand(walletLoadCmd)
	walletCmd.AddCommand(walletInfoCmd)
	rootcmd.AddCommand(walletCmd)
}