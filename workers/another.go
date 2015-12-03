package workers

import (
	"encoding/json"
	"github.com/Rabbit-Go/Base"
	"github.com/Rabbit-Go/Mongo"
	"github.com/Rabbit-Go/Redis"
	"github.com/Rabbit-Go/logger"
	"github.com/Sirupsen/logrus"
	"github.com/streadway/amqp"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type anotherMsg struct {
	Date      time.Time     `json:"date"`
	ExampleId bson.ObjectId `json:"example"`
}

var Another = Base.WorkerConf{
	QName:     "AnotherQ",
	RoutingKs: []string{"example.another"},
	AutoAck:   false,
	Handler:   anotherHanlder,
	Desc:      "This is another consumer.",
}

func anotherHanlder(msg amqp.Delivery) {
	relateLog := logger.Log.WithFields(logrus.Fields{
		"tag":    "worker",
		"worker": "anotherWorker",
	})
	var err error
	r := Redis.Pool.Get()
	defer r.Close()
	defer msg.Ack(true)

	var aMsg anotherMsg
	if err = json.Unmarshal(msg.Body, &aMsg); err != nil {
		relateLog.WithError(err).Error("json:unmarshal")
		return
	}

	session := Mongo.S.Clone()
	defer session.Close()

	_, err = session.DB("").C(Mongo.ColAnother).Insert(aMsg)

	if err != nil {
		relateLog.WithFields(logrus.Fields{
			"payload": aMsg,
		}).WithError(err).Error("mgo:insert")
		return
	}

	relateLog.Info(aMsg)
}
