package qa

import (
	"github.com/donyori/cqa/qa/qm"
)

func Init() {
	qm.Init()
}

func Shutdown() {
	<-qm.Exit(qm.ExitModeGracefully)
}

func ShutdownImmediately() {
	<-qm.Exit(qm.ExitModeImmediately)
}

func Exit() {
	<-qm.Exit(qm.ExitModeForcedly)
}
