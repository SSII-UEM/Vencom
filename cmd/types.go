package main

import (
	"math"
	"sync"
	"time"
)

type Product struct {
	Id                string  `json:"id"`
	Name              string  `json:"name"`
	Quantity          int     `json:"quantity"`
	HealthRestoration float64 `json:"healthRestoration"`
	Price             float64 `json:"price"`
	/*productLifespan time.Duration*/
}

type Coordinates struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

func (from *Coordinates) ComputeDistance(to *Coordinates) float64 {
	return math.Sqrt(math.Exp2(from.X-to.X) + math.Exp2(from.Y-to.Y) + math.Exp2(from.Z-to.Z))
}

type RetailerCustomData struct {
	AgentId     int                 `json:"agentId"`
	Coordinates Coordinates         `json:"coordinates"`
	Products    map[string]Product `json:"products"`
}

type BuyerCustomData struct {
	AgentId                    int           `json:"agentId"`
	IncomeAmount               float64       `json:"incomeAmount"`
	HealthReduction            float64       `json:"healthReduction"`
	IncomePeriodicity          time.Duration `json:"incomePeriodicity"`
	HealthReductionPeriodicity time.Duration `json:"healthReductionPeriodicity"`
	CurrentBalance             float64       `json:"currentBalance"`
	CurrentHealth              float64       `json:"currentHealth"`
	Coordinates                Coordinates   `json:"coordinates"`
}

type AgentsGlobalData struct {
	agents map[int]AgentData
	mutex sync.Mutex
}

func NewAgentsGlobalData() *AgentsGlobalData {
	return &AgentsGlobalData{
		agents: make(map[int]AgentData),
		mutex: sync.Mutex{},
	}
}

func (agentsGlobalData *AgentsGlobalData) registerBuyerData(id int, data BuyerCustomData) *AgentsGlobalData {
	agentsGlobalData.mutex.Lock()
	agentsGlobalData.agents[id] = AgentData{
		isRetailer: false,
		buyerCustomData: data,
		coordinates: (data).Coordinates,
	}
	defer agentsGlobalData.mutex.Unlock()
	return agentsGlobalData
}

func (agentsGlobalData *AgentsGlobalData) registerRetailerData(id int, data RetailerCustomData) *AgentsGlobalData {
	agentsGlobalData.mutex.Lock()
	agentsGlobalData.agents[id] = AgentData{
		isRetailer: true,
		retailerCustomData: data,
		coordinates: (data).Coordinates,
	}
	defer agentsGlobalData.mutex.Unlock()
	return agentsGlobalData
}

func (agentsGlobalData *AgentsGlobalData) getBuyerData(id int) (*BuyerCustomData, bool) {
	agentsGlobalData.mutex.Lock()
	temp, has := agentsGlobalData.agents[id]
	defer agentsGlobalData.mutex.Unlock()
	return &temp.buyerCustomData, has
}

func (agentsGlobalData *AgentsGlobalData) getRetailerData(id int) (*RetailerCustomData, bool) {
	agentsGlobalData.mutex.Lock()
	temp, has := agentsGlobalData.agents[id]
	defer agentsGlobalData.mutex.Unlock()
	return &temp.retailerCustomData, has
}

func (agentsGlobalData *AgentsGlobalData) getCoordinates(id int) (*Coordinates, bool) {
	agentsGlobalData.mutex.Lock()
	temp, has := agentsGlobalData.agents[id]
	defer agentsGlobalData.mutex.Unlock()
	return &temp.coordinates, has
}

type AgentData struct {
	isRetailer bool
	retailerCustomData RetailerCustomData
	buyerCustomData BuyerCustomData
	coordinates Coordinates
}

type Performative struct {
	code int
	name string
}

type Protocol struct {
	code int
	name string
	performatives []Performative
}

type Transaction struct {
	RetailerId      int             `json:"retailerId"`
	ProductId       string          `json:"productId"`
	Amount          int             `json:"amount"`
	BuyerCustomData BuyerCustomData `json:"buyerCustomData"`
}

type ProductMetrics struct {
	retailerId int
	product Product
	healthRestorationScore, buyerDistance float64
}