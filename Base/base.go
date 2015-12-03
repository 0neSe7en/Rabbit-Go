package Base

import (
	"github.com/Rabbit-Go/logger"
	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

type WorkerConf struct {
	QName     string              // Queue Name
	RoutingKs []string            // RoutingKey for binding
	AutoAck   bool                // AutoAck
	Handler   func(amqp.Delivery) // handler for message
	Desc      string              // Description for cli.
}

/* Work 接受三个参数. Exchange的配置信息从config.json中获取.分别是
   exchangeName
   exchangeType
   exchangeOptions.durable  -> default: true

Queue Default Config
   durable -> default: true

*/
func Work(conn *amqp.Connection, conf WorkerConf, errCh chan error) {
	baselog := logger.Log.WithFields(logrus.Fields{
		"tag":    "base:work",
		"worker": conf.QName,
	})
	ch, err := conn.Channel()
	err = ch.ExchangeDeclare(
		viper.GetString("exchangeName"),
		viper.GetString("exchangeType"),
		viper.GetBool("exchangeOptions.durable"),
		false, // autodelete
		false, // internal
		false, // noWait
		nil,
	)
	if err != nil {
		baselog.WithError(err).Fatal("create channel failed")
		errCh <- err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		conf.QName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,
	)
	if err != nil {
		baselog.WithError(err).Fatal("declare queue failed")
		errCh <- err
	}
	for _, s := range conf.RoutingKs {
		err = ch.QueueBind(q.Name, s, viper.GetString("exchangeName"), false, nil)
		if err != nil {
			baselog.Fatal("bind queue failed")
			errCh <- err
		}
	}
	baselog.WithFields(logrus.Fields{"routingKey": conf.RoutingKs}).Info("Binding queue")

	msgs, err := ch.Consume(
		q.Name,
		"",           // consumer string. if it is empty, rabbitmq generate a random string.
		conf.AutoAck, // autoAck
		false,        // exclusive
		false,        // noLocal
		false,        // noLocal
		nil,
	)
	if err != nil {
		baselog.WithError(err).Fatal("create consume failed")
		errCh <- err
	}

	forever := make(chan bool)

	for i := 0; i < viper.GetInt("workerCount"); i++ {
		baselog.Infof("[%d] Worker [%s] Started.", i, q.Name)
		go func() {
			for msg := range msgs {
				conf.Handler(msg)
			}
		}()
	}

	baselog.Info("[*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
