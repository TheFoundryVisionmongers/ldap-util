package main

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func init() {
	pwd, err := os.Executable()
	if err != nil {
		panic(err)
	}
	configPath = filepath.Join(filepath.Dir(pwd), "config.yml")
}

var (
	configPath string

	loadedConfig config
)

type config struct {
	// The ldap authentication configuration
	LDAPConfig ldapConfig `yaml:"ldap,omitempty"`
}

// ldapConfig is the configuration struct for LDAP Authentication settings
type ldapConfig struct {
	Base            string                `yaml:"base"`
	BindPassword    string                `yaml:"bind_pass"`
	BindUser        string                `yaml:"bind_user"`
	GroupSearch     ldapGroupSearchConfig `yaml:"group_search"`
	Host            string                `yaml:"host"`
	Port            int                   `yaml:"port"`
	SelfAuth        bool                  `yaml:"self_auth"`
	UseLDAP         bool                  `yaml:"use_ldap"`
	UserSearch      ldapUserSearchConfig  `yaml:"user_search"`
	UseSSL          bool                  `yaml:"use_ssl"`
	Roles           []ldapRolesConfig     `yaml:"roles"`
	AllowedUsers    []string              `yaml:"allowed_users"`
	ForbiddenUsers  []string              `yaml:"blocked_users"`
	RequiredGroups  []string              `yaml:"required_groups"`
	ForbiddenGroups []string              `yaml:"forbidden_groups"`
}

type ldapGroupSearchConfig struct {
	// DN to start the search from. For example "cn=Groups,dc=example,dc=com"
	DN string `yaml:"dn"`
	// Optional filter to apply when searching the directory. For example "(objectClass=posixGroup)"
	Filter string `yaml:"filter"`
	// These two fields are use to match a user to a group.
	//
	// It adds an additional requirement to the filter that an attribute in the group
	// match the user's attribute value. For example that the "members" attribute of
	// a group matches the "uid" of the user. The exact filter being added is:
	//
	//   (<groupAttr>=<userAttr value>)
	//
	UserAttr  string `yaml:"user_attr"`
	GroupAttr string `yaml:"group_attr"`
	// The attribute of the group that represents its name.
	NameAttr    string `yaml:"name_attr"`
	GroupPrefix string `yaml:"group_prefix"`
	GroupSuffix string `yaml:"group_suffix"`
}

type ldapRolesConfig struct {
	// GroupPrefix specifies the prefix of the group to match.
	GroupPrefix string `yaml:"group_prefix"`
	// GroupSuffix specifies the suffix of the group to match.
	GroupSuffix string `yaml:"group_suffix"`
	// NameAttr is the attribute of the group to match. If empty, defaults to NameAttr of the group search config.
	NameAttr string `yaml:"name_attr"`
	// Role is the role to use for users in groups matched by Pattern.
	Role string `yaml:"role"`
}

type ldapUserSearchConfig struct {
	// DN to start the search from. For example "cn=People,dc=example,dc=com"
	DN string `yaml:"dn"`
	// Optional filter to apply when searching the directory. For example "(objectClass=organizationalPerson)"
	Filter string `yaml:"filter"`
	// UserAttr is the attribute of the filter to match users by
	UserAttr string `yaml:"user_attr"`
	// NameAttr is the attribute of the user that represents its name.
	NameAttr string `yaml:"name_attr"`
}

func processYML() error {
	f, err := os.OpenFile(configPath, os.O_RDONLY, 0444)
	if err != nil {
		return err
	}
	defer f.Close()

	return yaml.NewDecoder(f).Decode(&loadedConfig)
}
