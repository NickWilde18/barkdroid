package push

import (
	"context"

	v1 "barkdroid/api/push/v1"
)

// IPushV1 defines the push notification API interface.
type IPushV1 interface {
	PushBody(ctx context.Context, req *v1.PushBodyReq) (res *v1.PushRes, err error)
	PushTitleBody(ctx context.Context, req *v1.PushTitleBodyReq) (res *v1.PushRes, err error)
	PushPost(ctx context.Context, req *v1.PushPostReq) (res *v1.PushRes, err error)
	RegisterDevice(ctx context.Context, req *v1.RegisterDeviceReq) (res *v1.RegisterDeviceRes, err error)
}
