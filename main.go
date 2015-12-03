package main

import (
	"fmt"
	"github.com/Rabbit-Go/Base"
	"github.com/Rabbit-Go/Mongo"
	"github.com/Rabbit-Go/Redis"
	"github.com/Rabbit-Go/logger"
	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
	"gopkg.in/mgo.v2"
	"os"
)

var runNumber int

var runWorker string
var mainlog = logger.Log.WithFields(logrus.Fields{
	"tag": "main",
})

func main() {
	cmdRun := &cobra.Command{
		Use:   "run",
		Short: "Run workers",
		Long:  `run all worker, Or run specific worker and specific worker number`,
		Run:   run,
	}
	cmdRun.Flags().IntVarP(&runNumber, "number", "n", 1, "specific worker number (default 1)")
	cmdRun.Flags().StringVarP(&runWorker, "worker", "w", "", "specific a worker to run")
	cmdLs := &cobra.Command{
		Use:   "ls",
		Short: "list all registered workers",
		Run:   ls,
	}

	rootCmd := &cobra.Command{Use: "rabbit"}
	rootCmd.AddCommand(cmdRun)
	rootCmd.AddCommand(cmdLs)
	rootCmd.Execute()
}

func run(cmd *cobra.Command, args []string) {
	var conf Base.WorkerConf
	var exists bool
	if runWorker != "" {
		if conf, exists = workerReg[runWorker]; !exists {
			fmt.Printf("%s doesn't exists.", runWorker)
			os.Exit(1)
		}
	}
	initConf(runNumber)

	conn, err := amqp.Dial(viper.GetString("amqpUrl"))
	if err != nil {
		mainlog.WithError(err).Fatal("Failed to connect to RabbitMQ")
	}
	defer conn.Close()

	if Mongo.S, err = mgo.Dial(viper.GetString("mongodb")); err != nil {
		mainlog.WithError(err).Fatal("Failed to connect to MongoDB")
	}
	defer Mongo.S.Close()

	// Redis.Pool 不会返回错误.
	Redis.Pool = Redis.NewPool(viper.GetString("redisServer"), viper.GetString("redisPassword"))
	defer Redis.Pool.Close()

	errCh := make(chan error)

	if runWorker == "" {
		for name, conf := range workerReg {
			mainlog.Infof("Worker [%s] is starting...", name)
			go Base.Work(conn, conf, errCh)
		}
	} else {
		mainlog.Infof("Worker [%s] is Starting...", runWorker)
		go Base.Work(conn, conf, errCh)
	}

	err = <-errCh
	mainlog.WithError(err).Fatal("Failed to init worker")
}

func ls(cmd *cobra.Command, args []string) {
	for name, conf := range workerReg {
		fmt.Printf("- %s: %s\n", name, conf.Desc)
	}
}
