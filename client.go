package eureka

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type HttpClient interface {
	Register(info *InstanceInfo) error
	Deregister(appName, instanceId string) error
	SendHeartBeat(appName, instanceId string, info *InstanceInfo, overriddenStatus InstanceStatus) error
	UpdateStatus(appName, instanceId string, newStatus InstanceStatus, info *InstanceInfo) error
	GetApplication(appName string) (*Application, error)
	GetInstanceByAppNameAndInstanceId(appName, instanceId string) (*InstanceInfo, error)
	GetInstanceByInstanceId(instanceId string) (*InstanceInfo, error)
	GetApplications(regions ...string) (*Applications, error)
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
	instanceResource := &InstanceResource{
		InstanceInfo: info,
	}

	resp, err := httpClient.makeRequest(http.MethodPost,
		httpClient.serviceUrl+"apps/"+info.AppName,
		instanceResource,
		map[string]string{
			"Accept-Encoding": "gzip",
			"Content-Type":    "application/json",
		},
	)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return httpClient.getError(resp)
	}

	return nil
}

func (httpClient DefaultHttpClient) Deregister(appName, instanceId string) error {
	resp, err := httpClient.makeRequest(http.MethodDelete,
		httpClient.serviceUrl+"apps/"+appName+"/"+instanceId,
		nil,
		nil)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return httpClient.getError(resp)
	}

	return nil
}

func (httpClient DefaultHttpClient) SendHeartBeat(appName, instanceId string, info *InstanceInfo, overriddenStatus InstanceStatus) error {
	heartBeatUrl, err := url.Parse(httpClient.serviceUrl + "apps/" + appName + "/" + instanceId)

	if err != nil {
		return err
	}

	query := heartBeatUrl.Query()
	query.Add("status", string(info.Status))
	query.Add("lastDirtyTimestamp", info.LastDirtyTimestamp)
	query.Add("overriddenstatus", string(overriddenStatus))

	heartBeatUrl.RawQuery = query.Encode()

	var resp *http.Response
	resp, err = httpClient.makeRequest(http.MethodPut,
		heartBeatUrl.String(),
		nil,
		nil)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return httpClient.getError(resp)
	}

	return nil
}

func (httpClient DefaultHttpClient) UpdateStatus(appName, instanceId string, newStatus InstanceStatus, info *InstanceInfo) error {
	updateStatusUrl, err := url.Parse(httpClient.serviceUrl + "apps/" + appName + "/" + instanceId + "/status")

	if err != nil {
		return err
	}

	query := updateStatusUrl.Query()

	query.Add("value", string(newStatus))
	query.Add("lastDirtyTimestamp", info.LastDirtyTimestamp)

	updateStatusUrl.RawQuery = query.Encode()

	var resp *http.Response
	resp, err = httpClient.makeRequest(http.MethodPut,
		updateStatusUrl.String(),
		nil,
		nil)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return httpClient.getError(resp)
	}

	return nil
}

func (httpClient DefaultHttpClient) GetApplication(appName string) (*Application, error) {
	resp, err := httpClient.makeRequest(http.MethodGet,
		httpClient.serviceUrl+"apps/"+appName,
		nil,
		map[string]string{
			"Accept": "application/json",
		},
	)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, httpClient.getError(resp)
	}

	applicationResource := &ApplicationResource{}
	err = httpClient.bindResponse(resp, applicationResource)

	if err != nil {
		return nil, err
	}

	return applicationResource.Application, nil
}

func (httpClient DefaultHttpClient) GetInstanceByAppNameAndInstanceId(appName, instanceId string) (*InstanceInfo, error) {
	resp, err := httpClient.makeRequest(http.MethodGet,
		httpClient.serviceUrl+"apps/"+appName+"/"+instanceId,
		nil,
		map[string]string{
			"Accept": "application/json",
		},
	)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, httpClient.getError(resp)
	}

	instanceResource := &InstanceResource{}
	err = httpClient.bindResponse(resp, instanceResource)

	if err != nil {
		return nil, err
	}

	return instanceResource.InstanceInfo, nil
}

func (httpClient DefaultHttpClient) GetInstanceByInstanceId(instanceId string) (*InstanceInfo, error) {
	resp, err := httpClient.makeRequest(http.MethodGet,
		httpClient.serviceUrl+"instances/"+instanceId,
		nil,
		map[string]string{
			"Accept": "application/json",
		},
	)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, httpClient.getError(resp)
	}

	instanceResource := &InstanceResource{}
	err = httpClient.bindResponse(resp, instanceResource)
	if err != nil {
		return nil, err
	}
	return instanceResource.InstanceInfo, nil
}

func (httpClient DefaultHttpClient) GetApplications(regions ...string) (*Applications, error) {
	applicationsUrl, err := url.Parse(httpClient.serviceUrl + "apps")

	if err != nil {
		return nil, err
	}

	query := applicationsUrl.Query()

	regionParameter := ""

	for _, region := range regions {
		regionParameter = regionParameter + "," + region
	}

	if regionParameter != "" {
		query.Add("regions", regionParameter[1:])
		applicationsUrl.RawQuery = query.Encode()
	}

	var resp *http.Response
	resp, err = httpClient.makeRequest(http.MethodGet,
		applicationsUrl.String(),
		nil,
		map[string]string{
			"Accept": "application/json",
		},
	)

	if err != nil {
		return nil, err
	}

	if strings.Trim(resp.Status, " ") != strconv.Itoa(http.StatusOK) {
		return nil, nil
	}

	applicationsResource := &ApplicationsResource{}
	err = httpClient.bindResponse(resp, applicationsResource)
	if err != nil {
		return nil, err
	}
	return applicationsResource.Applications, nil
}

func (httpClient DefaultHttpClient) getError(resp *http.Response) error {
	errorResponse := &Error{}

	err := httpClient.bindResponse(resp, errorResponse)
	if err != nil {
		return err
	}

	return errors.New(errorResponse.Message)
}

func (httpClient DefaultHttpClient) makeRequest(method string, url string, requestBodyObj interface{}, header map[string]string) (*http.Response, error) {
	body, err := json.Marshal(requestBodyObj)

	if err != nil {
		return nil, err
	}

	var req *http.Request
	req, err = http.NewRequest(method, url, bytes.NewReader(body))
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
