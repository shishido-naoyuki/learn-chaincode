/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type Account struct {
	ID          string  `json:"id"`
	PASSWORD          string  `json:"password"`
	CashBalance float64 `json:"cashBalance"`
}

type AccountList struct {
	LIST          string  `json:"list"`
	UserID        []string  `json:"userid"`
}

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 0 {
		fmt.Println("initdayo")
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
    return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

    // Handle different functions
    if function == "init" {
        return t.Init(stub, "init", args)
    } else if function == "createAccount" {
        return t.createAccount(stub, args)
    } else if function == "moneyTransfer" {
        return t.moneyTransfer(stub, args)
    }
    fmt.Println("invoke did not find func: " + function)

    return nil, errors.New("Received unknown function invocation: " + function)
}

func (t *SimpleChaincode) createAccount(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var username, password string
	var err error
	fmt.Println("running createAccount()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}
	
	//user create
	username = args[0]                            //rename for fun
	password = args[1]
	var account = Account{ID: username, PASSWORD: password, CashBalance: 10000000.0}
	accountBytes, err := json.Marshal(&account)
	if err != nil {
		fmt.Println("error creating account" + account.ID)
		return nil, errors.New("Error creating account " + account.ID)
	} else {
		fmt.Println("No existing account found for " + account.ID + ", initializing account.")
		err = stub.PutState(account.ID, accountBytes)
	}
	//uesr's list make
	var accountlist AccountList
	AccountListBytes, err := stub.GetState("LIST")
	err = json.Unmarshal(AccountListBytes, &accountlist)
	if err != nil {
		accountlist = AccountList{LIST: "LIST", UserID: []string{username}}
		accountBytes, err := json.Marshal(accountlist)
		if err !=  nil {
				fmt.Println("accountBytesError")
				return nil, errors.New("accountBytesError" + account.ID)
		}
		err = stub.PutState("LIST", accountBytes)
	} else {
		accountlist.UserID = append(accountlist.UserID, username)
		accountUpdataBytes, err := json.Marshal(&accountlist)
		if err !=  nil {
			fmt.Println("accountUpdataBytesError")
			return nil, errors.New("accountUpdataBytesError" + account.ID)
		} else{
			err = stub.PutState(accountlist.LIST, accountUpdataBytes)
		}
	}
	return nil, nil
}

func (t *SimpleChaincode) moneyTransfer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var fromUsername, toUsername string
	var money float64
	var err error
	
	fmt.Println("running moneyTransfer()")

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}
	
	fromUsername = args[0]                            //rename for fun
	toUsername = args[1]
	money, err = strconv.ParseFloat(args[2], 64)
	if err != nil {
		fmt.Println("moneyParseFloatError")
		return nil, errors.New("moneyParseFloatError")
	}
	
	var fromAccount Account
	var toAccount Account
	
	fromAccountBytes, err := stub.GetState(fromUsername)
	err = json.Unmarshal(fromAccountBytes, &fromAccount)
	if err != nil {
		fmt.Println("fromAccountBytesError")
		return nil, errors.New("fromAccountBytesError")
	}

	toAccountBytes, err := stub.GetState(toUsername)
	err = json.Unmarshal(toAccountBytes, &toAccount)
	if err != nil {
		fmt.Println("toAccountBytesError")
		return nil, errors.New("toAccountBytesError")
	}
	
	fromAccount.CashBalance -= money
	toAccount.CashBalance += money
	
	fromAccountUpdataBytes, err := json.Marshal(&fromAccount)
	if err !=  nil {
		fmt.Println("fromAccountUpdataBytesError")
			return nil, errors.New("fromAccountUpdataBytesError")
	} else{
		err = stub.PutState(fromAccount.ID, fromAccountUpdataBytes)
	}
	
	toAccountUpdataBytes, err := json.Marshal(&toAccount)
	if err !=  nil {
		fmt.Println("toAccountUpdataBytesError")
		return nil, errors.New("toAccountUpdataBytesError")
	} else{
		err = stub.PutState(toAccount.ID, toAccountUpdataBytes)
	}
	
	return nil, nil
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" {											//read a variable
		return t.read(stub, args)
	} else if function == "listRead" {
        return t.listRead(stub, args)
	}
	fmt.Println("query did not find func: " + function)						//error

	return nil, errors.New("Received unknown function query: " + function)
}

func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var username, password, jsonResp string
	var err error
	var account Account
    if len(args) != 2 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
    }

    username = args[0]
    password = args[1]
    
    accountBytes, err := stub.GetState(username)
    err = json.Unmarshal(accountBytes, &account)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + username + "\"}"
        return nil, errors.New(jsonResp)
    }
    if account.PASSWORD != password {
		jsonResp = "{\"Error\":\"login error \"}"
		return nil, errors.New(jsonResp)
    }
	
    return accountBytes, nil
    
    
}

func (t *SimpleChaincode) listRead(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var list, jsonResp string
	var err error
	var accountlist AccountList
    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
    }

    list = args[0]
    accountListBytes, err := stub.GetState("LIST")
    err = json.Unmarshal(accountListBytes, &accountlist)
    if err != nil {
        jsonResp = "{\"Error\":\"Failed to get state for " + list + "\"}"
        return nil, errors.New(jsonResp)
    }
    return accountListBytes, nil
}