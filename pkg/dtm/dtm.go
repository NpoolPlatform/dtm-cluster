package dtm

import (
	"context"
	"fmt"
	"net"

	constant "github.com/NpoolPlatform/dtm-cluster/pkg/const"

	"github.com/NpoolPlatform/go-service-framework/pkg/config"
	grpc2 "github.com/NpoolPlatform/go-service-framework/pkg/grpc"

	apimgrcli "github.com/NpoolPlatform/api-manager/pkg/client"

	"github.com/dtm-labs/dtmgrpc"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type uri struct {
	action string
	revert string
}

type Action struct {
	ServiceName string
	Action      string
	Revert      string
	Param       protoreflect.ProtoMessage
	uri         uri
}

func (act *Action) ConstructURI(ctx context.Context) error {
	api, err := apimgrcli.GetServiceMethodAPI(ctx, act.ServiceName, act.Action)
	if err != nil || api == nil {
		return fmt.Errorf("fail get service method api: %v", err)
	}
	act.uri.action = act.ServiceName + "/" + api.Path

	api, err = apimgrcli.GetServiceMethodAPI(ctx, act.ServiceName, act.Revert)
	if err != nil || api == nil {
		return fmt.Errorf("fail get service method api: %v", err)
	}
	act.uri.revert = act.ServiceName + "/" + api.Path

	return nil
}

func WithSage(ctx context.Context, actions []*Action, pre, post func(ctx context.Context) error) error {
	if pre != nil {
		if err := pre(ctx); err != nil {
			return fmt.Errorf("fail run pre: %v", err)
		}
	}

	svc, err := config.PeekService(constant.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return fmt.Errorf("fail peek dtm service: %v", err)
	}

	host := net.JoinHostPort(svc.Address, fmt.Sprintf("%v", svc.Port))
	gid := dtmgrpc.MustGenGid(host)

	saga := dtmgrpc.NewSagaGrpc(host, gid)
	for _, act := range actions {
		if err := act.ConstructURI(ctx); err != nil {
			return fmt.Errorf("fail construct action uri: %v", err)
		}
		saga = saga.Add(act.uri.action, act.uri.revert, act.Param)
	}
	saga.WaitResult = true

	err = saga.Submit()
	if err != nil {
		return fmt.Errorf("fail run saga: %v", err)
	}

	if post != nil {
		return post(ctx)
	}
	return nil
}
