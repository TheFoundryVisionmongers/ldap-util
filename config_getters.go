package main

// LDAPBase will return the LDAPBase set in the config file, CLI or environment variables.
func LDAPBase() string {
	return loadedConfig.LDAPConfig.Base
}

// LDAPBindPassword will return the LDAPBindPassword set in the config file, CLI or environment variables.
func LDAPBindPassword() string {
	return loadedConfig.LDAPConfig.BindPassword
}

// LDAPBindUser will return the LDAPBindUser set in the config file, CLI or environment variables.
func LDAPBindUser() string {
	return loadedConfig.LDAPConfig.BindUser
}

// LDAPGroupPrefix will return the group prefix set in the config file.
func LDAPGroupPrefix() string {
	return loadedConfig.LDAPConfig.GroupSearch.GroupPrefix
}

// LDAPGroupSuffix will return the group suffix set in the config file.
func LDAPGroupSuffix() string {
	return loadedConfig.LDAPConfig.GroupSearch.GroupSuffix
}

// LDAPHost will return the LDAPHost set in the config file, CLI or environment variables.
func LDAPHost() string {
	return loadedConfig.LDAPConfig.Host
}

// LDAPPort will return the LDAPPort set in the config file, CLI or environment variables.
func LDAPPort() int {
	return loadedConfig.LDAPConfig.Port
}

// LDAPSelfAuth will return whether the user should self-auth with the LDAP server; that is to user the users own
// credentials to bind on the LDAP server.
func LDAPSelfAuth() bool {
	return loadedConfig.LDAPConfig.SelfAuth
}

// LDAPUse will return the LDAPUse set in the config file, CLI or environment variables.
func LDAPUse() bool {
	return loadedConfig.LDAPConfig.UseLDAP
}

// LDAPUseSSL will return the LDAPUseSSL set in the config file, CLI or environment variables.
func LDAPUseSSL() bool {
	return loadedConfig.LDAPConfig.UseSSL
}

// LDAPGroupSearchBaseDN will return the LDAP Group Search BaseDN set in the config file
func LDAPGroupSearchDN() string {
	return loadedConfig.LDAPConfig.GroupSearch.DN
}

// LDAPGroupSearchFilter will return the LDAPGroupSearchFilter set in the config file
func LDAPGroupSearchFilter() string {
	return loadedConfig.LDAPConfig.GroupSearch.Filter
}

// LDAPGroupSearchGroupAttr will return the LDAPGroupSearchGroupAttr set in the config file
func LDAPGroupSearchGroupAttr() string {
	return loadedConfig.LDAPConfig.GroupSearch.GroupAttr
}

// LDAPGroupSearchUserAttr will return the LDAPGroupSearchUserAttr set in the config file
func LDAPGroupSearchUserAttr() string {
	return loadedConfig.LDAPConfig.GroupSearch.UserAttr
}

// LDAPGroupSearchNameAttr will return the LDAPGroupSearchNameAttr set in the config file
func LDAPGroupSearchNameAttr() string {
	return loadedConfig.LDAPConfig.GroupSearch.NameAttr
}

// LDAPUserSearchDN is the DN location for users in the directory
func LDAPUserSearchDN() string {
	return loadedConfig.LDAPConfig.UserSearch.DN
}

// LDAPUserSearchUserAttr the attribute for the username in the active directory
func LDAPUserSearchUserAttr() string {
	return loadedConfig.LDAPConfig.UserSearch.UserAttr
}

// LDAPUserSearchUserFilter is the filter for filtering the types of objectClass
// for users
func LDAPUserSearchUserFilter() string {
	return loadedConfig.LDAPConfig.UserSearch.Filter
}

// LDAPUserSearchNameAttr is the name attribute for obtaining the users full name
func LDAPUserSearchNameAttr() string {
	return loadedConfig.LDAPConfig.UserSearch.NameAttr
}
