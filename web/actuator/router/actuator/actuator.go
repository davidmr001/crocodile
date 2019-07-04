package actuator

import (
	"context"
	"crocodile/common/bind"
	"crocodile/common/cfg"
	"crocodile/common/e"
	"crocodile/common/registry"
	"crocodile/common/response"
	pbactuator "crocodile/service/actuator/proto/actuator"
	"github.com/gin-gonic/gin"
	"github.com/labulaka521/logging"
	"github.com/micro/go-micro/client"
	"time"
)

var (
	ActuatorClient pbactuator.ActuatorService
)

func Init() {
	c := client.NewClient(
		client.Retries(3),
		client.Registry(registry.Etcd(cfg.EtcdConfig.Endpoints...)),
	)
	ActuatorClient = pbactuator.NewActuatorService("crocodile.srv.actuator", c)

}

type QueryActuat struct {
	Name string `json:"name" validate:"required"`
}
type Actuat struct {
	Name      string `json:"name" validate:"required"`
	Address   []Addr `json:"address" validate:"required"`
	Createdby string `json:"createdby" validate:"required"`
}

type Addr struct {
	Ip string `json:"ip" validate:"required"`
}

// ""
// {
//    "name": "",
// 	  "address": [
// 	  		{"ip": "ip1"}
// 	  	]
// }
func CreateActuator(c *gin.Context) {
	var (
		app         response.Gin
		ctx         context.Context
		err         error
		loginuser   string
		exists      bool
		code        int32
		resp        *pbactuator.Response
		reqactuator pbactuator.Actuat
	)

	ctx, _ = context.WithTimeout(context.TODO(), time.Duration(cfg.MysqlConfig.MaxQueryTime)*time.Second)
	app = response.Gin{c}

	reqactuator = pbactuator.Actuat{}
	if err = bind.BindJson(c, &reqactuator); err != nil {
		code = e.ERR_BAD_REQUEST
		app.Response(code, nil)
		return
	}
	if loginuser, exists = c.Keys["user"].(string); !exists {
		code = e.ERR_TOKEN_INVALID
		app.Response(code, nil)
		return
	}

	reqactuator.Createdby = loginuser

	resp, err = ActuatorClient.CreateActuator(ctx, &reqactuator)
	if err != nil {
		logging.Errorf("CreateActuator Err: %v", err)
		code = e.ERR_CREATE_ACTUAT_FAIL
		app.Response(code, nil)
		return
	}
	app.Response(resp.Code, nil)
}
func DeleteActuator(c *gin.Context) {
	var (
		app            response.Gin
		deleteactuator pbactuator.Actuat
		ctx            context.Context
		err            error
		code           int32
		resp           *pbactuator.Response
		exits          bool
	)

	ctx, _ = context.WithTimeout(context.TODO(), time.Duration(cfg.MysqlConfig.MaxQueryTime)*time.Second)
	app = response.Gin{c}
	deleteactuator, exits = c.Keys["data"].(pbactuator.Actuat)
	if !exits {
		code = e.ERR_BAD_REQUEST
		app.Response(code, nil)
		return
	}

	resp, err = ActuatorClient.DeleteActuator(ctx, &deleteactuator)
	if err != nil {
		logging.Errorf("DeleteActuator Err: %v", err)
		code = e.ERR_DELETE_ACTUAT_FAIL
		app.Response(code, nil)
		return
	}
	app.Response(resp.Code, nil)
}

func ChangeActuator(c *gin.Context) {
	var (
		app            response.Gin
		changeactuator pbactuator.Actuat
		ctx            context.Context
		err            error
		code           int32
		resp           *pbactuator.Response
		exits          bool
	)

	ctx, _ = context.WithTimeout(context.TODO(), time.Duration(cfg.MysqlConfig.MaxQueryTime)*time.Second)
	app = response.Gin{c}

	changeactuator = pbactuator.Actuat{}
	changeactuator, exits = c.Keys["data"].(pbactuator.Actuat)
	if !exits {
		logging.Errorf("Not Exits Actuator Data")
		code = e.ERR_BAD_REQUEST
		app.Response(code, nil)
		return
	}

	resp, err = ActuatorClient.ChangeActuator(ctx, &changeactuator)
	if err != nil {
		logging.Errorf("CreateActuator Err: %v", err)
		code = e.ERR_CHANGE_ACTUAT_FAIL
		app.Response(code, nil)
		return
	}
	app.Response(resp.Code, nil)
}

func GetActuator(c *gin.Context) {
	var (
		app          response.Gin
		queryctuator pbactuator.Actuat
		ctx          context.Context
		err          error
		code         int32
		rsp          *pbactuator.Response
	)
	queryctuator = pbactuator.Actuat{}
	ctx, _ = context.WithTimeout(context.TODO(), time.Duration(cfg.MysqlConfig.MaxQueryTime)*time.Second)
	app = response.Gin{c}
	if err = bind.BindQuery(c, &queryctuator); err != nil {
		code = e.ERR_BAD_REQUEST
		app.Response(code, nil)
		return
	}

	rsp, err = ActuatorClient.GetActuator(ctx, &queryctuator)
	if err != nil {
		logging.Errorf("CreateActuator Err: %v", err)
		code = e.ERR_CHANGE_ACTUAT_FAIL
		app.Response(code, nil)
		return
	}
	app.Response(rsp.Code, rsp.Actuators)
}

func GetALLExecuteIP(c *gin.Context) {
	var (
		app  response.Gin
		ctx  context.Context
		err  error
		code int32
		rsp  *pbactuator.Response
	)
	ctx, _ = context.WithTimeout(context.TODO(), time.Duration(cfg.MysqlConfig.MaxQueryTime)*time.Second)
	app = response.Gin{c}
	rsp, err = ActuatorClient.GetAllExecutorIP(ctx, new(pbactuator.Actuat))
	if err != nil {
		logging.Errorf("GetAllExecutorIP Err:%v", err)
		code = e.ERR_GET_EXECUTOR_IP_FAIL
		app.Response(code, nil)
		return
	}
	app.Response(rsp.Code, rsp.ExecutorIps)
}