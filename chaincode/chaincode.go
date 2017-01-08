
package main

import (
    "encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("patients", []byte(args[0]))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke isur entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "createVital" {
		return t.createVital(stub, args)
	}else if function == "createPatient" {
		return t.createPatient(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// write - invoke function to write key/value pair
func (t *SimpleChaincode) createPatient(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}
	var vitals []string
	key = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(key, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	vitalsBytes, err := json.Marshal(&vitals)
	if err != nil {
			fmt.Println("Error marshalling vitals")
			return nil, errors.New("Error create patient")
		}
	err = stub.PutState(key + "vitals", vitalsBytes) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}


func (t *SimpleChaincode) createVital(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, vital,timestamp string
	var err error
	fmt.Println("running write()")

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3. name of the key and value to set")
	}
	var vitals []string
	key = args[0] //rename for funsies
	vital = args[1]
	timestamp = args[2]
	err = stub.PutState(key + "vital" + timestamp, []byte(vital)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	vitalBytes, err := stub.GetState(key + "vitals")
		err = json.Unmarshal(vitalBytes, &vitals)
		if err != nil {
			fmt.Println("Error unmarshel keys")
			return nil, errors.New("Error unmarshalling vitals ")
		}
		vitals = append(vitals, key + "vital" + timestamp)
	vitalsBytes, err := json.Marshal(&vitals)
	if err != nil {
			fmt.Println("Error marshalling vitals")
			return nil, errors.New("Error create patient")
		}
	err = stub.PutState(key + "vitals", vitalsBytes) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}
// read - query function to read key/value pair
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}
