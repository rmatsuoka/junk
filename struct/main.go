package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
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
	p := parser{args: args}
	s, ok := p.next()
	if !ok {
		return nil, errUnexpectedEOF
	}
	return p.value(s, typ{name: typeAny})
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

	value, err = p.value(v, typ)
	return key, value, err
}

func (p *parser) value(v string, typ typ) (any, error) {
	if typ.name == typeAny {
		switch v {
		case "{":
			return p.object()
		case "[":
			return p.array()
		default:
			return v, nil
		}
	} else {
		switch typ.name {
		case typeNumber, typeString, typeBoolean:
			return unmarshal[typ.name](v)
		case typeObject:
			if v != "{" {
				return nil, errTypeMismatch
			}
			return p.object()
		case typeArray:
			if v != "[" {
				return nil, errTypeMismatch
			}
			return p.fixedArray(typ.typeVar, typ.length)
		default:
			panic("unreachable")
		}
	}
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
		if _, ok := unmarshal[typeVar]; !ok {
			panic(errUnknownType)
			// return typ{}, errUnknownType
		}
		return typ{name: typeArray, typeVar: typeVar, length: l}, nil
	}
	name := typeName(s)
	if name == "" {
		name = typeAny
	}
	_, ok := unmarshal[name]
	if !ok {
		panic(errUnknownType)
		// return typ{}, errUnknownType
	}
	return typ{name: name}, nil
}

func unmarshalString[T any](s string) (any, error) {
	var x T
	err := json.Unmarshal([]byte(s), &x)
	return x, err
}

var unmarshal = map[typeName]func(s string) (any, error){
	typeNumber:  unmarshalString[float64],
	typeBoolean: unmarshalString[bool],
	typeString:  func(s string) (any, error) { return s, nil },
	typeObject:  func(s string) (any, error) { panic("not reachable") },
	typeArray:   func(s string) (any, error) { panic("not reachable") },
	typeAny:     func(s string) (any, error) { panic("not reachble") },
}

func (p *parser) array() ([]any, error) {
	var (
		ret []any
	)
	for s, ok := p.next(); ok; s, ok = p.next() {
		if s == "]" {
			return ret, nil
		}
		v, err := p.value(s, typ{name: typeAny})
		if err != nil {
			panic(err)
		}
		ret = append(ret, v)
	}
	panic(errUnexpectedEOF)
	// return nil, errUnexpectedEOF
}

func (p *parser) fixedArray(typeVar typeName, length int) ([]any, error) {
	var (
		ret []any
	)
	for i := 0; i < length; i++ {
		s, ok := p.next()
		if !ok {
			panic(errUnexpectedEOF)
		}
		value, err := p.value(s, typ{name: typeVar})
		if err != nil {
			return nil, err
		}
		ret = append(ret, value)
	}
	s, _ := p.next()
	if s != "]" {
		return nil, errTypeMismatch
	}
	return ret, nil
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
