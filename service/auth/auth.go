package auth

import (
	"bufio"
	"context"
	"google.golang.org/grpc"
	"io"
	"io/ioutil"
	"learn_together/api/authcontrol"
	"learn_together/initer"
	"os"
	"time"
)

var Enforcer *authcontrol.Enforcer

func AuthInit(config *initer.Config) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	//远程过程调用 创建一个客户端
	cc, err := authcontrol.NewClient(ctx, config.Auth.AuthAddr, config.Auth.AuthServiceName, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	models, err := ioutil.ReadFile(config.Auth.ModelPath)
	if err != nil {
		panic(err)
	}

	//新建一个执行器
	Enforcer, err = cc.NewEnforcer(ctx, authcontrol.Config{
		ModelText:     string(models),
		ConnectString: config.Auth.Connection,
		DriverName:    config.Auth.Driver,
		DbSpecified:   config.Auth.DbSpecified,
	})
	if err != nil {
		panic(err)
	}
	if config.Auth.IsLoadPolicy {
		policy, err := os.Open(config.Auth.PolicyPath)
		if err != nil {
			panic(config.Auth.PolicyPath + "config open error")
		}
		policies := getLines(policy)
		//添加策略
		_, err = Enforcer.AddPolicies(ctx, policies)
		if err != nil {
			panic("add policy error" + err.Error())
		}
		err = Enforcer.SavePolicy(ctx)
		if err != nil {
			panic("save policy error")
		}
	}

}
func getLines(file *os.File) []string {
	strs := make([]string, 0)
	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			return strs
		}
		strs = append(strs, string(line))
	}
}
