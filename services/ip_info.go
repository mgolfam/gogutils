package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mgolfam/gogutils/httpclient"

	"github.com/mgolfam/gogutils/glog"

	"github.com/mgolfam/gogutils/dto"
)

func GetIpInfo(ip string) *dto.IpInfo {
	if ip == "" || ip == "127.0.0" || ip == "0.0.0.0" {
		return nil
	}

	// Make a request to ip-api.com
	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)

	conf := httpclient.HttpConfig{
		Method:        "GET",
		URL:           url,
		Cache:         true,
		CacheTtl:      24 * 30 * 6 * 3600, // 6 month
		RetrieveCache: true,
		Timeout:       time.Second * 10,
	}

	response, err := httpclient.SendRequest(conf)
	if err != nil {
		glog.LogL(glog.ERROR, "Error getting IP info:", err)
		return nil
	}

	// Parse JSON response
	var ipInfo dto.IpInfo
	err = json.Unmarshal([]byte(response.Body), &ipInfo)
	if err != nil {
		glog.LogL("Error decoding JSON:", err)
		return nil
	}

	if ipInfo.Status == "success" {
		return &ipInfo
	}

	return nil

}
