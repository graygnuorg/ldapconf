# ldapconf

This package provides Go API for unmarshalling `ldap.conf` configuration
files into Go objects.

# Declaring Go structures

Use `ldap:"kw"` tags to mark fields that should be unmarshalled from
certain LDAP keywords, e.g.:

```golang
type LDAPCONF struct {
	Base	string	`ldap:"base"`
	URI	string	`ldap:"uri"`
	...
}
```

Keywords in `ldap.conf` file are case-insensitive.

# Reading from a file.

Use `ldap.Parse` function as shown in example below:

```golang
import "github.com/graygnuorg/ldapconf"

var cfg LDAPCONF
err := ldapconf.Parse(`ldap.conf`, &cfg)
```

# Unmarshalling a string

If the entire `ldap.conf` content is available in a string (or byte
array) variable, use the traditional __golang__ approach:

```golang
var cfg LDAPCONF
err := ldapconf.Unmarshal(string(`URI ldap://`), &cfg)
```
