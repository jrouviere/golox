package parser

type Env struct {
	parent *Env
	values map[string]interface{}
}

func NewEnv(parent *Env) *Env {
	return &Env{
		parent: parent,
		values: make(map[string]interface{}),
	}
}

func (e *Env) Define(name string, value interface{}) {
	e.values[name] = value
}

func (e *Env) Get(name string) (interface{}, error) {
	if v, ok := e.values[name]; ok {
		return v, nil
	}
	if e.parent != nil {
		return e.parent.Get(name)
	}
	return nil, &RuntimeError{Msg: "undefined variable " + name}
}

func (e *Env) Set(name string, value interface{}) error {
	if _, ok := e.values[name]; ok {
		e.values[name] = value
		return nil
	}
	if e.parent != nil {
		return e.parent.Set(name, value)
	}
	return &RuntimeError{Msg: "undefined variable " + name}
}
