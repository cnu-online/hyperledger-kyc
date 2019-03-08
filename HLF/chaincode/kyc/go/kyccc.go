package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

//Define the Smart Contract structure ...
type SmartContract struct {
}

//KycRecord ...
type KycRecord struct {
	KycID            string `json:"kycId"`
	LastName         string `json:"lastName"`
	FirstName        string `json:"firstName"`
	FatherSpouseName string `json:"fatherSpouseName"`

	Gender                 Gender            `json:"gender"`
	MaritalStatus          MaritalStatus     `json:"maritalStatus"`
	Citizenship            Citizenship       `json:"citizenship"`
	ResidentialStatus      ResidentialStatus `json:"residentialStatus"`
	OccupationType         OccupationType    `json:"occupationType"`
	ProofOfIdentityType    ProofOfIdentity   `json:"proofOfIdentityType"`
	ProofOfIdentity        string            `json:"proofOfIdentity"`
	AddressType            AddressType       `json:"addressType"`
	ProofOfAddressType     string            `json:"proofOfAddressType"`
	ProofOfAddress         ProofOfAddress    `json:"proofOfAddress"`
	ProofOfAddressExpiryDt time.Time         `json:"ProofOfAddressTypeExpiryDt"`
	PAddressLine1          string            `json:"pAddressLine1"`
	PAddressLine2          string            `json:"pAddressLine2"`
	PAddressLine3          string            `json:"pAddressLine3"`
	PAddressDistrict       string            `json:"pAddressDistrict"`
	PAddressTown           string            `json:"pAddressTown"`
	PAddressPinCode        string            `json:"pAddressPinCode"`

	CAddressLine1    string `json:"cAddressLine1"`
	CAddressLine2    string `json:"cAddressLine2"`
	CAddressLine3    string `json:"cAddressLine3"`
	CAddressDistrict string `json:"cAddressDistrict"`
	CAddressTown     string `json:"cAddressTown"`
	CAddressPinCode  string `json:"cAddressPinCode"`

	Telphone     string `json:"phone"`
	MobileNumber string `json:"mobileNumber"`
	EmailAddress string `json:"emailAddress"`
}

//AddressType ...
type AddressType int

const (
	AddressTypeResidential AddressType = 0
	AddressTypeBusiness    AddressType = 1
	AddressTypeUnspecified AddressType = 2
)

//ProofOfAddress ...
type ProofOfAddress int

// ...
const (
	POAPassport                  ProofOfAddress = 0
	POAVoterIDCard               ProofOfAddress = 1
	POAPANCard                   ProofOfAddress = 2
	POADrivingLicence            ProofOfAddress = 3
	POAAadhaarUID                ProofOfAddress = 4
	POANREGAJobCard              ProofOfAddress = 5
	POAOthers                    ProofOfAddress = 6
	POASimplifiedMeasuresAccount ProofOfAddress = 7
)

//ProofOfIdentity ...
type ProofOfIdentity int

const (
	POIPassport                  ProofOfIdentity = 0
	POIVoterIDCard               ProofOfIdentity = 1
	POIPANCard                   ProofOfIdentity = 2
	POIDrivingLicence            ProofOfIdentity = 3
	POIAadhaarUID                ProofOfIdentity = 4
	POINREGAJobCard              ProofOfIdentity = 5
	POIOthers                    ProofOfIdentity = 6
	POISimplifiedMeasuresAccount ProofOfIdentity = 7
)

//Citizenship ...
type Citizenship int

const (
	Indian           Citizenship = 0
	OtherCitizenship Citizenship = 1
)

//Gender ...
type Gender int

const (
	Male        Gender = 0
	Female      Gender = 1
	Transgender Gender = 2
)

//MaritalStatus ...
type MaritalStatus int

const (
	Married            MaritalStatus = 0
	Unmarried          MaritalStatus = 1
	OtherMaritalStatus MaritalStatus = 2
)

//ResidentialStatus ...
type ResidentialStatus int

const (
	ResidentIndividual   ResidentialStatus = 0
	NonResidentIndian    ResidentialStatus = 1
	ForeignNational      ResidentialStatus = 2
	PersonofIndianOrigin ResidentialStatus = 3
)

//OccupationType ...
type OccupationType int

const (
	PrivateSector    OccupationType = 0
	PublicSector     OccupationType = 1
	GivernmentSector OccupationType = 2
	Professional     OccupationType = 3
	SelfEmployed     OccupationType = 4
	Retired          OccupationType = 5
	Housewife        OccupationType = 6
	Student          OccupationType = 7
	Business         OccupationType = 8
	NotCategorised   OccupationType = 9
)

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createKycRecord" {
		return s.createKycRecord(APIstub, args)
	} else if function == "getKycs" {
		return s.getKycs(APIstub, args)
	} else if function == "getKyc" {
		return s.getKyc(APIstub, args)
	}else if function == "queryKycs" {
		return s.queryKycs(APIstub, args)
	}
	return shim.Error("Invalid Smart Contract function name.")
}
func (s *SmartContract) createKycRecord(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	KycID := args[0]
	firstName := args[1]
	lastName := args[2]

	kyc := KycRecord{KycID: KycID, FirstName: firstName, LastName: lastName}
	kycAsBytes, err := json.Marshal(kyc)
	if err != nil {
		return shim.Error("Invalid Value")
	}
	APIstub.PutState(KycID, kycAsBytes)
	return shim.Success(nil)
}
func (s *SmartContract) getKycs(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	org := strings.ToLower(args[0])
	orgtype := strings.ToLower(args[1])
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"kycRecord\",\"%s\":\"%s\"}}", orgtype, org)

	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)

}

// get Trade finance record
func (s *SmartContract) getKyc(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	KycId := args[0]
	tfAsBytes, _ := APIstub.GetState(KycId)
	return shim.Success(tfAsBytes)
}

//load two initial record
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	kycs := []KycRecord{
		KycRecord{KycID: time.Now().Format("20060102150405") + "1", FirstName: "Sreenivas", LastName: "Chinni"},
		KycRecord{KycID: time.Now().Format("20060102150405") + "2", FirstName: "Krishna Charan", LastName: "Chinni"},
	}

	rec1, err := json.Marshal(kycs[0])
	if err != nil {
		return shim.Error("error parsing record one")
	}
	APIstub.PutState(time.Now().Format("20060102150405")+"1", rec1)
	rec2, err1 := json.Marshal(kycs[1])
	if err1 != nil {
		return shim.Error("error parsing record two")
	}
	APIstub.PutState(time.Now().Format("20060102150405")+"2", rec2)

	return shim.Success(nil)
}
func (s *SmartContract) queryKycs(APIstub shim.ChaincodeStubInterface, args []string) sc.Response { 

	//   0
	// "queryString"
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	queryString := args[0]

	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

// ===========================================================================================
// constructQueryResponseFromIterator constructs a JSON array containing query results from
// a given result iterator
// ===========================================================================================
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return &buffer, nil
}



// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}




//peer chaincode query -C commonchannel -n kyccc -c '{"Args":["queryKeycs","{\"selector\":{\"firstname\":\"sreenivas\"}}"]}'