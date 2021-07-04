module meepo

go 1.16

replace utility-go => ../utility-go

require (
	github.com/AsynkronIT/protoactor-go v0.0.0-20210505180410-df90efd4b2b4
	github.com/StackExchange/wmi v0.0.0-20210224194228-fe8f1750fd46 // indirect
	github.com/c-bata/go-prompt v0.2.6
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/gogo/protobuf v1.3.2
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
	github.com/google/uuid v1.2.0
	github.com/hashicorp/golang-lru v0.5.4
	github.com/hashicorp/memberlist v0.2.4
	github.com/pborman/uuid v1.2.1
	github.com/pkg/errors v0.9.1
	github.com/shirou/gopsutil v3.21.4+incompatible
	github.com/vadv/gopher-lua-libs v0.1.2
	github.com/yuin/gopher-lua v0.0.0-20200816102855-ee81675732da
	go.uber.org/zap v1.16.0
	gopkg.in/yaml.v2 v2.4.0
	layeh.com/gopher-lfs v0.0.0-20201124131141-d5fb28581d14
	utility-go v0.0.0-00010101000000-000000000000
)
