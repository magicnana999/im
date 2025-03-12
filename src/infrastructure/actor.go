package infrastructure

import (
	"context"
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/asynkron/protoactor-go/cluster/clusterproviders/etcd"
	"github.com/asynkron/protoactor-go/cluster/identitylookup/disthash"
	"github.com/asynkron/protoactor-go/remote"
	"github.com/magicnana999/im/errors"
	"github.com/magicnana999/im/pb"
	imerror "github.com/magicnana999/im/pkg/error"
	"github.com/magicnana999/im/pkg/logger"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type ActorHandler interface {
	IsSupport(ctx context.Context, any any) bool
	Handle(ctx context.Context, any any) (any, error)
}

type ActorConfig struct {
	clientv3.Config
	ClusterName string `yaml:"clusterName"`
	KindName    string `yaml:"kindName"`
	ActorName   string `yaml:"actorName"`
	Port        int    `yaml:"port"`
}

type ActorCluster struct {
	ctx      context.Context
	cluster  *cluster.Cluster
	handlers []ActorHandler
}

func (s *ActorCluster) Receive(ctx actor.Context) {

	m := ctx.Message()
	for _, handler := range s.handlers {
		if handler.IsSupport(s.ctx, m) {
			ret, err := handler.Handle(s.ctx, m)
			if err != nil {
				ctx.Respond(response(ret, err))
			}
		}
	}

	ctx.Respond(response(nil, errors.NoHandlerSupport))
}

func (s *ActorCluster) Start(ctx context.Context) {

	s.cluster.StartMember()

	for {
		select {
		case <-ctx.Done():
			s.cluster.Shutdown(true)
		}
	}
}

func (s *ActorCluster) Shutdown(graceful bool) {
	s.cluster.Shutdown(graceful)
}

func (s *ActorCluster) AddHandler(h ...ActorHandler) {
	s.handlers = append(s.handlers, h...)
}
func response(ret any, err error) *pb.ActorResult {

	if err != nil {
		e := imerror.Format(err)
		return &pb.ActorResult{
			Code:    int32(e.Code),
			Message: e.Message,
		}
	}

	return &pb.ActorResult{
		Code:    int32(0),
		Message: "",
	}

}

func InitCluster(cfg *ActorConfig, handlers ...ActorHandler) *ActorCluster {
	system := actor.NewActorSystem()
	remoteConfig := remote.Configure("0.0.0.0", cfg.Port)

	clusterKind := cluster.NewKind(
		cfg.KindName,
		actor.PropsFromProducer(func() actor.Actor {
			return &ActorCluster{}
		}))

	clusterProvider, err := etcd.NewWithConfig("/actor-cluster", cfg.Config)
	if err != nil {
		logger.Fatalf("init etcd register failed: %v", err)
	}

	lookup := disthash.New()
	clusterConfig := cluster.Configure(cfg.ClusterName, clusterProvider, lookup, remoteConfig, cluster.WithKinds(clusterKind))
	c := cluster.New(system, clusterConfig)

	ac := &ActorCluster{
		cluster: c,
	}

	ac.AddHandler(handlers...)
	return ac
}
