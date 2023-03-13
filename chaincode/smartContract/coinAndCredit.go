package sc

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

//访问控制的信誉度标准
// var rank float64 =0.3

func (s *CCC) AddVehicle(stub shim.ChaincodeStubInterface, args []string, vehicles map[string]*Vehicle) peer.Response {
	if len(args) != 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1")
    }

    pkUser := args[0]
    if _, exists := vehicles[pkUser]; exists {
        return shim.Error(fmt.Sprintf("Vehicle with ID %s already exists", pkUser))
    }

    vehicle := &Vehicle{
    	PkUser:   pkUser,
    	CoinNum:  "",
    	Credit:   0,
    	Role:     "",
    	LastPing: time.Now().Unix(),
    	Activity: 0,
    }

    // 将车辆信息加入map中
    vehicles[pkUser] = vehicle

    return shim.Success(nil)
}

func (s *CCC) Ping(stub shim.ChaincodeStubInterface, args []string, vehicles map[string]*Vehicle) peer.Response {
    if len(args) != 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1")
    }

    pkUser := args[0]
    vehicle, exists := vehicles[pkUser]
    if !exists {
        return shim.Error(fmt.Sprintf("Vehicle with pkUser %s does not exist", pkUser))
    }

   // 更新 LastPing 和 PingCount 字段
   now := time.Now().Unix()
   freq :=0.0
   if vehicle.LastPing > 0 {
	   timeDelta := now - vehicle.LastPing
	   freq = 60.0 / float64(timeDelta)
   }
   vehicle.PingCount++
   vehicle.LastPing = now

   // 计算活跃度
   vehicle.Activity = vehicle.PingCount * int64(freq)

   // 将更新后的车辆信息保存回map中
   vehicles[pkUser] = vehicle
   
   return shim.Success(nil)
}

// func(s *CCC) changeStrategy(stub shim.ChaincodeStubInterface, args []string) peer.Response{
// 	pkUser := args[0]
// 	if len(args) != 2 {
// 		return shim.Error("Incorrect number of arguments. Expecting 2")
// 	}
// 	//管理员权限检查
// 	if pkUser !="admin" {
// 		return shim.Error("Do not have the authority")
// 	}
// 	temp, err := strconv.ParseFloat(args[1], 64); 
// 	if err != nil {
// 		return shim.Error("Type miss matched. Expecting number ")
// 	}
// 	rank =float64(temp)
// 	return shim.Success([]byte("update credit success"))
// }

// //一次的通信
// func (s *CCC) recordCommunication(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
//     if len(args) != 3 {
//         return shim.Error("Incorrect number of arguments. Expecting 3")
//     }
//     vehicleID := args[0]
//     serverID := args[1]
//     success := args[2]

//     vehicleBytes, err := APIstub.GetState(vehicleID)
//     if err != nil {
//         return shim.Error(err.Error())
//     }
//     if vehicleBytes == nil {
//         return shim.Error("Vehicle not found")
//     }
//     var vehicle CoinAndCredit
//     err = json.Unmarshal(vehicleBytes, &vehicle)
//     if err != nil {
//         return shim.Error(err.Error())
//     }

//     serverBytes, err := APIstub.GetState(serverID)
//     if err != nil {
//         return shim.Error(err.Error())
//     }
//     if serverBytes == nil {
//         return shim.Error("Server not found")
//     }
//     var server Server
//     err = json.Unmarshal(serverBytes, &server)
//     if err != nil {
//         return shim.Error(err.Error())
//     }
// 	//通讯成功
//     if success == "true" {
//         vehicle.SuccessNum++
//         server.SuccessNum++
//     }

//     vehicleBytes, err = json.Marshal(vehicle)
//     if err != nil {
//         return shim.Error(err.Error())
//     }
//     err = APIstub.PutState(vehicleID, vehicleBytes)
//     if err != nil {
//         return shim.Error(err.Error())
//     }

//     serverBytes, err = json.Marshal(server)
//     if err != nil {
//         return shim.Error(err.Error())
//     }
//     err = APIstub.PutState(serverID, serverBytes)
//     if err != nil {
//         return shim.Error(err.Error())
//     }

//     return shim.Success(nil)
// }

	//替换为addvehicle
// func (s *CCC) AddNewUserItem(stub shim.ChaincodeStubInterface, args []string) peer.Response {
// 	if len(args) != 3 {
// 		return shim.Error("Incorrect number of arguments. Expecting 3")
// 	}
// 	pkUser := args[0]
//     coinNum := args[1]
//     credit := args[2]

// 	IsExist := s.FeelExist(stub, pkUser)
// 	if (IsExist == "Yes") {
// 		return shim.Error("user has existed, please use the modify function")
// 	}

// 	newItem := CoinAndCredit{
//         PkUser: pkUser,
//         CoinNum: coinNum,
//         Credit: credit,
// 	}
// 	itembytes, err := json.Marshal(newItem)
// 	if err != nil {
// 		return shim.Error("json marshal shim.Error")
// 	}
// 	err = stub.PutState(pkUser, itembytes)
// 	if err != nil {
// 		return shim.Error(fmt.Sprintf("PutState shim.Error: %s", err))
// 	}

// 	return shim.Success([]byte("Add new user success"))
// }

func (s *CCC) UpdateCredit(stub shim.ChaincodeStubInterface, args[]string) peer.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	pkUser := args[0]
	action := args[1]    //action定义 0越权，1查询，2上传

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
	switch actionType {
    case "view":
        his_credit += 1.0
    case "upload":
        his_credit += 2.0
    case "violation":
        his_credit -= 5.0
    default:
        return shim.Error("unknown action type")
    }

	var credit float64
	credit=his_credit + commu_credit
	//防止越界
	if credit > 100 {
		credit=100
	}
	if credit <0  {
		credit =0
	}
		//注销
	if credit < rank {
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
    var tmp Vehicle
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
    var coinAndCredit Vehicle
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
    var coinAndCredit Vehicle
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
    var tmp Vehicle
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
