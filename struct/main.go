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
	return p.value(s, typ{kind: kindAny})
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
	if typ.kind == kindAny {
		switch v {
		case "{":
			return p.object()
		case "[":
			return p.array()
		default:
			return v, nil
		}
	} else {
		switch typ.kind {
		case kindNumber, kindString, kindBoolean:
			return unmarshal[typ.kind](v)
		case kindObject:
			if v != "{" {
				return nil, errTypeMismatch
			}
			return p.object()
		case kindArray:
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
	kind    typeKind
	typeVar typeKind // for array
	length  int      // for array
}

type typeKind string

const (
	kindNumber  typeKind = "number"
	kindString  typeKind = "string"
	kindBoolean typeKind = "boolean"
	kindObject  typeKind = "object"
	kindArray   typeKind = "array"
	kindAny     typeKind = "any"
)

var errUnknownType = errors.New("unknown type")

func parseType(s string) (typ, error) {
	if len(s) == 0 {
		return typ{kind: kindAny}, nil
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
		typeVar := typeKind(s[:i])
		if typeVar == "" {
			typeVar = kindAny
		}
		if _, ok := unmarshal[typeVar]; !ok {
			panic(errUnknownType)
			// return typ{}, errUnknownType
		}
		return typ{kind: kindArray, typeVar: typeVar, length: l}, nil
	}
	name := typeKind(s)
	if name == "" {
		name = kindAny
	}
	_, ok := unmarshal[name]
	if !ok {
		panic(errUnknownType)
		// return typ{}, errUnknownType
	}
	return typ{kind: name}, nil
}

func unmarshalString[T any](s string) (any, error) {
	var x T
	err := json.Unmarshal([]byte(s), &x)
	return x, err
}

var unmarshal = map[typeKind]func(s string) (any, error){
	kindNumber:  unmarshalString[float64],
	kindBoolean: unmarshalString[bool],
	kindString:  func(s string) (any, error) { return s, nil },
	kindObject:  func(s string) (any, error) { panic("not reachable") },
	kindArray:   func(s string) (any, error) { panic("not reachable") },
	kindAny:     func(s string) (any, error) { panic("not reachble") },
}

func (p *parser) array() ([]any, error) {
	var (
		ret []any
	)
	for s, ok := p.next(); ok; s, ok = p.next() {
		if s == "]" {
			return ret, nil
		}
		v, err := p.value(s, typ{kind: kindAny})
		if err != nil {
			panic(err)
		}
		ret = append(ret, v)
	}
	panic(errUnexpectedEOF)
	// return nil, errUnexpectedEOF
}

func (p *parser) fixedArray(typeVar typeKind, length int) ([]any, error) {
	var (
		ret []any
	)
	for i := 0; i < length; i++ {
		s, ok := p.next()
		if !ok {
			panic(errUnexpectedEOF)
		}
		value, err := p.value(s, typ{kind: typeVar})
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
