package ioc

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
)

// Definiendo DependencyType como un tipo basado en string.
type dependencyType string

// Declarando constantes de tipo DependencyType.
const (
	installation    dependencyType = "Installation"
	outBoundAdapter dependencyType = "OutBound"
	useCase         dependencyType = "UseCase"
	inboundAdapter  dependencyType = "Inbound"
)

var errs []error

type dependency any

type constructor any

type container struct {
	Error                 error
	dependencyType        dependencyType
	constructor           constructor
	constructorParameters []constructor
	dependency            dependency
}

var installations = make(map[string]container)
var outboundAdapters = make(map[string]container)
var usecases = make(map[string]container)
var inboundAdapters = make(map[string]container)

func dependencyContainerFactory(t dependencyType) map[string]container {

	if t == installation {
		return installations
	}

	if t == outBoundAdapter {
		return outboundAdapters
	}

	if t == useCase {
		return usecases
	}

	if t == inboundAdapter {
		return inboundAdapters
	}
	fmt.Println("unknown dependency type")
	os.Exit(1)
	return nil
}

func getConstructorKey(constructor constructor) (string, error) {

	// Usando reflexión para obtener información sobre la función.
	funcValue := reflect.ValueOf(constructor)

	// Usando la función runtime.FuncForPC para obtener información sobre la función.
	if funcPtr := funcValue.Pointer(); funcPtr != 0 {
		funcForPC := runtime.FuncForPC(funcPtr)
		if funcForPC != nil {
			funcName := funcForPC.Name()

			// Separando el nombre completo para obtener el paquete y el nombre de la función.
			// El formato usual es "path/to/package.FuncName".
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

func registry(t dependencyType, c constructor, constructorParameters ...constructor) error {

	constructorKey, err := getConstructorKey(c)

	if err != nil {
		errs = append(errs, err)
		return err
	}

	if dependencyContainerFactory(t)[constructorKey].constructor != nil {
		fmt.Println("dependency already loaded")
		return nil
	}

	// Obtener el tipo reflect.Type del primer constructor.
	constructorType := reflect.TypeOf(c)
	if constructorType.Kind() != reflect.Func {
		errs = append(errs, err)
		return errors.New("provided constructor is not a function")
	}
	var constructorParameterKeys []string

	for i := 0; i < len(constructorParameters); i++ {
		constructorKey, err := getConstructorKey(constructorParameters[i])
		if err != nil {
			errs = append(errs, err)
			return err
		}
		constructorParameterKeys = append(constructorParameterKeys, constructorKey)
	}

	dependencyContainerFactory(t)[constructorKey] = container{
		Error:                 err,
		dependencyType:        t,
		constructor:           c,
		constructorParameters: constructorParameters,
	}

	return nil
}

func Installation(c constructor) error {
	return registry(installation, c)
}

func OutBoundAdapter(c constructor, args ...constructor) error {
	return registry(outBoundAdapter, c, args...)
}

func Get(c constructor) (dependency, error) {
	constructorKey, err := getConstructorKey(c)
	if err != nil {
		return nil, err
	}
	dependency := dependencyContainerFactory(installation)[constructorKey].dependency
	if dependency != nil {
		return dependency, nil
	}
	dependency = dependencyContainerFactory(outBoundAdapter)[constructorKey].dependency
	if dependency != nil {
		return dependency, nil
	}
	dependency = dependencyContainerFactory(useCase)[constructorKey].dependency
	if dependency != nil {
		return dependency, nil
	}
	dependency = dependencyContainerFactory(inboundAdapter)[constructorKey].dependency
	if dependency != nil {
		return dependency, nil
	}
	return nil, errors.New("dependency is not present")
}

func LoadDependencies() error {

	for _, v := range errs {
		return v
	}

	for key, ctnr := range installations {
		value := reflect.ValueOf(ctnr.constructor)
		if err := dependencyRulesGuardClause(key, value); err != nil {
			return err
		}
		// Obtiene el número de valores de retorno
		numOut := value.Type().NumOut()
		result := value.Call(nil)
		if err := resultRulesGuardClause(key, numOut, result); err != nil {
			return err
		}
		container := installations[key]
		container.dependency = result[0].Interface()
		installations[key] = container
	}

	for key, ctnr := range outboundAdapters {
		value := reflect.ValueOf(ctnr.constructor)
		if err := dependencyRulesGuardClause(key, value); err != nil {
			return err
		}
		numOut := value.Type().NumOut()

		var args []reflect.Value
		for _, constructorParameter := range ctnr.constructorParameters {
			dependency, err := Get(constructorParameter)
			if err != nil {
				return err
			}
			args = append(args, reflect.ValueOf(dependency))
		}

		result := value.Call(args)
		if err := resultRulesGuardClause(key, numOut, result); err != nil {
			return err
		}
	}

	return nil
}

func dependencyRulesGuardClause(key string, value reflect.Value) error {
	// Obtiene el número de valores de retorno
	numOut := value.Type().NumOut()

	// Comprobar que haya al menos 1 valor de retorno y no más de 2
	if numOut < 1 {
		return errors.New(key + ": " + "the function must have at least 1 return value")
	}
	if numOut > 2 {
		return errors.New(key + ": " + "the function must have no more than 2 return values")
	}

	// Si hay un segundo valor de retorno, comprobar si es de tipo error
	if numOut == 2 && value.Type().Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
		return errors.New(key + ": " + "The second return value must be of error type")
	}

	return nil
}

func resultRulesGuardClause(key string, numOut int, result []reflect.Value) error {
	// Si hay dos valores de retorno y el segundo es un error, verificar si es nil
	if numOut == 2 {
		if err, ok := result[1].Interface().(error); ok {
			if err != nil {
				return err
			}
		} else {
			return errors.New(key + ": " + "the second return value must be of type error")
		}
	}
	return nil
}

/*
func LoadDependencies() error {
	for _, v := range dependencyContainer {
	}
	return nil
}
*/

/*
var installations = make(map[string]container)

var outboundAdapters = make(map[string]container)

var usecases = make(map[string]container)

var inboundAdapters = make(map[string]container)



/*
var installations = make(map[string]loadableDependency)

type loadableDependency func() (any, error)

func Installation(d loadableDependency, l ...loadableDependency) {
	installations[uuid.NewString()] = d
	//adapter := container[T]{loadableDependency: loadableDependency}
	//installations[uuid.NewString()] = &adapter
	//return &adapter
}*/

/*
type container[T any] struct {
	loadableDependency func() (T, error)
	isLoaded           bool
	Dependency         T
}

func (c *container[T]) Load() (any, error) {
	if c.isLoaded {
		return nil, errors.New("dependency already loaded")
	}
	instance, err := c.loadableDependency()
	c.Dependency = instance
	c.isLoaded = true
	return instance, err
}

type loadable[T any] interface {
	Load() (any, error)
}



var installations = make(map[string]loadable[any])

func Installation[T any](loadableDependency func() (T, error)) *container[T] {
	adapter := container[T]{loadableDependency: loadableDependency}
	installations[uuid.NewString()] = &adapter
	return &adapter
}*/

/*
var useCases = make(map[string]loadable[any])

func UseCase[T any](loadableDependency func() (T, error)) *container[T] {
	adapter := container[T]{loadableDependency: loadableDependency}
	useCases[uuid.NewString()] = &adapter
	return &adapter
}

var inboundAdapters = make(map[string]loadable[any])

func InboundAdapter[T any](loadableDependency func() (T, error)) *container[T] {
	adapter := container[T]{loadableDependency: loadableDependency}
	inboundAdapters[uuid.NewString()] = &adapter
	return &adapter
}

var outBoundAdapters = make(map[string]loadable[any])

func OutBoundAdapter[T any](loadableDependency func() (T, error)) *container[T] {
	adapter := container[T]{loadableDependency: loadableDependency}
	outBoundAdapters[uuid.NewString()] = &adapter
	return &adapter
}

func LoadDependencies() error {
	for _, v := range installations {
		_, err := v.Load()
		if err != nil {
			return err
		}
	}
	for _, v := range outBoundAdapters {
		_, err := v.Load()
		if err != nil {
			return err
		New(text)
	}
	for _, v := range useCases {
		_, err := v.Load()
		if err != nil {
			return err
		}
	}
	for _, v := range inboundAdapters {
		_, err := v.Load()
		if err != nil {
			return err
		}
	}
	return nil
}*/
