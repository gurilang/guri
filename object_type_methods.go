package tengo

import (
	"fmt"
)

type TypeMethod struct {
	ObjectImpl
	Tags  map[string]Object
	Value ToMethodConverter
}

func (c *TypeMethod) TypeName() string {
	return "type_method"
}

func (c *TypeMethod) String() string {
	return fmt.Sprintf("<type_method %s>", c.ToMap(true))
}

// Copy returns a copy of the type.
func (o TypeMethod) Copy() Object {
	return &o
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *TypeMethod) Equals(x Object) bool {
	return o == x
}

// IndexGet returns an element at a given index.
func (o *TypeMethod) IndexGet(_ *VM, key Object) (Object, error) {
	s, ok := key.(*String)
	if !ok {
		return nil, ErrInvalidIndexType
	}
	switch s.Value {
	case "tags":
		return &Map{Value: o.Tags}, nil
	case "value":
		return o.Value, nil
	default:
		return nil, ErrInvalidIndex
	}
}

// IndexSet sets an element at a given index.
func (o *TypeMethod) Set(key string, value Object) (err error) {
	switch key {
	case "value":
		if !value.CanCall() {
			return fmt.Errorf("value of %q isn't callable", key)
		}
		if conv, ok := value.(ToMethodConverter); ok {
			o.Value = conv
		} else {
			return fmt.Errorf("value of %q isn't convertible to method", key)
		}
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

func (o *TypeMethod) ToMap(bool) *Map {
	m := make(map[string]Object, 2)
	m["value"] = o.Value
	if len(o.Tags) > 0 {
		m["tags"] = &Map{Value: o.Tags}
	}
	return &Map{Value: m}
}

func (o *TypeMethod) FromMap(m *Map) error {
	value := m.Value["value"]
	if !value.CanCall() {
		return fmt.Errorf("value isn't callable")
	}
	if conv, ok := value.(ToMethodConverter); ok {
		o.Value = conv
	} else {
		return fmt.Errorf("value isn't convertible to method")
	}
	o.Tags = nil
	if tags, ok := m.Value["tags"]; ok {
		if m, ok := tags.(ToMapConverter); ok {
			o.Tags = m.ToMap(false).Value
		} else if tags != UndefinedValue {
			return fmt.Errorf("\"tags\" value isn't map")
		}
	}
	return nil
}

func (o *TypeMethod) CanCall() bool {
	return true
}

// IndexSet sets an element at a given index.
func (o *TypeMethod) IndexSet(_ *VM, key, value Object) (err error) {
	s, ok := key.(*String)
	if !ok {
		return ErrInvalidIndexType
	}
	switch s.Value {
	case "__map__":
		if toM, ok := value.(ToMapConverter); ok {
			for key, value := range toM.ToMap(false).Value {
				if err = o.Set(key, value); err != nil {
					return
				}
			}
			return
		}
		return fmt.Errorf("value of %q isn't convertible to map", s.Value)
	default:
		return o.Set(s.Value, value)
	}
}

type TypeMethods struct {
	ObjectImpl
	Value map[string]*TypeMethod
}

func (c *TypeMethods) TypeName() string {
	return "type_methods"
}

func (c *TypeMethods) String() string {
	return fmt.Sprintf("<type_methods %s>", c.ToMap(false))
}

// Copy returns a copy of the type.
func (o *TypeMethods) Copy() Object {
	value := make(map[string]*TypeMethod, len(o.Value))
	for k, v := range o.Value {
		value[k] = v.Copy().(*TypeMethod)
	}
	return &TypeMethods{Value: value}
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *TypeMethods) Equals(x Object) bool {
	return o == x
}

// IndexGet returns an element at a given index.
func (o *TypeMethods) IndexGet(_ *VM, key Object) (Object, error) {
	s, ok := key.(*String)
	if !ok {
		return nil, ErrInvalidIndexType
	}
	return o.Value[s.Value], nil
}

// IndexGet returns an element at a given index.
func (o *TypeMethods) IndexDel(_ *VM, key ...Object) error {
	for _, key := range key {
		s, ok := key.(*String)
		if !ok {
			return ErrInvalidIndexType
		}
		delete(o.Value, s.Value)
	}
	return nil
}

// Set sets an element at a given index.
func (o *TypeMethods) Set(key string, value Object) (err error) {
	switch key {
	case "":
		return ErrInvalidIndex
	default:
		switch t := value.(type) {
		case *Map:
			f := &TypeMethod{}
			if err = f.FromMap(t); err != nil {
				return err
			}
			o.Value[key] = f
			return
		default:
			if !t.CanCall() {
				return fmt.Errorf("value of %q isn't callable", key)
			}
			if conv, ok := value.(ToMethodConverter); ok {
				o.Value[key] = &TypeMethod{Value: conv}
			} else {
				return fmt.Errorf("value of %q isn't convertible to method", key)
			}
			return
		}
	}
}

func (o *TypeMethods) FromMap(m *Map) (err error) {
	for key, value := range m.Value {
		if err = o.Set(key, value); err != nil {
			return
		}
	}
	return
}

// IndexSet sets an element at a given index.
func (o *TypeMethods) IndexSet(_ *VM, key, value Object) error {
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

func (o *TypeMethods) ToMap(deep bool) *Map {
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
