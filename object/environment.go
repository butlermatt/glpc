package object

type Environment struct {
	parent *Environment
	m map[string]Object
}

const __GlobalEnv__ string = "__global__"

var fileScopes map[string]*Environment = make(map[string]*Environment)

func NewEnvironment(filename string) *Environment {
	env := &Environment{m: make(map[string]Object)}
	fileScopes[filename] = env
	return env
}

func NewEnclosedEnvironment(enclosing *Environment) *Environment {
	return &Environment{parent: enclosing, m: make(map[string]Object)}
}

func GetGlobal() *Environment {
	if env, ok := fileScopes[__GlobalEnv__]; ok {
		return env
	}

	env := NewEnvironment(__GlobalEnv__)
	// TODO: Populate builtin functions.
	return env
}
