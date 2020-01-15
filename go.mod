module SLGGAME

go 1.13

replace github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0

require (
	github.com/coreos/etcd v3.3.18+incompatible
	github.com/golang/protobuf v1.3.2
	github.com/sirupsen/logrus v1.4.2
	github.com/yaice-rx/yaice v0.0.0-20200115032937-d242ec87a366
	go.uber.org/zap v1.13.0
)
