package ioc

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/heimdalr/dag"
)

var (
	graph                  = dag.NewDAG()
	errs                   []error
	dependencyContainerMap = make(map[string]container)
	orderedDependencyKeys  []string
	atEndConstructor       *container
)

type dependency any
type constructor any

type container struct {
	id                    string
	constructor           constructor
	constructorParameters []constructor
	dependency            dependency
}

// visitor is used to walk the graph in a topological order
type visitor struct {
	orderedDependencyKeys *[]string
}

func (v visitor) Visit(vertex dag.Vertexer) {
	_, key := vertex.Vertex()
	*v.orderedDependencyKeys = append(*v.orderedDependencyKeys, key.(string))
}

// Registry registers a constructor and its dependencies.
func Registry(vertex constructor, edges ...constructor) {
	constructorKey, err := getConstructorKey(vertex)

	if err != nil {
		errs = append(errs, err)
		return
	}

	if dependencyContainerMap[constructorKey].constructor != nil {
		panic(fmt.Errorf("constructor already registered: %v", constructorKey))
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

// RegistryAtEnd registers a constructor to be initialized at the end of all other dependencies.
func RegistryAtEnd(vertex constructor, edges ...constructor) {
	if atEndConstructor != nil {
		panic("RegistryAtEnd can only be called once")
	}

	constructorKey, err := getConstructorKey(vertex)
	if err != nil {
		panic(err)
	}

	atEndConstructor = &container{
		id:                    constructorKey,
		constructor:           vertex,
		constructorParameters: edges,
	}
}

// LoadDependencies initializes all registered dependencies in the correct order.
func LoadDependencies() error {
	// Primero, procesar cualquier error acumulado
	for _, v := range errs {
		return v
	}

	// Construir el grafo de dependencias y ordenar las claves de los constructores
	for _, v := range dependencyContainerMap {
		for _, z := range v.constructorParameters {
			ctnr := getContainer(z)
			graph.AddEdge(v.id, ctnr.id)
		}
	}
	graph.OrderedWalk(visitor{orderedDependencyKeys: &orderedDependencyKeys})

	// Ejecutar los constructores en el orden correcto
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

	// Ejecutar el constructor registrado al final si existe
	if atEndConstructor != nil {
		endValue := reflect.ValueOf(atEndConstructor.constructor)
		var endArgs []reflect.Value

		for _, constructorParameter := range atEndConstructor.constructorParameters {
			dependency, err := get(constructorParameter)
			if err != nil {
				return err
			}
			endArgs = append(endArgs, reflect.ValueOf(dependency))
		}

		if err := validArgumentsGuardClause(atEndConstructor.id, endValue, endArgs); err != nil {
			return err
		}
		result := endValue.Call(endArgs)
		if err := resultRulesGuardClause(atEndConstructor.id, result); err != nil {
			return err
		}
		if len(result) != 0 {
			atEndConstructor.dependency = result[0].Interface()
		}
	}

	return nil
}

func getContainer(constructor constructor) container {
	constructorKey, _ := getConstructorKey(constructor)
	return dependencyContainerMap[constructorKey]
}

func dependencyRulesGuardClause(key string, value reflect.Value) error {
	funcType := value.Type()

	numOut := funcType.NumOut()

	if numOut > 2 {
		return errors.New(key + ": the function must have no more than 2 return values")
	}

	if numOut == 2 {
		if funcType.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
			return errors.New(key + ": the second return value must be of type error")
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

		if expectedType.Kind() == reflect.Interface && !argType.Implements(expectedType) {
			return fmt.Errorf("error in %s: argument %d does not implement the expected interface %v, but got %v", key, i, expectedType, argType)
		}

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
			return errors.New(key + ": the second return value must be of type error")
		}

		if err := secondVal.Interface(); err != nil {
			if actualErr, isErr := err.(error); isErr && actualErr != nil {
				return actualErr
			}
		}
	}
	return nil
}

func get(constructor constructor) (dependency, error) {
	constructorKey, err := getConstructorKey(constructor)
	if err != nil {
		return nil, err
	}
	dependency := dependencyContainerMap[constructorKey].dependency
	if dependency != nil {
		return dependency, nil
	}
	return nil, errors.New(constructorKey + " dependency is not present")
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
	return "", errors.New("constructor key can't be empty")
}
