package dbkit

import (
	//"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"testing"
	"time"
)

var sslUrl string = "mongodb://gumpmongodev:y04eTIOHkDTXsTMuXWNtlDtiWSasByuIQ5owemMbwdxJCNcUDurJMDhXO6zXzzYDjILolp3yLT31Dk9ETuKSJQ==@gumpmongodev.documents.azure.cn:10250/"

type Person struct {
	Name    string
	Phone   string
	SortIdx int
}

func TestMongoUpsertDoc(t *testing.T) {
	//InitMongoDB("mongodb://localhost:27017", "gumpcome")
	InitMongoDBWithSSL(sslUrl, "gumpcome")
	_, isOk, err := MongoUpsertDoc(&MongoSearch{
		Collection: "people",
		Key:        "phone",
		Value:      "123456789",
	}, &Person{"guanyu", "123456789", 10})
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(isOk)
}

func TestMongoFindDoc(t *testing.T) {
	//InitMongoDB("mongodb://localhost:27017", "gumpcome")
	InitMongoDBWithSSL(sslUrl, "gumpcome")
	result := Person{}
	MongoFindHandler("people", func(c *mgo.Collection) {
		c.Find(bson.M{"phone": "123456789"}).One(&result)
	})
	fmt.Println(result)
}

func TestMongoRemoveAllDoc(t *testing.T) {
	InitMongoDB("mongodb://localhost:27017", "gumpcome", 10)
	isOk, err := MongoRemoveAllDoc(&MongoSearch{
		Collection: "people",
		Key:        "phone",
		Value:      "123456789",
	})
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(isOk)
}

func TestBetchFindMongo(t *testing.T) {
	InitMongoDBWithSSL(sslUrl, "gumpcome")

	for i := 0; i < 1000; i++ {
		go func(num int) {
			result := Person{}
			MongoFindHandler("people", func(c *mgo.Collection) {
				c.Find(bson.M{"sortidx": num}).One(&result)
			})
			fmt.Printf("%v-%v\n", num, result)

			//err := MongoInsert("people", &Person{"guanyu", "123456789", num})
			//if err != nil {
			//	fmt.Printf("插入失败 %v\n", num)
			//}
		}(i)
	}
	time.Sleep(100 * time.Second)
}
