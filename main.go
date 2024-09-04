package main

import (
	"fmt"
	"log"
	"os"
	"rest-api-go/web"

	"github.com/joho/godotenv"
)

func main() {
	//Initialize setup for Org1
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the domain from the environment variable
	domain := os.Getenv("DOMAIN")
	if domain == "" {
		log.Fatal("DOMAIN environment variable is not set in the .env file")
	}
	cryptoPath := fmt.Sprintf("../../test-network/organizations/peerOrganizations/org1.%s", domain)

	orgConfig := web.OrgSetup{
		OrgName:      "Org1",
		MSPID:        "Org1MSP",
		CertPath:     fmt.Sprintf("%s/users/User1@org1.%s/msp/signcerts/cert.pem", cryptoPath, domain),
		KeyPath:      fmt.Sprintf("%s/users/User1@org1.%s/msp/keystore/", cryptoPath, domain),
		TLSCertPath:  fmt.Sprintf("%s/peers/peer0.org1.%s/tls/ca.crt", cryptoPath, domain),
		PeerEndpoint: "dns:///localhost:7051",
		GatewayPeer:  fmt.Sprintf("peer0.org1.%s", domain),
	}

	orgSetup, err := web.Initialize(orgConfig)
	if err != nil {
		fmt.Println("Error initializing setup for Org1: ", err)
	}
	web.Serve(web.OrgSetup(*orgSetup))
}
