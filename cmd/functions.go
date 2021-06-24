package main

import (
	"encoding/json"
	"fmt"
	"git.rwth-aachen.de/acs/public/cloud/mas/clonemap/pkg/agency"
	"git.rwth-aachen.de/acs/public/cloud/mas/clonemap/pkg/schemas"
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

	buyerData, _ := agentsGlobalData.getBuyerData(agent.GetAgentID())
	transaction := &Transaction{
		ProductId:       productId,
		BuyerCustomData: *buyerData,
		Amount:          amount,
	}
	str, _ := json.Marshal(&transaction)
	retailerResponse, _ := sendMessage(agent, retailerId, BUY_PRODUCT_PROTOCOL, schemas.FIPAPerfCallForProposal, string(str), true)
	if retailerResponse.Performative == schemas.FIPAPerfAcceptProposal /*OK*/ {
		json.Unmarshal([]byte(retailerResponse.Content), boughtProducts)
	}

	return
}

func sendMessage(localAgent *agency.Agent, remoteAgentId int, protocol int, performative int, content string, waitResponse bool) (message schemas.ACLMessage, err error) {
	var sentMessage, returnedMessage schemas.ACLMessage
	fmt.Printf("Sending message from %d to %d: %s\n", localAgent.GetAgentID(), remoteAgentId, content)
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