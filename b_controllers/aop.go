package b_controllers

import (
	"context"
	"do-global.com/beego-basic/b_globals"
	"do-global.com/beego-basic/b_logger"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"math/rand"
	"net/http"
	"reflect"
	"strings"
	"time"
	"unsafe"
)

func (t *BasicController) BPrepare() {
	t.startTime = time.Now().UnixNano()

	t.Ctx.Input.SetData("_____t", t)
	// 设置request生命周期的参数
	ctx := context.TODO()
	ctx = b_globals.WithRequestID(ctx, rand.Intn(time.Now().Second()+1))
	t.ReqCtx = ctx
}

func (t *BasicController) BAccessLog() {
	// access log
	spentTime := (time.Now().UnixNano() - t.startTime) / 1e6
	paramsStr, _ := json.Marshal(t.Ctx.Request.Form)
	reqPath := t.Ctx.Request.URL.Path
	now := time.Now()
	status := t.Ctx.Input.GetData("___status")
	if status == nil || status.(int) == 0 {
		status = t.Ctx.ResponseWriter.Status
	}
	if status == 0 {
		status = http.StatusOK
	}
	accessInfo := fmt.Sprintf("%v\001%v\001%v\001%v\001%v\001%v", now.Format("20060102"), now.UnixNano()/1e6, reqPath, string(paramsStr),
		status, spentTime)
	b_logger.AccessLogger.Info(accessInfo)
}

func (t *BasicController) BClearCtx()  {
	// 删掉reqId相关资源
	b_globals.RemoveOrmer(t.ReqCtx)
	t.ReqCtx.Done()
}

func (t *BasicController) BCheckAccess(accessWhiteList []string, tokenHeaderName string, checkTokenFunc func(token string) bool) {
	needCheck := beego.AppConfig.DefaultBool("secure.control_check", true)
	if !needCheck {
		return
	}
	controllerType, methods, isFind := t.GetRequestControllerAndMethods()
	if !isFind {
		t.Error403()
		return
	}

	controllerName := controllerType.String()
	var methodName string
	for _, v := range methods {
		methodName = v
		break
	}
	controllerAndMethod := controllerName + "." + methodName
	for _, whiteItem := range accessWhiteList {
		if strings.ToLower(controllerAndMethod) == strings.ToLower(whiteItem) {
			return
		}
	}

	token := t.Ctx.Request.Header.Get(tokenHeaderName)
	if token == "" {
		t.Error403()
	}
	if !checkTokenFunc(token) {
		t.Error403()
	}
}

func (t *BasicController) GetRequestControllerAndMethods() (reflect.Type, map[string]string, bool) {
	cInfo, isFind := beego.BeeApp.Handlers.FindRouter(t.Ctx)
	if isFind {
		controllerInfoV := reflect.ValueOf(cInfo).Elem()

		controllerTypeV := controllerInfoV.Field(1)
		controllerTypeV = reflect.NewAt(controllerTypeV.Type(), unsafe.Pointer(controllerTypeV.UnsafeAddr()))
		controllerType := controllerTypeV.Interface().(*reflect.Type)

		methodsV := controllerInfoV.Field(2)
		methodsV = reflect.NewAt(methodsV.Type(), unsafe.Pointer(methodsV.UnsafeAddr()))
		methods := methodsV.Interface().(*map[string]string)

		return *controllerType, *methods, true
	}
	return nil, nil, false
}
