package main

import (
	"errors"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.Level(6))
	logrus.SetReportCaller(true)
	logrus.SetOutput(io.MultiWriter(os.Stdout))

	type student struct {
		Age  int    `json:"age"`
		Name string `json:"name"`
	}
	stu := &student{
		Age:  11,
		Name: "zhnagsan",
	}
	logrus.WithFields(logrus.Fields{
		"err":     errors.New("Unmarshal failure"),
		"student": stu,
	}).Error("Unmarshal failure")
}
