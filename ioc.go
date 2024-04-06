package ioc

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"

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

func Registry(vertex constructor, edges ...constructor) {
	constructorKey, err := getConstructorKey(vertex)

	if err != nil {
		errs = append(errs, err)
		return
	}

	if dependencyContainerMap[constructorKey].constructor != nil {
		return
	}

	constructorType := reflect.TypeOf(vertex)
	if constructorType.Kind() != reflect.Func {
		errs = append(errs, err)
		return
	}

	id, err := graph.AddVertex(constructorKey)
	if err != nil {
		errs = append(errs, err)
		return
	}

	dependencyContainerMap[constructorKey] = container{
		id:                    id,
		constructor:           vertex,
		constructorParameters: edges,
	}
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
			dependency, err := get(constructorParameter)
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
		if len(result) != 0 {
			container.dependency = result[0].Interface()
		}
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

	if numOut > 2 {
		return errors.New(key + ": " + "the function must have no more than 2 return values")
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
		expectedType := funcType.In(i)
		argType := args[i].Type()

		// Si el tipo esperado es una interfaz, comprueba si el tipo del argumento la implementa
		if expectedType.Kind() == reflect.Interface && !argType.Implements(expectedType) {
			return fmt.Errorf("error in %s: argument %d does not implement the expected interface %v, but got %v", key, i, expectedType, argType)
		}

		// Si no es interfaz, comprueba la igualdad de tipos
		if expectedType.Kind() != reflect.Interface && expectedType != argType {
			return fmt.Errorf("error in %s: incorrect argument type for parameter %d, expected %v, but got %v", key, i, expectedType, argType)
		}
	}

	return nil
}

func resultRulesGuardClause(key string, result []reflect.Value) error {

	if len(result) == 1 {
		firstVal := result[0]
		if firstVal.Kind() == reflect.Ptr || firstVal.Kind() == reflect.Slice || firstVal.Kind() == reflect.Map ||
			firstVal.Kind() == reflect.Chan || firstVal.Kind() == reflect.Func || firstVal.Kind() == reflect.Interface {
			if !firstVal.IsNil() {
				if err, ok := firstVal.Interface().(error); ok {
					return err
				}
			}
		}
	}

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

func get(c constructor) (dependency, error) {
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

func Get[T any](c constructor) T {
	constructorKey, err := getConstructorKey(c)
	if err != nil {
		panic(fmt.Errorf("failed to get constructor key: %v", err))
	}
	dependency := dependencyContainerMap[constructorKey].dependency
	if dependency == nil {
		panic(fmt.Errorf(constructorKey + " dependency is not present"))
	}
	return dependency.(T)
}

var mu map[string]*sync.Mutex = make(map[string]*sync.Mutex)
var muLock sync.Mutex

type mockBehaviour[T any] struct {
	constructorKey string
	originalValue  T
}

func NewMockBehaviourForTesting[T any](c constructor, mock T) mockBehaviour[T] {
	constructorKey, err := getConstructorKey(c)
	if err != nil {
		panic(fmt.Errorf("failed to get constructor key: %v", err))
	}

	muLock.Lock()
	m, ok := mu[constructorKey]
	if !ok {
		m = &sync.Mutex{}
		mu[constructorKey] = m
	}
	muLock.Unlock()

	m.Lock()

	originalValue := dependencyContainerMap[constructorKey].dependency

	dependencyRefMap := dependencyContainerMap[constructorKey]
	dependencyRefMap.dependency = mock
	dependencyContainerMap[constructorKey] = dependencyRefMap

	return mockBehaviour[T]{
		constructorKey: constructorKey,
		originalValue:  originalValue.(T),
	}
}

func (mb mockBehaviour[T]) Release() {
	if m, ok := mu[mb.constructorKey]; ok {
		mockedRef := dependencyContainerMap[mb.constructorKey]
		mockedRef.dependency = mb.originalValue
		dependencyContainerMap[mb.constructorKey] = mockedRef
		m.Unlock()
	}
}
