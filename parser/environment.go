package parser

type Env struct {
	values map[string]interface{}
}

func NewEnv() *Env {
	return &Env{
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
	return nil, &RuntimeError{Msg: "undefined variable " + name}
}

func (e *Env) Set(name string, value interface{}) error {
	if _, ok := e.values[name]; ok {
		e.values[name] = value
		return nil
	}
	return &RuntimeError{Msg: "undefined variable " + name}
}
