package parser

type LoxFunction struct {
	Declaration *FunStmt
}

func (l *LoxFunction) Arity() int {
	return len(l.Declaration.params)
}

func (l *LoxFunction) Call(env *Env, args []interface{}) (interface{}, error) {
	fnEnv := NewEnv(env.Root())

	for i := range args {
		fnEnv.Define(l.Declaration.params[i].Lexeme, args[i])
	}

	err := l.Declaration.body.Evaluate(fnEnv)
	if err != nil {
		if rv, ok := err.(*ReturnValue); ok {
			return rv.val, nil
		}
		return nil, err
	}
	return nil, nil
}

func (l *LoxFunction) String() string {
	return "<fn " + l.Declaration.name.Lexeme + ">"
}
