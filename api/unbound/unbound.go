package unbound

import (
	"errors"
	"fmt"
	coreapi "github.com/sys-int/opnsense-api/api"
	"net/http"
)

//goland:noinspection GoNameStartsWithPackageName
type UnboundApi struct {
	*coreapi.OPNsense
}

type HostOverride struct {
	// 0 for disabled, 1 for enabled
	Enabled string `json:"enabled"`
	Host    string `json:"hostname"`
	Domain  string `json:"domain"`
	Ip      string `json:"server"`
	//	Rr          RR     `json:"rr,omitempty"`          // A, MX, CNAME...
	Mxprio      string `json:"mxprio,omitempty"`      // 10, 20
	Mx          string `json:"mx,omitempty"`          // mail.domain.tld ...
	Description string `json:"description,omitempty"` // any arbitrary text
	Uuid        string `json:"-"`
}

type RR struct {
	A struct {
		Value    string `json:"value"`
		Selected int    `json:"selected"`
	} `json:"A"`
	AAAA struct {
		Value    string `json:"value"`
		Selected int    `json:"selected"`
	} `json:"AAAA"`
	MX struct {
		Value    string `json:"value"`
		Selected int    `json:"selected"`
	} `json:"MX"`
}

type HostContainer struct {
	HostOverride HostOverride `json:"host,omitempty"`
}

type HostsContainer struct {
	HostOverrides []HostOverride `json:"rows,omitempty"`
}

func (opn *UnboundApi) HostOverrideCreateOrUpdate(hostOverride HostOverride) (string, error) {
	if hostOverride.Uuid == "" { // no uuid given use host / domain based fuzzy search
		var searchResult, _ = opn.HostEntryGetByFQDN(hostOverride.Host, hostOverride.Domain)

		if searchResult.Uuid != "" {
			fmt.Println(fmt.Sprintf("Found entry with same FQDN, doing update with uuid: %s", searchResult.Uuid))
			hostOverride.Uuid = searchResult.Uuid
		}
	}

	if hostOverride.Uuid == "" {
		return opn.HostOverrideCreate(hostOverride)
	} else {
		return opn.HostOverrideUpdate(hostOverride)
	}
}

func (opn *UnboundApi) HostOverrideUpdate(hostOverride HostOverride) (string, error) {
	// endpoint
	var endpoint = opn.EndpointForPluginControllerMethod(coreapi.Unbound, coreapi.Settings, "setHostOverride")
	var fullPath = fmt.Sprintf("%s/%s", endpoint, hostOverride.Uuid)

	var container HostContainer

	container.HostOverride = hostOverride

	response, _ := opn.Client().
		SetError(&coreapi.ServerError{}).
		SetResult(&coreapi.ServerResult{}).
		SetBody(container).
		Post(fullPath)

	if response.StatusCode() == http.StatusOK {
		srvResult := response.Result().(*coreapi.ServerResult)
		return srvResult.Uuid, nil
	} else {
		srvError := response.Error().(*coreapi.ServerError)
		return "", errors.New(srvError.Message)
	}
}

func (opn *UnboundApi) HostOverrideCreate(hostOverride HostOverride) (string, error) {
	// endpoint
	var endpoint = opn.EndpointForPluginControllerMethod(coreapi.Unbound, coreapi.Settings, "addHostOverride")

	var container HostContainer

	container.HostOverride = hostOverride

	response, err := opn.Client().
		SetError(&coreapi.ServerError{}).
		SetResult(&coreapi.ServerResult{}).
		SetBody(container).
		Post(endpoint)

	if response.StatusCode() == http.StatusOK {
		srvResult := response.Result().(*coreapi.ServerResult)
		return srvResult.Uuid, nil
	} else {
		srvError := response.Error().(*coreapi.ServerError)
		return "", errors.New(fmt.Sprintf("%s:%s", srvError.Message, err))
	}
}

func (opn *UnboundApi) HostEntryGetByFQDN(host string, domain string) (HostOverride, error) {
	var endpoint = opn.EndpointForPluginControllerMethod(coreapi.Unbound, coreapi.Settings, "searchHostOverride")
	var reqUrl = fmt.Sprintf("%s?searchPhrase=%s", endpoint, host)
	response, _ := opn.Client().
		SetError(&coreapi.ServerError{}).
		SetResult(&HostsContainer{}).
		Get(reqUrl)

	if response.StatusCode() == 200 {
		container := response.Result().(*HostsContainer)

		var allWithMatchingDomain = Filter(container.HostOverrides, func(override HostOverride) bool {
			return override.Domain == domain
		})

		if len(allWithMatchingDomain) > 1 {
			return HostOverride{}, &coreapi.TooManyFoundError{
				Err:  nil,
				Name: "found more then one entry",
			}
		}

		if len(allWithMatchingDomain) == 0 {
			return HostOverride{}, &coreapi.NotFoundError{
				Err:  nil,
				Name: "hostentry",
			}
		}

		return allWithMatchingDomain[0], nil
	} else if response.StatusCode() == 404 {
		return HostOverride{}, &coreapi.NotFoundError{
			Err:  nil,
			Name: "hostentry",
		}
	} else {
		srvError := response.Error().(*coreapi.ServerError)
		return HostOverride{}, errors.New(fmt.Sprintf("%s", srvError.Message))
	}
}

func (opn *UnboundApi) HostEntryGetByUuid(uuid string) (HostOverride, error) {
	var endpoint = opn.EndpointForPluginControllerMethod(coreapi.Unbound, coreapi.Settings, "getHostOverride")
	var fullPath = fmt.Sprintf("%s/%s", endpoint, uuid)

	response, err := opn.Client().
		SetResult(&HostContainer{}).
		SetError(&coreapi.ServerError{}).
		Get(fullPath)

	if response.StatusCode() == http.StatusOK && err == nil {
		srvResult := response.Result().(*HostContainer)
		return srvResult.HostOverride, err
	} else {
		return HostOverride{}, &coreapi.NotFoundError{
			Err:  nil,
			Name: "hostentry",
		}
	}
}

func (opn *UnboundApi) HostOverrideList() ([]HostOverride, error) {
	// endpoint
	var endpoint = opn.EndpointForPluginControllerMethod(coreapi.Unbound, coreapi.Settings, "searchHostOverride")

	response, err := opn.Client().
		SetResult(&HostsContainer{}).
		SetError(&coreapi.ServerError{}).
		Get(endpoint)

	if response.StatusCode() == 200 {
		container := response.Result().(*HostsContainer)
		return container.HostOverrides, nil
	} else {
		srvError := response.Error().(*coreapi.ServerError)
		return []HostOverride{}, errors.New(fmt.Sprintf("%s:%s", srvError.Message, err))
	}
}

func (opn *UnboundApi) HostEntryRemove(uuid string) error {
	var endpoint = opn.EndpointForPluginControllerMethod(coreapi.Unbound, coreapi.Settings, "delHostOverride")
	var fullPath = fmt.Sprintf("%s/%s", endpoint, uuid)

	response, err := opn.Client().
		SetError(&coreapi.ServerError{}).
		SetResult(&coreapi.ServerResult{}).
		Post(fullPath)

	if response.StatusCode() == 200 {
		return nil
	} else {
		srvError := response.Error().(*coreapi.ServerError)
		return errors.New(fmt.Sprintf("%s:%s", srvError.Message, err))
	}
}

func (opn *UnboundApi) HostEntryExists(host string, domain string) (bool, error) {
	if _, err := opn.HostEntryGetByFQDN(host, domain); err != nil {
		switch err.(type) {
		case *coreapi.NotFoundError:
			return false, nil
		default:
			return true, err
		}
	}
	return true, nil
}

func Filter[T any](vs []T, f func(T) bool) []T {
	filtered := make([]T, 0)
	for _, v := range vs {
		if f(v) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}
