package wallet

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"ocf/internal/common"
)

type WalletKey struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

type WalletManager struct {
	walletKey  *WalletKey
	walletPath string
}

func NewWalletManager() *WalletManager {
	walletPath := viper.GetString("account.wallet")
	if walletPath == "" {
		homeDir, _ := os.UserHomeDir()
		walletPath = filepath.Join(homeDir, ".ocf", "wallet.json")
	}

	return &WalletManager{
		walletPath: walletPath,
	}
}

func (wm *WalletManager) CreateWallet() error {
	common.Logger.Info("Creating new wallet...")

	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return fmt.Errorf("failed to generate keypair: %w", err)
	}

	wm.walletKey = &WalletKey{
		PublicKey:  base64.StdEncoding.EncodeToString(pubKey),
		PrivateKey: base64.StdEncoding.EncodeToString(privKey),
	}

	if err := wm.saveWallet(); err != nil {
		return fmt.Errorf("failed to save wallet: %w", err)
	}

	common.Logger.Infof("Wallet created successfully. Public key: %s", wm.walletKey.PublicKey)
	return nil
}

func (wm *WalletManager) LoadWallet() error {
	if _, err := os.Stat(wm.walletPath); os.IsNotExist(err) {
		return fmt.Errorf("wallet file not found at %s", wm.walletPath)
	}

	data, err := os.ReadFile(wm.walletPath)
	if err != nil {
		return fmt.Errorf("failed to read wallet file: %w", err)
	}

	// For simplicity, we'll just store the private key as base64
	privateKeyStr := string(data)
	privateKeyBytes, err := base64.StdEncoding.DecodeString(privateKeyStr)
	if err != nil {
		return fmt.Errorf("failed to decode private key: %w", err)
	}

	if len(privateKeyBytes) != ed25519.PrivateKeySize {
		return fmt.Errorf("invalid private key size")
	}

	privKey := ed25519.PrivateKey(privateKeyBytes)
	pubKey := privKey.Public().(ed25519.PublicKey)

	wm.walletKey = &WalletKey{
		PublicKey:  base64.StdEncoding.EncodeToString(pubKey),
		PrivateKey: base64.StdEncoding.EncodeToString(privKey),
	}

	common.Logger.Infof("Wallet loaded successfully. Public key: %s", wm.walletKey.PublicKey)
	return nil
}

func (wm *WalletManager) saveWallet() error {
	if wm.walletKey == nil {
		return fmt.Errorf("no wallet to save")
	}

	walletDir := filepath.Dir(wm.walletPath)
	if err := os.MkdirAll(walletDir, 0700); err != nil {
		return fmt.Errorf("failed to create wallet directory: %w", err)
	}

	if err := os.WriteFile(wm.walletPath, []byte(wm.walletKey.PrivateKey), 0600); err != nil {
		return fmt.Errorf("failed to write wallet file: %w", err)
	}

	return nil
}

func (wm *WalletManager) GetPublicKey() string {
	if wm.walletKey == nil {
		return ""
	}
	return wm.walletKey.PublicKey
}

func (wm *WalletManager) GetPrivateKey() string {
	if wm.walletKey == nil {
		return ""
	}
	return wm.walletKey.PrivateKey
}

func (wm *WalletManager) GetWalletPath() string {
	return wm.walletPath
}

func (wm *WalletManager) WalletExists() bool {
	_, err := os.Stat(wm.walletPath)
	return !os.IsNotExist(err)
}

func (wm *WalletManager) Initialize() error {
	if wm.WalletExists() {
		if err := wm.LoadWallet(); err != nil {
			common.Logger.Warnf("Failed to load existing wallet: %v", err)
			if err := wm.CreateWallet(); err != nil {
				return fmt.Errorf("failed to create new wallet after load failure: %w", err)
			}
		}
	} else {
		if err := wm.CreateWallet(); err != nil {
			return fmt.Errorf("failed to create new wallet: %w", err)
		}
	}
	return nil
}

func InitializeWallet() (*WalletManager, error) {
	wm := NewWalletManager()
	if err := wm.Initialize(); err != nil {
		return nil, err
	}
	return wm, nil
}