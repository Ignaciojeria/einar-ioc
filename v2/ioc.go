package ioc

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/heimdalr/dag"
)

type dependency any
type constructor any

type iocContainer struct {
	graph                  *dag.DAG
	errs                   []error
	dependencyContainerMap map[string]container
	orderedDependencyKeys  []string
	atEndConstructor       *container
}

type container struct {
	id                    string
	constructor           constructor
	constructorParameters []constructor
	dependency            dependency
}

// New creates a new instance of the Container.
func New() *iocContainer {
	return &iocContainer{
		graph:                  dag.NewDAG(),
		errs:                   []error{},
		dependencyContainerMap: make(map[string]container),
		orderedDependencyKeys:  []string{},
		atEndConstructor:       nil,
	}
}

type visitor struct {
	orderedDependencyKeys *[]string
}

func (v visitor) Visit(vertex dag.Vertexer) {
	_, key := vertex.Vertex()
	*v.orderedDependencyKeys = append(*v.orderedDependencyKeys, key.(string))
}

// Registry registers a constructor and its dependencies.
func (c *iocContainer) Registry(vertex constructor, edges ...constructor) {
	constructorKey, err := getConstructorKey(vertex)

	if err != nil {
		c.errs = append(c.errs, err)
		return
	}

	if c.dependencyContainerMap[constructorKey].constructor != nil {
		panic(fmt.Errorf("constructor already registered: %v", constructorKey))
	}

	constructorType := reflect.TypeOf(vertex)
	if constructorType.Kind() != reflect.Func {
		c.errs = append(c.errs, err)
		return
	}

	id, err := c.graph.AddVertex(constructorKey)
	if err != nil {
		c.errs = append(c.errs, err)
		return
	}

	c.dependencyContainerMap[constructorKey] = container{
		id:                    id,
		constructor:           vertex,
		constructorParameters: edges,
	}
}

// RegistryAtEnd registers a constructor to be initialized at the end of all other dependencies.
func (c *iocContainer) RegistryAtEnd(vertex constructor, edges ...constructor) {
	if c.atEndConstructor != nil {
		panic("RegistryAtEnd can only be called once")
	}

	constructorKey, err := getConstructorKey(vertex)
	if err != nil {
		panic(err)
	}

	c.atEndConstructor = &container{
		id:                    constructorKey,
		constructor:           vertex,
		constructorParameters: edges,
	}
}

// LoadDependencies initializes all registered dependencies in the correct order.
func (c *iocContainer) LoadDependencies() error {
	// Primero, procesar cualquier error acumulado
	for _, v := range c.errs {
		return v
	}

	// Construir el grafo de dependencias y ordenar las claves de los constructores
	for _, v := range c.dependencyContainerMap {
		for _, z := range v.constructorParameters {
			ctnr := c.getContainer(z)
			c.graph.AddEdge(v.id, ctnr.id)
		}
	}
	c.graph.OrderedWalk(visitor{orderedDependencyKeys: &c.orderedDependencyKeys})

	// Ejecutar los constructores en el orden correcto
	for i := len(c.orderedDependencyKeys) - 1; i >= 0; i-- {
		key := c.orderedDependencyKeys[i]
		ctnr := c.dependencyContainerMap[key]
		value := reflect.ValueOf(ctnr.constructor)
		if err := dependencyRulesGuardClause(key, value); err != nil {
			return err
		}
		var args []reflect.Value

		for _, constructorParameter := range ctnr.constructorParameters {
			dependency, err := c.get(constructorParameter)
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

		container := c.dependencyContainerMap[key]
		if len(result) != 0 {
			container.dependency = result[0].Interface()
		}
		c.dependencyContainerMap[key] = container
	}

	// Ejecutar el constructor registrado al final si existe
	if c.atEndConstructor != nil {
		endValue := reflect.ValueOf(c.atEndConstructor.constructor)
		var endArgs []reflect.Value

		for _, constructorParameter := range c.atEndConstructor.constructorParameters {
			dependency, err := c.get(constructorParameter)
			if err != nil {
				return err
			}
			endArgs = append(endArgs, reflect.ValueOf(dependency))
		}

		if err := validArgumentsGuardClause(c.atEndConstructor.id, endValue, endArgs); err != nil {
			return err
		}
		result := endValue.Call(endArgs)
		if err := resultRulesGuardClause(c.atEndConstructor.id, result); err != nil {
			return err
		}
		if len(result) != 0 {
			c.atEndConstructor.dependency = result[0].Interface()
		}
	}

	return nil
}

func (c *iocContainer) getContainer(constructor constructor) container {
	constructorKey, _ := getConstructorKey(constructor)
	return c.dependencyContainerMap[constructorKey]
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

func (c *iocContainer) get(constructor constructor) (dependency, error) {
	constructorKey, err := getConstructorKey(constructor)
	if err != nil {
		return nil, err
	}
	dependency := c.dependencyContainerMap[constructorKey].dependency
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
