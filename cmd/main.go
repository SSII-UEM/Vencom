package main

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/cloud/mas/clonemap/pkg/agency"
	"git.rwth-aachen.de/acs/public/cloud/mas/clonemap/pkg/schemas"
)



func main() {
	err := agency.StartAgency(task)
	if err != nil {
		fmt.Println(err)
	}
}

func task(ag *agency.Agent) (err error) {
	_type, _ := ag.GetAgentType()
	service := schemas.Service{
		Desc: _type,
		AgentID: ag.GetAgentID(),
	}
	_ , err = ag.DF.RegisterService(service) //se registran los agentes como servicios en el DF

	if err != nil {
		fmt.Println(err)
	} else {
		_type, _subtype := ag.GetAgentType()
		ag.Logger.NewLog(app, fmt.Sprintf("Initializing agent id=%d, type=%s, subtype=%s, custom=%s", ag.GetAgentID(), _type, _subtype, ag.GetCustomData()), "")

		if _type == "buyer" {
			id, data := ag.GetAgentID(), deserializeBuyerCustomData([]byte(ag.GetCustomData()))
			agentsGlobalData.registerBuyerData(id, data) //se registra la información del agente de manera global
			configureBuyerBehaviour(ag) //se asignan los comportamientos a los agentes compradores
		} else {
			id, data := ag.GetAgentID(), deserializeRetailerCustomData([]byte(ag.GetCustomData()))
			agentsGlobalData.registerRetailerData(id, data) //se registra la información del agente de manera global
			configureRetailerBehaviour(ag) //se asignan los comportamientos a los agentes vendedores
		}
	}

	return
}
