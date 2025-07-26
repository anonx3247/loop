package env

type Env struct {
	vars map[string]Var
}

type Var struct {
	Const bool
	Name  string
	Value Value
}

func NewEnv() *Env {
	return &Env{vars: make(map[string]Var)}
}

func (e *Env) Set(name string, value Value, isConst bool) {
	e.vars[name] = Var{
		Const: isConst,
		Name:  name,
		Value: value,
	}
}

func (e *Env) Get(name string) (Value, bool) {
	value, ok := e.vars[name]
	return value.Value, ok
}
