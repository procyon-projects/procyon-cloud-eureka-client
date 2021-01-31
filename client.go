package eureka

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type HttpClient interface {
	Register(info *InstanceInfo) error
	Deregister(appName, instanceID string) error
	SendHeartBeat(appName, instanceID string, info *InstanceInfo, overriddenStatus InstanceStatus) error
	UpdateStatus(appName, instanceID string, newStatus InstanceStatus, info *InstanceInfo) error
	GetInstances(appName string) (*ApplicationResponse, error)
	GetInstanceByAppAndInstanceId(appName, instanceID string) (*InstanceInfoResponse, error)
	GetInstanceByServiceId(instanceID string) (*InstanceInfoResponse, error)
	GetApplications(regions ...string) (*ApplicationsResponse, error)
}

type DefaultHttpClient struct {
	client     *http.Client
	serviceUrl string
}

func newDefaultHttpClient(serviceUrl string) DefaultHttpClient {
	return DefaultHttpClient{
		client:     http.DefaultClient,
		serviceUrl: serviceUrl,
	}
}

func (httpClient DefaultHttpClient) Register(info *InstanceInfo) error {
	_, err := httpClient.makeRequest(http.MethodPost,
		httpClient.serviceUrl+"apps/"+info.AppName,
		nil,
		nil)
	if err != nil {
		return err
	}
	return nil
}

func (httpClient DefaultHttpClient) Deregister(appName, instanceID string) error {
	_, err := httpClient.makeRequest(http.MethodDelete,
		httpClient.serviceUrl+"apps/"+appName+"/"+instanceID,
		nil,
		nil)
	if err != nil {
		return err
	}
	return nil
}

func (httpClient DefaultHttpClient) SendHeartBeat(appName, instanceID string, info *InstanceInfo, overriddenStatus InstanceStatus) error {
	heartBeatUrl, err := url.Parse(httpClient.serviceUrl + "apps/" + appName + "/" + instanceID)
	if err != nil {
		return err
	}
	heartBeatUrl.Query().Add("status", string(info.Status))
	heartBeatUrl.Query().Add("lastDirtyTimestamp", info.LastDirtyTimestamp)
	heartBeatUrl.Query().Add("overriddenstatus", string(overriddenStatus))

	_, err = httpClient.makeRequest(http.MethodPut,
		heartBeatUrl.String(),
		nil,
		nil)
	if err != nil {
		return err
	}
	return nil
}

func (httpClient DefaultHttpClient) UpdateStatus(appName, instanceID string, newStatus InstanceStatus, info *InstanceInfo) error {
	updateStatusUrl, err := url.Parse(httpClient.serviceUrl + "apps/" + appName + "/" + instanceID + "/status")
	if err != nil {
		return err
	}
	updateStatusUrl.Query().Add("value", string(info.Status))
	updateStatusUrl.Query().Add("lastDirtyTimestamp", info.LastDirtyTimestamp)

	_, err = httpClient.makeRequest(http.MethodPut,
		updateStatusUrl.String(),
		nil,
		nil)
	if err != nil {
		return err
	}
	return nil
}

func (httpClient DefaultHttpClient) GetInstances(appName string) (*ApplicationResponse, error) {
	resp, err := httpClient.makeRequest(http.MethodGet,
		httpClient.serviceUrl+"apps/"+appName,
		nil,
		map[string]string{
			"Accept": "application/json",
		})
	if err != nil {
		return nil, err
	}

	if strings.Trim(resp.Status, " ") != strconv.Itoa(http.StatusOK) {
		return nil, nil
	}

	applicationResponse := &ApplicationResponse{}
	err = httpClient.bindResponse(resp, applicationResponse)
	if err != nil {
		return nil, err
	}
	return applicationResponse, nil
}

func (httpClient DefaultHttpClient) GetInstanceByAppAndInstanceId(appName, instanceID string) (*InstanceInfoResponse, error) {
	resp, err := httpClient.makeRequest(http.MethodGet,
		httpClient.serviceUrl+"apps/"+appName+"/"+instanceID,
		nil,
		map[string]string{
			"Accept": "application/json",
		})
	if err != nil {
		return nil, err
	}

	if strings.Trim(resp.Status, " ") != strconv.Itoa(http.StatusOK) {
		return nil, nil
	}

	instanceInfoResponse := &InstanceInfoResponse{}
	err = httpClient.bindResponse(resp, instanceInfoResponse)
	if err != nil {
		return nil, err
	}
	return instanceInfoResponse, nil
}

func (httpClient DefaultHttpClient) GetInstanceByServiceId(instanceID string) (*InstanceInfoResponse, error) {
	resp, err := httpClient.makeRequest(http.MethodGet,
		httpClient.serviceUrl+"instances/"+instanceID,
		nil,
		map[string]string{
			"Accept": "application/json",
		})
	if err != nil {
		return nil, err
	}

	if strings.Trim(resp.Status, " ") != strconv.Itoa(http.StatusOK) {
		return nil, nil
	}

	instanceInfoResponse := &InstanceInfoResponse{}
	err = httpClient.bindResponse(resp, instanceInfoResponse)
	if err != nil {
		return nil, err
	}
	return instanceInfoResponse, nil
}

func (httpClient DefaultHttpClient) GetApplications(regions ...string) (*ApplicationsResponse, error) {
	resp, err := httpClient.makeRequest(http.MethodGet,
		httpClient.serviceUrl+"apps/",
		nil,
		map[string]string{
			"Accept": "application/json",
		})
	if err != nil {
		return nil, err
	}

	if strings.Trim(resp.Status, " ") != strconv.Itoa(http.StatusOK) {
		return nil, nil
	}

	applications := &ApplicationsResponse{}
	err = httpClient.bindResponse(resp, applications)
	if err != nil {
		return nil, err
	}
	return applications, nil
}

func (httpClient DefaultHttpClient) makeRequest(method string, url string, body io.Reader, header map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if header != nil {
		for headerKey, headerValue := range header {
			req.Header.Set(headerKey, headerValue)
		}
	}
	return httpClient.client.Do(req)
}

func (httpClient DefaultHttpClient) bindResponse(resp *http.Response, responseObject interface{}) error {
	responseArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseArray, responseObject)
	if err != nil {
		return err
	}
	return nil
}
