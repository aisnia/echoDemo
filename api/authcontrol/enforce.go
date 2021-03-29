package authcontrol

import (
	"context"
	"encoding/json"
	"fmt"
	pb "learn_together/api/authcontrol/proto"
	"log"
	"reflect"
)

// Config contains data needed to create an enforcer.
type Config struct {
	DriverName    string
	ConnectString string
	ModelText     string
	DbSpecified   bool
}

// Enforcer is the main interface for authorization enforcement and policy management.
type Enforcer struct {
	handler int32
	client  *Client
}

// NewEnforcer creates an enforcer via file or DB.
// File:
// e := casbin.NewEnforcer("path/to/basic_model.conf", "path/to/basic_policy.csv")
// MySQL DB:
// a := mysqladapter.NewDBAdapter("mysql", "mysql_username:mysql_password@tcp(127.0.0.1:3306)/")
// e := casbin.NewEnforcer("path/to/basic_model.conf", a)

func (c *Client) NewEnforcer(ctx context.Context, config Config) (*Enforcer, error) {
	var adapterHandler int32 = -1
	enforcer := &Enforcer{client: c}

	// Maybe it does not need NewAdapter.  //新建一个adapter  默认是文件

	adapterReply, err := c.remoteClient.NewAdapter(ctx, &pb.NewAdapterRequest{
		DriverName:    config.DriverName,
		ConnectString: config.ConnectString,
		DbSpecified:   config.DbSpecified,
	})
	
	if err != nil {
		return enforcer, err
	}
	adapterHandler = adapterReply.Handler

	e, err := c.remoteClient.NewEnforcer(ctx, &pb.NewEnforcerRequest{
		ModelText:     config.ModelText,
		AdapterHandle: adapterHandler,
	})
	if err != nil {
		return enforcer, err
	}
	enforcer.handler = e.Handler
	return enforcer, nil
}
func makeABAC(obj interface{}) (string, error) {
	data, err := json.Marshal(&obj)
	if err != nil {
		return "", err
	}
	return "ABAC::" + string(data), nil
}

// Enforce decides whether a "subject" can access a "object" with the operation "action", input parameters are usually: (sub, obj, act).
func (e *Enforcer) Enforce(ctx context.Context, params ...interface{}) (bool, error) {
	var data []string
	for _, item := range params {
		var value string
		var err error
		if reflect.TypeOf(item).Kind() == reflect.Struct {
			value, err = makeABAC(data)
			if err != nil {
				return false, err
			}
		} else {
			value = fmt.Sprintf("%v", item)
		}
		data = append(data, value)
	}

	res, err := e.client.remoteClient.Enforce(ctx, &pb.EnforceRequest{
		EnforcerHandler: e.handler,
		Params:          data,
	})
	log.Println(res, err)
	return res.Res, err
}

// LoadPolicy reloads the policy from file/database.
func (e *Enforcer) LoadPolicy(ctx context.Context) error {
	_, err := e.client.remoteClient.LoadPolicy(ctx, &pb.EmptyRequest{Handler: e.handler})
	return err
}

// SavePolicy saves the current policy (usually after changed with Casbin API) back to file/database.
func (e *Enforcer) SavePolicy(ctx context.Context) error {
	_, err := e.client.remoteClient.SavePolicy(ctx, &pb.EmptyRequest{Handler: e.handler})
	return err
}
