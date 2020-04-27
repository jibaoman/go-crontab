package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"time"
)

func main() {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		err     error
		kv      clientv3.KV
		//putResp *clientv3.PutResponse
		getResp *clientv3.GetResponse
	)
	config = clientv3.Config{
		Endpoints:   []string{"172.29.3.100:2379"},
		DialTimeout: 5 * time.Second,
	}

	//建立一个客户端
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	//用于读写etcd
	kv = clientv3.NewKV(client)

	//
	kv.Put(context.TODO(),"/cron/jobs/job4","job4")

	if getResp,err = kv.Get(context.TODO(),"/cron/jobs/",clientv3.WithPrefix());err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(getResp.Kvs)
	}
}
