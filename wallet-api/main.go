package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"github.com/google/uuid"
	"github.com/tyler-smith/go-bip39"
	"github.com/tyler-smith/go-bip32"
	"strconv"
)

type Wallet struct {
	ID      string `json:"id"`
	Address string `json:"address"`
}

var (
	wallets = make(map[string]Wallet)
	mu      sync.Mutex
)

func createWallet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	mu.Lock()
	defer mu.Unlock()

	id := uuid.New().String()
	wallet := Wallet{
		ID:      id,
		Address: "0x" + id[:8], // fake address for now
	}

	wallets[id] = wallet

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wallet)
}

func getWallet(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	mu.Lock()
	defer mu.Unlock()

	wallet, exists := wallets[id]
	if !exists {
		http.Error(w, "Wallet not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(wallet)
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func generateMnemonic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		http.Error(w, "Failed to generate entropy", http.StatusInternalServerError)
		return
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		http.Error(w, "Failed to generate mnemonic", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"mnemonic": mnemonic,
	})
}

func deriveWallet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	mnemonic := r.URL.Query().Get("mnemonic")
	indexStr := r.URL.Query().Get("index")

	if mnemonic == "" || indexStr == "" {
		http.Error(w, "Missing mnemonic or index", http.StatusBadRequest)
		return
	}

	if !bip39.IsMnemonicValid(mnemonic) {
		http.Error(w, "Invalid mnemonic", http.StatusBadRequest)
		return
	}

	index, err := strconv.Atoi(indexStr)
	if err != nil {
		http.Error(w, "Invalid index", http.StatusBadRequest)
		return
	}

	seed := bip39.NewSeed(mnemonic, "")

	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		http.Error(w, "Failed to create master key", http.StatusInternalServerError)
		return
	}

	childKey, err := masterKey.NewChildKey(uint32(index))
	if err != nil {
		http.Error(w, "Failed to derive child key", http.StatusInternalServerError)
		return
	}

	address := childKey.PublicKey().String()

	json.NewEncoder(w).Encode(map[string]string{
		"address": address,
	})
}

func main() {
	http.HandleFunc("/wallet", createWallet)
	http.HandleFunc("/get", getWallet)
	http.HandleFunc("/health", health)
	http.HandleFunc("/mnemonic", generateMnemonic)
	http.HandleFunc("/derive", deriveWallet)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}