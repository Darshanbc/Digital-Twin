package main

import (
	// "bytes"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/golang/protobuf/proto"
	//"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	mspprotos "github.com/hyperledger/fabric/protos/msp"
	sc "github.com/hyperledger/fabric/protos/peer"
	"log"
	strconv2 "strconv"
	"strings"
)

var logger = shim.NewLogger("farmer_agricert")

type SmartContract struct {
}

type Farmer struct {
	KYC_ID string
	Plots  []Plot       `json:"plot"`
	Store  []Fertilizer `json:"store"`
}
type CertAgent struct {
	KYC_ID              string
	certificateRequests [] string
	fertilizerRequests  []string
}

type KYC struct {
	First_Name string `json:"first name"`
	Last_Name  string `json:"last name"`
	Age        string `json:"age"`
	Address    string `json:"address"`
	Email      string `json:"email"`
	MobileNo   string `json:"mobile_no"`
}

type Plot struct {
	PlotId       string
	Co_ordinates []Latlong `json:"co_ordinates"`
	Survey_no    string    `json:"survey_no"`
	Crop_history []string  `json:"crop_history"`
	SoilType     string    `json:"soil_type"`
	Current_crop string    `json:"current_crop"`
}

type Latlong struct {
	Lattitude float64 `json:"lattitude"`
	Longitude float64 `json:"longitude"`
}

type Fertilizer struct {
	FertlizerName string  `json:"fertlizer_name"`
	FertlizerID   string  `json:"fertlizer_id"`
	Quantity      float32 `json:"quantity"`
}

type Crop struct {
	CropName        string       `json:"crop_name"`
	Type            string       `json:"type"`
	CropCycle       []Cycle      `json:"crop_cycle"`
	FertilzerUsed   []Fertilizer `json:"fertilzer_used"`
	FertilzerReq    []Fertilizer `json:"fertilzer_req"`
	CertRequestTxID string       `json:"cert_request_tx_id"`
	Cert            Certificate  `json:"cert"`
}

type Cycle struct {
	FromMonth string `json:"from_month"`
	ToMonth   string `json:"to_month"`
}

//0-----------------------------------Cert Agency-----
//

type Requests struct {
	Asset interface{}`json:"asset"` //cert or fert
	Status string `json:"status"`
}

//type AgriRequests struct {
//
//	FertilizerRequest struct {
//		Request []Request
//		status string
//	} `json:"fertilizer_request"`
//
//	CertRequest       struct {
//		Request []Request
//		status string
//	} `json:"cert_request"`
//}
//
//type Request struct {
//	TrxId string `json:"txId"`
//	AprrovalStatus string `json:"aprroval_status"`
//}

type Certificate struct {
	CertficateID     string `json:"certficate_id"`
	DigitalSignature string `json:"digital_signature"`
	IssueDate        string `json:"issue_date"`
	Type             string `json:"type"`
	//OrganicPercentage float32
}
type IDs struct {
	KYC_ID  string
	Crop_ID string
	PlotID  string
}

//----------------------

//var logger = shim.NewLogger("example_cc0")

func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) sc.Response {
	fmt.Println("-----------Instantiating---------------------------")
	id := []IDs{
		IDs{
			KYC_ID:  "KYC1001",
			Crop_ID: "CROP1001",
			PlotID:  "Plot1001",
		},
	}
	idsAsBytes, _ := json.Marshal(id[0])
	stub.PutState("ids", idsAsBytes)

	kyc := KYC{First_Name: "Criyagen", Address: "Bangalore", Email: "agricert@criyagen.com"}
	idsAsBytes, _ = json.Marshal(kyc)
	stub.PutState("1000", idsAsBytes)

	agent := CertAgent{KYC_ID: "1000"}
	idsAsBytes, _ = json.Marshal(agent)
	stub.PutState("Criyagen", idsAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) sc.Response {
	logger.Info("-------------------Invoke----------")
	// Retrieve the requested Smart Contract function and arguments
	function, args := stub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	logger.Info("fucntion")
	logger.Info(function)
	logger.Info("args")
	logger.Info(args)
	if function == "newuser" {
		return s.KYCRegistration(stub, args)
	} else if function == "PlotRegisteration" {
		return s.PlotRegisteration(stub)
	} else if function == "query" {
		return s.query(stub, args)
	} else if function == "CropDetails" {
		return s.CreateCrop(stub)
	} else if function == "addFertilizerToCrop" {
		return s.addFertilizerToCrop(stub)
	} else if function == "addFertilizerToStore" {
		return s.addFertilizerToStore(stub)
	} else if function == "ApproveOrDenyFertilizer" {
		return s.ApproveOrDenyFertilizer(stub, args)
	} else {
		return shim.Error("Invalid Smart contract function aname")
	}

}

func (s *SmartContract) KYCRegistration(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	logger.Info("########### KYCRegistration ###########")

	creator, err := stub.GetCreator()
	if err != nil {
		return shim.Error(err.Error())
	}

	id := &mspprotos.SerializedIdentity{}
	err = proto.Unmarshal(creator, id)
	block, _ := pem.Decode(id.GetIdBytes())
	cert, err := x509.ParseCertificate(block.Bytes)
	enrollID := cert.Subject.CommonName

	userAsbytes, _ := stub.GetState(enrollID)

	if userAsbytes != nil {
		return shim.Error("User already exists")
	}
	fmt.Printf("enrollID: %s", enrollID)
	//mspID := id.GetMspid()

	newuser := KYC{First_Name: args[0], Last_Name: args[1], Age: args[2], Address: args[3], Email: args[4], MobileNo: args[5]}
	//generating user ID
	idsAsBytes, _ := stub.GetState("ids")
	ids := IDs{}
	json.Unmarshal(idsAsBytes, &ids)
	logger.Debug("ids are ", ids)
	kycID := ids.KYC_ID
	logger.Info("kycID ", kycID)
	j := string([]rune(ids.KYC_ID)[3:])
	logger.Debug("value of J in KYC_ID is ", j)
	tempKYCno, _ := strconv2.Atoi(j)
	tempKYCno = tempKYCno + 1
	ids.KYC_ID = "KYC" + strconv2.Itoa(tempKYCno)
	idsAsBytes, _ = json.Marshal(ids)
	stub.PutState("ids", idsAsBytes)
	// fmt.Println("enrollID", enrollID)
	fmt.Println("user object ", newuser)
	//------------------------
	newuserAsBytes, _ := json.Marshal(newuser)
	stub.PutState(kycID, newuserAsBytes)

	newfarmer := Farmer{}
	newfarmer.KYC_ID = kycID
	newfarmerAsBytes, _ := json.Marshal(newfarmer)
	stub.PutState(enrollID, newfarmerAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) PlotRegisteration(stub shim.ChaincodeStubInterface) sc.Response {
	// {"co_ordinates":[{"lattitude":45.15445,"longitude":32.45455},{"lattitude":45.15445,"longitude":32.45455},{"lattitude":45.15445,"longitude":32.45455},{"lattitude":45.15445,"longitude":32.45455}],
	// "survey_no":"Eg456-456-5677",
	// "soil_type":"sasa"}
	logger.Info("-----------------------PlotRegistration--------------------")
	//retreiving the args
	args := stub.GetArgs()

	//-------------------------------
	// creatinng plot object
	tempPlot := Plot{}

	if err := json.Unmarshal([]byte(args[1]), &tempPlot); err != nil {
		log.Fatal(err)
		return shim.Error(err.Error())
	}

	//-----------------------------------

	enrollID := getEnrollID(stub)
	userAsbytes, _ := stub.GetState(enrollID)
	tempfarmer := Farmer{}
	json.Unmarshal(userAsbytes, &tempfarmer)
	//---------------------------------------

	// PlotId geneartion and assigning to Plot
	length := len(tempfarmer.Plots)
	tempPlot.PlotId = "Plot" + strconv2.Itoa(length)

	tempfarmer.Plots = append(tempfarmer.Plots, tempPlot)
	userAsbytes, _ = json.Marshal(tempfarmer)
	stub.PutState(enrollID, userAsbytes)
	return shim.Success(nil)
}

func getEnrollID(stub shim.ChaincodeStubInterface) string {
	creator, _ := stub.GetCreator()
	id := &mspprotos.SerializedIdentity{}
	_ = proto.Unmarshal(creator, id)
	block, _ := pem.Decode(id.GetIdBytes())
	cert, _ := x509.ParseCertificate(block.Bytes)
	enrollID := cert.Subject.CommonName
	return enrollID
}

//fill crop history

func (s *SmartContract) CreateCrop(stub shim.ChaincodeStubInterface) sc.Response {
	//["plotID","","certRegistration"]
	fmt.Print("-------------------cropDetails------------------")
	args := stub.GetArgs()
	plotID := args[1]
	logger.Info("plotId", string(plotID))
	enrollID := getEnrollID(stub)
	tempfarmer := Farmer{}
	farmerAsbytes, _ := stub.GetState(enrollID)
	json.Unmarshal(farmerAsbytes, &tempfarmer)
	flag := 0
	var plot Plot
	var index int
	for index, plot = range tempfarmer.Plots {
		if plot.PlotId == string(plotID) {
			flag = 1
			break
		}
	}
	if flag == 0 {
		return shim.Error("Invalid plot ID")
	}
	logger.Info("plot is ", plot)
	idsAsBytes, _ := stub.GetState("ids")
	ids := IDs{}
	json.Unmarshal(idsAsBytes, &ids)
	logger.Debug("ids are ", ids)
	cropID := ids.Crop_ID
	logger.Info("crop", cropID)
	j := string([]rune(ids.Crop_ID)[4:])
	logger.Debug("value of J in cropID is ", j)
	tempcropno, _ := strconv2.Atoi(j)
	tempcropno = tempcropno + 1
	ids.Crop_ID = "CROP" + strconv2.Itoa(tempcropno)
	idsAsBytes, _ = json.Marshal(ids)
	stub.PutState("ids", idsAsBytes)

	crop := Crop{}
	if err := json.Unmarshal([]byte(args[2]), &crop); err != nil {
		log.Fatal(err)
		return shim.Error(err.Error())
	}
	logger.Info("crop deatils sent is ", crop)

	if string(args[3]) == "certRegistration" {
		logger.Info("Arg[3] pased is ", string(args[3]))
		crop.CertRequestTxID = stub.GetTxID()
		tempfarmer.Plots[index].Current_crop = cropID
		req := Requests{}
		req.Status="Pending"
		reqAsBytes,_:=json.Marshal(req)
		stub.PutState(stub.GetTxID(),reqAsBytes)


	} else {
		tempfarmer.Plots[index].Crop_history = append(tempfarmer.Plots[index].Crop_history, cropID)

	}

	logger.Info("crop before map ", args[2])
	logger.Info("Crop details ", crop)
	cropAsBytes, _ := json.Marshal(crop)
	stub.PutState(cropID, cropAsBytes)

	logger.Info("after adding crop history", plot)
	farmerAsbytes, _ = json.Marshal(tempfarmer)
	stub.PutState(enrollID, farmerAsbytes)
	return shim.Success(nil)

}

func (s *SmartContract) query(stub shim.ChaincodeStubInterface, args []string) sc.Response {
	logger.Info("args", args[0])
	result, _ := stub.GetState(args[0])

	return shim.Success(result)
}

func (s *SmartContract) addFertilizerToStore(stub shim.ChaincodeStubInterface) sc.Response {
	//[{fertlizer_id:"jhjkh",fertlizer_name:"ghjghg",quantity:4}]
	args := stub.GetArgs()
	fert := Fertilizer{}
	if err := json.Unmarshal([]byte(args[1]), &fert); err != nil {
		log.Fatal(err)
		return shim.Error(err.Error())
	}
	enrollID := getEnrollID(stub)
	idAsBytes, _ := stub.GetState(enrollID)
	farmer := Farmer{}
	json.Unmarshal(idAsBytes, &farmer)
	farmer.Store = append(farmer.Store, fert)
	idAsBytes, _ = json.Marshal(farmer)
	stub.PutState(enrollID, idAsBytes)
	return shim.Success(nil)

}
func (s *SmartContract) addFertilizerToCrop(stub shim.ChaincodeStubInterface) sc.Response {
	//[plotId,{fertlizer_name:"name",fertlizer_id:"asd",quantity:2.56},"pending"]
	//get plot
	enrollID := getEnrollID(stub)
	idAsBytes, _ := stub.GetState(enrollID)
	farmer := Farmer{}
	json.Unmarshal(idAsBytes, &farmer)
	args := stub.GetArgs()
	plotId := string(args[1])
	var index int
	flag := 0
	var cropid string
	for index, _ = range farmer.Plots {

		if farmer.Plots[index].PlotId == plotId {
			flag = 1
			break;
		}
	}
	//check if plot exists and get current crop id
	if flag == 1 {
		cropid = farmer.Plots[index].Current_crop
	} else {
		return shim.Error("plot doesnot belong to you")
	}

	logger.Info("current crop", cropid)
	cropAsBytes, _ := stub.GetState(cropid)
	crop := Crop{}
	json.Unmarshal(cropAsBytes, &crop)
	certAgentAsBytes, _ := stub.GetState("Criyagen")
	agent := CertAgent{}
	json.Unmarshal(certAgentAsBytes, &agent)
	fert := Fertilizer{}
	if err := json.Unmarshal([]byte(args[2]), &fert); err != nil {
		log.Fatal(err)
		return shim.Error(err.Error())
	}

	flag = 0
	for index, _ = range farmer.Store {

		if farmer.Store[index].FertlizerID == fert.FertlizerID && farmer.Store[index].Quantity <= fert.Quantity {
			flag = 1
			break;
		}

	}
	if flag == 1 {
		txid := stub.GetTxID()
		req := Requests{Asset:fert, Status: "pending"}
		reqAsBytes,_:=json.Marshal(req)
		agent.fertilizerRequests = append(agent.fertilizerRequests, txid)
		stub.PutState(txid, reqAsBytes)
		cropAsBytes, _ = json.Marshal(crop)
		stub.PutState(cropid, cropAsBytes)
		idAsBytes, _ = json.Marshal(agent)
		stub.PutState("Criyagen", idAsBytes)
		return shim.Success(nil)
	} else {
		return shim.Error("Fertilizer not present in the stocl or insufficient quantity")
	}

}
func (s *SmartContract) ApproveOrDenyFertilizer(stub shim.ChaincodeStubInterface, args [] string) sc.Response {
	//	[ 'CROP1002',
	//	'DarshanBC6',
	//	'dfd44d11d76d4f2bcfce965b31f7a8ddee2c89a9f7e689b39a31ce6fcf688f2d',
	//'yes' ]


	req := Requests{}
	fert := Fertilizer{}
	reqAsbytes, _ := stub.GetState(args[2])//txid
	json.Unmarshal(reqAsbytes, &req)
	req.Status=args[3]//status
	if args[3]=="yes" {
		farmer := Farmer{}
		idsAsBytes, _ := stub.GetState(args[1])//enrollid
		json.Unmarshal(idsAsBytes, &farmer)
		var index int
		var flag int
		for index, _ = range farmer.Store {
			if farmer.Store[index].FertlizerID == fert.FertlizerID && farmer.Store[index].Quantity <= fert.Quantity {
				flag = 1
				farmer.Store[index].Quantity = farmer.Store[index].Quantity - fert.Quantity
				break;
			}
		}
		if flag == 0 {
			return shim.Error("Insufficient Stock")
		}
		certAgentAsBytes, _ := stub.GetState("Criyagen")
		agent := CertAgent{}
		json.Unmarshal(certAgentAsBytes, &agent)

		for index, element := range agent.fertilizerRequests {
			if strings.Compare(strings.ToLower(element), args[2]) == 0 {//txid
				agent.fertilizerRequests = append(agent.fertilizerRequests[:index], agent.fertilizerRequests[index+1:]...)
			}
		}
		//upload agent
		certAgentAsBytes,_=json.Marshal(agent)
		stub.PutState("Criyagen",certAgentAsBytes)
		//upload farmer
		idsAsBytes,_=json.Marshal(farmer)
		stub.PutState(args[1],idsAsBytes)
	}
	//read Crop
	crop,_:=stub.GetState(args[0])
	//write crop
	stub.PutState(args[0],crop)
	return shim.Success(nil)


}
func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		logger.Error("Error starting Simple chaincode: %s", err)
	}
}

//----------------------fertilizer agency--------------

//type Fertilizer struct {
//	FertlizerName string
//	FertlizerID   string
//	Quantity      float32
//}
