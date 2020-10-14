package RabbitMQ

import (
	"fmt"
	"github.com/streadway/amqp"
	"time"
)

type IMessageQueue interface {
	ConnectMQ(url string,exchangeName string ,exchangeKind string)
	QueueDeclare(routerKey string,QueueName string)
	Produce(message string,routerKey string)
	Consume(queueName string,key string) <-chan amqp.Delivery
	Destory()
}

type MessageQueue struct {
	conn  *amqp.Connection
	channel *amqp.Channel
	queue  			amqp.Queue
	exchangeName 	string
	exchangeKind	string
	AuthAndRouterMap 	map[int]string
	GameAndRouterMap 	map[int]string
	CommonAndRouterMap 	map[int]string
}

func NewMessageQueue()IMessageQueue{
	return &MessageQueue{}
}

func (this *MessageQueue)ConnectMQ(url string,exchangeName string ,exchangeKind string){
	conn,err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil{
		panic(err)
	}
	this.conn = conn
	this.channel, err = conn.Channel()
	if err != nil {
		panic(err)
	}
	this.exchangeName = exchangeName
	this.exchangeKind = exchangeKind
	err = this.channel.ExchangeDeclare("test","topic",false, false, false, true, nil)
	if err != nil {
		panic(err)
	}
	mq,err := this.channel.QueueDeclare("test",false,false,false,true,nil)
	if nil != err{
		println(err.Error())
	}
	this.queue = mq
	//创建完成后，需要把queue绑定到对应的exchange中
	err = this.channel.QueueBind("test","","test",true,nil)
	if nil != err{
		println(err.Error())
	}

}

/**
 * 创建消息队列
 */
func (this *MessageQueue)QueueDeclare(routerKey string,QueueName string){

}

func (this *MessageQueue)Produce(message string,routerKey string){
	for ; ; {
		fmt.Println(this.channel.Publish("test","test",false,false,amqp.Publishing{Body:[]byte("Go Go AMQP!")}))
		time.Sleep(2*time.Second)
	}
}

func (this *MessageQueue)Consume(queueName string,key string) <-chan amqp.Delivery {
	data,err := this.channel.Consume(queueName,key,true,false,false,true,nil)
	if err== nil{
		return  data
	} else{
		return nil
	}
}

func (this *MessageQueue)Destory(){
	//this.channel.Close()
	//this.conn.Close()
}
