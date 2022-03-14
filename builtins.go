package tengo

import (
	"fmt"
	"strconv"
)

var builtinFuncs = []*BuiltinFunction{
	{
		Name:  "len",
		Value: builtinLen,
	},
	{
		Name:  "copy",
		Value: builtinCopy,
	},
	{
		Name:  "append",
		Value: builtinAppend,
	},
	{
		Name:  "delete",
		Value: builtinDelete,
	},
	{
		Name:  "splice",
		Value: builtinSplice,
	},
	{
		Name:  "string",
		Value: builtinString,
	},
	{
		Name:  "int",
		Value: builtinInt,
	},
	{
		Name:  "bool",
		Value: builtinBool,
	},
	{
		Name:  "float",
		Value: builtinFloat,
	},
	{
		Name:  "char",
		Value: builtinChar,
	},
	{
		Name:  "bytes",
		Value: builtinBytes,
	},
	{
		Name:  "time",
		Value: builtinTime,
	},
	{
		Name:  "is_int",
		Value: builtinIsInt,
	},
	{
		Name:  "is_float",
		Value: builtinIsFloat,
	},
	{
		Name:  "is_string",
		Value: builtinIsString,
	},
	{
		Name:  "is_bool",
		Value: builtinIsBool,
	},
	{
		Name:  "is_char",
		Value: builtinIsChar,
	},
	{
		Name:  "is_bytes",
		Value: builtinIsBytes,
	},
	{
		Name:  "is_array",
		Value: builtinIsArray,
	},
	{
		Name:  "is_immutable_array",
		Value: builtinIsImmutableArray,
	},
	{
		Name:  "is_map",
		Value: builtinIsMap,
	},
	{
		Name:  "is_immutable_map",
		Value: builtinIsImmutableMap,
	},
	{
		Name:  "is_iterable",
		Value: builtinIsIterable,
	},
	{
		Name:  "is_time",
		Value: builtinIsTime,
	},
	{
		Name:  "is_error",
		Value: builtinIsError,
	},
	{
		Name:  "is_undefined",
		Value: builtinIsUndefined,
	},
	{
		Name:  "is_function",
		Value: builtinIsFunction,
	},
	{
		Name:  "is_callable",
		Value: builtinIsCallable,
	},
	{
		Name:  "type_name",
		Value: builtinTypeName,
	},
	{
		Name:  "format",
		Value: builtinFormat,
	},
	{
		Name:  "range",
		Value: builtinRange,
	},
	{
		Name:  "map",
		Value: builtinMap,
	},
	{
		Name:     "get_methods",
		Value:    builtinGetMethods,
		IsMethod: true,
	},
	{
		Name:  "type",
		Value: builtinType,
	},
	{
		Name:  "typeof",
		Value: builtinTypeOf,
	},
	{
		Name:  "field",
		Value: builtinField,
	},
	{
		Name:  "fields",
		Value: builtinFields,
	},
	{
		Name:  "method",
		Value: builtinMethod,
	},
	{
		Name:  "methods",
		Value: builtinMethods,
	},
	{
		Name:  "property",
		Value: builtinProperty,
	},
	{
		Name:  "properties",
		Value: builtinProperties,
	},
}

// GetAllBuiltinFunctions returns all builtin function objects.
func GetAllBuiltinFunctions() []*BuiltinFunction {
	return append([]*BuiltinFunction{}, builtinFuncs...)
}

func builtinTypeName(ctx *CallContext) (Object, error) {
	if len(ctx.Args) != 1 {
		return nil, ErrWrongNumArguments
	}
	return &String{Value: ctx.Args[0].TypeName()}, nil
}

func builtinIsString(ctx *CallContext) (Object, error) {
	if len(ctx.Args) != 1 {
		return nil, ErrWrongNumArguments
	}
	if _, ok := ctx.Args[0].(*String); ok {
		return TrueValue, nil
	}
	return FalseValue, nil
}

func builtinIsInt(ctx *CallContext) (Object, error) {
	if len(ctx.Args) != 1 {
		return nil, ErrWrongNumArguments
	}
	if _, ok := ctx.Args[0].(*Int); ok {
		return TrueValue, nil
	}
	return FalseValue, nil
}

func builtinIsFloat(ctx *CallContext) (Object, error) {
	if len(ctx.Args) != 1 {
		return nil, ErrWrongNumArguments
	}
	if _, ok := ctx.Args[0].(*Float); ok {
		return TrueValue, nil
	}
	return FalseValue, nil
}

func builtinIsBool(ctx *CallContext) (Object, error) {
	if len(ctx.Args) != 1 {
		return nil, ErrWrongNumArguments
	}
	if _, ok := ctx.Args[0].(*Bool); ok {
		return TrueValue, nil
	}
	return FalseValue, nil
}

func builtinIsChar(ctx *CallContext) (Object, error) {
	if len(ctx.Args) != 1 {
		return nil, ErrWrongNumArguments
	}
	if _, ok := ctx.Args[0].(*Char); ok {
		return TrueValue, nil
	}
	return FalseValue, nil
}

func builtinIsBytes(ctx *CallContext) (Object, error) {
	if len(ctx.Args) != 1 {
		return nil, ErrWrongNumArguments
	}
	if _, ok := ctx.Args[0].(*Bytes); ok {
		return TrueValue, nil
	}
	return FalseValue, nil
}

func builtinIsArray(ctx *CallContext) (Object, error) {
	if len(ctx.Args) != 1 {
		return nil, ErrWrongNumArguments
	}
	if _, ok := ctx.Args[0].(*Array); ok {
		return TrueValue, nil
	}
	return FalseValue, nil
}

func builtinIsImmutableArray(ctx *CallContext) (Object, error) {
	if len(ctx.Args) != 1 {
		return nil, ErrWrongNumArguments
	}
	if _, ok := ctx.Args[0].(*ImmutableArray); ok {
		return TrueValue, nil
	}
	return FalseValue, nil
}

func builtinIsMap(ctx *CallContext) (Object, error) {
	if len(ctx.Args) != 1 {
		return nil, ErrWrongNumArguments
	}
	if _, ok := ctx.Args[0].(*Map); ok {
		return TrueValue, nil
	}
	return FalseValue, nil
}

func builtinIsImmutableMap(ctx *CallContext) (Object, error) {
	if len(ctx.Args) != 1 {
		return nil, ErrWrongNumArguments
	}
	if _, ok := ctx.Args[0].(*ImmutableMap); ok {
		return TrueValue, nil
	}
	return FalseValue, nil
}

func builtinIsTime(ctx *CallContext) (Object, error) {
	if len(ctx.Args) != 1 {
		return nil, ErrWrongNumArguments
	}
	if _, ok := ctx.Args[0].(*Time); ok {
		return TrueValue, nil
	}
	return FalseValue, nil
}

func builtinIsError(ctx *CallContext) (Object, error) {
	if len(ctx.Args) != 1 {
		return nil, ErrWrongNumArguments
	}
	if _, ok := ctx.Args[0].(*Error); ok {
		return TrueValue, nil
	}
	return FalseValue, nil
}

func builtinIsUndefined(ctx *CallContext) (Object, error) {
	if len(ctx.Args) != 1 {
		return nil, ErrWrongNumArguments
	}
	if ctx.Args[0] == UndefinedValue {
		return TrueValue, nil
	}
	return FalseValue, nil
}

func builtinIsFunction(ctx *CallContext) (Object, error) {
	if len(ctx.Args) != 1 {
		return nil, ErrWrongNumArguments
	}
	switch ctx.Args[0].(type) {
	case *CompiledFunction:
		return TrueValue, nil
	}
	return FalseValue, nil
}

func builtinIsCallable(ctx *CallContext) (Object, error) {
	if len(ctx.Args) != 1 {
		return nil, ErrWrongNumArguments
	}
	if ctx.Args[0].CanCall() {
		return TrueValue, nil
	}
	return FalseValue, nil
}

func builtinIsIterable(ctx *CallContext) (Object, error) {
	if len(ctx.Args) != 1 {
		return nil, ErrWrongNumArguments
	}
	if ctx.Args[0].CanIterate() {
		return TrueValue, nil
	}
	return FalseValue, nil
}

// len(obj object) => int
func builtinLen(ctx *CallContext) (Object, error) {
	if len(ctx.Args) != 1 {
		return nil, ErrWrongNumArguments
	}
	switch arg := ctx.Args[0].(type) {
	case *Array:
		return &Int{Value: int64(len(arg.Value))}, nil
	case *ImmutableArray:
		return &Int{Value: int64(len(arg.Value))}, nil
	case *String:
		return &Int{Value: int64(len(arg.Value))}, nil
	case *Bytes:
		return &Int{Value: int64(len(arg.Value))}, nil
	case *Map:
		return &Int{Value: int64(len(arg.Value))}, nil
	case *ImmutableMap:
		return &Int{Value: int64(len(arg.Value))}, nil
	default:
		return nil, ErrInvalidArgumentType{
			Name:     "first",
			Expected: "array/string/bytes/map",
			Found:    arg.TypeName(),
		}
	}
}

// range(start, stop[, step])
func builtinRange(ctx *CallContext) (Object, error) {
	numArgs := len(ctx.Args)
	if numArgs < 2 || numArgs > 3 {
		return nil, ErrWrongNumArguments
	}
	var start, stop, step *Int

	for i, arg := range ctx.Args {
		v, ok := ctx.Args[i].(*Int)
		if !ok {
			var name string
			switch i {
			case 0:
				name = "start"
			case 1:
				name = "stop"
			case 2:
				name = "step"
			}

			return nil, ErrInvalidArgumentType{
				Name:     name,
				Expected: "int",
				Found:    arg.TypeName(),
			}
		}
		if i == 2 && v.Value <= 0 {
			return nil, ErrInvalidRangeStep
		}
		switch i {
		case 0:
			start = v
		case 1:
			stop = v
		case 2:
			step = v
		}
	}

	if step == nil {
		step = &Int{Value: int64(1)}
	}

	return buildRange(start.Value, stop.Value, step.Value), nil
}

func buildRange(start, stop, step int64) *Array {
	array := &Array{}
	if start <= stop {
		for i := start; i < stop; i += step {
			array.Value = append(array.Value, &Int{
				Value: i,
			})
		}
	} else {
		for i := start; i > stop; i -= step {
			array.Value = append(array.Value, &Int{
				Value: i,
			})
		}
	}
	return array
}

func builtinFormat(ctx *CallContext) (Object, error) {
	numArgs := len(ctx.Args)
	if numArgs == 0 {
		return nil, ErrWrongNumArguments
	}
	format, ok := ctx.Args[0].(*String)
	if !ok {
		return nil, ErrInvalidArgumentType{
			Name:     "format",
			Expected: "string",
			Found:    ctx.Args[0].TypeName(),
		}
	}
	if numArgs == 1 {
		// okay to return 'format' directly as String is immutable
		return format, nil
	}
	s, err := Format(format.Value, ctx.Args[1:]...)
	if err != nil {
		return nil, err
	}
	return &String{Value: s}, nil
}

func builtinCopy(ctx *CallContext) (Object, error) {
	if len(ctx.Args) != 1 {
		return nil, ErrWrongNumArguments
	}
	return ctx.Args[0].Copy(), nil
}

func builtinString(ctx *CallContext) (Object, error) {
	argsLen := len(ctx.Args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, ErrWrongNumArguments
	}
	if _, ok := ctx.Args[0].(*String); ok {
		return ctx.Args[0], nil
	}
	v, ok := ToString(ctx.Args[0])
	if ok {
		if len(v) > MaxStringLen {
			return nil, ErrStringLimit
		}
		return &String{Value: v}, nil
	}
	if argsLen == 2 {
		return ctx.Args[1], nil
	}
	return UndefinedValue, nil
}

func builtinInt(ctx *CallContext) (Object, error) {
	argsLen := len(ctx.Args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, ErrWrongNumArguments
	}
	if _, ok := ctx.Args[0].(*Int); ok {
		return ctx.Args[0], nil
	}
	v, ok := ToInt64(ctx.Args[0])
	if ok {
		return &Int{Value: v}, nil
	}
	if argsLen == 2 {
		return ctx.Args[1], nil
	}
	return UndefinedValue, nil
}

func builtinFloat(ctx *CallContext) (Object, error) {
	argsLen := len(ctx.Args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, ErrWrongNumArguments
	}
	if _, ok := ctx.Args[0].(*Float); ok {
		return ctx.Args[0], nil
	}
	v, ok := ToFloat64(ctx.Args[0])
	if ok {
		return &Float{Value: v}, nil
	}
	if argsLen == 2 {
		return ctx.Args[1], nil
	}
	return UndefinedValue, nil
}

func builtinBool(ctx *CallContext) (Object, error) {
	if len(ctx.Args) != 1 {
		return nil, ErrWrongNumArguments
	}
	if _, ok := ctx.Args[0].(*Bool); ok {
		return ctx.Args[0], nil
	}
	v, ok := ToBool(ctx.Args[0])
	if ok {
		if v {
			return TrueValue, nil
		}
		return FalseValue, nil
	}
	return UndefinedValue, nil
}

func builtinChar(ctx *CallContext) (Object, error) {
	argsLen := len(ctx.Args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, ErrWrongNumArguments
	}
	if _, ok := ctx.Args[0].(*Char); ok {
		return ctx.Args[0], nil
	}
	v, ok := ToRune(ctx.Args[0])
	if ok {
		return &Char{Value: v}, nil
	}
	if argsLen == 2 {
		return ctx.Args[1], nil
	}
	return UndefinedValue, nil
}

func builtinBytes(ctx *CallContext) (Object, error) {
	argsLen := len(ctx.Args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, ErrWrongNumArguments
	}

	// bytes(N) => create a new bytes with given size N
	if n, ok := ctx.Args[0].(*Int); ok {
		if n.Value > int64(MaxBytesLen) {
			return nil, ErrBytesLimit
		}
		return &Bytes{Value: make([]byte, int(n.Value))}, nil
	}
	v, ok := ToByteSlice(ctx.Args[0])
	if ok {
		if len(v) > MaxBytesLen {
			return nil, ErrBytesLimit
		}
		return &Bytes{Value: v}, nil
	}
	if argsLen == 2 {
		return ctx.Args[1], nil
	}
	return UndefinedValue, nil
}

func builtinTime(ctx *CallContext) (Object, error) {
	argsLen := len(ctx.Args)
	if !(argsLen == 1 || argsLen == 2) {
		return nil, ErrWrongNumArguments
	}
	if _, ok := ctx.Args[0].(*Time); ok {
		return ctx.Args[0], nil
	}
	v, ok := ToTime(ctx.Args[0])
	if ok {
		return &Time{Value: v}, nil
	}
	if argsLen == 2 {
		return ctx.Args[1], nil
	}
	return UndefinedValue, nil
}

// append(arr, items...)
func builtinAppend(ctx *CallContext) (Object, error) {
	if len(ctx.Args) < 2 {
		return nil, ErrWrongNumArguments
	}
	switch arg := ctx.Args[0].(type) {
	case *Array:
		return &Array{Value: append(arg.Value, ctx.Args[1:]...)}, nil
	case *ImmutableArray:
		return &Array{Value: append(arg.Value, ctx.Args[1:]...)}, nil
	default:
		return nil, ErrInvalidArgumentType{
			Name:     "first",
			Expected: "array",
			Found:    arg.TypeName(),
		}
	}
}

// builtinDelete deletes Object keys
// usage: delete(map, "key") or delete(map, "key1", "key2", "keyN")
// key must be a string
func builtinDelete(ctx *CallContext) (Object, error) {
	argsLen := len(ctx.Args)
	if argsLen < 2 {
		return nil, ErrWrongNumArguments
	}
	if err := ctx.Args[0].IndexDel(ctx.VM, ctx.Args[1:]...); err == nil {
		return UndefinedValue, nil
	} else {
		return nil, err
	}
}

// builtinSplice deletes and changes given Array, returns deleted items.
// usage:
// deleted_items := splice(array[,start[,delete_count[,item1[,item2[,...]]]])
func builtinSplice(ctx *CallContext) (Object, error) {
	argsLen := len(ctx.Args)
	if argsLen == 0 {
		return nil, ErrWrongNumArguments
	}

	array, ok := ctx.Args[0].(*Array)
	if !ok {
		return nil, ErrInvalidArgumentType{
			Name:     "first",
			Expected: "array",
			Found:    ctx.Args[0].TypeName(),
		}
	}
	arrayLen := len(array.Value)

	var startIdx int
	if argsLen > 1 {
		arg1, ok := ctx.Args[1].(*Int)
		if !ok {
			return nil, ErrInvalidArgumentType{
				Name:     "second",
				Expected: "int",
				Found:    ctx.Args[1].TypeName(),
			}
		}
		startIdx = int(arg1.Value)
		if startIdx < 0 || startIdx > arrayLen {
			return nil, ErrIndexOutOfBounds
		}
	}

	delCount := len(array.Value)
	if argsLen > 2 {
		arg2, ok := ctx.Args[2].(*Int)
		if !ok {
			return nil, ErrInvalidArgumentType{
				Name:     "third",
				Expected: "int",
				Found:    ctx.Args[2].TypeName(),
			}
		}
		delCount = int(arg2.Value)
		if delCount < 0 {
			return nil, ErrIndexOutOfBounds
		}
	}
	// if count of to be deleted items is bigger than expected, truncate it
	if startIdx+delCount > arrayLen {
		delCount = arrayLen - startIdx
	}
	// delete items
	endIdx := startIdx + delCount
	deleted := append([]Object{}, array.Value[startIdx:endIdx]...)

	head := array.Value[:startIdx]
	var items []Object
	if argsLen > 3 {
		items = make([]Object, 0, argsLen-3)
		for i := 3; i < argsLen; i++ {
			items = append(items, ctx.Args[i])
		}
	}
	items = append(items, array.Value[endIdx:]...)
	array.Value = append(head, items...)

	// return deleted items
	return &Array{Value: deleted}, nil
}

// builtinMap make new map merging args of map and kwargs
// Usage: map([map...]...[,key=value,keyN=value])
// Examples:
// map(a=1,b=2) => {"a":1,"b":2}
// map({"a":1},{"b":2};c=3) => {"a":1,"b":2,"c":3}
func builtinMap(ctx *CallContext) (Object, error) {
	res := ctx.Kwargs
	if res == nil {
		res = make(map[string]Object)
	}

	for i, v := range ctx.Args {
		switch v {
		case nil:
		case UndefinedValue:
		default:
			switch t := v.(type) {
			case *Map:
				for key, value := range t.Value {
					res[key] = value
				}
			default:
				return nil, ErrInvalidArgumentType{
					Name:     "arg #" + strconv.Itoa(i),
					Expected: "map",
					Found:    t.TypeName(),
				}
			}
		}
	}
	return &Map{Value: res}, nil
}

// get_methods
func builtinGetMethods(ctx *CallContext) (Object, error) {
	if !ctx.This.CanIterate() {
		return nil, fmt.Errorf("not iterable: %s", ctx.This.TypeName())
	}
	var (
		items []Object
		item  Object
	)
	it := ctx.This.Iterate()
	for it.Next() {
		if item = it.Value(); item.Method() {
			items = append(items, item)
		}
	}
	return &Array{Value: items}, nil
}

// builtinTypeOf return type of object instance
func builtinTypeOf(ctx *CallContext) (_ Object, err error) {
	if len(ctx.Args) != 1 {
		return nil, ErrWrongNumArguments
	}
	switch t := ctx.Args[0].(type) {
	case ObjectInstancer:
		return t.InstanceType(), nil
	default:
		return nil, ErrInvalidArgumentType{
			Name:     "arg #0",
			Expected: "instance",
			Found:    t.TypeName(),
		}
	}
}

// builtinType define new type
func builtinType(ctx *CallContext) (_ Object, err error) {
	if len(ctx.Args) < 1 {
		return nil, ErrWrongNumArguments
	}

	var (
		name   string
		args   = ctx.Args
		kwargs = ctx.Kwargs
	)

start:
	switch n := args[0].(type) {
	case *String:
		name = n.Value
	case *Map:
		if len(args) != 1 {
			return nil, ErrWrongNumArguments
		}
		if len(kwargs) > 0 {
			return nil, ErrUnexpectedKwargs
		}
		kwargs = n.Value
		if nameValue, ok := kwargs["name"]; ok {
			args = []Object{nameValue}
		} else {
			return nil, fmt.Errorf("name of type is undefined")
		}
		goto start
	default:
		return nil, ErrInvalidArgumentType{
			Name:     "arg #0",
			Expected: "string|map",
			Found:    n.TypeName(),
		}
	}
	if name == "" {
		return nil, fmt.Errorf("name of type is empty")
	}

	t := &Type{
		Name:       name,
		Fields:     &TypeFields{Value: map[string]*TypeField{}},
		Methods:    &TypeMethods{Value: map[string]*TypeMethod{}},
		Properties: &TypeProperties{Value: map[string]*TypeProperty{}},
	}

	for i, arg := range args[1:] {
		switch at := arg.(type) {
		case *TypeFields:
			t.Fields = at
		case *TypeMethods:
			t.Methods = at
		case *TypeProperties:
			t.Properties = at
		default:
			if i == 0 && arg.Method() {
				t.New = arg
			}
		}
	}

	if fields, ok := kwargs["fields"]; ok {
		switch vt := fields.(type) {
		case *Map:
			if err = t.Fields.FromMap(vt); err != nil {
				return
			}
		case *TypeFields:
			t.Fields = vt
		default:
			return nil, fmt.Errorf("'fields' isn't map or <fields>")
		}
	}

	if val, ok := kwargs["methods"]; ok {
		switch vt := val.(type) {
		case *Map:
			if err = t.Methods.FromMap(vt); err != nil {
				return
			}
		case *TypeMethods:
			t.Methods = vt
		default:
			return nil, fmt.Errorf("'methods' isn't map or <methods>")
		}
	}

	if val, ok := kwargs["properties"]; ok {
		switch vt := val.(type) {
		case *Map:
			if err = t.Properties.FromMap(vt); err != nil {
				return
			}
		case *TypeProperties:
			t.Properties = vt
		default:
			return nil, fmt.Errorf("'properties' isn't map or <properties>")
		}
	}

	if value, ok := kwargs["new"]; ok {
		if !value.CanCall() {
			return nil, fmt.Errorf("'new' isn't callable")
		}
		t.New = value
	}

	return t, nil
}

// builtinField create new field
// usage: field(value) or field(value, tag_name="tag_value")
// value is any value
func builtinField(ctx *CallContext) (Object, error) {
	if len(ctx.Args) != 1 {
		return nil, ErrWrongNumArguments
	}
	return &TypeField{Value: ctx.Args[0], Tags: ctx.Kwargs}, nil
}

// builtinFields create new fields
// usage: fields(f1=value,f2=value)
// value is any value
func builtinFields(ctx *CallContext) (Object, error) {
	var fields = make(map[string]*TypeField, len(ctx.Kwargs))
	for name, value := range ctx.Kwargs {
		switch t := value.(type) {
		case *TypeField:
			fields[name] = t
		case *Default:
			fields[name] = &TypeField{Value: UndefinedValue, Tags: make(map[string]Object, 0)}
		default:
			fields[name] = &TypeField{Value: value, Tags: make(map[string]Object, 0)}
		}
	}
	return &TypeFields{Value: fields}, nil
}

// builtinMethod create new method
// usage: method(caller) or method(caller, tag_name="tag_value")
// caller is any callable value
func builtinMethod(ctx *CallContext) (Object, error) {
	if len(ctx.Args) != 1 {
		return nil, ErrWrongNumArguments
	}
	if conv, ok := ctx.Args[0].(ToMethodConverter); ok {
		return &TypeMethod{Value: conv, Tags: ctx.Kwargs}, nil
	} else {
		return nil, fmt.Errorf("value %s isn't convertible to method", ctx.Args[0].TypeName())
	}
}

// builtinMethods create new methods
// usage: methods(f1=callable, f2=method(callable), f3=method(callable, tag_name="tag_value")
// callable is callable value
func builtinMethods(ctx *CallContext) (Object, error) {
	var methods = make(map[string]*TypeMethod, len(ctx.Kwargs))
	for name, value := range ctx.Kwargs {
		switch t := value.(type) {
		case *TypeMethod:
			methods[name] = t
		default:
			if !value.CanCall() {
				return nil, fmt.Errorf("method %q: value isn't callable", name)
			}
			if conv, ok := t.(ToMethodConverter); ok {
				methods[name] = &TypeMethod{Value: conv, Tags: make(map[string]Object, 0)}
			} else {
				return nil, fmt.Errorf("value %s isn't convertible to method", t.TypeName())
			}
		}
	}
	return &TypeMethods{Value: methods}, nil
}

// builtinProperty create new property
// usage: property(getter, setter) or property(getter, setter, tag_name="tag_value")
//        or method(get=getter, set=setter, tag_name="tag_value")
// caller is any callable value
func builtinProperty(ctx *CallContext) (Object, error) {
	var (
		p TypeProperty
	)
	switch len(ctx.Args) {
	case 0:
		p.Getter, p.Setter = ctx.Kwargs["get"], ctx.Kwargs["set"]
		delete(ctx.Kwargs, "get")
		delete(ctx.Kwargs, "set")
	case 2:
		p.Getter, p.Setter = ctx.Args[0], ctx.Args[1]
	default:
		return nil, ErrWrongNumArguments
	}
	p.Tags = ctx.Kwargs
	return &p, nil
}

// builtinProperties create new properties
// usage: properties(f1={get:getter}, f2=property(getter, setter))
func builtinProperties(ctx *CallContext) (_ Object, err error) {
	var properties = make(map[string]*TypeProperty, len(ctx.Kwargs))
	for name, value := range ctx.Kwargs {
		switch t := value.(type) {
		case *TypeProperty:
			properties[name] = t
		case *Map:
			var p TypeProperty
			if err = p.FromMap(t); err != nil {
				return
			}
			properties[name] = &p
		default:
			return nil, fmt.Errorf("bad property %q value type", name)
		}
	}
	return &TypeProperties{Value: properties}, nil
}
