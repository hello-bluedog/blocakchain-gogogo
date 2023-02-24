package sc

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

func (s *CCC) AddNewUserItem(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	pkUser := args[0]
    coinNum := args[1]
    credit := args[2]

	IsExist := s.FeelExist(stub, pkUser)
	if (IsExist == "Yes") {
		return shim.Error("user has existed, please use the modify function")
	}

	newItem := CoinAndCredit{
        PkUser: pkUser,
        CoinNum: coinNum,
        Credit: credit,
	}
	itembytes, err := json.Marshal(newItem)
	if err != nil {
		return shim.Error("json marshal shim.Error")
	}
	err = stub.PutState(pkUser, itembytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("PutState shim.Error: %s", err))
	}

	return shim.Success([]byte("Add new user success"))
}
func (s *CCC) UpdateCredit(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	pkUser := args[0]
	credit := args[1]

    itembytes, err := stub.GetState(pkUser)
    if(len(itembytes) == 0){
        return shim.Error("no such user in ledger!!!")
    }
    var tmp CoinAndCredit
    _ = json.Unmarshal(itembytes, &tmp)
    tmp.Credit = credit
    itembytes, err = json.Marshal(tmp)
	if err != nil {
		return shim.Error("json marshal shim.Error")
	}

	err = stub.PutState(pkUser, itembytes)

	return shim.Success([]byte("update credit success"))
}

func (s *CCC) UpdateCoinNum(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	pkUser := args[0]
	coinNum := args[1]

    itembytes, _ := stub.GetState(pkUser)
    if(len(itembytes) == 0){
        return shim.Error("no such user in ledger!!!")
    }
    var tmp CoinAndCredit
    _ = json.Unmarshal(itembytes, &tmp)
    tmp.CoinNum = coinNum
    itembytes, err := json.Marshal(tmp)
	if err != nil {
		return shim.Error("json marshal shim.Error")
	}

	err = stub.PutState(pkUser, itembytes)

	return shim.Success([]byte("update credit success"))
}

func (s *CCC) DeleteItem(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	// check if exists
	pkUser := args[0]
	itembytes, err := stub.GetState(pkUser)
	if (err != nil || len(itembytes) == 0) {
		return shim.Error("No such User!")
	}
	//do deletion
	err = stub.DelState(pkUser)
	if err != nil {
		return shim.Error("Delete State Error")
	}
	return shim.Success([]byte("Delete User Success"))
}

func (s * CCC) QueryCoinNum(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
    pkUser := args[0]
    itemBytes, _ := stub.GetState(pkUser)
    if(len(itemBytes) == 0) {
        return shim.Error("no such user in ledger!")
    }
    var coinAndCredit CoinAndCredit
    _ = json.Unmarshal(itemBytes, &coinAndCredit)
    return shim.Success([]byte(coinAndCredit.CoinNum))
}

func (s * CCC) QueryCredit(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
    pkUser := args[0]
    itemBytes, _ := stub.GetState(pkUser)
    if(len(itemBytes) == 0) {
        return shim.Error("no such user in ledger!")
    }
    var coinAndCredit CoinAndCredit
    _ = json.Unmarshal(itemBytes, &coinAndCredit)
    return shim.Success([]byte(coinAndCredit.Credit))
}

func (s * CCC)ChangeRole(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
    pkUser := args[0]
    role := args[1]
    itembytes, err := stub.GetState(pkUser)
    if(len(itembytes) == 0){
        return shim.Error("no such user in ledger!!!")
    }
    var tmp CoinAndCredit
    _ = json.Unmarshal(itembytes, &tmp)
    tmp.Role = role
    itembytes, err = json.Marshal(tmp)
	if err != nil {
		return shim.Error("json marshal shim.Error")
    }
	err = stub.PutState(pkUser, itembytes)

	return shim.Success([]byte("update role success"))
}

/*func (s * CCC)FeelExist(stub shim.ChaincodeStubInterface, target string) string {
	ResBytes, _ := stub.GetState(target)
	if(len(ResBytes) == 0) {
		return "No"
	} else  {
		return "Yes"
	}
}*/
