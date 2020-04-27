package main


import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"time"
)

type TimePoint struct {
	StartTime int64 `bson:"startTime"`
	EndTime int64 `bson:"endTime"`
}

type LogRecord struct {
	JobName string `bson:"jobName"`
	Command string `bson:"command"`
	Err string `bson:"err"`
	Content string `bson:"content"`
	TimePoint TimePoint `bson:"timePoint"`
}
func main() {
	var (
		client *mongo.Client
		err error
		database *mongo.Database
		collection *mongo.Collection
		record *LogRecord
		result *mongo.InsertManyResult
		docId objectid.ObjectID
		logArr []interface{}
		insertId interface{}
	)
	if client,err = mongo.Connect(context.TODO(),"mongodb://172.29.3.100:27017",clientopt.ConnectTimeout(5 * time.Second)); err != nil {
		fmt.Println(err)
		return
	}

	database = client.Database("cron")

	collection = database.Collection("log")

	//插入记录
	record = &LogRecord{
		JobName:   "job10",
		Command:   "echo hello",
		Err:       "",
		Content:   "hello",
		TimePoint: TimePoint{StartTime:time.Now().Unix(),EndTime:time.Now().Unix()+5},
	}

	logArr = []interface{}{record,record,record}

	if result,err = collection.InsertMany(context.TODO(),logArr); err != nil {
		fmt.Println(err)
		return
	}

	for _,insertId =  range result.InsertedIDs {
		docId = insertId.(objectid.ObjectID)
		fmt.Println("自增ID:",docId.Hex())
	}



}

