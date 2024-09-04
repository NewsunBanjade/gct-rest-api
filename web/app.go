package web

import (
	"fmt"
	"net/http"
	"os"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// OrgSetup contains organization's config to interact with the network.
type OrgSetup struct {
	OrgName      string
	MSPID        string
	CryptoPath   string
	CertPath     string
	KeyPath      string
	TLSCertPath  string
	PeerEndpoint string
	GatewayPeer  string
	Gateway      client.Gateway
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func apiKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("Authorization")
		if apiKey != os.Getenv("API_KEY") {
			writeError(w, "Forbidden: Not Authorized", http.StatusUnauthorized)
			return

		}
		next.ServeHTTP(w, r)
	})
}

func Serve(setups OrgSetup) {

	// Initialize your handlers
	http.Handle("/createTransaction", apiKeyMiddleware(http.HandlerFunc(setups.CreateTransaction)))
	http.Handle("/transactions", apiKeyMiddleware(http.HandlerFunc(setups.ReadAllTransactions)))
	http.Handle("/updateTransactionProcess", apiKeyMiddleware(http.HandlerFunc(setups.UpdateTransactionProcess)))
	http.Handle("/transaction", apiKeyMiddleware(http.HandlerFunc(setups.GetTransactionById)))
	http.Handle("/transactionHistory", apiKeyMiddleware(http.HandlerFunc(setups.GetTransactionHistoryById)))
	http.Handle("/addReceiptPayment", apiKeyMiddleware(http.HandlerFunc(setups.AddReceiptPayment)))
	http.Handle("/addMemberPayment", apiKeyMiddleware(http.HandlerFunc(setups.AddMemberPayment)))

	//Server
	fmt.Println("Listening (http://localhost:3264/)...")
	if err := http.ListenAndServe(":3264", nil); err != nil {
		fmt.Println(err)
	}
}
