package remote

import (
	"crypto/tls"
	"go-brunel/internal/pkg/server/notify"
	"go-brunel/internal/pkg/server/store"
	"go-brunel/internal/pkg/shared/remote"
	"net/rpc"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
)

func Server(
	jr store.JobStore,
	lr store.LogStore,
	cr store.ContainerStore,
	rr store.RepositoryStore,
	sr store.StageStore,
	notify notify.Notify,
	credentials remote.Credentials,
	listen string,
) error {
	service := &RPC{
		JobStore:        jr,
		LogStore:        lr,
		ContainerStore:  cr,
		RepositoryStore: rr,
		StageStore:      sr,
		Notify:          notify,
	}
	err := rpc.Register(service)
	if err != nil {
		return errors.Wrap(err, "error registering remote RPC service")
	}

	tlsConfig, err := credentials.ServerConfig()
	if err != nil {
		return err
	}

	log.Info("listening for RPC agent connections on ", listen)
	l, err := tls.Listen("tcp", listen, tlsConfig)
	if err != nil {
		return errors.Wrap(err, "error listening for remote RPC")
	}
	go rpc.Accept(l)
	return nil
}
