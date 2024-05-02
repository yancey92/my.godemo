package model

type ReqBody struct {
	Name string `json:"name" example:"张三"` // 名字
	Age  int    `json:"age" example:"20"`  // 年龄
}
