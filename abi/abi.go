package abi

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// ABI
type ABI struct {
	Constructor *Function
	Functions   map[string]Function
	Events      map[string]Event
	Errors      map[string]Error
}

func (abi *ABI) UnmarshalJSON(data []byte) error {
	type Value struct {
		Name       string          `json:"name"`
		Type       string          `json:"type"`
		Components []FunctionValue `json:"components"`
		Indexed    bool            `json:"indexed"`
		Anonymous  bool            `json:"anonymous"`
	}
	type Interface struct {
		Type            string           `json:"type"`
		Name            string           `json:"name"`
		Inputs          []Value          `json:"inputs"`
		Outputs         []Value          `json:"outputs"`
		StateMutability *StateMutability `json:"stateMutability"`
		Payable         bool             `json:"payable"`
		Constant        bool             `json:"constant"`
	}
	var interfaces []Interface

	if err := json.Unmarshal(data, &interfaces); err != nil {
		return err
	}

	abi.Functions = make(map[string]Function)
	abi.Events = make(map[string]Event)
	abi.Errors = make(map[string]Error)

	setToFunction := func(i *Interface, f *Function) error {
		if i.StateMutability == nil {
			return errors.New("stateMutability is required")
		}

		f.Type = FunctionType(i.Type)
		f.Name = i.Name
		f.StateMutability = *i.StateMutability
		f.Payable = i.Payable
		f.Constant = i.Constant

		f.Inputs = make([]FunctionValue, len(i.Inputs))
		for idx, input := range i.Inputs {
			f.Inputs[idx].Name = input.Name
			f.Inputs[idx].Type = input.Type
			f.Inputs[idx].Components = input.Components
		}

		f.Outputs = make([]FunctionValue, len(i.Outputs))
		for idx, output := range i.Outputs {
			f.Outputs[idx].Name = output.Name
			f.Outputs[idx].Type = output.Type
			f.Outputs[idx].Components = output.Components
		}

		return nil
	}

	setToEvent := func(i *Interface, e *Event) error {
		e.Type = EventType(i.Type)
		e.Name = i.Name
		e.Inputs = make([]EventValue, len(i.Inputs))
		for idx, input := range i.Inputs {
			e.Inputs[idx].Name = input.Name
			e.Inputs[idx].Type = input.Type
			e.Inputs[idx].Components = input.Components
			e.Inputs[idx].Indexed = input.Indexed
			e.Inputs[idx].Anonymous = input.Anonymous
		}
		return nil
	}

	setToError := func(i *Interface, e *Error) error {
		e.Type = ErrorType(i.Type)
		e.Name = i.Name
		e.Inputs = make([]ErrorValue, len(i.Inputs))
		for idx, input := range i.Inputs {
			e.Inputs[idx].Name = input.Name
			e.Inputs[idx].Type = input.Type
			e.Inputs[idx].Components = input.Components
		}
		return nil
	}

	for _, i := range interfaces {
		switch i.Type {
		case string(FunctionTypeFunction), string(FunctionTypeReceive):
			if _, ok := abi.Functions[i.Name]; ok {
				return fmt.Errorf("function %s is already set", i.Name)
			}
			f := Function{}
			if err := setToFunction(&i, &f); err != nil {
				return err
			}
			selector, err := CalculateSelector(&f)
			if err != nil {
				return err
			}
			f.Selector = selector
			abi.Functions[i.Name] = f
			break
		case string(FunctionTypeConstructor):
			if abi.Constructor != nil {
				return errors.New("constructor is already set")
			}
			f := Function{}
			if err := setToFunction(&i, &f); err != nil {
				return err
			}
			abi.Constructor = &f
			break
		case string(EventTypeEvent):
			if _, ok := abi.Events[i.Name]; ok {
				return fmt.Errorf("event %s is already set", i.Name)
			}
			e := Event{}
			if err := setToEvent(&i, &e); err != nil {
				return err
			}
			abi.Events[i.Name] = e
			break
		case string(ErrorTypeError):
			if _, ok := abi.Errors[i.Name]; ok {
				return fmt.Errorf("error %s is already set", i.Name)
			}
			e := Error{}
			if err := setToError(&i, &e); err != nil {
				return err
			}
			abi.Errors[i.Name] = e
			break
		}
	}
	return nil
}

// Function
type Function struct {
	Type            FunctionType
	Name            string
	Inputs          []FunctionValue
	Outputs         []FunctionValue
	StateMutability StateMutability
	Payable         bool
	Constant        bool
	// Additional Field
	Selector [4]byte
}

// Event
type Event struct {
	Type   EventType
	Name   string
	Inputs []EventValue
}

// Error
type Error struct {
	Name   string
	Type   ErrorType
	Inputs []ErrorValue
}

// ValueTypes
type IValue interface {
	GetName() string
	GetType() string
}

type BaseValue struct {
	Name       string
	Type       string
	Components []BaseValue
}

func (b *BaseValue) GetName() string {
	return b.Name
}

func (b *BaseValue) GetType() string {
	return b.Type
}

type FunctionValue = BaseValue

type EventValue struct {
	BaseValue
	Indexed   bool
	Anonymous bool
}

type ErrorValue = BaseValue

// Enums
type FunctionType string

const (
	FunctionTypeFunction    FunctionType = "function"
	FunctionTypeConstructor FunctionType = "constructor"
	FunctionTypeReceive     FunctionType = "receive"
)

func (f *FunctionType) UnmarshalJSON(b []byte) error {
	s := FunctionType(b)
	switch s {
	case FunctionTypeFunction, FunctionTypeConstructor, FunctionTypeReceive:
		*f = s
		return nil
	}
	return errors.New("invalid function type")
}

type StateMutability string

const (
	StateMutabilityPure       StateMutability = "pure"
	StateMutabilityView       StateMutability = "view"
	StateMutabilityNonPayable StateMutability = "nonpayable"
	StateMutabilityPayable    StateMutability = "payable"
)

func (m *StateMutability) UnmarshalJSON(b []byte) error {
	s := StateMutability(strings.Trim(string(b), "\""))
	switch s {
	case StateMutabilityPure, StateMutabilityView, StateMutabilityNonPayable, StateMutabilityPayable:
		*m = s
		return nil
	}
	return errors.New("invalid state mutability: " + string(s))
}

type EventType string

const (
	EventTypeEvent EventType = "event"
)

func (e *EventType) UnmarshalJSON(b []byte) error {
	s := EventType(strings.Trim(string(b), "\""))
	switch s {
	case EventTypeEvent:
		*e = s
		return nil
	}
	return errors.New("invalid event type")
}

type ErrorType string

const (
	ErrorTypeError ErrorType = "error"
)

func (e *ErrorType) UnmarshalJSON(b []byte) error {
	s := ErrorType(strings.Trim(string(b), "\""))
	switch s {
	case ErrorTypeError:
		*e = s
		return nil
	}
	return errors.New("invalid error type")
}
