package main

import (
	"encoding/json"
	"fmt"
	"git.rwth-aachen.de/acs/public/cloud/mas/clonemap/pkg/agency"
	"git.rwth-aachen.de/acs/public/cloud/mas/clonemap/pkg/schemas"
	"time"
)

func configureBuyerBehaviour(agent *agency.Agent) {
	data, _ := agentsGlobalData.getBuyerData(agent.GetAgentID())

	behaviour1, _ := agent.NewPeriodicBehavior(data.IncomePeriodicity, func() error {
		buyerData, _ := agentsGlobalData.getBuyerData(agent.GetAgentID())
		buyerData.CurrentBalance += buyerData.IncomeAmount
		agentsGlobalData.registerBuyerData(agent.GetAgentID(), *buyerData)
		return nil
	})
	behaviour1.Start()

	behaviour2, _ := agent.NewPeriodicBehavior(data.HealthReductionPeriodicity, func() error {
		buyerData, _ := agentsGlobalData.getBuyerData(agent.GetAgentID())
		buyerData.CurrentHealth -= buyerData.HealthReduction
		agentsGlobalData.registerBuyerData(agent.GetAgentID(), *buyerData)
		return nil
	})
	behaviour2.Start()

	behaviour3, _ := agent.NewPeriodicBehavior(10 * time.Second, func() error {
		buyerData, _ := agentsGlobalData.getBuyerData(agent.GetAgentID())
		fmt.Printf("Agent#%d current hp = %v\n", agent.GetAgentID(), buyerData.CurrentHealth)
		if buyerData.CurrentHealth < 50.0 {
			retailers := getRetailers(agent)
			//TODO: inefficient. redo
			for len(retailers) > 0 && buyerData.CurrentHealth < 90 {
				productMetrics := sortProductMetrics(generateProductMetricsFromRetailers(agent, retailers))
				fmt.Println(productMetrics)
				if len(productMetrics) > 0 {
					productMetric := productMetrics[0]
					var boughtProducts []Product = buyProduct(agent, productMetric.retailerId, productMetric.product.Id, 1)
					for idx, product := range boughtProducts {
						fmt.Printf("Agent#%d BoughtProduct#%d: %v", agent.GetAgentID(), idx, product)
					}
					for _, boughtProduct := range boughtProducts {
						buyerData.CurrentHealth += boughtProduct.HealthRestoration
					}
				}
			}
			agentsGlobalData.registerBuyerData(agent.GetAgentID(), *buyerData)
		}
		return nil
	})
	behaviour3.Start()
}

func configureRetailerBehaviour(agent *agency.Agent) {
	data, _ := agentsGlobalData.getRetailerData(agent.GetAgentID())

	behaviour1, _ := agent.NewMessageBehavior(REQUEST_COORDINATES_PROTOCOL, map[int]func(schemas.ACLMessage)error {
		schemas.FIPAPerfRequest: func(message schemas.ACLMessage) error {
			str, _ := json.Marshal(data.Coordinates)
			sendMessage(agent, message.Sender, REQUEST_COORDINATES_PROTOCOL, schemas.FIPAPerfAgree, string(str), false)
			return nil
		},
	}, func(message schemas.ACLMessage) error {
		return nil
	})
	behaviour1.Start()

	behaviour2, _ := agent.NewMessageBehavior(REQUEST_PRODUCT_LIST_PROTOCOL, map[int]func(message schemas.ACLMessage)error{
		schemas.FIPAPerfRequest: func(message schemas.ACLMessage) error {
			str, _ := json.Marshal(data.Products)
			sendMessage(agent, message.Sender, REQUEST_PRODUCT_LIST_PROTOCOL, schemas.FIPAPerfAgree, string(str), false)
			return nil
		},
	}, func(message schemas.ACLMessage)error{
		return nil
	})
	behaviour2.Start()

	behaviour3, _ := agent.NewMessageBehavior(BUY_PRODUCT_PROTOCOL, map[int]func(schemas.ACLMessage) error{
		schemas.FIPAPerfCallForProposal: func(message schemas.ACLMessage) error {
			fmt.Printf("Agent id=%d requested: %s\n", message.Sender, message.Content)
			var transaction Transaction
			json.Unmarshal([]byte(message.Content), &transaction)
			if product, has := data.Products[transaction.ProductId]; has && transaction.Amount > 0 && transaction.Amount <= product.Quantity && transaction.BuyerCustomData.CurrentBalance >= product.Price* float64(transaction.Amount) {
				var boughtProducts []Product = make([]Product, transaction.Amount)
				for i := 0; i < transaction.Amount; i++ {
					boughtProducts[i] = product
				}
				str, _ := json.Marshal(&boughtProducts)
				sendMessage(agent, message.Sender, BUY_PRODUCT_PROTOCOL, schemas.FIPAPerfAcceptProposal, string(str), false)
				fmt.Printf("[VALID] Agent id=%d requested: %s\n", message.Sender, message.Content)
			} else {
				fmt.Printf("[INVALID] Agent id=%d requested: %s\n", message.Sender, message.Content)
			}
			return nil
		},
	}, func(message schemas.ACLMessage) error {
		return nil
	})
	behaviour3.Start()
}