package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

const (
	chaincodeName string = "transactionnews"
	channelId     string = "gctchain"
)

type TransactionSuccess struct {
	BlockId string `json:"blockId"`
}

func (setup *OrgSetup) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("Received Invoke request")
	if err := r.ParseForm(); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	function := "CreateTransaction"

	defer r.Body.Close()

	var transaction TransactionInitiate
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json, _ := json.Marshal(&transaction)
	fmt.Println(string(json))
	res, err := SubmitTransaction(setup, function, json)

	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeBlockSuccess(w, res)

}

func (setup *OrgSetup) ReadAllTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	function := "GetAllTransactions"

	res, err := EvaluateTransaction(setup, function)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set(applicationJson())
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "%s", res)
}

func (setup *OrgSetup) UpdateTransactionProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("Received Invoke request")
	if err := r.ParseForm(); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	function := "UpdateTransactionProcess"

	var transactionProcess UpdateTransactionProcess
	err := json.NewDecoder(r.Body).Decode(&transactionProcess)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}
	json, _ := json.Marshal(&transactionProcess)
	fmt.Printf("\n ---- UpdateProcess JSON Received -------")
	fmt.Println(string(json))
	res, err := SubmitTransaction(setup, function, json)

	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeBlockSuccess(w, res)

}

func (setup *OrgSetup) GetTransactionHistoryById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	queryParams := r.URL.Query()

	id := queryParams.Get("id")
	function := "GetTransactionHistoryById"

	res, err := EvaluateTransaction(setup, function, id)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set(applicationJson())
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", res)
}

func (setup *OrgSetup) GetTransactionById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	queryParams := r.URL.Query()

	id := queryParams.Get("id")
	function := "ReadTransactionById"

	res, err := EvaluateTransaction(setup, function, id)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set(applicationJson())
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", res)

}

func (setup *OrgSetup) AddReceiptPayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("Received Invoke request")
	if err := r.ParseForm(); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	function := "AddRecipientPayment"

	defer r.Body.Close()

	var groupPayment RecipientGroupPayment
	err := json.NewDecoder(r.Body).Decode(&groupPayment)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json, _ := json.Marshal(&groupPayment)
	fmt.Printf("\n ---- UpdateProcess JSON Received -------")
	fmt.Println(string(json))
	res, err := SubmitTransaction(setup, function, json)

	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeBlockSuccess(w, res)

}

func (setup *OrgSetup) AddMemberPayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("Received  request")
	if err := r.ParseForm(); err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	function := "AddMemberPayment"
	var memberTransaction UpdateMemberTransaction
	err := json.NewDecoder(r.Body).Decode(&memberTransaction)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json, _ := json.Marshal(&memberTransaction)
	fmt.Printf("\n ---- UpdateProcess JSON Received -------")
	fmt.Println(string(json))

	res, err := SubmitTransaction(setup, function, json)

	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeBlockSuccess(w, res)

}

func SubmitTransaction(setup *OrgSetup, funcName string, data []byte) (string, error) {
	fmt.Printf("channel: %s, chaincode: %s, function: %s\n", channelId, chaincodeName, funcName)
	network := setup.Gateway.GetNetwork(channelId)
	contract := network.GetContract(chaincodeName)
	txn_proposal, err := contract.NewProposal(funcName, client.WithArguments(string(data)))
	if err != nil {
		return "", fmt.Errorf("error creating txn proposal: %s", err)

	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		return "", fmt.Errorf("error endorsing txn: %s", err)

	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		return "", fmt.Errorf("error submitting transaction: %s", err)

	}

	return txn_committed.TransactionID(), nil

}

func EvaluateTransaction(setup *OrgSetup, funcName string, arg ...string) ([]byte, error) {
	fmt.Printf("channel: %s, chaincode: %s, function: %s\n", channelId, chaincodeName, funcName)
	network := setup.Gateway.GetNetwork(channelId)
	contract := network.GetContract(chaincodeName)
	evaluateResponse, err := contract.EvaluateTransaction(funcName, arg...)
	if err != nil {
		return nil, err
	}
	return evaluateResponse, nil
}

func applicationJson() (key string, value string) {
	return "content-type", "application/json"
}

func writeError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set(applicationJson())
	w.WriteHeader(statusCode)

	errorResponse := ErrorResponse{
		Error: message,
	}

	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		writeError(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func writeBlockSuccess(w http.ResponseWriter, tnxID string) {
	w.Header().Set(applicationJson())
	w.WriteHeader(http.StatusOK)

	transactionSuccess := TransactionSuccess{
		BlockId: tnxID,
	}

	if err := json.NewEncoder(w).Encode(transactionSuccess); err != nil {
		writeError(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
