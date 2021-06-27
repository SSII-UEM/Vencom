package main

import (
	"encoding/json"
	"fmt"
	"git.rwth-aachen.de/acs/public/cloud/mas/clonemap/pkg/agency"
	"git.rwth-aachen.de/acs/public/cloud/mas/clonemap/pkg/schemas"
	"strconv"
	"time"
)

func configureBuyerBehaviour(agent *agency.Agent) {
	data, _ := agentsGlobalData.getBuyerData(agent.GetAgentID()) //se recupera la información del agente

	behaviour1, _ := agent.NewPeriodicBehavior(data.IncomePeriodicity, func() error {
		buyerData, _ := agentsGlobalData.getBuyerData(agent.GetAgentID())
		buyerData.CurrentBalance += buyerData.IncomeAmount //se aumenta el capital del agente
		agentsGlobalData.registerBuyerData(agent.GetAgentID(), *buyerData)
		agent.Logger.NewLog(app, fmt.Sprintf("Health = %f, balance = %f", buyerData.CurrentHealth, buyerData.CurrentBalance), strconv.FormatInt(time.Now().Unix(), 10))
		return nil
	})
	behaviour1.Start()

	behaviour2, _ := agent.NewPeriodicBehavior(data.HealthReductionPeriodicity, func() error {
		buyerData, _ := agentsGlobalData.getBuyerData(agent.GetAgentID())
		buyerData.CurrentHealth -= buyerData.HealthReduction //se reduce la vida del agente
		agentsGlobalData.registerBuyerData(agent.GetAgentID(), *buyerData)
		agent.Logger.NewLog(app, fmt.Sprintf("Health = %f, balance = %f", buyerData.CurrentHealth, buyerData.CurrentBalance), strconv.FormatInt(time.Now().Unix(), 10))
		return nil
	})
	behaviour2.Start()

	behaviour3, _ := agent.NewPeriodicBehavior(5 * time.Second, func() error {
		buyerData, _ := agentsGlobalData.getBuyerData(agent.GetAgentID())
		if buyerData.CurrentHealth < 50.0 {
			retailers := getRetailers(agent) //se obtienen los ids de los agentes vendedores
			//TODO: inefficient. redo
			if len(retailers) > 0 && buyerData.CurrentHealth < 90 {
				productMetrics := sortProductMetrics(generateProductMetricsFromRetailers(agent, retailers)) //se obtienen las métricas de los productos de manera ordenada
				agent.Logger.NewLog(app, fmt.Sprintf("Collected products from retailers: %v", productMetrics), "")
				if len(productMetrics) > 0 {
					var productMetric ProductMetrics = productMetrics[0]
					var boughtProducts []Product = buyProduct(agent, productMetric.retailerId, productMetric.product.Id, 1) //se compra 1 unidad del producto
					agent.Logger.NewLog(app, fmt.Sprintf("Bought %d product/s from retailer#%d. Products: %v", len(boughtProducts), productMetric.retailerId, boughtProducts), "")
					agent.Logger.NewLog(app, fmt.Sprintf("[BEFORE] Health = %f, balance = %f", buyerData.CurrentHealth, buyerData.CurrentBalance), strconv.FormatInt(time.Now().Unix(), 10))
					for _, boughtProduct := range boughtProducts {
						buyerData.CurrentBalance -= boughtProduct.Price
						buyerData.CurrentHealth += boughtProduct.HealthRestoration
					}
					agent.Logger.NewLog(app, fmt.Sprintf("[AFTER] Health = %f, balance = %f", buyerData.CurrentHealth, buyerData.CurrentBalance), strconv.FormatInt(time.Now().Unix(), 10))
				}
			}
			//TODO: this overrides an older state on new ones. Redo to safely support concurrent modifications
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
			sendMessage(agent, message.Sender, REQUEST_COORDINATES_PROTOCOL, schemas.FIPAPerfAgree, string(str), false) //se envía la respuesta con las coordenadas
			return nil
		},
	}, func(message schemas.ACLMessage) error {
		return nil
	})
	behaviour1.Start()

	behaviour2, _ := agent.NewMessageBehavior(REQUEST_PRODUCT_LIST_PROTOCOL, map[int]func(message schemas.ACLMessage)error{
		schemas.FIPAPerfRequest: func(message schemas.ACLMessage) error {
			str, _ := json.Marshal(data.Products)
			sendMessage(agent, message.Sender, REQUEST_PRODUCT_LIST_PROTOCOL, schemas.FIPAPerfAgree, string(str), false) //se envía la respuesta con la lista de productos
			return nil
		},
	}, func(message schemas.ACLMessage)error{
		return nil
	})
	behaviour2.Start()

	behaviour3, _ := agent.NewMessageBehavior(BUY_PRODUCT_PROTOCOL, map[int]func(schemas.ACLMessage) error{
		schemas.FIPAPerfCallForProposal: func(message schemas.ACLMessage) error {
			var transaction Transaction
			json.Unmarshal([]byte(message.Content), &transaction)
			agent.Logger.NewLog(app, "Received transaction", fmt.Sprintf("%v", transaction))
			product, has := data.Products[transaction.ProductId]
			agent.Logger.NewLog(app, fmt.Sprintf("Checking product#%s", (&transaction).ProductId), fmt.Sprintf("%v", product))
			//Se comprueba que la transacción es válida
			if has && (&transaction).Amount > 0 && (&transaction).Amount <= product.Quantity && transaction.BuyerCustomData.CurrentBalance >= product.Price * float64((&transaction).Amount) {
				var boughtProducts []Product = make([]Product, 0)
				for i := 0; i < (&transaction).Amount; i++ {
					boughtProducts = append(boughtProducts, product)
				}
				str, _ := json.Marshal(&boughtProducts)
				agent.Logger.NewLog(app,"Returning products", fmt.Sprintf("%v", &boughtProducts))
				sendMessage(agent, message.Sender, BUY_PRODUCT_PROTOCOL, schemas.FIPAPerfAgree, string(str), false) //se envía la lista de productos vendidos
			}
			return nil
		},
	}, func(message schemas.ACLMessage) error {
		return nil
	})
	behaviour3.Start()
}