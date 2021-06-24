package main

import (
	"encoding/json"
)

func deserializeBuyerCustomData(bytes []byte) (buyerCustomData BuyerCustomData) {
	json.Unmarshal(bytes, &buyerCustomData)
	return
}

func deserializeRetailerCustomData(bytes []byte) (retailerCustomData RetailerCustomData) {
	json.Unmarshal(bytes, &retailerCustomData)
	return
}

func deserializeTransaction(bytes []byte) (transaction Transaction) {
	json.Unmarshal(bytes, &transaction)
	return
}

func serializeCustomData(object struct{}) (serializedCustomData []byte) {
	serializedCustomData, _ = json.Marshal(object)
	return
}
