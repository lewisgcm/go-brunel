package remote

import (
	"crypto/tls"
	"github.com/pkg/errors"
	"go-brunel/internal/pkg/shared"
	"go-brunel/internal/pkg/shared/remote"
	"net/rpc"
)

type rpcClient struct {
	client *rpc.Client
}

func rpcError(err error) error {
	return errors.Wrap(err, "RPC error")
}

func NewRPCClient(credentials remote.Credentials, endpoint string) (Remote, error) {
	tlsConfig, err := remote.ClientConfig(credentials)
	if err != nil {
		return nil, errors.Wrap(err, "error generating TLS configuration for connecting to RPC endpoint")
	}

	conn, err := tls.Dial("tcp", endpoint, tlsConfig)
	if err != nil {
		return nil, errors.Wrap(err, "error dialing RPC endpoint")
	}
	return &rpcClient{
		client: rpc.NewClient(conn),
	}, nil
}

func (c *rpcClient) GetNextAvailableJob() (*shared.Job, error) {
	var reply remote.GetNextAvailableJobResponse
	err := c.client.Call("RPC.GetNextAvailableJob", &remote.Empty{}, &reply)
	if err != nil {
		return nil, rpcError(err)
	}
	return reply.Job, nil
}

func (c *rpcClient) SetJobState(id shared.JobID, state shared.JobState) error {
	return rpcError(
		c.client.Call("RPC.SetJobState", &remote.SetJobStateRequest{Id: id, State: state}, &remote.Empty{}),
	)
}

func (c *rpcClient) HasBeenCancelled(id shared.JobID) (bool, error) {
	var reply bool
	e := c.client.Call("RPC.HasBeenCancelled", &id, &reply)
	return reply, rpcError(e)
}

func (c *rpcClient) Log(id shared.JobID, message string, logType shared.LogType) error {
	return rpcError(
		c.client.Call("RPC.Log", &remote.LogRequest{Id: id, Message: message, LogType: logType}, &remote.Empty{}),
	)
}

func (c *rpcClient) AddContainer(id shared.JobID, containerID shared.ContainerID, meta shared.ContainerMeta, container shared.Container, state shared.ContainerState) error {
	return rpcError(
		c.client.Call("RPC.AddContainer", &remote.AddContainerRequest{
			Id:          id,
			ContainerID: containerID,
			Meta:        meta,
			Container:   container,
			State:       state,
		}, &remote.Empty{}),
	)
}

func (c *rpcClient) SetContainerState(id shared.ContainerID, state shared.ContainerState) error {
	return rpcError(
		c.client.Call("RPC.SetContainerState", &remote.SetContainerStateRequest{Id: id, State: state}, &remote.Empty{}),
	)
}

func (c *rpcClient) ContainerLog(id shared.ContainerID, message string, logType shared.LogType) error {
	return rpcError(
		c.client.Call("RPC.ContainerLog", &remote.ContainerLogRequest{Id: id, Message: message, LogType: logType}, &remote.Empty{}),
	)
}

func (c *rpcClient) SearchForValue(searchPath []string, name string) (string, error) {
	var reply string
	e := c.client.Call("RPC.SearchForValue", remote.SearchForXRequest{SearchPath: searchPath, Name: name}, &reply)
	return reply, rpcError(e)
}

func (c *rpcClient) SearchForSecret(searchPath []string, name string) (string, error) {
	var reply string
	e := c.client.Call("RPC.SearchForSecret", remote.SearchForXRequest{SearchPath: searchPath, Name: name}, &reply)
	return reply, rpcError(e)
}
