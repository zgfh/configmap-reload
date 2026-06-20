package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/shirou/gopsutil/process"
)

func mockProcess() (chan bool, chan bool) {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	result := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGUSR1)
	go func() {
		select {
		case sig := <-sigs:
			fmt.Println()
			fmt.Println("received signal:", sig)
			result <- true
		case <-done:
			fmt.Println("exiting")
			return
		}
	}()
	fmt.Println("awaiting signal")
	return result, done
}

func TestGetAvailableUpgrades(t *testing.T) {
	resultC, done := mockProcess()
	defer func() {
		done <- true
	}()

	p, _ := process.NewProcess(int32(os.Getpid()))
	processNmae, _ := p.Cmdline()
	fmt.Printf("process name: %s, pid: %d", processNmae, p.Pid)

	// don't know why, if not sleep,it can not get new process use 'process.Processes()'
	time.Sleep(5 * time.Second)

	err := signalHookExec(&signalHookParam{
		processName:  processNmae,
		signalNumber: 0x1e, //syscall.SIGUSR1
	})
	t.Logf("signalHookExec (%d : %s ) result: %+v", p.Pid, processNmae, err)
	result := false
	select {
	case <-resultC:
		result = true
	case <-time.After(5 * time.Second): // 等待 2 秒，给系统足够的时间投递和处理信号
		t.Logf("result not finish: %+v", err)
	}
	if !result {
		t.Errorf("signalHookExec err : %+v", err)
	}

}
