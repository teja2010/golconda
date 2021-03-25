package args

import (
	. "github.com/teja2010/golconda/src/any"

)

type arg struct {
	names []string
	val AnyValue
}

func _arg(names []string) arg {
	return arg{names, NoneValue()}
}

func getArgs() []arg {
	return []arg{
		_arg([]string{"--help", "-h"}),
	}
}

func ParseArgs () map[string]AnyValue {

	_ = getArgs()

	argMap := make(map[string]AnyValue)

	return argMap
}
