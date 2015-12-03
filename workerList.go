package main

import (
	"github.com/Rabbit-Go/Base"
	"github.com/Rabbit-Go/workers"
)

var workerReg = map[string]Base.WorkerConf{
	"Another": workers.Another,
	"Example": workers.Example,
}
