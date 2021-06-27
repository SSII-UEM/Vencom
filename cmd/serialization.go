package main

import (
	"encoding/json"
)

//se deserializa un string al tipo BuyerCustomData
func deserializeBuyerCustomData(bytes []byte) (buyerCustomData BuyerCustomData) {
	json.Unmarshal(bytes, &buyerCustomData)
	return
}

//se deserializa un string al tipo RetailerCustomData
func deserializeRetailerCustomData(bytes []byte) (retailerCustomData RetailerCustomData) {
	json.Unmarshal(bytes, &retailerCustomData)
	return
}

//Se deserializa un string al tipo Transaction
func deserializeTransaction(bytes []byte) (transaction Transaction) {
	json.Unmarshal(bytes, &transaction)
	return
}

//Se serializa un objecto a un array de bytes
func serializeCustomData(object struct{}) (serializedCustomData []byte) {
	serializedCustomData, _ = json.Marshal(object)
	return
}
