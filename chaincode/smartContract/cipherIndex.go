package sc
import (
    "encoding/json"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

func (s *CCC) AppendNewMessage(stub shim.ChaincodeStubInterface, args []string) peer.Response {
    if (len(args) != 8){
        return shim.Error("expect for 8 args!")
    }
    pkEnc := args[0]
    pkUser := args[1]
    roadName := args[2]
    messageLevel := args[3]
    weather := args[4]
    condition :=args[5]
    traffic := args[6]
    averageSpeed := args[7]

    IsExist := s.FeelExist(stub, pkEnc)
	if (IsExist == "Yes") {
		return shim.Error("Data Item has existed, please use the modify function")
	}

    var tmp = CipherIndex{
        PkEnc : pkEnc,
        PkUser : pkUser,
        RoadName : roadName,
        MessageLevel: messageLevel,
        Weather: weather,
        Condition: condition,
        Traffic: traffic,
        AverageSpeed: averageSpeed,
    }
    itemBytes, err := json.Marshal(tmp)
    if err != nil {
        return shim.Error("json marshal errror!")
    }
    err = stub.PutState(pkEnc, itemBytes)
    if err != nil {
        return shim.Error("json marshal errror!")
    }
    return shim.Success([]byte("Add new message success"))
}

func (s *CCC) DeleteMessage(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	// check if exists
	pkEnc := args[0]
	databytes, err := stub.GetState(pkEnc)
	if (err != nil || len(databytes) == 0) {
		return shim.Error("No target Message")
	}

	err = stub.DelState(pkEnc)
	if err != nil {
		return shim.Error("Delete State Error")
	}
	return shim.Success([]byte("Delete Message Success"))
}

func (s *CCC) GetMessageByRoadName(stub shim.ChaincodeStubInterface, args []string) peer.Response{
    if (len(args) != 1){
        return shim.Error("expect for 1 args")
    }
    messageIterator, err := stub.GetStateByRange("", "")
    if err != nil{
        return shim.Error("no message in log!")
    }
    defer messageIterator.Close()
    var messages []*CipherIndex
    for messageIterator.HasNext(){
        messageResponse, err := messageIterator.Next()
        if err != nil{
            return shim.Error("iterator error!")
        }
        var cipherIndex CipherIndex
        err = json.Unmarshal(messageResponse.Value, &cipherIndex)
        if err != nil{
            return shim.Error("json unmarshal error!")
        }
        if(cipherIndex.RoadName == args[0]){
            messages = append(messages, &cipherIndex)
        }

    }
    var returnBytes, _ = json.Marshal(messages)
    return shim.Success(returnBytes)
}
func (s * CCC)FeelExist(stub shim.ChaincodeStubInterface, target string) string {
	ResBytes, _ := stub.GetState(target)
	if(len(ResBytes) == 0) {
		return "No"
	} else  {
		return "Yes"
	}
}
