package main

import (
	"encoding/json"
	"git.rwth-aachen.de/acs/public/cloud/mas/clonemap/pkg/agency"
	"git.rwth-aachen.de/acs/public/cloud/mas/clonemap/pkg/schemas"
	"strconv"
)

func getRetailers(agent *agency.Agent) (retailersId []int) {
	services, _ := agent.DF.SearchForService("retailer")
	retailersId = make([]int, len(services))
	for idx, service := range services {
		retailersId[idx] = service.AgentID
	}
	return
}

func buyProduct(agent *agency.Agent, retailerId int, productId string, amount int) (boughtProducts []Product) {
	boughtProducts = make([]Product, 0)
	buyerData, _ := agentsGlobalData.getBuyerData(agent.GetAgentID())
	transaction := &Transaction{
		ProductId:       productId,
		BuyerCustomData: *buyerData,
		Amount:          amount,
	}
	str, _ := json.Marshal(&transaction)
	//retailerResponse, err := sendMessage(agent, retailerId, BUY_PRODUCT_PROTOCOL, schemas.FIPAPerfCallForProposal, string(str), true)
	message, _ := agent.ACL.NewMessage(retailerId, BUY_PRODUCT_PROTOCOL, schemas.FIPAPerfCallForProposal, string(str))
	agent.ACL.SendMessage(message)
	response, _ := agent.ACL.RecvMessageWait()
	agent.Logger.NewLog(app, strconv.Itoa(response.Performative), "")
	if response.Performative == schemas.FIPAPerfAgree /*OK*/ {
		var temp []Product
		json.Unmarshal([]byte(response.Content), &temp)
		boughtProducts = append(boughtProducts, temp...)
	}
	return
}

func sendMessage(localAgent *agency.Agent, remoteAgentId int, protocol int, performative int, content string, waitResponse bool) (message schemas.ACLMessage, err error) {
	var sentMessage, returnedMessage schemas.ACLMessage
	sentMessage, _ = localAgent.ACL.NewMessage(remoteAgentId, protocol, performative, content)
	err = localAgent.ACL.SendMessage(sentMessage)

	if err != nil && waitResponse {
		returnedMessage, err = localAgent.ACL.RecvMessageWait()
		if err != nil {
			message = returnedMessage
		}
	}

	return message, err
}