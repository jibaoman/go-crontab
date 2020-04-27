package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"time"
)

func main() {
	var (
		config clientv3.Config
		client *clientv3.Client
		err    error
		kv clientv3.KV
		getResp *clientv3.GetResponse
		watchStartRevision int64
		watcher clientv3.Watcher
		watchRespChan <- chan clientv3.WatchResponse
		watchResp clientv3.WatchResponse
		event *clientv3.Event
	)
	config = clientv3.Config{
		Endpoints:   []string{"172.29.3.100:2379"},
		DialTimeout: 5 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		fmt.Println(err)
		return
	}

	kv = clientv3.NewKV(client)

	go func() {
		for {
			kv.Put(context.TODO(),"/cron/jobs/job7","i am job7")

			kv.Delete(context.TODO(),"/cron/jobs/job7")

			time.Sleep(1 * time.Second)
		}
	}()

	if getResp,err = kv.Get(context.TODO(),"/cron/jobs/job7"); err != nil {
		fmt.Println(err)
		return
	}

	if len(getResp.Kvs) != 0 {
		fmt.Println("当前值：",string(getResp.Kvs[0].Value))
	}

	//当前ETCD集群事务ID，单调递增的
	watchStartRevision = getResp.Header.Revision + 1

	//创建一个watcher
	watcher = clientv3.NewWatcher(client)

	//启动监听
	fmt.Println("从该版本开始往后监听",watchStartRevision)

	ctx,cancelFunc := context.WithCancel(context.TODO())
	time.AfterFunc(5 * time.Second, func() {
		cancelFunc()
		fmt.Println("5秒时间已到")
	})

	watchRespChan = watcher.Watch(ctx,"/cron/jobs/job7",clientv3.WithRev(watchStartRevision))

	for watchResp =  range watchRespChan {
		for _,event =  range watchResp.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("修改为：",string(event.Kv.Value),"Revision:",event.Kv.CreateRevision,event.Kv.ModRevision)
			case mvccpb.DELETE:
				fmt.Println("删除了","Revision",event.Kv.ModRevision)

			}
		}
	}

}