package main

const app string = "app"

var agentsGlobalData AgentsGlobalData = *NewAgentsGlobalData()

const (
	REQUEST_COORDINATES_PROTOCOL = iota
	REQUEST_PRODUCT_LIST_PROTOCOL = iota
	BUY_PRODUCT_PROTOCOL = iota
)