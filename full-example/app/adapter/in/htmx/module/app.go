package module

import (
	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/heimdalr/dag"
)

var _ = ioc.Registry(NewAppModule)

type AppModule struct {
	d *dag.DAG
}

func NewAppModule() IModule {
	return AppModule{
		d: dag.NewDAG(),
	}
}

func (m AppModule) DAG() *dag.DAG {
	return m.d
}
