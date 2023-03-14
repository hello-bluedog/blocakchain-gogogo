package sc

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)


func (s *CCC) AddVehicle(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1")
    }

    pkUser := args[0]
	coinNumber := args[1]
	credit := args[2]
	role := args[3]

    vehicle := Vehicle{
    	PkUser:   pkUser,
    	CoinNum:  coinNumber,
    	Credit:   credit,
    	Role:     role,
    	LastPing: time.Now().Unix(),
    	Activity: 0.0,
    }

	//车辆信息上传
	itembytes, err := json.Marshal(vehicle)
	if err != nil {
		return shim.Error("json marshal shim.Error")
	}
	err = stub.PutState(pkUser, itembytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("PutState shim.Error: %s", err))
	}

	return shim.Success([]byte("Add new user success"))
}

func (s *CCC) Ping(stub shim.ChaincodeStubInterface, args []string) peer.Response {
    if len(args) != 1 {
        return shim.Error("Incorrect number of arguments. Expecting 1")
    }

    pkUser := args[0]
	itembytes, err := stub.GetState(pkUser)
    if(len(itembytes) == 0){
        return shim.Error("no such user in ledger!!!")
    }
    var vehicle Vehicle
    _ = json.Unmarshal(itembytes, &vehicle)
    

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
   vehicle.Activity = float64(vehicle.PingCount) * freq

   // 将更新后的车辆信息保存
	itembytes,err = json.Marshal(vehicle)
	if err != nil {
		return shim.Error("json marshal shim.Error")
	}
	err = stub.PutState(pkUser, itembytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("PutState shim.Error: %s", err))
	}
   
   return shim.Success([]byte("Ping success"))
}


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
    var myVehicle Vehicle
    _ = json.Unmarshal(itembytes, &myVehicle)
	
	//历史信誉
	var totalCredit float64
    var weightTotal float64
    var currentTimestamp = time.Now().Unix()

    historyIter, err := stub.GetHistoryForKey(pkUser)
    if err != nil {
        return shim.Error("failed to get history for pkUser")
    }
    defer historyIter.Close()

	for historyIter.HasNext() {
		historicValue, err := historyIter.Next()
		if err != nil {
			return shim.Error("failed to get next history value for pkUser")
		}
	
		var vehicle Vehicle
		err = json.Unmarshal(historicValue.Value, &vehicle)
		if err != nil {
			return shim.Error("failed to unmarshal historic vehicle value for pkUser")
		}
	
		historicCredit, err := strconv.ParseFloat(vehicle.Credit, 64)
		if err != nil {
			return shim.Error("failed to parse historic vehicle credit for pkUser")
		}
	
		// historicTimestamp := historicValue.Timestamp.Seconds * 1e9 + int64(historicValue.Timestamp.Nanos)
		historicTimestamp := vehicle.LastPing
		// 计算时间差，并计算权重系数
		timeDelta := currentTimestamp - historicTimestamp
		decayRate := 0.001
    	secondsInYear := 31536000.0 // 秒数
    	weight := math.Exp(-float64(timeDelta) * decayRate / secondsInYear)
	
		// 计算加权信誉度
		weightedCredit := historicCredit * weight
	
		// 将加权信誉度加入总和中
		totalCredit += weightedCredit
	
		// 统计历史版本数
		weightTotal++
	}
	
	// 计算平均信誉度
	his_credit := totalCredit / float64(weightTotal)
	if his_credit >100.0 {
		his_credit =100.0
	}

	//通信信誉
	var commu_credit float64
	if myVehicle.Activity>100.0 {
		commu_credit = 100.0
	} else {
		commu_credit = myVehicle.Activity
	}

	var final_credit float64
	final_credit = his_credit*0.58 + commu_credit*0.4
	
	//行为信誉
	switch action {
    case "1":
        final_credit += 1.0
    case "2":
        final_credit += 2.0
    case "3":
        final_credit -= 5.0
    default:
        return shim.Error("unknown action type")
    }

	if final_credit > 100.0 {
		final_credit = 100.0
	}else if final_credit <0.0 {
		final_credit =0.0
	}

	myVehicle.Credit=strconv.FormatFloat(final_credit,'f',2,64)
    itembytes, err = json.Marshal(myVehicle)
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
