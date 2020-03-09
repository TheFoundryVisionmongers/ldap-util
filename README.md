# LDAP Util

LDAP Util is a Go helper utility app, to test you can bind to an LDAP server using the given
username and password.

It should be run with your flix `config.yml` file in the same directory; alternatively, you can optionally provide the
full path to the config file (this second option is shown below).

It will use the LDAP settings in the config to
attempt to connect, bind, authenticate and search for groups.

## Usage

```bash
$ ./ldap-utils --config-file=/home/ben/flix/flix-server/config.yml
  This test application will read and output all LDAP data from your Flix config file.  This may include sensitive passwords.  Do you wish to continue? [y/N] y
  Loading config from /home/ben/flix/flix-server/config.yml
  Loaded config: {{dc=flix,dc=local admin uid=cgornea,ou=Users,dc=flix,dc=local {ou=Groups,dc=flix,dc=local (objectClass=posixGroup) uid memberUid cn flix_ } 10.0.145.222 389 false true {ou=Users,dc=flix,dc=local (objectClass=organizationalPerson) uid givenName} false}}
  Attempting to dial the LDAP server at 10.0.145.222:389...
  ✓ Connected to LDAP Server
  Please enter the username you'd like to log into the LDAP server with: ben
  Please enter the password you'd like to log into the LDAP server with. Note: this value will be printed to the console and may be output during the testing process: admin
  Will log into the LDAP server as ben:admin
  Attempting to bind on LDAP server
  Self auth is set to false, will attempt to bind with bind user set in config
  ✓ Established bind user - uid=cgornea,ou=Users,dc=flix,dc=local:admin
  ✓ Bind Complete
  Searching for LDAP Users
  	DN: ou=Users,dc=flix,dc=local
  	Filter: (&(objectClass=organizationalPerson)(uid=ben))
  	 Attributes [dn cn uid displayName mail ou givenName]
  Binding as returned user to verify credentials - uid=ben,ou=Users,dc=flix,dc=local:admin
  Got user details:
  	Username: ben
  	CN: Ben Cragg
  	Mail: 
  	Display name: 
  UserSearchNameAttr set, overriding display name to Ben
  Finding users groups, user search attr: ben
  Searching for LDAP groups:
  	BaseDN: ou=Groups,dc=flix,dc=local
  	Filter: (&(objectClass=posixGroup)(memberUid=ben))
  	Attributes: [cn]
  Found 3 groups for user
  ✓ Got group names: [foundry flix_foundry foundry_flix]
  Closing connection to LDAP server
  ✓ Finished.
```
