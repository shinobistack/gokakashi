package agent

import "fmt"

type Agent struct{}

func New() *Agent {
	return &Agent{}
}

func (a *Agent) Start() error {
	fmt.Println("TODO: implement v2 agent start logic")

	return nil
}
