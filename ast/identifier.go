package ast

import (
	"com.loop.anonx3247/env"
	"com.loop.anonx3247/utils"
)

type Identifier struct {
	source utils.String
}

func NewIdentifier(source utils.String) Identifier {
	return Identifier{source: source}
}

func (i Identifier) Source() utils.String {
	return i.source
}

func (i Identifier) Name() string {
	return i.source.String()
}

func (i Identifier) Eval(env *env.Env) (env.Value, error) {
	value, ok := env.Get(i.Name())
	if !ok {
		return nil, utils.Error{Source: i.Source(), Message: "variable not found"}
	}
	return value, nil
}
