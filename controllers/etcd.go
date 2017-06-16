package controllers

import (
	"context"
	"time"

	"encoding/json"

	log "chpkg.in/qiniu/log.v1"
	"github.com/astaxie/beego"
	"github.com/coreos/etcd/client"
)

type EtcdController struct {
	beego.Controller
}

type ObjectBody struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func GenerateKapi() (kapi client.KeysAPI) {

	var cfg = client.Config{
		Endpoints:               []string{"http://127.0.0.1:2379"},
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	c, err := client.New(cfg)
	if err != nil {
		log.Error(err)
		return
	}
	kapi = client.NewKeysAPI(c)

	return kapi
}

// @Title Create
// @Description create newkeyvalue
// @Param	key-value	string 	string	true		"The object content"
// @Success 200 {string} ob.key ob.value
// @Failure 403 body is empty
// @router / [post]
func (e *EtcdController) Post() {
	var ob ObjectBody
	err := json.Unmarshal(e.Ctx.Input.RequestBody, &ob)

	if err != nil {
		log.Error(err)
		return
	}

	kapi := GenerateKapi()
	resp, err := kapi.Set(context.Background(), ob.Key, ob.Value, nil)
	if err != nil {
		log.Error(err)
		return
	}
	log.Info("resp:", resp)
	e.Data["json"] = map[string]string{"key": ob.Key, "value": resp.Node.Value}

	e.ServeJSON()
}

// @Title Get
// @Description find value by key
// @Param	key		path 	string	true		"the key you want to get"
// @Success 200 value
// @Failure 403 :key is empty
// @router /:key [get]
func (e *EtcdController) Get() {
	kapi := GenerateKapi()
	key := e.Ctx.Input.Param(":key")
	if key != "" {
		v, err := kapi.Get(context.Background(), key, nil)
		if err != nil {
			e.Data["json"] = err.Error()
		} else {
			e.Data["json"] = v.Node.Value
		}
	}
	e.ServeJSON()
}

// @Title Delete
// @Description delete the object
// @Param	key		path 	string	true		"The key-value you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 key is empty
// @router /:key [delete]
func (e *EtcdController) Delete() {
	kapi := GenerateKapi()
	key := e.Ctx.Input.Param(":key")
	if key != "" {
		_, err := kapi.Delete(context.Background(), key, nil)
		if err != nil {
			e.Data["json"] = err.Error()
		} else {
			e.Data["json"] = "delete success"
		}
	}
	e.ServeJSON()
}
