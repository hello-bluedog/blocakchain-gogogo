package sc

import (
	"encoding/json"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)
//访问控制的信誉度标准
var rank float64 =0.3

func(s *CCC) changeStrategy(stub shim.ChaincodeStubInterface, args []string) peer.Response{
	pkUser := args[0]
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	//管理员权限检查
	if pkUser !="admin" {
		return shim.Error("Do not have the authority")
	}
	temp, err := strconv.ParseFloat(args[1], 64); 
	if err != nil {
		return shim.Error("Type miss matched. Expecting number ")
	}
	rank =float64(temp)
	return shim.Success([]byte("update credit success"))
}

//一次的通信
func (s *CCC) recordCommunication(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
    if len(args) != 3 {
        return shim.Error("Incorrect number of arguments. Expecting 3")
    }
    vehicleID := args[0]
    serverID := args[1]
    success := args[2]

    vehicleBytes, err := APIstub.GetState(vehicleID)
    if err != nil {
        return shim.Error(err.Error())
    }
    if vehicleBytes == nil {
        return shim.Error("Vehicle not found")
    }
    var vehicle CoinAndCredit
    err = json.Unmarshal(vehicleBytes, &vehicle)
    if err != nil {
        return shim.Error(err.Error())
    }

    serverBytes, err := APIstub.GetState(serverID)
    if err != nil {
        return shim.Error(err.Error())
    }
    if serverBytes == nil {
        return shim.Error("Server not found")
    }
    var server Server
    err = json.Unmarshal(serverBytes, &server)
    if err != nil {
        return shim.Error(err.Error())
    }
	//通讯成功
    if success == "true" {
        vehicle.SuccessNum++
        server.SuccessNum++
    }

    vehicleBytes, err = json.Marshal(vehicle)
    if err != nil {
        return shim.Error(err.Error())
    }
    err = APIstub.PutState(vehicleID, vehicleBytes)
    if err != nil {
        return shim.Error(err.Error())
    }

    serverBytes, err = json.Marshal(server)
    if err != nil {
        return shim.Error(err.Error())
    }
    err = APIstub.PutState(serverID, serverBytes)
    if err != nil {
        return shim.Error(err.Error())
    }

    return shim.Success(nil)
}


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
	id:=args[1]

    itembytes, err := stub.GetState(pkUser)
    if(len(itembytes) == 0){
        return shim.Error("no such user in ledger!!!")
    }
    var tmp CoinAndCredit
    _ = json.Unmarshal(itembytes, &tmp)
	//历史信誉
	var his_credit float64
	his_credit,err = strconv.ParseFloat(tmp.Credit, 64)
	if err != nil {
		return shim.Error("json marshal shim.Error")
	}
	//通信信誉
	var commu_credit float64
	var server Server
	item_server, err := stub.GetState(id)
	_ = json.Unmarshal(item_server, &server)
	commu_credit= float64(tmp.SuccessNum)/float64(server.SuccessNum)
	//行为信誉


	var credit float64
	credit=his_credit + commu_credit
	
		//注销
	if(credit < rank) {
		tmp.IsRevoked =true
	}
	cre:=strconv.FormatFloat(credit,'f',-1,64)
    tmp.Credit = cre
    itembytes, err = json.Marshal(tmp)
	if err != nil {
		return shim.Error("json marshal shim.Error")
	}
	err = stub.PutState(pkUser, itembytes)
	return shim.Success([]byte("update credit success"))
}

func (s *CCC) UpdateCredit2(stub shim.ChaincodeStubInterface, pkUser string) peer.Response {
    itembytes, err := stub.GetState(pkUser)
    if(len(itembytes) == 0){
        return shim.Error("no such user in ledger!!!")
    }
    var tmp CoinAndCredit
	var cre1 float64
	var cre2 float64 
    var cre3 float64
	_ = json.Unmarshal(itembytes, &tmp)
    cre1,err = strconv.ParseFloat(tmp.Credit, 64)
	if err != nil {
		return shim.Error("json marshal shim.Error")
	}
	// cre1 = queryCredit()		历史信誉
	// cre2 = countCredit()  	通信信誉
	// cre3 = mutualCredit()   	交互信誉
	creNow:=0.5*cre1 + 0.3*cre2 + 0.2*cre3
		//注销
	if(creNow < rank) {
		tmp.IsRevoked =true
	}
	cre:=strconv.FormatFloat(creNow,'f',-1,64)
	tmp.Credit=cre
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
		return shim.Error("Incorrect number of arguments. Expecting 2")
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

/*func (s * CCC)FeelExist(stub shim.ChaincodeStubInterface, target string) string {
	ResBytes, _ := stub.GetState(target)
	if(len(ResBytes) == 0) {
		return "No"
	} else  {
		return "Yes"
	}
}*/
