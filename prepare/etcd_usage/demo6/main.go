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
		lease clientv3.Lease
		leaseGranResp *clientv3.LeaseGrantResponse
		leaseId clientv3.LeaseID
		putResp *clientv3.PutResponse
		getResp *clientv3.GetResponse
		keepResp *clientv3.LeaseKeepAliveResponse
		keepRespChan <- chan *clientv3.LeaseKeepAliveResponse
		kv clientv3.KV
	)

	config = clientv3.Config{
		Endpoints:   []string{"172.29.3.100:2379"}, // 集群列表
		DialTimeout: 5 * time.Second,
	}

	// 建立一个客户端
	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	//申请租约
	lease = clientv3.NewLease(client)
	if leaseGranResp,err = lease.Grant(context.TODO(),10);err != nil {
		fmt.Println(err)
		return
	}

	leaseId = leaseGranResp.ID

	if keepRespChan,err = lease.KeepAlive(context.TODO(),leaseId); err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		for {
			select {
			case keepResp = <- keepRespChan:
				if keepRespChan == nil {
					fmt.Println("租约已经失效")
					goto END
				} else {
					fmt.Println("收到续约应答：",keepResp.ID)
				}

			}
		}

		END:
	}()

	kv = clientv3.NewKV(client)

	if putResp,err = kv.Put(context.TODO(),"/cron/jobs/job1","",clientv3.WithLease(leaseId));err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("写入成功：",putResp.Header.Revision)


	for {
		if getResp,err = kv.Get(context.TODO(),"/cron/jobs/job1"); err != nil {
			fmt.Println(err)
			return
		}

		if getResp.Count == 0 {
			fmt.Println("kv过期了")
			break
		}

		fmt.Println("还没过期",getResp.Kvs)
		time.Sleep(2 * time.Second)
	}
}