package RouterManager

import (
	"SLGGAME/CoreManager/RabbitMQ"
	"SLGGAME/Service"
)

type RouterService struct {
}

func NewServer(type_ string, serverGroup string) Service.IService {
	return &RouterService{}
}

var  rabbitMQ string = "localhost:5672"

func (r RouterService) RegisterProtoHandler() {

}

func (r RouterService) BeforeRunThreadHook() {

}

func (r RouterService) AfterRunThreadHook() {

}

func (r RouterService) Run() {
	rabbitMqMgr := RabbitMQ.NewMessageQueue()
	rabbitMqMgr.ConnectMQ(rabbitMQ,"test","topic")
	rabbitMqMgr.QueueDeclare("test","test")
	rabbitMqMgr.Produce("test","test")
}

func (r RouterService) ObserverPProf(addr string) {

}



