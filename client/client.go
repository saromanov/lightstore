package client

//Implementation of client to Lightstore

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"log"
)

type (
	KV map[string]string
)

type Client struct {
	addr string
}

func NewClient(addr string) *Client {
	cl := new(Client)
	cl.addr = addr
	_, err := cl.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return cl
}

func (cl *Client) Ping()(bool, error){
	resp, err := http.Get(fmt.Sprintf("%s/ping", cl.addr))
	if err != nil{
		return false, errors.New(fmt.Sprintf("Server by address %s is not ready", cl.addr))
	}

	resp.Body.Close()
	return true, nil
}

//Set basic key-value to lightstore
func (cl *Client) Set(key, value string)(int, error) {
	jsonStr := fmt.Sprintf(`{"%s":"%s"}`, key, value)
	return cl.set(jsonStr)
}

func (cl *Client) set(jsonStr string) (int, error) {
	url := fmt.Sprintf("%s/set", cl.addr)
	return cl.sendRequest(url, bytes.NewBuffer([]byte(jsonStr)))
}

func (cl *Client) sendRequest(url string, buff *bytes.Buffer)(int, error){
	req, err := cl.NewRequest("POST", url, buff)
	req.Header.Set("X-Custom-Header", "lightstore")
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if strings.HasPrefix(resp.Status, "200") {
		return 1, nil
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		return 0, errors.New(string(body))
	}
}

func (cl *Client) NewRequest(method string, url string, buff *bytes.Buffer)(*http.Request, error) {
	if buff == nil {
		return http.NewRequest("POST", url, nil)
	}

	return http.NewRequest("POST", url, buff)
}

//SetMap provides append to lightstore pairs key-value
func (cl *Client) SetMap(values KV) (int, error) {
	result := "{"
	c := 0
	for key, value := range values {
		if c > 0 {
			result += ","
		}
		result += fmt.Sprintf(`"%s":"%s"`, key, value)
		c += 1
	}
	result += "}"
	return cl.set(result)
}

//Get valur by key
func (cl *Client) Get(key string) string {
	return cl.get(fmt.Sprintf("http://127.0.0.1:8080/get/%s", key))
}

//Stat is not ready now
func (cl *Client) Stat() string {
	return cl.get("http://127.0.0.1:8080/_stat")
}

func (*Client) get(url string) string {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Custom-Header", "lightstore")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	if strings.HasPrefix(resp.Status, "200") {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		result := string(body)
		return result[1 : len(result)-1]
	} else {
		return ""
	}
}

//CreatePage provides create new page on the lightstore
func (cl *Client) CreatePage(pagename string) (int, error) {
	url := fmt.Sprintf("%s/create/%s", cl.addr, pagename)
	return cl.sendRequest(url, nil)
}

func (cl *Client) SetToPage(pagename, key, value string) (int, error) {
	url := fmt.Sprintf("%s/set/%s", cl.addr, pagename)
	jsonStr := fmt.Sprintf(`{"%s":"%s"}`, key, value)
	return cl.sendRequest(url, bytes.NewBuffer([]byte(jsonStr)))
}

func (cl *Client) GetFromPage(pagename, key string)string{
	url := fmt.Sprintf("%s/dbget/%s/%s", cl.addr, pagename, key)
	return cl.get(url)
}

//SerializeKey provides serialization key data
func (cl *Client) SerializeKey(key string, value string) string{
	url := fmt.Sprintf("%s/serialize/%s/%s", cl.addr, key, value)
	return cl.get(url)
}

func (cl *Client) GetHistory() string {
	url := fmt.Sprintf("%s/history", cl.addr)
	return cl.get(url)
}

func (cl *Client) Find(key string) string {
	url := fmt.Sprintf("%s/find/%s", cl.addr, key)
	return cl.get(url)
}

func (cl *Client) Subscribe(key string){
	url := fmt.Sprintf("%s/subscribe/%s", cl.addr, key)
	cl.get(url)
}