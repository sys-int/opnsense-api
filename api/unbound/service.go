package unbound

import (
	"errors"
	"fmt"
	coreapi "github.com/sys-int/opnsense-api/api"
)

func (opn *UnboundApi) ServiceRestart() error {
	var endpoint = opn.EndpointForPluginControllerMethod(coreapi.Unbound, coreapi.Service, "restart")

	response, err := opn.Client().
		SetError(&coreapi.ServerError{}).
		SetResult(&coreapi.ServerResult{}).
		Post(endpoint)

	if response.StatusCode() == 200 {
		return nil
	} else {
		srvError := response.Error().(*coreapi.ServerError)
		return errors.New(fmt.Sprintf("%s:%s", srvError.Message, err))
	}
}
