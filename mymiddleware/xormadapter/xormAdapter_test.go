// Copyright 2017 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package xormadapter

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/casbin/casbin/v2"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func testGetPolicy(t *testing.T, e *casbin.Enforcer, res [][]string) {
	t.Helper()
	myRes := e.GetPolicy()
	log.Print("Policy: ", myRes)

	m := make(map[string]bool, len(res))
	for _, value := range res {
		key := strings.Join(value, ",")
		m[key] = true
	}

	for _, value := range myRes {
		key := strings.Join(value, ",")
		if !m[key] {
			t.Error("Policy: ", myRes, ", supposed to be ", res)
			break
		}
	}
}
//初始化策略
func initPolicy(t *testing.T, driverName string, dataSourceName string, dbSpecified ...bool) {
	// Because the DB is empty at first, 英文数据库首先是空的
	// 我们需要从文件中加载一些策略到数据库中
	// so we need to load the policy from the file adapter (.CSV) first.
	e, _ := casbin.NewEnforcer("examples/rbac_model.conf", "examples/rbac_policy.csv")

	// 新建一个适配器 xorm的
	a, err := NewAdapter(driverName, dataSourceName, dbSpecified...)
	if err != nil {
		panic(err)
	}

	// 这是将当前策略保存到数据库的技巧。
	// 我们无法调用e.SavePolicy（），因为强制执行器中的适配器仍然是文件适配器。
	// 当前策略表示Casbin强制程序中的策略（也称为内存中的策略）).
	err = a.SavePolicy(e.GetModel())
	if err != nil {
		panic(err)
	}

	// Clear the current policy.
	e.ClearPolicy()
	testGetPolicy(t, e, [][]string{})

	// Load the policy from DB.
	err = a.LoadPolicy(e.GetModel())
	if err != nil {
		panic(err)
	}
	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}})
}

func testSaveLoad(t *testing.T, driverName string, dataSourceName string, dbSpecified ...bool) {
	// 初始化数据库的 一些策略
	initPolicy(t, driverName, dataSourceName, dbSpecified...)
	// Note: you don't need to look at the above code
	// if you already have a working DB with policy inside.

	// Now the DB has policy, so we can provide a normal use case.
	// Create an adapter and an enforcer.
	// NewEnforcer() will load the policy automatically.
	a, _ := NewAdapter(driverName, dataSourceName, dbSpecified...)
	e, _ := casbin.NewEnforcer("examples/rbac_model.conf", a)
	//使用rbac 则用以下的策略即可
	//p, data2_admin, data2, read
	//p  角色  资源  行为
	//g  用户 角色    这样就行了
	fmt.Println(e.GetAllRoles())   //[data2_admin]
	fmt.Println(e.GetAllSubjects()) //主体 传送过来的可以是角色 也 可以是 用户哦
	fmt.Println(e.GetAllObjects()) //资源 data1 与 data2
	fmt.Println(e.GetAllActions()) //权限 read，write
	fmt.Println(e.Enforce("alice","data2","read"))
	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}})
}

func testAutoSave(t *testing.T, driverName string, dataSourceName string, dbSpecified ...bool) {
	// 初始化策略数据库
	initPolicy(t, driverName, dataSourceName, dbSpecified...)
	// Note: you don't need to look at the above code
	// if you already have a working DB with policy inside.  已经有了策略数据库

	// 现在数据库有了策略，因此我们可以提供一个正常的用例。
	// Create an adapter and an enforcer. 创建一个适配器和一个执行器。
	// NewEnforcer() will load the policy automatically.     NewEnforcer（）将自动加载策略。
	a, _ := NewAdapter(driverName, dataSourceName, dbSpecified...)
	e, _ := casbin.NewEnforcer("examples/rbac_model.conf", a)

	// AutoSave is enabled by default. 自动保存是默认的
	// Now we disable it.  现在我们先禁止它
	e.EnableAutoSave(false)

	var err error
	logErr := func(action string) {
		if err != nil {
			t.Fatalf("test action[%s] failed, err: %v", action, err)
		}
	}

	// Because AutoSave is disabled, the policy change only affects the policy in Casbin enforcer,
	// 由于已禁用自动保存功能，因此策略更改仅会影响Casbin执行器中的策略，
	// it doesn't affect the policy in the storage.  不影响存储的策略
	_, err = e.AddPolicy("alice", "data1", "write")
	logErr("AddPolicy")

	// 从存储中重新加载 看看是否影响了
	// Reload the policy from the storage to see the effect.
	err = e.LoadPolicy()
	logErr("LoadPolicy")
	//  This is still the original policy.
	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}})

	// Now we enable the AutoSave. 打开自动保存
	e.EnableAutoSave(true)

	//启动自动保存 都会影响的
	// Because AutoSave is enabled, the policy change not only affects the policy in Casbin enforcer,
	// but also affects the policy in the storage.
	_, err = e.AddPolicy("alice", "data1", "write")
	logErr("AddPolicy2")
	// Reload the policy from the storage to see the effect.
	err = e.LoadPolicy()
	logErr("LoadPolicy2")
	// The policy has a new rule: {"alice", "data1", "write"}.
	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}, {"alice", "data1", "write"}})

	// Remove the added rule.
	_, err = e.RemovePolicy("alice", "data1", "write")
	logErr("RemovePolicy")
	err = e.LoadPolicy()
	logErr("LoadPolicy3")
	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}})

	//删除两个策略规则
	// Remove "data2_admin" related policy rules via a filter. 通过过滤器删除与“ data2_admin”相关的策略规则。
	// Two rules: {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"} are deleted.
	_, err = e.RemoveFilteredPolicy(0, "data2_admin")
	logErr("RemoveFilteredPolicy")
	err = e.LoadPolicy()
	logErr("LoadPolicy4")

	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}})
}

func testFilteredPolicy(t *testing.T, driverName string, dataSourceName string, dbSpecified ...bool) {
	// Initialize some policy in DB.
	initPolicy(t, driverName, dataSourceName, dbSpecified...)
	// Note: you don't need to look at the above code
	// if you already have a working DB with policy inside.

	// Now the DB has policy, so we can provide a normal use case.
	// Create an adapter and an enforcer.
	// NewEnforcer() will load the policy automatically.
	a, _ := NewAdapter(driverName, dataSourceName, dbSpecified...)
	e, _ := casbin.NewEnforcer("examples/rbac_model.conf")
	// Now set the adapter 这里也能设置 adapter
	e.SetAdapter(a)

	var err error
	logErr := func(action string) {
		if err != nil {
			t.Fatalf("test action[%s] failed, err: %v", action, err)
		}
	}

	// Load only alice's policies    根据 sub获取去策略 policy
	err = e.LoadFilteredPolicy(Filter{V0: []string{"alice"}})
	logErr("LoadFilteredPolicy")
	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}})

	// Load only bob's policies
	err = e.LoadFilteredPolicy(Filter{V0: []string{"bob"}})
	logErr("LoadFilteredPolicy2")
	testGetPolicy(t, e, [][]string{{"bob", "data2", "write"}})

	// Load policies for data2_admin
	err = e.LoadFilteredPolicy(Filter{V0: []string{"data2_admin"}})
	logErr("LoadFilteredPolicy3")
	testGetPolicy(t, e, [][]string{{"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}})

	// Load policies for alice and bob  可以传多个角色 然后获取策略
	err = e.LoadFilteredPolicy(Filter{V0: []string{"alice", "bob"}})
	logErr("LoadFilteredPolicy4")
	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}})
}

func testRemovePolicies(t *testing.T, driverName string, dataSourceName string, dbSpecified ...bool) {
	// Initialize some policy in DB.
	initPolicy(t, driverName, dataSourceName, dbSpecified...)
	// Note: you don't need to look at the above code
	// if you already have a working DB with policy inside.

	// Now the DB has policy, so we can provide a normal use case.
	// Create an adapter and an enforcer.
	// NewEnforcer() will load the policy automatically.
	a, _ := NewAdapter(driverName, dataSourceName, dbSpecified...)
	e, _ := casbin.NewEnforcer("examples/rbac_model.conf")

	// Now set the adapter
	e.SetAdapter(a)

	var err error
	logErr := func(action string) {
		if err != nil {
			t.Fatalf("test action[%s] failed, err: %v", action, err)
		}
	}

	err = a.AddPolicies("p", "p", [][]string{{"max", "data2", "read"}, {"max", "data1", "write"}, {"max", "data1", "delete"}})
	logErr("AddPolicies")

	// Load policies for max
	err = e.LoadFilteredPolicy(Filter{V0: []string{"max"}})
	logErr("LoadFilteredPolicy")

	testGetPolicy(t, e, [][]string{{"max", "data2", "read"}, {"max", "data1", "write"}, {"max", "data1", "delete"}})

	// Remove policies
	err = a.RemovePolicies("p", "p", [][]string{{"max", "data2", "read"}, {"max", "data1", "write"}})
	logErr("RemovePolicies")

	// Reload policies for max
	err = e.LoadFilteredPolicy(Filter{V0: []string{"max"}})
	logErr("LoadFilteredPolicy2")

	testGetPolicy(t, e, [][]string{{"max", "data1", "delete"}})
}

func testAddPolicies(t *testing.T, driverName string, dataSourceName string, dbSpecified ...bool) {
	// Initialize some policy in DB.
	initPolicy(t, driverName, dataSourceName, dbSpecified...)
	// Note: you don't need to look at the above code
	// if you already have a working DB with policy inside.

	// Now the DB has policy, so we can provide a normal use case.
	// Create an adapter and an enforcer.
	// NewEnforcer() will load the policy automatically.
	a, _ := NewAdapter(driverName, dataSourceName, dbSpecified...)
	e, _ := casbin.NewEnforcer("examples/rbac_model.conf")

	// Now set the adapter
	e.SetAdapter(a)

	var err error
	logErr := func(action string) {
		if err != nil {
			t.Fatalf("test action[%s] failed, err: %v", action, err)
		}
	}

	//为某个角色添加 资源的权限
	err = a.AddPolicies("p", "p", [][]string{{"max", "data2", "read"}, {"max", "data1", "write"}})
	logErr("AddPolicies")

	// Load policies for max  然后根据过滤器 获取max的策略
	err = e.LoadFilteredPolicy(Filter{V0: []string{"max"}})
	logErr("LoadFilteredPolicy")

	testGetPolicy(t, e, [][]string{{"max", "data2", "read"}, {"max", "data1", "write"}})
}

func testUpdatePolicies(t *testing.T, driverName string, dataSourceName string, dbSpecified ...bool) {
	// Initialize some policy in DB.
	initPolicy(t, driverName, dataSourceName, dbSpecified...)
	// Note: you don't need to look at the above code
	// if you already have a working DB with policy inside.

	// Now the DB has policy, so we can provide a normal use case.
	// Create an adapter and an enforcer.
	// NewEnforcer() will load the policy automatically.
	a, _ := NewAdapter(driverName, dataSourceName, dbSpecified...)
	e, _ := casbin.NewEnforcer("examples/rbac_model.conf")

	// Now set the adapter
	e.SetAdapter(a)

	var err error
	logErr := func(action string) {
		if err != nil {
			t.Fatalf("test action[%s] failed, err: %v", action, err)
		}
	}

	err = a.UpdatePolicy("p", "p", []string{"bob", "data2", "write"}, []string{"alice", "data2", "write"})
	logErr("UpdatePolicy")

	testGetPolicy(t, e, [][]string{{"alice", "data1", "read"}, {"alice", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}})
}

func TestAdapters(t *testing.T) {
	// 您还可以使用以下方式来使用现有的数据库“ abc”：  里面会有表
	//testSaveLoad(t, "mysql", "root:@tcp(127.0.0.1:3306)/abc", true)

	//testSaveLoad(t, "mysql", "root:123456@tcp(127.0.0.1:3306)/")

	//自动保存哦
	testAutoSave(t, "mysql", "root:123456@tcp(127.0.0.1:3306)/")

	//过滤器 获取对应的 policy
	//testFilteredPolicy(t, "mysql", "root:123456@tcp(127.0.0.1:3306)/")

	//添加策略
	//testAddPolicies(t, "mysql", "root:@tcp(127.0.0.1:3306)/")

	//testRemovePolicies(t, "mysql", "root:@tcp(127.0.0.1:3306)/")

	//testUpdatePolicies(t, "mysql", "root:@tcp(127.0.0.1:3306)/")
}
