package module

import "github.com/heimdalr/dag"

type IModule interface {
	DAG() *dag.DAG
}
