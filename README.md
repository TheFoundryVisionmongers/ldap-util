# LDAP Util

LDAP Util is a Go helper utility app, to test you can bind to an LDAP server using the given
username and password.

## Usage

```bash
./ldap-util -h=<hostname> -bindUser=uid=username,ou=Users,dc=company,dc=com -bindPass=<password>
```