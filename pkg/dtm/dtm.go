package dtm

import (
	"context"
	"fmt"
	"net"
	"sync"

	constant "github.com/NpoolPlatform/dtm-cluster/pkg/const"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	grpc2 "github.com/NpoolPlatform/go-service-framework/pkg/grpc"

	apimwcli "github.com/NpoolPlatform/basal-middleware/pkg/client/api"
	apimgrpb "github.com/NpoolPlatform/message/npool/basal/mw/v1/api"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonpb "github.com/NpoolPlatform/message/npool/basetypes/v1"

	"github.com/dtm-labs/dtm/client/dtmcli/dtmimp"
	"github.com/dtm-labs/dtm/client/dtmgrpc"
	"google.golang.org/protobuf/proto"
)

type uri struct {
	action string
	revert string
}

type Action struct {
	ServiceName string
	Action      string
	Revert      string
	Args        proto.Message
	uri         uri
}

type SagaDispose struct {
	TransOptions dtmimp.TransOptions
	Actions      []*Action
}

func NewSagaDispose(options dtmimp.TransOptions) *SagaDispose {
	return &SagaDispose{
		TransOptions: options,
	}
}

func (sd *SagaDispose) Add(svcName, action, revert string, args proto.Message) {
	sd.Actions = append(sd.Actions, &Action{
		ServiceName: svcName,
		Action:      action,
		Revert:      revert,
		Args:        args,
	})
}

var apiMap sync.Map
var hostMap sync.Map

func (act *Action) apiKey(action string) string {
	return fmt.Sprintf("%v:%v", act.ServiceName, action)
}

func (act *Action) constructURI(ctx context.Context) (err error) {
	api, ok := apiMap.Load(act.apiKey(act.Action))
	if !ok {
		_api, err := apimwcli.GetAPIOnly(ctx, &apimgrpb.Conds{
			ServiceName: &commonpb.StringVal{Op: cruder.EQ, Value: act.ServiceName},
			Path:        &commonpb.StringVal{Op: cruder.EQ, Value: act.Action},
		})
		if err != nil || _api == nil {
			return fmt.Errorf("service %v api %v: %v", act.ServiceName, act.Action, err)
		}
		if _api.Path == "" {
			return fmt.Errorf("invalid api path: %v, %v", act.ServiceName, act.Action)
		}
		api = _api
		apiMap.Store(act.apiKey(act.Action), api)
	}

	host, ok := hostMap.Load(act.ServiceName)
	if !ok {
		svc, err := config.PeekService(act.ServiceName, grpc2.GRPCTAG)
		if err != nil {
			return fmt.Errorf("service %v: %v", act.ServiceName, err)
		}
		host = net.JoinHostPort(svc.Address, fmt.Sprintf("%v", svc.Port))
		hostMap.Store(act.ServiceName, host)
	}

	act.uri.action = host.(string) + "/" + api.(*apimgrpb.API).Path
	if act.Revert == "" {
		return nil
	}

	api, ok = apiMap.Load(act.apiKey(act.Revert))
	if !ok {
		_api, err := apimwcli.GetAPIOnly(ctx, &apimgrpb.Conds{
			ServiceName: &commonpb.StringVal{Op: cruder.EQ, Value: act.ServiceName},
			Path:        &commonpb.StringVal{Op: cruder.EQ, Value: act.Revert},
		})
		if err != nil || _api == nil {
			return fmt.Errorf("service %v api %v: %v", act.ServiceName, act.Revert, err)
		}
		if _api.Path == "" {
			return fmt.Errorf("invalid api path: %v, %v", act.ServiceName, act.Revert)
		}
		api = _api
		apiMap.Store(act.apiKey(act.Revert), api)
	}
	act.uri.revert = host.(string) + "/" + api.(*apimgrpb.API).Path

	return nil
}

func WithSaga(ctx context.Context, dispose *SagaDispose) error {
	host, ok := hostMap.Load(constant.ServiceName)
	if !ok {
		svc, err := config.PeekService(constant.ServiceName, grpc2.GRPCTAG)
		if err != nil {
			return fmt.Errorf("service %v: %v", constant.ServiceName, err)
		}
		host = net.JoinHostPort(svc.Address, fmt.Sprintf("%v", svc.Port))
		hostMap.Store(constant.ServiceName, host)
	}

	gid := dtmgrpc.MustGenGid(host.(string))
	saga := dtmgrpc.NewSagaGrpc(host.(string), gid)
	for _, act := range dispose.Actions {
		if err := act.constructURI(ctx); err != nil {
			return fmt.Errorf("construct uri: %v", err)
		}
		saga = saga.Add(act.uri.action, act.uri.revert, act.Args)
	}
	saga.TransOptions = dispose.TransOptions

	err := saga.Submit()
	if err != nil {
		return fmt.Errorf("saga: %v", err)
	}

	return nil
}
