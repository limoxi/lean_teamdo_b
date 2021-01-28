package account

import (
	"context"
	"github.com/bitly/go-simplejson"
	"github.com/kfchen81/beego/vanilla/middleware"
	"net/http"
)

var gInstance *ContextFactory

type ContextFactory struct {
}

// NewContext 构造含有corp的Context
func (this *ContextFactory) NewContext(ctx context.Context, request *http.Request, userId int, jwtToken string, rawData *simplejson.Json) context.Context {
	user := new(User)
	user.Model = nil
	user.Id = 1
	user.RawData = rawData
	user.Ctx = ctx
	ctx = context.WithValue(ctx, "account", user)
	return ctx
}

func init() {
	gInstance = &ContextFactory{}
	middleware.SetBusinessContextFactory(gInstance)
}
