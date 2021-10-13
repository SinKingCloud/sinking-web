package sinking_sdk_go

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// rpcRequestBuild rpc消息构建
type rpcRequestBuild struct {
	name     string
	url      string
	method   string
	header   map[string]string
	service  *Service
	param    *Param
	register *Register
}

// Rpc 构建远程调用服务名
func (r *Register) Rpc(name string) *rpcRequestBuild {
	return &rpcRequestBuild{register: r, name: name}
}

// Header 构建远程调用header
func (r *rpcRequestBuild) Header(header map[string]string) *rpcRequestBuild {
	r.header = header
	return r
}

// Method 构建远程调用Method
func (r *rpcRequestBuild) Method(method string) *rpcRequestBuild {
	r.method = method
	return r
}

// Call 远程调用
func (r *rpcRequestBuild) Call(url string, param *Param) (string, error) {
	r.service, _ = r.register.GetService(r.name)
	if r.service == nil {
		return "", errors.New("not found online service")
	}
	r.url = url
	r.param = param
	body, err := r.sendRequest()
	if err != nil {
		//下线服务
		r.register.changeServerStatus(r.service.ServiceHash, 1)
		return "", err
	}
	return body, nil
}

// sendRequest 发送请求
func (r *rpcRequestBuild) sendRequest() (string, error) {
	if r.method == "" {
		r.method = http.MethodPost
	}
	req, err := http.NewRequest(r.method, "http://"+r.service.Addr+r.url, strings.NewReader(toJson(r.param)))
	if err != nil {
		return "", err
	}
	req.Header.Set("content-type", "application/json")
	for k, v := range r.header {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp != nil && resp.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(resp.Body)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
