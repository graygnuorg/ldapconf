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
	"reflect"
	"testing"
)

func TestText(t *testing.T) {
	inputText := `# Initial comment

  # One more comment
base dc=example,dc=com



use_tls	true
port	330
uri	ldap://a	ldapi://b
`
	type Input struct {
		Base	string	`ldap:"base"`
		UseTLS	bool	`ldap:"use_tls"`
		Port	int16	`ldap:"port"`
		URI	[]string	`ldap:"uri"`
	}
	var in Input
	if err := Unmarshal([]byte(inputText), &in); err != nil {
		t.Fatalf("error: %v", err)
	}
	if !reflect.DeepEqual(in,
		Input{
			Base: `dc=example,dc=com`,
			UseTLS: true,
			Port: 330,
			URI: []string{`ldap://a`, `ldapi://b`},
		}) {
		t.Fatalf("expectation failed: %#v\n", in)
	}
}
