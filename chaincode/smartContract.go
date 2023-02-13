package main

import (
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric/scripts/fabric-samples/dataSharing/smartContract"
)


func main() {
	err := shim.Start(new(sc.CCC))
	if err != nil {
		fmt.Sprintf("Error creating new Chaincode: %s", err)
	}
}
