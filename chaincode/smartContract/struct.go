package sc

import (
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type CCC struct {
}

type CipherIndex struct{
    PkEnc string `json:"pkEnc"`
	PkUser string `json:"pkUser"`
    RoadName string `json:"roadName"`
	MessageLevel    string    `json:"messageLevel"`
    Weather string `json:"weather"`
    Condition string `json:"condition"`
    Traffic string `json:"traffic"`
    AverageSpeed string `json:"avarageSpeed"`
}

type Vehicle struct {
	PkUser string `json:"pkUser"`
	CoinNum string `json:"coinNum"`
    Credit string `json:"credit"`
	Role string `json:"role"`
	PingCount int64 `json:"pingCount"`
	LastPing  int64 `json:"lastPing"`
	Activity float64 `json:"avtivity"`
}

func (s *CCC) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success([]byte("success"))
}

func (s *CCC) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	function, args := stub.GetFunctionAndParameters()

	// coin/credit ledger
	if function == "AddVehicle" {    //传参初始化
		return s.AddVehicle(stub, args)
	} else if function == "Ping" {     //交易后ping一下更新活跃度
        return s.Ping(stub, args) 
	} else if function == "UpdateCredit" { //ping完再update更新信誉度
		return s.UpdateCredit(stub, args)
	} else if function == "UpdateCoinNum" {
		return s.UpdateCoinNum(stub, args)
	} else if function == "DeleteItem" {
		return s.DeleteItem(stub, args)
	} else if function == "QueryCoinNum" {
		return s.QueryCoinNum(stub, args)
    } else if function == "QueryCredit" {
        return s.QueryCredit(stub, args)
    } else if function == "ChangeRole" {
        return s.ChangeRole(stub, args)
    }
	
	// cipherIndex ledger
	if function == "AppendNewMessage" {
		return s.AppendNewMessage(stub, args)
	} else if function == "DeleteMessage" {
		return s.DeleteMessage(stub, args)
	} else if function == "GetMessageByRoadName" {
		return s.GetMessageByRoadName(stub, args)
    }


	return shim.Error(fmt.Sprintf("Invalid Chaincode function name: %s", function))
}
