package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	var (
		config clientv3.Config
		client *clientv3.Client
		err error
		kv clientv3.KV
		putOp clientv3.Op
		getOp clientv3.Op
		opResp clientv3.OpResponse
	)
	config  =  clientv3.Config{
		Endpoints:            []string{"172.29.3.100:2379"},
		DialTimeout:          5 * time.Second,
	}

	if client,err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	kv = clientv3.NewKV(client)

	putOp = clientv3.OpPut("/cron/jobs/job8","123456")

	if opResp,err = kv.Do(context.TODO(),putOp);err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("写入Revision:", opResp.Put().Header.Revision)

	getOp = clientv3.OpGet("/cron/jobs/job8")

	if opResp,err = kv.Do(context.TODO(),getOp);err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("数据Revision:",opResp.Get().Kvs[0].ModRevision)
	fmt.Println("数据value:",string(opResp.Get().Kvs[0].Value))
}
