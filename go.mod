module learn_together

go 1.14

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/casbin/casbin/v2 v2.25.5
	github.com/coreos/etcd v3.3.13+incompatible
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-playground/validator/v10 v10.4.1
	github.com/go-redis/redis/v8 v8.8.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/uuid v1.2.0 // indirect
	github.com/hashicorp/go-uuid v1.0.1
	github.com/jordan-wright/email v4.0.1-0.20210109023952-943e75fe5223+incompatible
	github.com/labstack/echo/v4 v4.2.1
	github.com/lib/pq v1.7.0
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/spf13/viper v1.7.1
	go.etcd.io/etcd v3.3.25+incompatible
	golang.org/x/crypto v0.0.0-20210314154223-e6e6c4f2bb5b // indirect
	golang.org/x/net v0.0.0-20210316092652-d523dce5a7f4 // indirect
	golang.org/x/sys v0.0.0-20210316092937-0b90fd5c4c48 // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/grpc v1.26.0
	google.golang.org/protobuf v1.26.0
	xorm.io/xorm v1.0.7
)

replace (
	github.com/coreos/etcd => github.com/ozonru/etcd v3.3.20-grpc1.27-origmodule+incompatible
	google.golang.org/grpc => google.golang.org/grpc v1.27.0
)
