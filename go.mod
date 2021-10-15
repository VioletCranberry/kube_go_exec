module kube_go_exec

go 1.16

require (
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/viper v1.8.1
	kube_go_exec/utils v0.0.0-00010101000000-000000000000
)

replace kube_go_exec/utils => ./utils
