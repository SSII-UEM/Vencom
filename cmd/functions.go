package main

import (
	"encoding/json"
	"git.rwth-aachen.de/acs/public/cloud/mas/clonemap/pkg/agency"
	"git.rwth-aachen.de/acs/public/cloud/mas/clonemap/pkg/schemas"
	"strconv"
)

//devuelve una lista con los ids de los agentes vendedores
func getRetailers(agent *agency.Agent) (retailersId []int) {
	services, _ := agent.DF.SearchForService("retailer")
	retailersId = make([]int, len(services))
	for idx, service := range services {
		retailersId[idx] = service.AgentID
	}
	return
}

//compra una cantidad determinada de producto
func buyProduct(agent *agency.Agent, retailerId int, productId string, amount int) (boughtProducts []Product) {
	boughtProducts = make([]Product, 0)
	buyerData, _ := agentsGlobalData.getBuyerData(agent.GetAgentID())
	transaction := &Transaction{
		ProductId:       productId,
		BuyerCustomData: *buyerData,
		Amount:          amount,
	} // se crea la transacción
	str, _ := json.Marshal(&transaction) //se serializa
	//retailerResponse, err := sendMessage(agent, retailerId, BUY_PRODUCT_PROTOCOL, schemas.FIPAPerfCallForProposal, string(str), true)
	message, _ := agent.ACL.NewMessage(retailerId, BUY_PRODUCT_PROTOCOL, schemas.FIPAPerfCallForProposal, string(str))
	agent.ACL.SendMessage(message) //se envía la transacción al vendedor
	response, _ := agent.ACL.RecvMessageWait() //se espera la respuesta del vendedor
	agent.Logger.NewLog(app, strconv.Itoa(response.Performative), "")
	if response.Performative == schemas.FIPAPerfAgree /*OK*/ {
		var temp []Product
		json.Unmarshal([]byte(response.Content), &temp) //se deserializa la lista de productos comprados
		boughtProducts = append(boughtProducts, temp...)
	}
	return
}

// envía un mensaje
func sendMessage(localAgent *agency.Agent, remoteAgentId int, protocol int, performative int, content string, waitResponse bool) (message schemas.ACLMessage, err error) {
	var sentMessage, returnedMessage schemas.ACLMessage
	sentMessage, _ = localAgent.ACL.NewMessage(remoteAgentId, protocol, performative, content) //se crea el mensaje
	err = localAgent.ACL.SendMessage(sentMessage) //se envía el mensaje

	if err != nil && waitResponse {
		returnedMessage, err = localAgent.ACL.RecvMessageWait() //se bloquea el thread hasta que se reciba un mensaje
		if err != nil {
			message = returnedMessage
		}
	}

	return message, err
}