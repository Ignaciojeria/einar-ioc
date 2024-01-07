package ioc

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/heimdalr/dag"
)

var graph = dag.NewDAG()

var errs []error

type dependency any

type constructor any

type container struct {
	id                    string
	constructor           constructor
	constructorParameters []constructor
	dependency            dependency
}

type visitor struct {
}

func (v visitor) Visit(vertex dag.Vertexer) {
	_, key := vertex.Vertex()
	orderedDependencyKeys = append(orderedDependencyKeys, key.(string))
}

var dependencyContainerMap = make(map[string]container)

var orderedDependencyKeys []string

func Registry(vertex constructor, edges ...constructor) error {
	constructorKey, err := getConstructorKey(vertex)

	if err != nil {
		errs = append(errs, err)
		return err
	}

	if dependencyContainerMap[constructorKey].constructor != nil {
		return nil
	}

	constructorType := reflect.TypeOf(vertex)
	if constructorType.Kind() != reflect.Func {
		errs = append(errs, err)
		return errors.New("provided constructor is not a function")
	}
	var constructorParameterKeys []string

	for i := 0; i < len(edges); i++ {
		constructorKey, err := getConstructorKey(edges[i])
		if err != nil {
			errs = append(errs, err)
			return err
		}
		constructorParameterKeys = append(constructorParameterKeys, constructorKey)
	}

	id, err := graph.AddVertex(constructorKey)
	if err != nil {
		return err
	}

	dependencyContainerMap[constructorKey] = container{
		id:                    id,
		constructor:           vertex,
		constructorParameters: edges,
	}
	return nil
}

func getConstructorKey(constructor constructor) (string, error) {
	funcValue := reflect.ValueOf(constructor)

	if funcPtr := funcValue.Pointer(); funcPtr != 0 {
		funcForPC := runtime.FuncForPC(funcPtr)
		if funcForPC != nil {
			funcName := funcForPC.Name()

			parts := strings.Split(funcName, "/")
			lastPart := parts[len(parts)-1]
			subParts := strings.SplitN(lastPart, ".", 2)

			packageName := strings.Join(parts[:len(parts)-1], "/") + "/" + subParts[0]
			functionName := subParts[1]

			return packageName + "." + functionName, nil
		}
	}
	return "", errors.New("constructor key cant be empty")
}

func LoadDependencies() error {
	for _, v := range errs {
		return v
	}
	for _, v := range dependencyContainerMap {
		for _, z := range v.constructorParameters {
			c := getContainer(z)
			graph.AddEdge(v.id, c.id)
		}
	}
	graph.OrderedWalk(visitor{})

	for i := len(orderedDependencyKeys) - 1; i >= 0; i-- {
		key := orderedDependencyKeys[i]
		ctnr := dependencyContainerMap[key]
		value := reflect.ValueOf(ctnr.constructor)
		if err := dependencyRulesGuardClause(key, value); err != nil {
			return err
		}
		var args []reflect.Value

		for _, constructorParameter := range ctnr.constructorParameters {
			dependency, err := Get(constructorParameter)
			if err != nil {
				return err
			}
			args = append(args, reflect.ValueOf(dependency))
		}

		if err := validArgumentsGuardClause(key, value, args); err != nil {
			return err
		}
		result := value.Call(args)
		if err := resultRulesGuardClause(key, result); err != nil {
			return err
		}

		container := dependencyContainerMap[key]
		container.dependency = result[0].Interface()
		dependencyContainerMap[key] = container
	}
	return nil
}

func getContainer(c constructor) container {
	constructorKey, _ := getConstructorKey(c)
	return dependencyContainerMap[constructorKey]
}

func dependencyRulesGuardClause(key string, value reflect.Value) error {
	funcType := value.Type()

	numOut := funcType.NumOut()

	if numOut < 1 {
		return errors.New(key + ": " + "the function must have at least 1 return value")
	}

	if numOut == 2 {
		if funcType.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
			return errors.New(key + ": " + "the second return value must be of type error")
		}
	}

	return nil
}

func validArgumentsGuardClause(key string, value reflect.Value, args []reflect.Value) error {
	funcType := value.Type()

	if funcType.NumIn() != len(args) {
		return fmt.Errorf("error in %s: incorrect number of arguments, expected %d, but got %d", key, funcType.NumIn(), len(args))
	}

	for i := 0; i < len(args); i++ {
		if funcType.In(i) != args[i].Type() {
			return fmt.Errorf("error in %s: incorrect argument type for parameter %d, expected %v, but got %v", key, i, funcType.In(i), args[i].Type())
		}
	}

	return nil
}

func resultRulesGuardClause(key string, result []reflect.Value) error {
	if len(result) == 2 {
		secondVal := result[1]

		if secondVal.IsNil() {
			return nil
		}

		if _, ok := secondVal.Interface().(error); !ok {
			return errors.New(key + ": " + "the second return value must be of type error")
		}

		if err := secondVal.Interface(); err != nil {
			if actualErr, isErr := err.(error); isErr && actualErr != nil {
				return actualErr
			}
		}
	}
	return nil
}

func Get(c constructor) (dependency, error) {
	constructorKey, err := getConstructorKey(c)
	if err != nil {
		return nil, err
	}
	dependency := dependencyContainerMap[constructorKey].dependency
	if dependency != nil {
		return dependency, nil
	}
	return nil, errors.New(constructorKey + "dependency is not present")
}
