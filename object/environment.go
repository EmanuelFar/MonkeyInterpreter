package object


type Environment struct {
  store   map[string]Object
  outer   *Environment
}

func NewEnvironment() *Environment {
  s := make(map[string]Object)
  return &Environment{store: s, outer: nil}
}

func (e *Environment) Get(name string) (Object, bool) {
  obj, ok := e.store[name]
  if !ok && e.outer != nil {
    // if Object isn't in current environment, search in outer environment.
    obj, ok = e.outer.Get(name)
  }
  return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
  e.store[name] = val
  return val
}

// We create a new environment for functions in order to prevent data override.
func NewEnclosedEnvironment(outer *Environment) *Environment {
  env := NewEnvironment()
  env.outer = outer
  return env
}
