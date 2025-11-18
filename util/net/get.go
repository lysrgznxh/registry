package net

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

/*
*
测试连通性
*/
func TestConnectivity(url string, timeout time.Duration) bool {
	client := http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   true, //true:不同HTTP请求之间TCP连接的重用将被阻止（http1.1默认为长连接，此处改为短连接）
			MaxIdleConnsPerHost: 512,  //控制每个主机下的最大闲置连接数目
		},
		Timeout: timeout, //Client请求的时间限制,该超时限制包括连接时间、重定向和读取response body时间;Timeout为零值表示不设置超时
	}
	req, err := http.NewRequest("GET", url, bytes.NewReader([]byte("")))
	if err != nil {
		return false
	}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	//log.Debug("探测cosmos",resp.Body)
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	return true
}

// 发起http请求
func HttpGet(url string, timeout time.Duration) (content string, err error) {
	client := http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   true, //true:不同HTTP请求之间TCP连接的重用将被阻止（http1.1默认为长连接，此处改为短连接）
			MaxIdleConnsPerHost: 512,  //控制每个主机下的最大闲置连接数目
		},
		Timeout: timeout, //Client请求的时间限制,该超时限制包括连接时间、重定向和读取response body时间;Timeout为零值表示不设置超时
	}
	req, err := http.NewRequest("GET", url, bytes.NewReader([]byte("")))
	if err != nil {
		return content, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return content, err
	}
	//log.Debug("探测cosmos",resp.Body)
	defer resp.Body.Close()
	contentByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return content, err
	}
	return string(contentByte), nil
}

// 发起http请求
func HttpGetByte(url string, timeout time.Duration) (content []byte, err error) {
	client := http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   true, //true:不同HTTP请求之间TCP连接的重用将被阻止（http1.1默认为长连接，此处改为短连接）
			MaxIdleConnsPerHost: 512,  //控制每个主机下的最大闲置连接数目
		},
		Timeout: timeout, //Client请求的时间限制,该超时限制包括连接时间、重定向和读取response body时间;Timeout为零值表示不设置超时
	}
	req, err := http.NewRequest("GET", url, bytes.NewReader([]byte("")))
	if err != nil {
		return content, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return content, err
	}
	//log.Debug("探测cosmos",resp.Body)
	defer resp.Body.Close()
	contentByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return content, err
	}
	return contentByte, nil
}
