package unbound

import (
	"errors"
	"fmt"
	coreapi "github.com/sys-int/opnsense-api/api"
)

func (opn *UnboundApi) ServiceRestart(uuid string) error {
	var endpoint = opn.EndpointForPluginControllerMethod(coreapi.Unbound, coreapi.Service, "restart")
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
