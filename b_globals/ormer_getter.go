package b_globals

import (
	"context"
	"errors"
	"github.com/astaxie/beego/orm"
)

var (
	ormerMapper = make(map[int]orm.Ormer)
)

func GetOrmer(ctx context.Context) (orm.Ormer, error) {
	reqId, exists := GetRequestID(ctx)
	if !exists {
		return nil, errors.New("get ormer error")
	}
	ormer, exists := ormerMapper[reqId]
	if !exists {
		ormer = orm.NewOrm()
	}
	return ormer, nil
}

func RemoveOrmer(ctx context.Context) {
	reqId, exists := GetRequestID(ctx)
	if !exists {
		return
	}
	delete(ormerMapper, reqId)
}
