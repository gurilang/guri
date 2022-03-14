package tengo

import (
	"fmt"
)

type TypeField struct {
	ObjectImpl
	Tags  map[string]Object
	Value Object
}

func (c *TypeField) TypeName() string {
	return "type_field"
}

func (c *TypeField) String() string {
	return fmt.Sprintf("<type_field %s>", c.ToMap(true))
}

// Copy returns a copy of the type.
func (o TypeField) Copy() Object {
	return &o
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *TypeField) Equals(x Object) bool {
	return o == x
}

// IndexGet returns an element at a given index.
func (o *TypeField) IndexGet(_ *VM, key Object) (Object, error) {
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
func (o *TypeField) Set(key string, value Object) (err error) {
	switch key {
	case "value":
		o.Value = value
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

func (o *TypeField) ToMap(bool) *Map {
	m := make(map[string]Object, 2)
	m["value"] = o.Value
	if o.Tags != nil && len(o.Tags) > 0 {
		m["tags"] = &Map{Value: o.Tags}
	}
	return &Map{Value: m}
}

func (o *TypeField) FromMap(m *Map) error {
	o.Value = m.Value["value"]
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

// IndexSet sets an element at a given index.
func (o *TypeField) IndexSet(_ *VM, key, value Object) (err error) {
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

type TypeFields struct {
	ObjectImpl
	Value map[string]*TypeField
}

func (c *TypeFields) TypeName() string {
	return "type_fields"
}

func (c *TypeFields) String() string {
	return fmt.Sprintf("<type_fields %s>", c.ToMap(true))
}

// Copy returns a copy of the type.
func (o *TypeFields) Copy() Object {
	value := make(map[string]*TypeField, len(o.Value))
	for k, v := range o.Value {
		value[k] = v.Copy().(*TypeField)
	}
	return &TypeFields{Value: value}
}

// Equals returns true if the value of the type is equal to the value of
// another object.
func (o *TypeFields) Equals(x Object) bool {
	return o == x
}

// IndexGet returns an element at a given index.
func (o *TypeFields) IndexGet(_ *VM, key Object) (Object, error) {
	s, ok := key.(*String)
	if !ok {
		return nil, ErrInvalidIndexType
	}
	return o.Value[s.Value], nil
}

// IndexGet returns an element at a given index.
func (o *TypeFields) IndexDel(_ *VM, key ...Object) error {
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
func (o *TypeFields) Set(key string, value Object) (err error) {
	switch key {
	case "":
		return ErrInvalidIndex
	default:
		switch t := value.(type) {
		case *TypeField:
			o.Value[key] = t
		case *Map:
			var f TypeField
			if err = f.FromMap(t); err != nil {
				return
			}
			o.Value[key] = &f
		case *Default:
			o.Value[key] = &TypeField{Value: UndefinedValue, Tags: make(map[string]Object, 0)}
		default:
			o.Value[key] = &TypeField{Value: value, Tags: make(map[string]Object, 0)}
		}
	}
	return nil
}

func (o *TypeFields) FromMap(m *Map) (err error) {
	for key, value := range m.Value {
		if err = o.Set(key, value); err != nil {
			return
		}
	}
	return
}

// IndexSet sets an element at a given index.
func (o *TypeFields) IndexSet(_ *VM, key, value Object) error {
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

func (o *TypeFields) ToMap(deep bool) *Map {
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
