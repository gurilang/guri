package tengo

func ToMapMethod(o ToMapConverter) *UserFunctionCtx {
	return &UserFunctionCtx{Value: func(ctx *CallContext) (ret Object, err error) {
		switch len(ctx.Args) {
		case 0:
			return o.ToMap(false), nil
		case 1:
			return o.ToMap(!ctx.Args[0].IsFalsy()), nil
		default:
			return nil, ErrWrongNumArguments
		}
	}}
}

func FromMapMethod(o FromMapLoader) *UserFunctionCtx {
	return &UserFunctionCtx{Value: func(ctx *CallContext) (ret Object, err error) {
		switch len(ctx.Args) {
		case 1:
			if m, ok := ctx.Args[0].(ToMapConverter); ok {
				return o.(Object), o.FromMap(m.ToMap(true))
			}
			return nil, ErrInvalidArgumentType{
				Name:     "input",
				Expected: "map",
				Found:    ctx.Args[0].TypeName(),
			}
		default:
			return nil, ErrWrongNumArguments
		}
	}}
}
