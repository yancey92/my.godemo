package jpushkit

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

//Client 推送客户端
type Client struct {
	AppKey    string
	Secret    string
	pushURL   string
	deviceURL string
}

//NewClient 新建客户端
func NewClient(appKey, secret string) *Client {
	pClient := &Client{AppKey: appKey, Secret: secret}
	pClient.pushURL = CONSTANT_PUSH_URL
	pClient.deviceURL = CONSTANT_DEVICE_URL
	return pClient
}

func (c *Client) request(method, url string, body io.Reader, group bool) (*Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.getAuthor(group))
	req.Header.Set("Charset", "UTF-8")
	req.Header.Set("Content-Type", "application/json")

	cl := &http.Client{}
	resp, err := cl.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &Response{data: buf}, nil
}

func (c *Client) getAuthor(group bool) string {
	auth := c.AppKey + ":" + c.Secret
	if group {
		auth = CONSTANT_GROUP + "-" + auth
	}
	buf := []byte(auth)
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString(buf))
}

//Push 推送
func (c *Client) Push(req *PushRequest) (map[string]interface{}, error) {
	url := CONSTANT_V3_PUSH
	buf, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := c.request("POST", url, bytes.NewReader(buf), false)
	if err != nil {
		return nil, err
	}
	return resp.Map()
}

//Validate 推送格式是否正确
func (c *Client) Validate(req *PushRequest) (map[string]interface{}, error) {
	url := c.pushURL + "/v3/push/validate"
	buf, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	resp, err := c.request("POST", url, bytes.NewReader(buf), false)
	if err != nil {
		return nil, err
	}
	return resp.Map()
}
