/* This file is part of ldapconf
   Copyright (C) 2026 Sergey Poznyakoff

   Ldapconf is free software; you can redistribute it and/or modify it
   under the terms of the GNU General Public License as published by the
   Free Software Foundation; either version 3 of the License, or (at your
   option) any later version.

   Ldapconf is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License along
   with ldapconf. If not, see <http://www.gnu.org/licenses/>. */

package ldapconf

import (
	"bufio"
	"bytes"
	"encoding"
	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	wsRe = regexp.MustCompile(`\s+`)
)

func parseConfig(input []byte) map[string]string {
	cf := map[string]string{}
	scanner := bufio.NewScanner(bytes.NewReader(input))
	for scanner.Scan() {
		s := strings.Trim(scanner.Text(), " \t")
		if strings.HasPrefix(s, "#") {
			continue
		}
		a := wsRe.Split(s, 2)
		if len(a) != 2 {
			continue
		}
		cf[strings.ToLower(a[0])] = a[1]
	}
	return cf
}

func stringToBool(s string) bool {
	switch strings.ToLower(s) {
	case `0`, `false`, `off`, `no`:
		return false
	default:
		return true
	}
}

func unmarshalValue(kind reflect.Kind, fval reflect.Value, cval string) error {
	var pval reflect.Value
	if fval.Kind() != reflect.Ptr {
		pval = fval.Addr()
	} else {
		if fval.IsNil() {
			fval.Set(reflect.New(fval.Type().Elem()))
		} 
		pval = fval
	}
	if pval.CanInterface() {
		if um, ok := pval.Interface().(encoding.TextUnmarshaler); ok {
			return um.UnmarshalText([]byte(cval))
		}
	}

	switch kind {
	case reflect.Bool:
		fval.SetBool(stringToBool(cval))
	case reflect.Int:
		if n, err := strconv.ParseInt(cval, 10, 0); err != nil {
			return err
		} else {
			fval.SetInt(n)
		}

	case reflect.Int8:
		if n, err := strconv.ParseInt(cval, 10, 8); err != nil {
			return err
		} else {
			fval.SetInt(n)
		}

	case reflect.Int16:
		if n, err := strconv.ParseInt(cval, 10, 16); err != nil {
			return err
		} else {
			fval.SetInt(n)
		}

	case reflect.Int32:
		if n, err := strconv.ParseInt(cval, 10, 32); err != nil {
			return err
		} else {
			fval.SetInt(n)
		}

	case reflect.Int64:
		if n, err := strconv.ParseInt(cval, 10, 64); err != nil {
			return err
		} else {
			fval.SetInt(n)
		}

	case reflect.Uint:
		if n, err := strconv.ParseUint(cval, 10, 0); err != nil {
			return err
		} else {
			fval.SetUint(n)
		}

	case reflect.Uint8:
		if n, err := strconv.ParseUint(cval, 10, 8); err != nil {
			return err
		} else {
			fval.SetUint(n)
		}

	case reflect.Uint16:
		if n, err := strconv.ParseUint(cval, 10, 16); err != nil {
			return err
		} else {
			fval.SetUint(n)
		}

	case reflect.Uint32:
		if n, err := strconv.ParseUint(cval, 10, 32); err != nil {
			return err
		} else {
			fval.SetUint(n)
		}

	case reflect.Uint64:
		if n, err := strconv.ParseUint(cval, 10, 64); err != nil {
			return err
		} else {
			fval.SetUint(n)
		}

	case reflect.Float32:
		if n, err := strconv.ParseFloat(cval, 32); err != nil {
			return err
		} else {
			fval.SetFloat(n)
		}

	case reflect.Float64:
		if n, err := strconv.ParseFloat(cval, 64); err != nil {
			return err
		} else {
			fval.SetFloat(n)
		}

	case reflect.String:
		fval.SetString(cval)

	case reflect.Slice:
		if err := unmarshalSlice(fval, cval); err != nil {
			return err
		}
		
	default:
		return fmt.Errorf("unsupported field type: %v", kind)

		/* Following types are not supported:

			Struct
			Array
			Uintptr
			Complex64
			Complex128
			Chan
			Func
			Interface
			Map
			Pointer
			UnsafePointer
			*/
	}
	return nil
}

func unmarshalSlice(fval reflect.Value, cval string) error {
	eltype := fval.Type().Elem()
	kind := eltype.Kind()
	fields := strings.Fields(cval)
	l := len(fields)
	result := reflect.MakeSlice(reflect.SliceOf(eltype), l, l)
	
	for i := 0; i < l; i++ {
                elt := reflect.New(eltype).Elem()
		if err := unmarshalValue(kind, elt, fields[i]); err != nil {
        		return err
		}
                result.Index(i).Set(elt)
        }
	fval.Set(result)
	return nil
}

func Unmarshal(in []byte, out any) error {
	cf := parseConfig(in)
	if out == nil {
		return fmt.Errorf("passed object is nil")
	}
	v := reflect.ValueOf(out)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return fmt.Errorf("passed object is nil")
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("passed object is not a pointer to structure")
	}

	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fld := t.Field(i)
		tag := fld.Tag.Get(`ldap`)
		if tag == "" || tag == "-" {
			continue
		}
		fval := v.FieldByIndex([]int{i})
		if cval, ok := cf[strings.ToLower(tag)]; ok {
			if err := unmarshalValue(fld.Type.Kind(), fval, cval); err != nil {
				fmt.Errorf("error parsing %s: %v", fld.Name, err)
			}
		}
	}
	return nil
}

// Parse reads configuration file `filename` and unmarshals it to
// the output object `out`.
func Parse(filename string, out any) error {
	content, err := ioutil.ReadFile(filename)
	if err == nil {
		return Unmarshal(content, out)
	}
	return err
}
