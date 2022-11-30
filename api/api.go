package api

import (
	"crypto/tls"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"net/url"
)

type OPNsense struct {
	BaseUrl     url.URL
	ApiKey      string
	ApiSecret   string
	NoSslVerify bool
}

func (opn *OPNsense) Client() *resty.Request {

	if opn.NoSslVerify {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	client := resty.New()
	client.SetBasicAuth(opn.ApiKey, opn.ApiSecret)

	return client.R()
}

type NotFoundError struct {
	Name string
	Err  error
}

type TooManyFoundError struct {
	Name string
	Err  error
}

type ServerError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ServerResult struct {
	ResultStatus string `json:"result"`
	Uuid         string `json:"uuid,omitempty"`
}

func (f *NotFoundError) Error() string {
	return fmt.Sprintf("not found: %s", f.Name)
}

func (f *TooManyFoundError) Error() string {
	return fmt.Sprintf("too many found: %s", f.Name)
}

// EndpointForModule so basically api/<plugin>
func (opn *OPNsense) EndpointForModule(module Module) string {
	return fmt.Sprintf("%s/api/%s", opn.BaseUrl.String(), module.String())
}

// EndpointForModuleController so basically api/<plugin>/<controller>
func (opn *OPNsense) EndpointForModuleController(module Module, controller Controller) string {
	return fmt.Sprintf("%s/%s", opn.EndpointForModule(module), controller.String())
}

// EndpointForPluginControllerMethod so basically api/<plugin>/<controller>/<method>
func (opn *OPNsense) EndpointForPluginControllerMethod(module Module, controller Controller, method string) string {
	return fmt.Sprintf("%s/%s", opn.EndpointForModuleController(module, controller), method)
}
