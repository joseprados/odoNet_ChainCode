package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//TestReadingAsset_Init
func TestReadingAsset_Init(t *testing.T) {
	reading := new(ReadingAsset)
	stub := shim.NewMockStub("reading", reading)
	checkInit(t, stub, [][]byte{[]byte("init")})
	checkState(t, stub, "readingIDIndex", getExpectedReadingIDIndex("init"))
}

//TestReadingAsset_InvokeUnknownFunction
func TestReadingAsset_InvokeUnknownFunction(t *testing.T) {
	reading := new(ReadingAsset)
	stub := shim.NewMockStub("reading", reading)
	checkInit(t, stub, [][]byte{[]byte("init")})
	checkInvokeUnknownFunction(t, stub, [][]byte{[]byte("myFunction"), []byte("docType:Asset")})
}

//TestReadingAsset_Invoke_addNewReading
func TestReadingAsset_Invoke_addNewReadingOK(t *testing.T) {
	reading := new(ReadingAsset)
	stub := shim.NewMockStub("reading", reading)
	checkInit(t, stub, [][]byte{[]byte("init")})
	checkInvoke(t, stub, getFirstReadingAssetForTesting())
	newReadingID := "100001"
	checkState(t, stub, newReadingID, getNewReadingExpected())
	checkState(t, stub, "readingIDIndex", getExpectedReadingIDIndex("addNewReading"))
}

//TestReadingAsset_Invoke_addNewReadingUnknownField
func TestReadingAsset_Invoke_addNewReadingUnknownField(t *testing.T) {
	reading := new(ReadingAsset)
	stub := shim.NewMockStub("reading", reading)
	checkInit(t, stub, [][]byte{[]byte("init")})
	res := stub.MockInvoke("1", getReadingAssetWithUnknownFieldForTesting())
	if res.Status != shim.OK {
		checkError(t, "Reading Data is Corrupted", res.Message)
	} else {
		fmt.Println("Unknown Field Error was expected, but not raised")
		t.FailNow()
	}
}

//TestReadingAsset_Invoke_addNewReading
func TestReadingAsset_Invoke_addNewReadingDuplicate(t *testing.T) {
	reading := new(ReadingAsset)
	stub := shim.NewMockStub("reading", reading)
	checkInit(t, stub, [][]byte{[]byte("init")})
	checkInvoke(t, stub, getFirstReadingAssetForTesting())
	newReadingID := "100001"
	checkState(t, stub, newReadingID, getNewReadingExpected())
	checkState(t, stub, "readingIDIndex", getExpectedReadingIDIndex("addNewReading"))
	res := stub.MockInvoke("1", getFirstReadingAssetForTesting())
	if res.Status != shim.OK {
		checkError(t, "This Reading already exists: 100001", res.Message)
	} else {
		fmt.Println("Error was expected, but not raised")
		t.FailNow()
	}
}

//TestReadingAsset_Invoke_updateReadingOK  //change template
func TestReadingAsset_Invoke_updateReadingOK(t *testing.T) {
	reading := new(ReadingAsset)
	stub := shim.NewMockStub("reading", reading)
	checkInit(t, stub, [][]byte{[]byte("init")})
	checkInvoke(t, stub, getFirstReadingAssetForTesting())
	res := stub.MockInvoke("1", getUpdateReadingAssetForOKTesting())
	if res.Status != shim.OK {
		fmt.Println("Invoke", "failed", string(res.Message))
		t.FailNow()
	}
	checkState(t, stub, "100001", getUpdatedReadingExpected())
}

//TestReadingAsset_Invoke_updateReadingValueNOK  //change template
func TestReadingAsset_Invoke_updateReadingValueNOK(t *testing.T) {
	reading := new(ReadingAsset)
	stub := shim.NewMockStub("reading", reading)
	checkInit(t, stub, [][]byte{[]byte("init")})
	checkInvoke(t, stub, getFirstReadingAssetForTesting())
	res := stub.MockInvoke("1", getUpdateReadingAssetForValueNOKTesting())
	if res.Status != shim.OK {
		checkError(t, "updateReading: New Reading is less than Current Reading - cannot update", res.Message)
	} else {
		fmt.Println("Error was expected, but not raised")
		t.FailNow()
	}
	checkState(t, stub, "100001", getNewReadingExpected())
}

//TestReadingAsset_Invoke_updateReadingDateNOK  //change template
func TestReadingAsset_Invoke_updateReadingDateNOK(t *testing.T) {
	reading := new(ReadingAsset)
	stub := shim.NewMockStub("reading", reading)
	checkInit(t, stub, [][]byte{[]byte("init")})
	checkInvoke(t, stub, getFirstReadingAssetForTesting())
	res := stub.MockInvoke("1", getUpdateReadingAssetForDateNOKTesting())
	if res.Status != shim.OK {
		checkError(t, "updateReading: New Date is earlier than Current Date - cannot update", res.Message)
	} else {
		fmt.Println("Error was expected, but not raised")
		t.FailNow()
	}
	checkState(t, stub, "100001", getNewReadingExpected())
}

//TestReadingAsset_Invoke_removeAllReadingsOK  //change template
func TestReadingAsset_Invoke_removeAllReadingsOK(t *testing.T) {
	reading := new(ReadingAsset)
	stub := shim.NewMockStub("reading", reading)
	checkInit(t, stub, [][]byte{[]byte("init")})
	checkInvoke(t, stub, getFirstReadingAssetForTesting())
	checkInvoke(t, stub, getSecondReadingAssetForTesting())
	checkReadAllReadingsOK(t, stub)
	checkState(t, stub, "readingIDIndex", getExpectedReadingIDIndex("beforeRemoveReading"))
	checkInvoke(t, stub, getRemoveAllReadingAssetsForTesting())
	checkState(t, stub, "readingIDIndex", getExpectedReadingIDIndex(""))
}

//TestReadingAsset_Invoke_removeReadingNOK  //change template
func TestReadingAsset_Invoke_removeAllReadingsNOK(t *testing.T) {
	reading := new(ReadingAsset)
	stub := shim.NewMockStub("reading", reading)
	checkInit(t, stub, [][]byte{[]byte("init")})
	res := stub.MockInvoke("1", getRemoveAllReadingAssetsForTesting())
	if res.Status != shim.OK {
		checkError(t, "removeAllReadings: No readings to remove", res.Message)
	} else {
		fmt.Println("Error was expected, but not raised")
		t.FailNow()
	}
	checkState(t, stub, "readingIDIndex", getExpectedReadingIDIndex(""))
}

//TestReadingAsset_Query_readReading
func TestReadingAsset_Query_readReading(t *testing.T) {
	reading := new(ReadingAsset)
	stub := shim.NewMockStub("reading", reading)
	checkInit(t, stub, [][]byte{[]byte("init")})
	readingID := "100001"
	checkInvoke(t, stub, getFirstReadingAssetForTesting())
	checkReadReadingOK(t, stub, readingID)
	checkReadReadingNOK(t, stub, "")
}

//TestReadingAsset_Query_readAllReadings
func TestReadingAsset_Query_readAllReadings(t *testing.T) {
	reading := new(ReadingAsset)
	stub := shim.NewMockStub("reading", reading)
	checkInit(t, stub, [][]byte{[]byte("init")})
	checkInvoke(t, stub, getFirstReadingAssetForTesting())
	checkInvoke(t, stub, getSecondReadingAssetForTesting())
	checkReadAllReadingsOK(t, stub)
}

/*
*
*	Helper Functions
*
 */
//Get first ReadingAsset for testing
func getFirstReadingAssetForTesting() [][]byte {
	return [][]byte{[]byte("addNewReading"),
		[]byte("{\"vehicleID\":\"100001\",\"docType\":\"Asset.Reading\",\"reading\":\"50\",\"creationDate\":\"12/01/2017\"}")}
}

//Get ReadingAsset with unknown field for testing
func getReadingAssetWithUnknownFieldForTesting() [][]byte {
	return [][]byte{[]byte("addNewReading"),
		[]byte("{\"vehicleID\":\"100001\",\"docuType\":\"Asset.Reading\",\"reading\":\"50\",\"creationDate\":\"12/01/2017\"}")}
}

//Get second ReadingAsset for testing
func getSecondReadingAssetForTesting() [][]byte {
	return [][]byte{[]byte("addNewReading"),
		[]byte("{\"vehicleID\":\"100002\",\"docType\":\"Asset.Reading\",\"reading\":\"70\",\"creationDate\":\"12/01/2017\"}")}
}

//Get update ReadingAsset for OK testing
func getUpdateReadingAssetForOKTesting() [][]byte {
	return [][]byte{[]byte("updateReading"),
		[]byte("{\"vehicleID\":\"100001\",\"docType\":\"Asset.Reading\",\"reading\":\"100\",\"creationDate\":\"12/20/2017\"}")}
}

//Get update ReadingAsset for reading value NOK testing
func getUpdateReadingAssetForValueNOKTesting() [][]byte {
	return [][]byte{[]byte("updateReading"),
		[]byte("{\"vehicleID\":\"100001\",\"docType\":\"Asset.Reading\",\"reading\":\"20\",\"creationDate\":\"12/01/2017\"}")}
}

//Get update ReadingAsset for reading Date NOK testing
func getUpdateReadingAssetForDateNOKTesting() [][]byte {
	return [][]byte{[]byte("updateReading"),
		[]byte("{\"vehicleID\":\"100001\",\"docType\":\"Asset.Reading\",\"reading\":\"100\",\"creationDate\":\"11/20/2017\"}")}
}

//Get remove all ReadingAssets for testing //change template
func getRemoveAllReadingAssetsForTesting() [][]byte {
	return [][]byte{[]byte("removeAllReadings")}
}

//Get an expected value for testing
func getNewReadingExpected() []byte {
	var reading Reading
	reading.VehicleID = "100001"
	reading.ObjectType = "Asset.Reading"
	reading.Reading = "50"
	reading.CreationDate = "12/01/2017"
	readingJSON, err := json.Marshal(reading)
	if err != nil {
		fmt.Println("Error converting a Reading record to JSON")
		return nil
	}
	return []byte(readingJSON)
}

//Get an expected value for testing
func getUpdatedReadingExpected() []byte {
	var reading Reading
	reading.VehicleID = "100001"
	reading.ObjectType = "Asset.Reading"
	reading.Reading = "100"
	reading.CreationDate = "12/20/2017"
	readingJSON, err := json.Marshal(reading)
	if err != nil {
		fmt.Println("Error converting a Reading record to JSON")
		return nil
	}
	return []byte(readingJSON)
}

//Get expected values of Readings for testing
func getExpectedReadings() []byte {
	var readings []Reading
	var reading Reading
	reading.VehicleID = "100001"
	reading.ObjectType = "Asset.Reading"
	reading.Reading = "50"
	reading.CreationDate = "12/01/2017"
	readings = append(readings, reading)
	reading.VehicleID = "100002"
	reading.ObjectType = "Asset.Reading"
	reading.Reading = "70"
	reading.CreationDate = "12/01/2017"
	readings = append(readings, reading)
	readingJSON, err := json.Marshal(readings)
	if err != nil {
		fmt.Println("Error converting reading records to JSON")
		return nil
	}
	return []byte(readingJSON)
}

func getExpectedReadingIDIndex(funcName string) []byte {
	var readingIDIndex ReadingIDIndex
	switch funcName {
	case "addNewReading":
		readingIDIndex.VehicleIDs = append(readingIDIndex.VehicleIDs, "100001")
		readingIDIndexBytes, err := json.Marshal(readingIDIndex)
		if err != nil {
			fmt.Println("Error converting ReadingIDIndex to JSON")
			return nil
		}
		return readingIDIndexBytes
	case "beforeRemoveReading":
		readingIDIndex.VehicleIDs = append(readingIDIndex.VehicleIDs, "100001")
		readingIDIndex.VehicleIDs = append(readingIDIndex.VehicleIDs, "100002")
		readingIDIndexBytes, err := json.Marshal(readingIDIndex)
		if err != nil {
			fmt.Println("Error converting ReadingIDIndex to JSON")
			return nil
		}
		return readingIDIndexBytes
	case "afterRemoveReading":
		readingIDIndex.VehicleIDs = append(readingIDIndex.VehicleIDs, "100001")
		readingIDIndexBytes, err := json.Marshal(readingIDIndex)
		if err != nil {
			fmt.Println("Error converting ReadingIDIndex to JSON")
			return nil
		}
		return readingIDIndexBytes
	default:
		readingIDIndexBytes, err := json.Marshal(readingIDIndex)
		if err != nil {
			fmt.Println("Error converting ReadingIDIndex to JSON")
			return nil
		}
		return readingIDIndexBytes
	}
}

//checkInit - helper to check the Initialization of chaincode: ReadingAsset
func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
}

//checkState - helper for checking the chaincode state for a given stateKey afgainst an expected value
func checkState(t *testing.T, stub *shim.MockStub, stateKey string, expectedState []byte) {
	actualState := stub.State[stateKey]
	if actualState == nil {
		fmt.Println("State for ", stateKey, ": failed to get value")
		t.FailNow()
	}

	if bytes.Compare(actualState, expectedState) != 0 {
		fmt.Println("Incorrect State for " + stateKey + ": \nExpected: " + string(expectedState) + "\nActual  : " + string(actualState))
		t.FailNow()
	}
}

//checkInvoke - helper for checking Invoke of chaincode
func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		t.FailNow()
	}
}

func checkInvokeUnknownFunction(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		expectedErr := "Received unknown function invocation"
		actualErr := string(res.Message)
		if !(strings.Contains(actualErr, expectedErr)) {
			fmt.Println("Invoke with unknown function must throw error: ", expectedErr, "\n Actual Error :", actualErr)
			t.FailNow()
		}
	} else {
		t.FailNow()
	}
}

//checkReadReadingOK - helper for positive test readReading
func checkReadReadingOK(t *testing.T, stub *shim.MockStub, readingID string) {
	res := stub.MockInvoke("1", [][]byte{[]byte("readReading"), []byte(readingID)})
	if res.Status != shim.OK {
		fmt.Println("func readReading with ID: ", readingID, " failed"+string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("func readReading with ID: ", readingID, "failed to get value")
		t.FailNow()
	}
	if bytes.Compare(getNewReadingExpected(), []byte(res.Payload)) != 0 {
		fmt.Println("func readReading with ID: ", readingID, "Expected:", string(getNewReadingExpected()), "Actual:", string(res.Payload))
		t.FailNow()
	}
}

//checkReadReadingNOK - helper for negative testing of readReading
func checkReadReadingNOK(t *testing.T, stub *shim.MockStub, readingID string) {
	//with no readingID
	res := stub.MockInvoke("1", [][]byte{[]byte("readReading"), []byte("")})
	if res.Status != shim.OK {
		expectedErr := "retrieveReading: Corrupt reading record "
		actualErr := string(res.Message)
		if !(strings.Contains(actualErr, expectedErr)) {
			fmt.Println("func readReading negative test: ", "Expected Error:", expectedErr, "Actual Error", actualErr)
			t.FailNow()
		}
	} else {
		t.FailNow()
	}
}

func checkReadAllReadingsOK(t *testing.T, stub *shim.MockStub) {
	res := stub.MockInvoke("1", [][]byte{[]byte("readAllReadings")})
	if res.Status != shim.OK {
		fmt.Println("func readAllReadings failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("func readAllReadings failed to get value")
		t.FailNow()
	}
	if bytes.Compare(getExpectedReadings(), []byte(res.Payload)) != 0 {
		fmt.Println("func readAllReadings Expected:\n", string(getExpectedReadings()), "\nActual:\n", string(res.Payload))
		t.FailNow()
	}
}

func checkError(t *testing.T, exp string, act string) {
	if strings.Compare(exp, act) != 0 {
		fmt.Println("Unexpected Error! Expecting ", exp, "\n Actual :", act)
		t.FailNow()
	}
}
