package tengo

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/d5/tengo/v2/parser"
	"github.com/d5/tengo/v2/token"
)

var (
	// TrueValue represents a true value.
	TrueValue Object = &Bool{value: true}

	// FalseValue represents a false value.
	FalseValue Object = &Bool{value: false}

	// UndefinedValue represents an undefined value.
	UndefinedValue Object = &Undefined{}

	// DefaultValue represents an undefined kwarg value.
	DefaultValue Object = &Default{}
)

type (
	// ToMapConverter converts object to map.
	ToMapConverter interface {
		Object
		ToMap(deep bool) *Map
	}

	// FromMap reset object values from map.
	FromMapLoader interface {
		Object
		FromMap(m *Map) error
	}

	// ToArrayConverter converts object to array
	ToArrayConverter interface {
		Object
		ToArray(deep bool) *Array
	}

	// ToMethodConverter converts callabled object to method of instance
	ToMethodConverter interface {
		Object
		// ToMethod converts this caller to method of instance
		ToMethodOf(of Object) ToMethodConverter
	}
)

// Object represents an object in the VM.
type Object interface {
	// TypeName should return the name of the type.
	TypeName() string

	// String should return a string representation of the type's value.
	String() string

	// BinaryOp should return another object that is the result of a given
	// binary operator and a right-hand side object. If BinaryOp returns an
	// error, the VM will treat it as a run-time error.
	BinaryOp(op token.Token, rhs Object) (Object, error)

	// IsFalsy should return true if the value of the type should be considered
	// as falsy.
	IsFalsy() bool

	// Equals should return true if the value of the type should be considered
	// as equal to the value of another object.
	Equals(another Object) bool

	// Copy should return a copy of the type (and its value). Copy function
	// will be used for copy() builtin function which is expected to deep-copy
	// the values generally.
	Copy() Object

	// IndexGet should take an index Object and return a result Object or an
	// error for indexable objects. Indexable is an object that can take an
	// index and return an object. If error is returned, the runtime will treat
	// it as a run-time error and ignore returned value. If Object is not
	// indexable, ErrNotIndexable should be returned as error. If nil is
	// returned as value, it will be converted to UndefinedToken value by the
	// runtime.
	IndexGet(_ *VM, index Object) (value Object, err error)

	// IndexSet should take an index Object and a value Object for index
	// assignable objects. Index assignable is an object that can take an index
	// and a value on the left-hand side of the assignment statement. If Object
	// is not index assignable, ErrNotIndexAssignable should be returned as
	// error. If an error is returned, it will be treated as a run-time error.
	IndexSet(_ *VM, index, value Object) error

	// IndexDel should delete an index Object and a value Object for index
	// assignable objects. Index deletable is an object that can delete an index.
	// If Object is not index deletable, ErrNotIndexDeletable should be returned as
	// error. If an error is returned, it will be treated as a run-time error.
	IndexDel(_ *VM, _ ...Object) error

	// Iterate should return an Iterator for the type.
	Iterate() Iterator

	// CanIterate should return whether the Object can be Iterated.
	CanIterate() bool

	// Call should take an arbitrary number of arguments and returns a return
	// value and/or an error, which the VM will consider as a run-time error.
	Call(ctx *CallContext) (ret Object, err error)

	// CanCall should return whether the Object can be Called.
	CanCall() bool

	// Method should return if caller is method
	Method() bool

	// MethodTo should return instance of this method
	MethodTo() Object
}

// ObjectImpl represents a default Object Implementation. To defined a new
// value type, one can embed ObjectImpl in their type declarations to avoid
// implementing all non-significant methods. TypeName() and String() methods
// still need to be implemented.
type ObjectImpl struct {
}

// TypeName returns the name of the type.
func (o *ObjectImpl) TypeName() string {
	panic(ErrNotImplemented)
}

func (o *ObjectImpl) String() string {
	panic(ErrNotImplemented)
}

// BinaryOp returns another object that is the result of a given binary
// operator and a right-hand side object.
func (o *ObjectImpl) BinaryOp(_ token.Token, _ Object) (Object, error) {
	return nil, ErrInvalidOperator
}

// Copy returns a copy of the type.
func (o *ObjectImpl) Copy() Object {
	return nil
}

// IsFalsy returns true if the value of the type is falsy.
func (o *ObjectImpl) IsFalsy() bool {
	return false
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *ObjectImpl) Equals(x Object) bool {
	return o == x
}

// IndexGet returns an element at a given index.
func (o *ObjectImpl) IndexGet(_ *VM, _ Object) (res Object, err error) {
	return nil, ErrNotIndexable
}

// IndexSet sets an element at a given index.
func (o *ObjectImpl) IndexSet(_ *VM, _, _ Object) (err error) {
	return ErrNotIndexAssignable
}

// IndexDel deletes an element at a given index.
func (o *ObjectImpl) IndexDel(_ *VM, _ ...Object) (err error) {
	return ErrNotIndexDeletable
}

// Iterate returns an iterator.
func (o *ObjectImpl) Iterate() Iterator {
	return nil
}

// CanIterate returns whether the Object can be Iterated.
func (o *ObjectImpl) CanIterate() bool {
	return false
}

// Call takes an arbitrary number of arguments and returns a return value
// and/or an error.
func (o *ObjectImpl) Call(*CallContext) (ret Object, err error) {
	return nil, nil
}

// CanCall returns whether the Object can be Called.
func (o *ObjectImpl) CanCall() bool {
	return false
}

// Method shold return if caller is method
func (o *ObjectImpl) Method() bool {
	return false
}

// MethodTo should return instance of this method
func (*ObjectImpl) MethodTo() Object {
	return nil
}

// Array represents an array of objects.
type Array struct {
	ObjectImpl
	Value []Object
}

// TypeName returns the name of the type.
func (o *Array) TypeName() string {
	return "array"
}

func (o *Array) String() string {
	var elements []string
	for _, e := range o.Value {
		elements = append(elements, e.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

// BinaryOp returns another object that is the result of a given binary
// operator and a right-hand side object.
func (o *Array) BinaryOp(op token.Token, rhs Object) (Object, error) {
	if rhs, ok := rhs.(*Array); ok {
		switch op {
		case token.Add:
			if len(rhs.Value) == 0 {
				return o, nil
			}
			return &Array{Value: append(o.Value, rhs.Value...)}, nil
		}
	}
	return nil, ErrInvalidOperator
}

// Copy returns a copy of the type.
func (o *Array) Copy() Object {
	var c []Object
	for _, elem := range o.Value {
		c = append(c, elem.Copy())
	}
	return &Array{Value: c}
}

// IsFalsy returns true if the value of the type is falsy.
func (o *Array) IsFalsy() bool {
	return len(o.Value) == 0
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *Array) Equals(x Object) bool {
	var xVal []Object
	switch x := x.(type) {
	case *Array:
		xVal = x.Value
	case *ImmutableArray:
		xVal = x.Value
	default:
		return false
	}
	if len(o.Value) != len(xVal) {
		return false
	}
	for i, e := range o.Value {
		if !e.Equals(xVal[i]) {
			return false
		}
	}
	return true
}

// IndexGet returns an element at a given index.
func (o *Array) IndexGet(_ *VM, index Object) (res Object, err error) {
	intIdx, ok := index.(*Int)
	if !ok {
		err = ErrInvalidIndexType
		return
	}
	idxVal := int(intIdx.Value)
	if idxVal < 0 || idxVal >= len(o.Value) {
		res = UndefinedValue
		return
	}
	res = o.Value[idxVal]
	return
}

// IndexSet sets an element at a given index.
func (o *Array) IndexSet(_ *VM, index, value Object) (err error) {
	intIdx, ok := ToInt(index)
	if !ok {
		err = ErrInvalidIndexType
		return
	}
	if intIdx < 0 || intIdx >= len(o.Value) {
		err = ErrIndexOutOfBounds
		return
	}
	o.Value[intIdx] = value
	return nil
}

// Iterate creates an array iterator.
func (o *Array) Iterate() Iterator {
	return &ArrayIterator{
		v: o.Value,
		l: len(o.Value),
	}
}

// CanIterate returns whether the Object can be Iterated.
func (o *Array) CanIterate() bool {
	return true
}

// ToArrayConverter converts object to array. If not deep, returns this object,
// other else, return new array with deeped values
func (o *Array) ToArray(deep bool) *Array {
	if deep {
		ret := make([]Object, len(o.Value))
		for i, v := range o.Value {
			switch t := v.(type) {
			case ToArrayConverter:
				ret[i] = t.ToArray(true)
			case ToMapConverter:
				ret[i] = t.ToMap(true)
			default:
				ret[i] = t
			}
		}
		return &Array{Value: ret}
	}
	return o
}

// Bool represents a boolean value.
type Bool struct {
	ObjectImpl

	// this is intentionally non-public to force using objects.TrueValue and
	// FalseValue always
	value bool
}

func (o *Bool) String() string {
	if o.value {
		return "true"
	}

	return "false"
}

// TypeName returns the name of the type.
func (o *Bool) TypeName() string {
	return "bool"
}

// Copy returns a copy of the type.
func (o *Bool) Copy() Object {
	return o
}

// IsFalsy returns true if the value of the type is falsy.
func (o *Bool) IsFalsy() bool {
	return !o.value
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *Bool) Equals(x Object) bool {
	return o == x
}

// GobDecode decodes bool value from input bytes.
func (o *Bool) GobDecode(b []byte) (err error) {
	o.value = b[0] == 1
	return
}

// GobEncode encodes bool values into bytes.
func (o *Bool) GobEncode() (b []byte, err error) {
	if o.value {
		b = []byte{1}
	} else {
		b = []byte{0}
	}
	return
}

// BuiltinFunction represents a builtin function.
type BuiltinFunction struct {
	ObjectImpl
	Name         string
	Value        CallableFuncCtx
	IsMethod     bool
	methodTarget Object
}

// TypeName returns the name of the type.
func (o *BuiltinFunction) TypeName() string {
	return "builtin-function:" + o.Name
}

func (o *BuiltinFunction) String() string {
	return "<builtin-function>"
}

// Copy returns a copy of the type.
func (o BuiltinFunction) Copy() Object {
	return &o
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *BuiltinFunction) Equals(_ Object) bool {
	return false
}

// Call executes a builtin function.
func (o *BuiltinFunction) Call(ctx *CallContext) (Object, error) {
	return o.Value(ctx)
}

// CanCall returns whether the Object can be Called.
func (o *BuiltinFunction) CanCall() bool {
	return true
}

// Method shold return if caller is method
func (o *BuiltinFunction) Method() bool {
	return o.IsMethod
}

// Target should return current method target
func (o *BuiltinFunction) MethodTo() Object {
	return o.methodTarget
}

// Target should set current method target
func (o BuiltinFunction) ToMethodOf(of Object) ToMethodConverter {
	o.methodTarget = of
	return &o
}

// BuiltinModule is an importable module that's written in Go.
type BuiltinModule struct {
	Attrs map[string]Object
}

// Import returns an immutable map for the module.
func (m *BuiltinModule) Import(moduleName string) (interface{}, error) {
	return m.AsImmutableMap(moduleName), nil
}

// AsImmutableMap converts builtin module into an immutable map.
func (m *BuiltinModule) AsImmutableMap(moduleName string) *ImmutableMap {
	attrs := make(map[string]Object, len(m.Attrs))
	for k, v := range m.Attrs {
		attrs[k] = v.Copy()
	}
	attrs["__module_name__"] = &String{Value: moduleName}
	return &ImmutableMap{Value: attrs}
}

// Bytes represents a byte array.
type Bytes struct {
	ObjectImpl
	Value []byte
}

func (o *Bytes) String() string {
	return string(o.Value)
}

// TypeName returns the name of the type.
func (o *Bytes) TypeName() string {
	return "bytes"
}

// BinaryOp returns another object that is the result of a given binary
// operator and a right-hand side object.
func (o *Bytes) BinaryOp(op token.Token, rhs Object) (Object, error) {
	switch op {
	case token.Add:
		switch rhs := rhs.(type) {
		case *Bytes:
			if len(o.Value)+len(rhs.Value) > MaxBytesLen {
				return nil, ErrBytesLimit
			}
			return &Bytes{Value: append(o.Value, rhs.Value...)}, nil
		}
	}
	return nil, ErrInvalidOperator
}

// Copy returns a copy of the type.
func (o *Bytes) Copy() Object {
	return &Bytes{Value: append([]byte{}, o.Value...)}
}

// IsFalsy returns true if the value of the type is falsy.
func (o *Bytes) IsFalsy() bool {
	return len(o.Value) == 0
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *Bytes) Equals(x Object) bool {
	t, ok := x.(*Bytes)
	if !ok {
		return false
	}
	return bytes.Equal(o.Value, t.Value)
}

// IndexGet returns an element (as Int) at a given index.
func (o *Bytes) IndexGet(_ *VM, index Object) (res Object, err error) {
	intIdx, ok := index.(*Int)
	if !ok {
		err = ErrInvalidIndexType
		return
	}
	idxVal := int(intIdx.Value)
	if idxVal < 0 || idxVal >= len(o.Value) {
		res = UndefinedValue
		return
	}
	res = &Int{Value: int64(o.Value[idxVal])}
	return
}

// Iterate creates a bytes iterator.
func (o *Bytes) Iterate() Iterator {
	return &BytesIterator{
		v: o.Value,
		l: len(o.Value),
	}
}

// CanIterate returns whether the Object can be Iterated.
func (o *Bytes) CanIterate() bool {
	return true
}

// Char represents a character value.
type Char struct {
	ObjectImpl
	Value rune
}

func (o *Char) String() string {
	return string(o.Value)
}

// TypeName returns the name of the type.
func (o *Char) TypeName() string {
	return "char"
}

// BinaryOp returns another object that is the result of a given binary
// operator and a right-hand side object.
func (o *Char) BinaryOp(op token.Token, rhs Object) (Object, error) {
	switch rhs := rhs.(type) {
	case *Char:
		switch op {
		case token.Add:
			r := o.Value + rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Char{Value: r}, nil
		case token.Sub:
			r := o.Value - rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Char{Value: r}, nil
		case token.Less:
			if o.Value < rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.Greater:
			if o.Value > rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.LessEq:
			if o.Value <= rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.GreaterEq:
			if o.Value >= rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		}
	case *Int:
		switch op {
		case token.Add:
			r := o.Value + rune(rhs.Value)
			if r == o.Value {
				return o, nil
			}
			return &Char{Value: r}, nil
		case token.Sub:
			r := o.Value - rune(rhs.Value)
			if r == o.Value {
				return o, nil
			}
			return &Char{Value: r}, nil
		case token.Less:
			if int64(o.Value) < rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.Greater:
			if int64(o.Value) > rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.LessEq:
			if int64(o.Value) <= rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.GreaterEq:
			if int64(o.Value) >= rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		}
	}
	return nil, ErrInvalidOperator
}

// Copy returns a copy of the type.
func (o *Char) Copy() Object {
	return &Char{Value: o.Value}
}

// IsFalsy returns true if the value of the type is falsy.
func (o *Char) IsFalsy() bool {
	return o.Value == 0
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *Char) Equals(x Object) bool {
	t, ok := x.(*Char)
	if !ok {
		return false
	}
	return o.Value == t.Value
}

type VarArgMode uint8

func (m VarArgMode) Pos() int {
	if m == VarArgNamed {
		return 1
	}
	return 0
}

const (
	VarArgNone = iota
	VarArgNamed
	VarArgAnonymous
)

// Variadic represents a variadic expression.
type Variadic struct {
	Valid bool
	Name  string
}

// CompiledFunction represents a compiled function.
type CompiledFunction struct {
	ObjectImpl
	Instructions     []byte
	NumLocals        int // number of local variables (including function parameters)
	NumArgs          int
	VarArgs          Variadic
	Kwargs           map[string]int
	KwargsNames      []string
	KwargsDefaults   []Object
	VarKwargs        Variadic
	SourceMap        map[int]parser.Pos
	IsMethod         bool // receive `this` as first arg
	methodTarget     Object
	Free             []*ObjectPtr
	vm               *VM
	vmConstantsCount int
	running          uint32
}

// TypeName returns the name of the type.
func (o *CompiledFunction) TypeName() string {
	return "compiled-function"
}

func (o *CompiledFunction) String() string {
	return "<compiled-function>"
}

func (o *CompiledFunction) Method() bool {
	return o.IsMethod
}

func (o *CompiledFunction) MethodTo() Object {
	return o.methodTarget
}

// Copy returns a copy of the type.
func (o *CompiledFunction) Copy() Object {
	return &CompiledFunction{
		Instructions:   append([]byte{}, o.Instructions...),
		NumLocals:      o.NumLocals,
		NumArgs:        o.NumArgs,
		VarArgs:        o.VarArgs,
		Kwargs:         o.Kwargs,
		KwargsNames:    o.KwargsNames,
		KwargsDefaults: o.KwargsDefaults,
		VarKwargs:      o.VarKwargs,
		IsMethod:       o.IsMethod,
		methodTarget:   o.methodTarget,
		Free:           append([]*ObjectPtr{}, o.Free...), // DO NOT Copy() of elements; these are variable pointers
	}
}

// ToMethod converts this caller to method of instance
func (o CompiledFunction) ToMethodOf(of Object) ToMethodConverter {
	o.methodTarget = of
	return &o
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *CompiledFunction) Equals(_ Object) bool {
	return false
}

// SourcePos returns the source position of the instruction at ip.
func (o *CompiledFunction) SourcePos(ip int) parser.Pos {
	for ip >= 0 {
		if p, ok := o.SourceMap[ip]; ok {
			return p
		}
		ip--
	}
	return parser.NoPos
}

// CanCall returns whether the Object can be Called.
func (o *CompiledFunction) CanCall() bool {
	return true
}

func (o *CompiledFunction) SetupVM(vm *VM) {
	constsOffset := len(vm.bc.Constants)

	// Load fn
	inst := MakeInstruction(parser.OpConstant, constsOffset)

	// Load args
	inst = append(inst,
		MakeInstruction(parser.OpConstant, constsOffset+1)...)
	// Load kwargs
	inst = append(inst,
		MakeInstruction(parser.OpConstant, constsOffset+2)...)

	// Call, set value to a global, stop
	inst = append(inst, MakeInstruction(parser.OpCall, 0, 1, 0, 1)...)
	inst = append(inst, MakeInstruction(parser.OpSuspend)...)

	bc := *vm.bc
	bc.Constants = append(bc.Constants, o, &Array{}, &Map{})

	bc.MainFunction = &CompiledFunction{
		Instructions: inst,
	}

	o.vm = NewVM(&bc, vm.globals, vm.maxAllocs)
	o.vm.context = vm.context
	o.vmConstantsCount = constsOffset + 3
}

// Call should take an arbitrary number of arguments and returns a return
// value and/or an error, which the VM will consider as a run-time error.
func (o *CompiledFunction) Call(ctx *CallContext) (ret Object, err error) {
	if atomic.LoadUint32(&o.running) == 1 {
		err = fmt.Errorf("compile function is running")
		return
	}

	atomic.StoreUint32(&o.running, 1)
	defer atomic.StoreUint32(&o.running, 0)

	if o.vm == nil {
		o.SetupVM(ctx.VM)
	}

	o.vm.bc.Constants[o.vmConstantsCount-2].(*Array).Value = ctx.Args
	o.vm.bc.Constants[o.vmConstantsCount-1].(*Map).Value = ctx.Kwargs

	if err = o.vm.Run(); err == nil {
		ret = o.vm.stack[0]
	}
	return
}

// Error represents an error value.
type Error struct {
	ObjectImpl
	Value Object
}

// TypeName returns the name of the type.
func (o *Error) TypeName() string {
	return "error"
}

func (o *Error) String() string {
	if o.Value != nil {
		return fmt.Sprintf("error: %s", o.Value.String())
	}
	return "error"
}

// IsFalsy returns true if the value of the type is falsy.
func (o *Error) IsFalsy() bool {
	return true // error is always false.
}

// Copy returns a copy of the type.
func (o *Error) Copy() Object {
	return &Error{Value: o.Value.Copy()}
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *Error) Equals(x Object) bool {
	return o == x // pointer equality
}

// IndexGet returns an element at a given index.
func (o *Error) IndexGet(_ *VM, index Object) (res Object, err error) {
	if strIdx, _ := ToString(index); strIdx != "value" {
		err = ErrInvalidIndexOnError
		return
	}
	res = o.Value
	return
}

// Float represents a floating point number value.
type Float struct {
	ObjectImpl
	Value float64
}

func (o *Float) String() string {
	return strconv.FormatFloat(o.Value, 'f', -1, 64)
}

// TypeName returns the name of the type.
func (o *Float) TypeName() string {
	return "float"
}

// BinaryOp returns another object that is the result of a given binary
// operator and a right-hand side object.
func (o *Float) BinaryOp(op token.Token, rhs Object) (Object, error) {
	switch rhs := rhs.(type) {
	case *Float:
		switch op {
		case token.Add:
			r := o.Value + rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Float{Value: r}, nil
		case token.Sub:
			r := o.Value - rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Float{Value: r}, nil
		case token.Mul:
			r := o.Value * rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Float{Value: r}, nil
		case token.Quo:
			r := o.Value / rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Float{Value: r}, nil
		case token.Less:
			if o.Value < rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.Greater:
			if o.Value > rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.LessEq:
			if o.Value <= rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.GreaterEq:
			if o.Value >= rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		}
	case *Int:
		switch op {
		case token.Add:
			r := o.Value + float64(rhs.Value)
			if r == o.Value {
				return o, nil
			}
			return &Float{Value: r}, nil
		case token.Sub:
			r := o.Value - float64(rhs.Value)
			if r == o.Value {
				return o, nil
			}
			return &Float{Value: r}, nil
		case token.Mul:
			r := o.Value * float64(rhs.Value)
			if r == o.Value {
				return o, nil
			}
			return &Float{Value: r}, nil
		case token.Quo:
			r := o.Value / float64(rhs.Value)
			if r == o.Value {
				return o, nil
			}
			return &Float{Value: r}, nil
		case token.Less:
			if o.Value < float64(rhs.Value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.Greater:
			if o.Value > float64(rhs.Value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.LessEq:
			if o.Value <= float64(rhs.Value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.GreaterEq:
			if o.Value >= float64(rhs.Value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		}
	}
	return nil, ErrInvalidOperator
}

// Copy returns a copy of the type.
func (o *Float) Copy() Object {
	return &Float{Value: o.Value}
}

// IsFalsy returns true if the value of the type is falsy.
func (o *Float) IsFalsy() bool {
	return math.IsNaN(o.Value)
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *Float) Equals(x Object) bool {
	t, ok := x.(*Float)
	if !ok {
		return false
	}
	return o.Value == t.Value
}

// ImmutableArray represents an immutable array of objects.
type ImmutableArray struct {
	ObjectImpl
	Value []Object
}

// TypeName returns the name of the type.
func (o *ImmutableArray) TypeName() string {
	return "immutable-array"
}

func (o *ImmutableArray) String() string {
	var elements []string
	for _, e := range o.Value {
		elements = append(elements, e.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

// BinaryOp returns another object that is the result of a given binary
// operator and a right-hand side object.
func (o *ImmutableArray) BinaryOp(op token.Token, rhs Object) (Object, error) {
	if rhs, ok := rhs.(*ImmutableArray); ok {
		switch op {
		case token.Add:
			return &Array{Value: append(o.Value, rhs.Value...)}, nil
		}
	}
	return nil, ErrInvalidOperator
}

// Copy returns a copy of the type.
func (o *ImmutableArray) Copy() Object {
	var c []Object
	for _, elem := range o.Value {
		c = append(c, elem.Copy())
	}
	return &Array{Value: c}
}

// IsFalsy returns true if the value of the type is falsy.
func (o *ImmutableArray) IsFalsy() bool {
	return len(o.Value) == 0
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *ImmutableArray) Equals(x Object) bool {
	var xVal []Object
	switch x := x.(type) {
	case *Array:
		xVal = x.Value
	case *ImmutableArray:
		xVal = x.Value
	default:
		return false
	}
	if len(o.Value) != len(xVal) {
		return false
	}
	for i, e := range o.Value {
		if !e.Equals(xVal[i]) {
			return false
		}
	}
	return true
}

// IndexGet returns an element at a given index.
func (o *ImmutableArray) IndexGet(_ *VM, index Object) (res Object, err error) {
	intIdx, ok := index.(*Int)
	if !ok {
		err = ErrInvalidIndexType
		return
	}
	idxVal := int(intIdx.Value)
	if idxVal < 0 || idxVal >= len(o.Value) {
		res = UndefinedValue
		return
	}
	res = o.Value[idxVal]
	return
}

// Iterate creates an array iterator.
func (o *ImmutableArray) Iterate() Iterator {
	return &ArrayIterator{
		v: o.Value,
		l: len(o.Value),
	}
}

// CanIterate returns whether the Object can be Iterated.
func (o *ImmutableArray) CanIterate() bool {
	return true
}

// ImmutableMap represents an immutable map object.
type ImmutableMap struct {
	ObjectImpl
	Value map[string]Object
}

// TypeName returns the name of the type.
func (o *ImmutableMap) TypeName() string {
	return "immutable-map"
}

func (o *ImmutableMap) String() string {
	var pairs []string
	for k, v := range o.Value {
		pairs = append(pairs, fmt.Sprintf("%s: %s", k, v.String()))
	}
	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}

// Copy returns a copy of the type.
func (o *ImmutableMap) Copy() Object {
	c := make(map[string]Object)
	for k, v := range o.Value {
		c[k] = v.Copy()
	}
	return &Map{Value: c}
}

// IsFalsy returns true if the value of the type is falsy.
func (o *ImmutableMap) IsFalsy() bool {
	return len(o.Value) == 0
}

// IndexGet returns the value for the given key.
func (o *ImmutableMap) IndexGet(_ *VM, index Object) (res Object, err error) {
	strIdx, ok := ToString(index)
	if !ok {
		err = ErrInvalidIndexType
		return
	}
	res, ok = o.Value[strIdx]
	if !ok {
		res = UndefinedValue
	}
	return
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *ImmutableMap) Equals(x Object) bool {
	var xVal map[string]Object
	switch x := x.(type) {
	case *Map:
		xVal = x.Value
	case *ImmutableMap:
		xVal = x.Value
	default:
		return false
	}
	if len(o.Value) != len(xVal) {
		return false
	}
	for k, v := range o.Value {
		tv := xVal[k]
		if !v.Equals(tv) {
			return false
		}
	}
	return true
}

// Iterate creates an immutable map iterator.
func (o *ImmutableMap) Iterate() Iterator {
	var keys []string
	for k := range o.Value {
		keys = append(keys, k)
	}
	return &MapIterator{
		v: o.Value,
		k: keys,
		l: len(keys),
	}
}

// CanIterate returns whether the Object can be Iterated.
func (o *ImmutableMap) CanIterate() bool {
	return true
}

// ToMap convert to map.
func (o *ImmutableMap) ToMap(deep bool) *Map {
	if deep {
		m := make(map[string]Object, len(o.Value))
		for k, v := range o.Value {
			if c, ok := v.(ToMapConverter); ok {
				m[k] = c.ToMap(true)
			} else {
				m[k] = v
			}
		}
		return &Map{Value: m}
	}
	return &Map{Value: o.Value}
}

// Int represents an integer value.
type Int struct {
	ObjectImpl
	Value int64
}

func (o *Int) String() string {
	return strconv.FormatInt(o.Value, 10)
}

// TypeName returns the name of the type.
func (o *Int) TypeName() string {
	return "int"
}

// BinaryOp returns another object that is the result of a given binary
// operator and a right-hand side object.
func (o *Int) BinaryOp(op token.Token, rhs Object) (Object, error) {
	switch rhs := rhs.(type) {
	case *Int:
		switch op {
		case token.Add:
			r := o.Value + rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Int{Value: r}, nil
		case token.Sub:
			r := o.Value - rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Int{Value: r}, nil
		case token.Mul:
			r := o.Value * rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Int{Value: r}, nil
		case token.Quo:
			r := o.Value / rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Int{Value: r}, nil
		case token.Rem:
			r := o.Value % rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Int{Value: r}, nil
		case token.And:
			r := o.Value & rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Int{Value: r}, nil
		case token.Or:
			r := o.Value | rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Int{Value: r}, nil
		case token.Xor:
			r := o.Value ^ rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Int{Value: r}, nil
		case token.AndNot:
			r := o.Value &^ rhs.Value
			if r == o.Value {
				return o, nil
			}
			return &Int{Value: r}, nil
		case token.Shl:
			r := o.Value << uint64(rhs.Value)
			if r == o.Value {
				return o, nil
			}
			return &Int{Value: r}, nil
		case token.Shr:
			r := o.Value >> uint64(rhs.Value)
			if r == o.Value {
				return o, nil
			}
			return &Int{Value: r}, nil
		case token.Less:
			if o.Value < rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.Greater:
			if o.Value > rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.LessEq:
			if o.Value <= rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.GreaterEq:
			if o.Value >= rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		}
	case *Float:
		switch op {
		case token.Add:
			return &Float{Value: float64(o.Value) + rhs.Value}, nil
		case token.Sub:
			return &Float{Value: float64(o.Value) - rhs.Value}, nil
		case token.Mul:
			return &Float{Value: float64(o.Value) * rhs.Value}, nil
		case token.Quo:
			return &Float{Value: float64(o.Value) / rhs.Value}, nil
		case token.Less:
			if float64(o.Value) < rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.Greater:
			if float64(o.Value) > rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.LessEq:
			if float64(o.Value) <= rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.GreaterEq:
			if float64(o.Value) >= rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		}
	case *Char:
		switch op {
		case token.Add:
			return &Char{Value: rune(o.Value) + rhs.Value}, nil
		case token.Sub:
			return &Char{Value: rune(o.Value) - rhs.Value}, nil
		case token.Less:
			if o.Value < int64(rhs.Value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.Greater:
			if o.Value > int64(rhs.Value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.LessEq:
			if o.Value <= int64(rhs.Value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.GreaterEq:
			if o.Value >= int64(rhs.Value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		}
	}
	return nil, ErrInvalidOperator
}

// Copy returns a copy of the type.
func (o *Int) Copy() Object {
	return &Int{Value: o.Value}
}

// IsFalsy returns true if the value of the type is falsy.
func (o *Int) IsFalsy() bool {
	return o.Value == 0
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *Int) Equals(x Object) bool {
	t, ok := x.(*Int)
	if !ok {
		return false
	}
	return o.Value == t.Value
}

// Map represents a map of objects.
type Map struct {
	ObjectImpl
	Value map[string]Object
}

// TypeName returns the name of the type.
func (o *Map) TypeName() string {
	return "map"
}

func (o *Map) String() string {
	var pairs []string
	for k, v := range o.Value {
		pairs = append(pairs, fmt.Sprintf("%s: %s", k, v.String()))
	}
	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}

// Copy returns a copy of the type.
func (o *Map) Copy() Object {
	c := make(map[string]Object)
	for k, v := range o.Value {
		c[k] = v.Copy()
	}
	return &Map{Value: c}
}

// IsFalsy returns true if the value of the type is falsy.
func (o *Map) IsFalsy() bool {
	return len(o.Value) == 0
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *Map) Equals(x Object) bool {
	var xVal map[string]Object
	switch x := x.(type) {
	case *Map:
		xVal = x.Value
	case *ImmutableMap:
		xVal = x.Value
	default:
		return false
	}
	if len(o.Value) != len(xVal) {
		return false
	}
	for k, v := range o.Value {
		tv := xVal[k]
		if !v.Equals(tv) {
			return false
		}
	}
	return true
}

// IndexGet returns the value for the given key.
func (o *Map) IndexGet(_ *VM, index Object) (res Object, err error) {
	strIdx, ok := ToString(index)
	if !ok {
		err = ErrInvalidIndexType
		return
	}
	res, ok = o.Value[strIdx]
	if !ok {
		res = UndefinedValue
	}
	return
}

// IndexSet sets the value for the given key.
func (o *Map) IndexSet(_ *VM, index, value Object) (err error) {
	strIdx, ok := ToString(index)
	if !ok {
		err = ErrInvalidIndexType
		return
	}
	o.Value[strIdx] = value
	return nil
}

// IndexDel deletes the value for the given key.
func (o *Map) IndexDel(_ *VM, index ...Object) (err error) {
	for _, index := range index {
		strIdx, ok := ToString(index)
		if !ok {
			err = ErrInvalidIndexType
			return
		}
		delete(o.Value, strIdx)
	}
	return nil
}

// Iterate creates a map iterator.
func (o *Map) Iterate() Iterator {
	var keys []string
	for k := range o.Value {
		keys = append(keys, k)
	}
	return &MapIterator{
		v: o.Value,
		k: keys,
		l: len(keys),
	}
}

// ToMap convert to map.
func (o *Map) ToMap(deep bool) *Map {
	if deep {
		m := make(map[string]Object, len(o.Value))
		for k, v := range o.Value {
			if c, ok := v.(ToMapConverter); ok {
				m[k] = c.ToMap(true)
			} else {
				m[k] = v
			}
		}
		return &Map{Value: m}
	}
	return o
}

// CanIterate returns whether the Object can be Iterated.
func (o *Map) CanIterate() bool {
	return true
}

// ObjectPtr represents a free variable.
type ObjectPtr struct {
	ObjectImpl
	Value *Object
}

func (o *ObjectPtr) String() string {
	return "free-var"
}

// TypeName returns the name of the type.
func (o *ObjectPtr) TypeName() string {
	return "<free-var>"
}

// Copy returns a copy of the type.
func (o *ObjectPtr) Copy() Object {
	return o
}

// IsFalsy returns true if the value of the type is falsy.
func (o *ObjectPtr) IsFalsy() bool {
	return o.Value == nil
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *ObjectPtr) Equals(x Object) bool {
	return o == x
}

// String represents a string value.
type String struct {
	ObjectImpl
	Value   string
	runeStr []rune
}

// TypeName returns the name of the type.
func (o *String) TypeName() string {
	return "string"
}

func (o *String) String() string {
	return strconv.Quote(o.Value)
}

// BinaryOp returns another object that is the result of a given binary
// operator and a right-hand side object.
func (o *String) BinaryOp(op token.Token, rhs Object) (Object, error) {
	switch op {
	case token.Add:
		switch rhs := rhs.(type) {
		case *String:
			if len(o.Value)+len(rhs.Value) > MaxStringLen {
				return nil, ErrStringLimit
			}
			return &String{Value: o.Value + rhs.Value}, nil
		default:
			rhsStr := rhs.String()
			if len(o.Value)+len(rhsStr) > MaxStringLen {
				return nil, ErrStringLimit
			}
			return &String{Value: o.Value + rhsStr}, nil
		}
	case token.Less:
		switch rhs := rhs.(type) {
		case *String:
			if o.Value < rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		}
	case token.LessEq:
		switch rhs := rhs.(type) {
		case *String:
			if o.Value <= rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		}
	case token.Greater:
		switch rhs := rhs.(type) {
		case *String:
			if o.Value > rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		}
	case token.GreaterEq:
		switch rhs := rhs.(type) {
		case *String:
			if o.Value >= rhs.Value {
				return TrueValue, nil
			}
			return FalseValue, nil
		}
	}
	return nil, ErrInvalidOperator
}

// IsFalsy returns true if the value of the type is falsy.
func (o *String) IsFalsy() bool {
	return len(o.Value) == 0
}

// Copy returns a copy of the type.
func (o *String) Copy() Object {
	return &String{Value: o.Value}
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *String) Equals(x Object) bool {
	t, ok := x.(*String)
	if !ok {
		return false
	}
	return o.Value == t.Value
}

// IndexGet returns a character at a given index.
func (o *String) IndexGet(_ *VM, index Object) (res Object, err error) {
	intIdx, ok := index.(*Int)
	if !ok {
		err = ErrInvalidIndexType
		return
	}
	idxVal := int(intIdx.Value)
	if o.runeStr == nil {
		o.runeStr = []rune(o.Value)
	}
	if idxVal < 0 || idxVal >= len(o.runeStr) {
		res = UndefinedValue
		return
	}
	res = &Char{Value: o.runeStr[idxVal]}
	return
}

// Iterate creates a string iterator.
func (o *String) Iterate() Iterator {
	if o.runeStr == nil {
		o.runeStr = []rune(o.Value)
	}
	return &StringIterator{
		v: o.runeStr,
		l: len(o.runeStr),
	}
}

// CanIterate returns whether the Object can be Iterated.
func (o *String) CanIterate() bool {
	return true
}

// Time represents a time value.
type Time struct {
	ObjectImpl
	Value time.Time
}

func (o *Time) String() string {
	return o.Value.String()
}

// TypeName returns the name of the type.
func (o *Time) TypeName() string {
	return "time"
}

// BinaryOp returns another object that is the result of a given binary
// operator and a right-hand side object.
func (o *Time) BinaryOp(op token.Token, rhs Object) (Object, error) {
	switch rhs := rhs.(type) {
	case *Int:
		switch op {
		case token.Add: // time + int => time
			if rhs.Value == 0 {
				return o, nil
			}
			return &Time{Value: o.Value.Add(time.Duration(rhs.Value))}, nil
		case token.Sub: // time - int => time
			if rhs.Value == 0 {
				return o, nil
			}
			return &Time{Value: o.Value.Add(time.Duration(-rhs.Value))}, nil
		}
	case *Time:
		switch op {
		case token.Sub: // time - time => int (duration)
			return &Int{Value: int64(o.Value.Sub(rhs.Value))}, nil
		case token.Less: // time < time => bool
			if o.Value.Before(rhs.Value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.Greater:
			if o.Value.After(rhs.Value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.LessEq:
			if o.Value.Equal(rhs.Value) || o.Value.Before(rhs.Value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		case token.GreaterEq:
			if o.Value.Equal(rhs.Value) || o.Value.After(rhs.Value) {
				return TrueValue, nil
			}
			return FalseValue, nil
		}
	}
	return nil, ErrInvalidOperator
}

// Copy returns a copy of the type.
func (o *Time) Copy() Object {
	return &Time{Value: o.Value}
}

// IsFalsy returns true if the value of the type is falsy.
func (o *Time) IsFalsy() bool {
	return o.Value.IsZero()
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *Time) Equals(x Object) bool {
	t, ok := x.(*Time)
	if !ok {
		return false
	}
	return o.Value.Equal(t.Value)
}

// Undefined represents an undefined value.
type Undefined struct {
	ObjectImpl
}

// TypeName returns the name of the type.
func (o *Undefined) TypeName() string {
	return "undefined"
}

func (o *Undefined) String() string {
	return "<undefined>"
}

// Copy returns a copy of the type.
func (o *Undefined) Copy() Object {
	return o
}

// IsFalsy returns true if the value of the type is falsy.
func (o *Undefined) IsFalsy() bool {
	return true
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *Undefined) Equals(x Object) bool {
	return o == x
}

// IndexGet returns an element at a given index.
func (o *Undefined) IndexGet(_ *VM, _ Object) (Object, error) {
	return UndefinedValue, nil
}

// Iterate creates a map iterator.
func (o *Undefined) Iterate() Iterator {
	return o
}

// CanIterate returns whether the Object can be Iterated.
func (o *Undefined) CanIterate() bool {
	return true
}

// Next returns true if there are more elements to iterate.
func (o *Undefined) Next() bool {
	return false
}

// Key returns the key or index value of the current element.
func (o *Undefined) Key() Object {
	return o
}

// Value returns the value of the current element.
func (o *Undefined) Value() Object {
	return o
}

// Default represents an undefined value of keyword argument.
type Default struct {
	ObjectImpl
}

// TypeName returns the name of the type.
func (o *Default) TypeName() string {
	return "default"
}

func (o *Default) String() string {
	return "<default>"
}

// Copy returns a copy of the type.
func (o *Default) Copy() Object {
	return o
}

// IsFalsy returns true if the value of the type is falsy.
func (o *Default) IsFalsy() bool {
	return true
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *Default) Equals(x Object) bool {
	return o == x
}

// IndexGet returns an element at a given index.
func (o *Default) IndexGet(_ *VM, _ Object) (Object, error) {
	return UndefinedValue, nil
}

// Iterate creates a map iterator.
func (o *Default) Iterate() Iterator {
	return o
}

// CanIterate returns whether the Object can be Iterated.
func (o *Default) CanIterate() bool {
	return true
}

// Next returns true if there are more elements to iterate.
func (o *Default) Next() bool {
	return false
}

// Key returns the key or index value of the current element.
func (o *Default) Key() Object {
	return o
}

// Value returns the value of the current element.
func (o *Default) Value() Object {
	return o
}

// UserFunction represents a user function.
type UserFunction struct {
	ObjectImpl
	Name     string
	Value    CallableFunc
	IsMethod bool
	Target   Object
}

// TypeName returns the name of the type.
func (o *UserFunction) TypeName() string {
	return "user-function:" + o.Name
}

func (o *UserFunction) String() string {
	return "<user-function>"
}

// Copy returns a copy of the type.
func (o UserFunction) Copy() Object {
	return &o
}

// Method shold return if caller is method
func (o *UserFunction) Method() bool {
	return o.IsMethod
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *UserFunction) Equals(_ Object) bool {
	return false
}

// Call invokes a user function.
func (o *UserFunction) Call(ctx *CallContext) (Object, error) {
	if len(ctx.Kwargs) > 0 {
		return nil, ErrUnexpectedKwargs
	}
	return o.Value(ctx.Args...)
}

// CanCall returns whether the Object can be Called.
func (o *UserFunction) CanCall() bool {
	return true
}

func (o *UserFunction) MethodTo() Object {
	return o.Target
}

func (o UserFunction) ToMethodOf(of Object) ToMethodConverter {
	o.Target = of
	return &o
}

// UserFunctionCtx represents a user function with call context.
type UserFunctionCtx struct {
	ObjectImpl
	Name     string
	Value    CallableFuncCtx
	IsMethod bool
	Target   Object
}

// TypeName returns the name of the type.
func (o *UserFunctionCtx) TypeName() string {
	return "user-function-ctx:" + o.Name
}

func (o *UserFunctionCtx) String() string {
	return "<user-function-ctx>"
}

// Copy returns a copy of the type.
func (o UserFunctionCtx) Copy() Object {
	return &o
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *UserFunctionCtx) Equals(_ Object) bool {
	return false
}

// Call invokes a user function.
func (o *UserFunctionCtx) Call(ctx *CallContext) (ret Object, err error) {
	return o.Value(ctx)
}

// CanCall returns whether the Object can be Called.
func (o *UserFunctionCtx) CanCall() bool {
	return true
}

// Method shold return if caller is method
func (o *UserFunctionCtx) Method() bool {
	return o.IsMethod
}

// Target should return current method target
func (o *UserFunctionCtx) MethodTo() Object {
	return o.Target
}

func (o UserFunctionCtx) ToMethodOf(of Object) ToMethodConverter {
	o.Target = of
	return &o
}

// VmContextObject represents the VmContext object
type VmContextObject struct {
	ObjectImpl
	Value *VmContext
}

// Copy returns a copy of the type.
func (c *VmContextObject) Copy() Object {
	return &VmContextObject{Value: c.Value}
}

func (c *VmContextObject) TypeName() string {
	return "context"
}

func (c *VmContextObject) String() string {
	return "<vm-context>"
}
