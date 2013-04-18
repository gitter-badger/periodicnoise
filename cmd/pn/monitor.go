package main

import (
	"log"
	"os/exec"
	"strings"
)

var monitoringCalls = map[monitoringResult]string{}
var monitoringEvent string

type monitoringResult string

const (
	monitorOk       monitoringResult = "OK"
	monitorCritical monitoringResult = "CRITICAL"
	monitorWarning  monitoringResult = "WARNING"
	monitorDebug    monitoringResult = "DEBUG"
	monitorUnknown  monitoringResult = "UNKNOWN"
)

var monitoringResults = []monitoringResult{monitorOk, monitorCritical,
	monitorWarning, monitorDebug, monitorUnknown}

func (m monitoringResult) String() string { return string(m) }

// Hook for passive monitoring solution
func monitor(state monitoringResult, message string) {
	// sanity check arguments
	switch state {
	case monitorOk, monitorCritical, monitorWarning, monitorDebug, monitorUnknown:
		// These are valid
	default:
		panic("unknown monitoring state")
	}

	call, exists := monitoringCalls[state]
	if !exists {
		return
	}

	call = strings.Replace(call, "%(event)", monitoringEvent, -1)
	call = strings.Replace(call, "%(state)", state.String(), -1)
	call = strings.Replace(call, "%(message)", message, -1)
	// do argument interpolation
	cmd := commander.Command("/bin/sh", "-c", call)
	err := cmd.Run()
	if err != nil {
		log.Fatalln("FATAL: Monitoring script failed, hope your monitoring system detects dead clients")
	}
}

// infrastructure for dependency injection for os.exec Command and run
type Executor interface {
	Run() error
}

type Commander interface {
	Command(name string, args ...string) Executor
}

var commander = execCommander{}

// default implementations
type execExecutor struct {
	*exec.Cmd
}

type execCommander struct{}

func (e execCommander) Command(name string, args ...string) Executor {
	return execExecutor{exec.Command(name, args...)}
}