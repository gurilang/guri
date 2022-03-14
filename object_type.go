package tengo

import (
	"fmt"
	"strings"
)

type ObjectInstancer interface {
	InstanceType() *Type
}

type TypeProperty struct {
	ObjectImpl
	Getter,
	Setter Object
	Tags map[string]Object
}

func (c *TypeProperty) TypeName() string {
	return "type_property"
}

func (c *TypeProperty) String() string {
	var s string
	if c.Getter != nil {
		s = "r"
	}
	if c.Setter != nil {
		s += "w"
	}
	return fmt.Sprintf("<type_property %s>", s)
}

// Copy returns a copy of the type.
func (o TypeProperty) Copy() Object {
	return &o
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *TypeProperty) Equals(x Object) bool {
	return o == x
}

// IndexGet returns an element at a given index.
func (o *TypeProperty) IndexGet(_ *VM, key Object) (Object, error) {
	s, ok := key.(*String)
	if !ok {
		return nil, ErrInvalidIndexType
	}
	switch s.Value {
	case "tags":
		return &Map{Value: o.Tags}, nil
	case "get":
		return o.Getter, nil
	case "set":
		return o.Setter, nil
	case "__map__":
		return ToMapMethod(o), nil
	default:
		return nil, ErrInvalidIndex
	}
}

// IndexSet sets an element at a given index.
func (o *TypeProperty) IndexSet(_ *VM, key, value Object) (err error) {
	s, ok := key.(*String)
	if !ok {
		return ErrInvalidIndexType
	}
	switch s.Value {
	case "__map__":
		if toM, ok := value.(ToMapConverter); ok {
			for key, value := range toM.ToMap(true).Value {
				if err = o.Set(key, value); err != nil {
					return
				}
			}
			return
		}
		return fmt.Errorf("value of %q isn't convertible to map", key.String())
	default:
		return o.Set(s.Value, value)
	}
}

// IndexSet sets an element at a given index.
func (o *TypeProperty) Set(key string, value Object) (err error) {
	switch key {
	case "get":
		if value == UndefinedValue {
			o.Getter = nil
			return nil
		}
		if !value.Method() {
			return fmt.Errorf("value of %q isn't method", key)
		}
		o.Getter = value
		return
	case "set":
		if value == UndefinedValue {
			o.Setter = nil
			return nil
		}
		if !value.Method() {
			return fmt.Errorf("value of %q isn't method", key)
		}
		o.Setter = value
		return
	case "tags":
		if toM, ok := value.(ToMapConverter); ok {
			o.Tags = toM.ToMap(false).Value
			return
		}
		return fmt.Errorf("value of %q isn't convertible to map", key)
	default:
		return ErrInvalidIndex
	}
}

func (o *TypeProperty) ToMap(bool) *Map {
	m := make(map[string]Object, 2)
	if o.Getter != nil {
		m["get"] = o.Getter
	}
	if o.Setter != nil {
		m["set"] = o.Setter
	}
	if o.Tags != nil {
		m["tags"] = &Map{Value: o.Tags}
	}
	return &Map{Value: m}
}

func (o *TypeProperty) FromMap(m *Map) (err error) {
	var get, set, tags = m.Value["get"], m.Value["set"], m.Value["tags"]
	if get != nil {
		if err = o.Set("get", get); err != nil {
			return
		}
	}
	if set != nil {
		if err = o.Set("set", set); err != nil {
			return
		}
	}
	if tags != nil {
		if err = o.Set("tags", tags); err != nil {
			return
		}
	}
	return
}

type TypeProperties struct {
	ObjectImpl
	Value map[string]*TypeProperty
}

func (c *TypeProperties) TypeName() string {
	return "properties"
}

func (c *TypeProperties) String() string {
	return "<properties>"
}

// Copy returns a copy of the type.
func (o TypeProperties) Copy() Object {
	value := make(map[string]*TypeProperty, len(o.Value))
	for k, v := range o.Value {
		vcopy := *v
		value[k] = &vcopy
	}
	return &o
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *TypeProperties) Equals(x Object) bool {
	return o == x
}

// IndexGet returns an element at a given index.
func (o *TypeProperties) IndexGet(_ *VM, key Object) (Object, error) {
	s, ok := key.(*String)
	if !ok {
		return nil, ErrInvalidIndexType
	}
	switch s.Value {
	case "__map__":
		return ToMapMethod(o), nil
	}
	return o.Value[s.Value], nil
}

// Set sets an element at a given index.
func (o *TypeProperties) Set(key string, value Object) (err error) {
	switch key {
	case "":
		return ErrInvalidIndex
	case "__map__":
		if toM, ok := value.(ToMapConverter); ok {
			for key, value := range toM.ToMap(true).Value {
				if err = o.Set(key, value); err != nil {
					return
				}
			}
			return
		}
		return fmt.Errorf("value of %q isn't convertible to map", key)
	default:
		if len(key) > 4 {
			switch key[0:4] {
			case "get_":
				key := key[4:]
				prop := o.Value[key]
				if prop == nil {
					prop = &TypeProperty{}
					o.Value[key] = prop
				}
				prop.Getter = value
				return nil
			case "set_":
				key := key[4:]
				prop := o.Value[key]
				if prop == nil {
					prop = &TypeProperty{}
					o.Value[key] = prop
				}
				prop.Setter = value
				return nil
			}
		}

		if value == UndefinedValue {
			delete(o.Value, key)
			return nil
		}

		switch t := value.(type) {
		case *Map:
			var g, s = t.Value["get"], t.Value["set"]
			if g == UndefinedValue {
				g = nil
			}
			if s == UndefinedValue {
				s = nil
			}
			if g == nil && s == nil {
				delete(o.Value, key)
			}
			prop := &TypeProperty{}
			prop.Getter = g
			prop.Setter = s
			o.Value[key] = prop
			return nil
		case *TypeProperty:
			o.Value[key] = t
			return nil
		default:
			return ErrInvalidIndexValueType
		}
	}
}

func (o *TypeProperties) FromMap(m *Map) (err error) {
	for key, value := range m.Value {
		if err = o.Set(key, value); err != nil {
			return
		}
	}
	return
}

// IndexSet sets an element at a given index.
func (o *TypeProperties) IndexSet(_ *VM, key, value Object) error {
	s, ok := key.(*String)
	if !ok {
		return ErrInvalidIndexType
	}
	if s.Value == "" {
		return ErrInvalidIndex
	}

	switch s.Value {
	case "__map__":
		if toM, ok := value.(ToMapConverter); ok {
			return o.FromMap(toM.ToMap(true))
		}
		return fmt.Errorf("value of %q isn't convertible to map", key.String())
	default:
		return o.Set(s.Value, value)
	}
}

func (o *TypeProperties) ToMap(deep bool) *Map {
	m := make(map[string]Object, len(o.Value))
	if deep {
		for k, v := range o.Value {
			m[k] = v.ToMap(true)
		}
	} else {
		for key, value := range o.Value {
			m[key] = value
		}
	}
	return &Map{Value: m}
}

type Type struct {
	ObjectImpl

	Name       string
	Tags       map[string]Object
	New        Object
	Methods    *TypeMethods
	Fields     *TypeFields
	Properties *TypeProperties
}

func (c *Type) TypeName() string {
	return "type"
}

func (c *Type) String() string {
	return fmt.Sprintf("<type: %s>", c.Name)
}

func (o *Type) Call(ctx *CallContext) (ret Object, err error) {
	if o.New == nil {
		var (
			obj = &Instance{
				Type:   o,
				Values: make(map[string]Object, len(ctx.Kwargs)),
			}
		)

		for name, value := range o.Fields.Value {
			obj.Values[name] = value.Value.Copy()
		}

		for name, value := range ctx.Kwargs {
			obj.Set(ctx.VM, name, value)
		}
		return obj, nil
	}

	switch t := o.New.(type) {
	case *CompiledFunction:
		var (
			fn  = t.Copy().(*CompiledFunction)
			obj = &Instance{
				Type:   o,
				Values: map[string]Object{},
			}
		)

		fn.methodTarget = obj

		for name, value := range o.Fields.Value {
			obj.Values[name] = value.Value.Copy()
		}

		if _, err = fn.Call(ctx); err == nil {
			ret = obj
		}
	default:
		ret, err = t.Call(&CallContext{VM: ctx.VM, This: o, Args: ctx.Args, Kwargs: ctx.Kwargs})
	}
	return
}

func (o *Type) CanCall() bool {
	return true
}

// Copy returns a copy of the type.
func (o *Type) Copy() Object {
	return o
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *Type) Equals(x Object) bool {
	return o == x
}

func (o *Type) ToMap(deep bool) *Map {
	v := map[string]Object{
		"name": &String{Value: o.Name},
	}

	if o.New != nil {
		v["new"] = o.New
	}

	if o.Tags != nil {
		v["tags"] = &Map{Value: o.Tags}
	}

	v["methods"] = o.Methods.ToMap(deep)
	v["fields"] = o.Fields.ToMap(deep)
	v["properties"] = o.Properties.ToMap(deep)

	return &Map{Value: v}
}

// IndexGet returns an element at a given index.
func (o *Type) IndexGet(_ *VM, key Object) (Object, error) {
	s, ok := key.(*String)
	if !ok {
		return nil, ErrInvalidIndexType
	}

	switch s.Value {
	case "tags":
		return &Map{Value: o.Tags}, nil
	case "new":
		return o.New, nil
	case "fields":
		return o.Fields, nil
	case "methods":
		return o.Methods, nil
	case "props":
		return o.Properties, nil
	case "__map__":
		return o.ToMap(true), nil
	default:
		return nil, ErrInvalidIndex
	}
}

// IndexSet sets an element at a given index.
func (o *Type) IndexSet(_ *VM, key, value Object) (err error) {
	s, ok := key.(*String)
	if !ok {
		return ErrInvalidIndexType
	}
	if s.Value == "" {
		return ErrInvalidIndex
	}

	switch s.Value {
	case "tags":
		if toM, ok := value.(ToMapConverter); ok {
			o.Tags = toM.ToMap(false).Value
			return
		}
		return fmt.Errorf("value of %q isn't convertible to map", key)
	case "new":
		if !value.CanCall() {
			return fmt.Errorf("value of %q isn't callable", key)
		}
		if !value.Method() {
			return fmt.Errorf("value of %q isn't method", key)
		}
		o.New = value
		return
	case "fields":
		if v, ok := value.(*TypeFields); ok {
			o.Fields = v
			return
		}
		return fmt.Errorf("value of %q isn't fields", key)
	case "methods":
		return
	case "props":
		return
	default:
		return ErrInvalidIndex
	}
}

func (o *Type) Value() interface{} {
	return o
}

// Instance represents a map of objects.
type Instance struct {
	ObjectImpl
	Type     *Type
	Callable bool
	Values   map[string]Object
	Methods  map[string]ToMethodConverter
}

// TypeName returns the name of the type.
func (o *Instance) TypeName() string {
	return "instance::" + o.Type.Name
}

func (o *Instance) String() string {
	var pairs []string
	for k, v := range o.Values {
		pairs = append(pairs, fmt.Sprintf("%s: %s", k, v.String()))
	}
	return fmt.Sprintf("<%s #%p {%s}>", o.Type.Name, o, strings.Join(pairs, ", "))
}

// Copy returns a copy of the type.
func (o *Instance) Copy() Object {
	c := make(map[string]Object)
	for k, v := range o.Values {
		c[k] = v.Copy()
	}
	return &Instance{Values: c}
}

// IsFalsy returns true if the value of the type is falsy.
func (o *Instance) IsFalsy() bool {
	return len(o.Values) == 0
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *Instance) Equals(x Object) bool {
	var xVal map[string]Object
	switch x := x.(type) {
	case *Instance:
		xVal = x.Values
	default:
		return false
	}
	if len(o.Values) != len(xVal) {
		return false
	}
	for k, v := range o.Values {
		tv := xVal[k]
		if !v.Equals(tv) {
			return false
		}
	}
	return true
}

// IndexGet returns the value for the given key.
func (o *Instance) IndexGet(Vm *VM, index Object) (res Object, err error) {
	name, ok := ToString(index)
	if !ok {
		err = ErrInvalidIndexType
		return
	}

	switch name {
	case "__map__":
		return &Map{Value: o.Values}, nil
	case "__type__":
		return o.Type, nil
	default:
		if res, ok = o.Values[name]; ok {
			return
		}
		if m, ok := o.Type.Methods.Value[name]; ok {
			if res, ok = o.Methods[name]; ok {
				return
			} else if o.Methods == nil {
				o.Methods = map[string]ToMethodConverter{}
			}
			m2 := m.Value.ToMethodOf(o)
			o.Methods[name] = m2
			return m2, nil
		}
		if prop := o.Type.Properties.Value[name]; prop != nil {
			if prop.Setter != nil {
				_, err = prop.Setter.Call(&CallContext{VM: Vm, This: o})
			}
		}
		return UndefinedValue, nil
	}
}

// IndexSet sets the value for the given key.
func (o *Instance) IndexSet(_ *VM, index, value Object) (err error) {
	strIdx, ok := ToString(index)
	if !ok {
		err = ErrInvalidIndexType
		return
	}
	o.Values[strIdx] = value
	return nil
}

// IndexSet sets the value for the given key.
func (o *Instance) Set(Vm *VM, name string, value Object) (err error) {
	if prop := o.Type.Properties.Value[name]; prop != nil {
		if prop.Setter != nil {
			_, err = prop.Setter.Call(&CallContext{VM: Vm, This: o, Args: []Object{value}})
		}
	} else {
		o.Values[name] = value
	}
	return
}

func (o *Instance) InstanceType() *Type {
	return o.Type
}

// Iterate creates a map iterator.
func (o *Instance) Iterate() Iterator {
	var keys []string
	for k := range o.Values {
		keys = append(keys, k)
	}
	return &MapIterator{
		v: o.Values,
		k: keys,
		l: len(keys),
	}
}

// CanIterate returns whether the Object can be Iterated.
func (o *Instance) CanIterate() bool {
	return true
}
