package main

import (
	"git.rwth-aachen.de/acs/public/cloud/mas/clonemap/pkg/agency"
	"git.rwth-aachen.de/acs/public/cloud/mas/clonemap/pkg/schemas"
	"sort"
)

func computeDistance(agent *agency.Agent, agentId int) (distance float64) {
	message, _ := sendMessage(agent, agentId, 0, schemas.FIPAPerfRequest, "", true)
	if message.Performative == schemas.FIPAPerfAgree {
		var coordinates1, coordinates2 *Coordinates
		coordinates1, _ = agentsGlobalData.getCoordinates(agent.GetAgentID())
		coordinates2, _ = agentsGlobalData.getCoordinates(agentId)
		distance = coordinates1.ComputeDistance(coordinates2)
	}
	return
}

func generateProductMetricsFromRetailers(agent *agency.Agent, retailersId []int) (productsMetrics []ProductMetrics) {
	//TODO: redo
	productsMetrics = make([]ProductMetrics, 0)
	for _, retailerId := range retailersId {
		retailerData, _ := agentsGlobalData.getRetailerData(retailerId)
		for _, product := range retailerData.Products {
			productsMetrics = append(productsMetrics, ProductMetrics{
				retailerId: retailerId,
				product: product,
				healthRestorationScore: product.HealthRestoration / product.Price,
				buyerDistance: computeDistance(agent, retailerId),
			})
		}
	}
	return
}

func sortProductMetrics(metrics []ProductMetrics) []ProductMetrics {
	if len(metrics) > 0 {
		sort.Slice(metrics, func(i,j int) bool {
			return metrics[i].healthRestorationScore > metrics[j].healthRestorationScore && metrics[i].buyerDistance > metrics[j].buyerDistance
		})
	}
	return metrics
}