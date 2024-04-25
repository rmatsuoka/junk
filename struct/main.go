package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	errUnexpectedEOF = errors.New("unexpected EOF")
)

type parser struct {
	args []string
	i    int
}

func (p *parser) next() (string, bool) {
	if p.i >= len(p.args) {
		return "", false
	}
	x := p.args[p.i]
	p.i++
	return x, true
}

func parse(args []string) (any, error) {
	var (
		ret any
		err error
		p   = parser{args: args}
	)
	s, ok := p.next()
	if !ok {
		return nil, errUnexpectedEOF
	}
	switch s {
	case "{":
		ret, err = p.object()
	case "[":
		ret, err = p.array()
	default:
		// return nil, fmt.Errorf("expected { or ( but got %s", s)
		panic(fmt.Sprintf("expected { or ( but got %s", s))
	}
	return ret, err
}

func (p *parser) object() (map[string]any, error) {
	var (
		ret = make(map[string]any)
		err error
	)
	for s, ok := p.next(); ok; s, ok = p.next() {
		switch s {
		case "}":
			return ret, err
		default:
			key, value, err := p.keyvalue(s)
			if err != nil {
				return nil, err
			}
			ret[key] = value
		}
	}
	panic(errUnexpectedEOF)
	// return nil, errUnexpectedEOF
}

var errTypeMismatch = errors.New("type mismatch")

func (p *parser) keyvalue(s string) (key string, value any, err error) {
	keytype, v, found := strings.Cut(s, "=")
	if !found {
		// return "", "", fmt.Errorf("no =")
		panic("no =: " + s)
	}

	key, t, _ := lastCut(keytype, ":")
	typ, err := parseType(t)
	if err != nil {
		panic(err)
	}

	if typ.name == typeAny {
		if len(v) != 1 {
			return key, v, nil
		}
		switch v[0] {
		case '{':
			value, err = p.object()
		case '[':
			value, err = p.array()
		default:
			value = v
		}
	} else {
		switch typ.name {
		case typeNumber, typeString, typeBoolean:
			value, err = types[typ.name](v)
		case typeObject:
			if v != "{" {
				return "", nil, errTypeMismatch
			}
			value, err = p.object()
		case typeArray:
			if v != "[" {
				return "", nil, errTypeMismatch
			}
			value, err = p.fixedArray(typ.typeVar, typ.length)
		default:
			panic("unreachable")
		}
	}
	return key, value, err
}

func lastCut(s, sep string) (before, after string, found bool) {
	i := strings.LastIndex(s, sep)
	if i == -1 {
		return s, "", false
	}
	return s[:i], s[i+len(sep):], true
}

type typ struct {
	name    typeName
	typeVar typeName // for array
	length  int      // for array
}

type typeName string

const (
	typeNumber  typeName = "number"
	typeString  typeName = "string"
	typeBoolean typeName = "boolean"
	typeObject  typeName = "object"
	typeArray   typeName = "array"
	typeAny     typeName = "any"
)

var errUnknownType = errors.New("unknown type")

func parseType(s string) (typ, error) {
	if len(s) == 0 {
		return typ{name: typeAny}, nil
	}
	if s[len(s)-1] == ']' {
		if len(s) == 1 {
			panic(errUnknownType)
			// return typ{}, errUnknownType
		}
		i := strings.LastIndexByte(s, '[')
		if i < 0 {
			panic(errUnknownType)
			// return typ{}, errUnknownType
		}
		var l int
		if i == len(s)-2 {
			l = 0
		} else {
			var err error
			l, err = strconv.Atoi(s[i+1 : len(s)-1])
			if err != nil {
				log.Panicf("s[%d:%d]: %v", i+1, len(s)-2, errUnknownType.Error())
				// return typ{}, errUnknownType
			}
		}
		typeVar := typeName(s[:i])
		if typeVar == "" {
			typeVar = typeAny
		}
		if _, ok := types[typeVar]; !ok {
			panic(errUnknownType)
			// return typ{}, errUnknownType
		}
		return typ{name: typeArray, typeVar: typeVar, length: l}, nil
	}
	name := typeName(s)
	if name == "" {
		name = typeAny
	}
	_, ok := types[name]
	if !ok {
		panic(errUnknownType)
		// return typ{}, errUnknownType
	}
	return typ{name: name}, nil
}

func unmarshal[T any](s string) (any, error) {
	var x T
	err := json.Unmarshal([]byte(s), &x)
	return x, err
}

var types = map[typeName]func(s string) (any, error){
	typeNumber:  unmarshal[float64],
	typeBoolean: unmarshal[bool],
	typeString:  func(s string) (any, error) { return s, nil },
	typeObject:  func(s string) (any, error) { panic("not reachable") },
	typeArray:   func(s string) (any, error) { panic("not reachable") },
	typeAny:     func(s string) (any, error) { panic("not reachble") },
}

func (p *parser) array() ([]any, error) {
	var (
		ret []any
		err error
	)
	for s, ok := p.next(); ok; s, ok = p.next() {
		if len(s) != 1 {
			ret = append(ret, s)
			continue
		}
		switch s {
		case "{":
			o, err := p.object()
			if err != nil {
				return nil, err
			}
			ret = append(ret, o)
		case "[":
			a, err := p.array()
			if err != nil {
				return nil, err
			}
			ret = append(ret, a)
		case "]":
			return ret, err
		default:
			ret = append(ret, s)
		}
	}
	panic(errUnexpectedEOF)
	// return nil, errUnexpectedEOF
}

func (p *parser) fixedArray(typeVar typeName, length int) ([]any, error) {
	var (
		ret []any
		err error
	)
	for i := 0; i < length; i++ {
		s, ok := p.next()
		if !ok {
			panic(errUnexpectedEOF)
		}
		var value any
		switch typeVar {
		case typeNumber, typeString, typeBoolean:
			value, err = types[typeVar](s)
		case typeObject:
			if s != "{" {
				return nil, errTypeMismatch
			}
			value, err = p.object()
		case typeArray:
			if s != "[" {
				return nil, errTypeMismatch
			}
			value, err = p.array()
		case typeAny:
			value = s
		default:
			panic("unreachable")
		}
		if err != nil {
			return nil, err
		}
		ret = append(ret, value)
	}
	s, _ := p.next()
	if s != "]" {
		return nil, errTypeMismatch
	}
	return ret, err
}

func main() {
	flag.Parse()

	x, err := parse(flag.Args())
	if err != nil {
		log.Fatal(err)
	}
	j, err := json.Marshal(x)
	if err != nil {
		log.Fatal(err)
	}
	b := new(bytes.Buffer)
	json.Indent(b, j, "", "  ")
	b.WriteByte('\n')
	b.WriteTo(os.Stdout)
}
