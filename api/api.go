package api

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/url"
)

type OPNsense struct {
	BaseUrl     url.URL
	ApiKey      string
	ApiSecret   string
	NoSslVerify bool
}

func (opn *OPNsense) Send(request *http.Request) (*http.Response, error) {
	var client = &http.Client{}

	certPool, _ := x509.SystemCertPool()
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: opn.NoSslVerify,
			RootCAs:            certPool,
		},
	}

	request.SetBasicAuth(opn.ApiKey, opn.ApiSecret)
	return client.Do(request)
}

type NotFoundError struct {
	Name string
	Err  error
}

type TooManyFoundError struct {
	Name string
	Err  error
}

func (f *NotFoundError) Error() string {
	return fmt.Sprintf("not found: %s", f.Name)
}

func (f *TooManyFoundError) Error() string {
	return fmt.Sprintf("too many found: %s", f.Name)
}

// EndpointForModule so basically api/<plugin>
func (opn *OPNsense) EndpointForModule(module string) string {
	return fmt.Sprintf("%s/api/%s", opn.BaseUrl.String(), module)
}

// EndpointForModuleController so basically api/<plugin>/<controller>
func (opn *OPNsense) EndpointForModuleController(module string, controller string) string {
	return fmt.Sprintf("%s/%s", opn.EndpointForModule(module), controller)
}

// EndpointForPluginControllerMethod so basically api/<plugin>/<controller>/<method>
func (opn *OPNsense) EndpointForPluginControllerMethod(module string, controller string, method string) string {
	return fmt.Sprintf("%s/%s", opn.EndpointForModuleController(module, controller), method)
}
