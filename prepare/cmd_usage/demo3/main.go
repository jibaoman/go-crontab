package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)


type result struct {
	err error
	output []byte
}

func main() {
	//执行一个CMD,让他执行一个协程
	var (
		ctx context.Context
		cancelFunc context.CancelFunc
		cmd *exec.Cmd
		resultChan chan *result
		res *result
	)
	//创建一个结果队列
	resultChan = make(chan *result,1000)
	ctx,cancelFunc = context.WithCancel(context.TODO())

	go func() {
		var (
			output []byte
			err error
		)
		cmd = exec.CommandContext(ctx,"C:\\cygwin64\\bin\\bash.exe","-c","sleep 2;echo hello;")
		output,err = cmd.CombinedOutput()

		resultChan <- &result{
			err:    err,
			output: output,
		}
	}()

	time.Sleep(1*time.Second)

	cancelFunc()

	res = <- resultChan
	fmt.Println(res.err,string(res.output))
}

