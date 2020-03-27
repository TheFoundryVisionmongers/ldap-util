package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"gopkg.in/ldap.v2"
)

// bold wraps the input string in special characters which make text bold in a bash terminal.
func bold(input string) string {
	// \033[1m starts printing in bold, \033[0m stops the bold
	return fmt.Sprintf("\033[1m%s\033[0m ", input)
}

// getResponse will get a response from the user, it uses the 'validate' function to check if the input is valid.  It
// will repeatedly ask the user for further input if the input is invalid.
func getResponse(prompt string, validate func(string) bool) string {
	var resp string

	for {
		fmt.Print(bold(prompt))
		n, err := fmt.Scanln(&resp)

		// If n==0, the user entered no input, ignore error and just validate the response
		if (err != nil && n != 0) || !validate(resp) {
			fmt.Println("Invalid input, please try again.")
			continue
		}
		break
	}

	return resp
}

func getUsernamePass() (string, string) {
	username := getResponse("Please enter the username you'd like to log into the LDAP server with:", func(s string) bool { return s != "" })
	password := getResponse("Please enter the password you'd like to log into the LDAP server with. Note: this value will be printed to the console and may be output during the testing process:", func(_ string) bool { return true })

	return username, password
}

func getAttrs(e ldap.Entry, name string) []string {
	for _, a := range e.Attributes {
		if a.Name != name {
			continue
		}
		return a.Values
	}
	if name == "dn" || name == "DN" {
		return []string{e.DN}
	}
	return nil
}

func getAttr(e ldap.Entry, name string) string {
	if a := getAttrs(e, name); len(a) > 0 {
		return a[0]
	}
	return ""
}

// groups will search the ldap for the given users groups membership
// using the configuration group search params. it will then return the found groups
// or an error
func groups(user ldap.Entry, l *ldap.Conn) ([]string, error) {
	var groups []*ldap.Entry
	var groupNames []string

	for _, attr := range getAttrs(user, LDAPGroupSearchUserAttr()) {
		fmt.Printf("Finding users groups, user search attr: %s\n", attr)

		filter := fmt.Sprintf("(%s=%s)", LDAPGroupSearchGroupAttr(), ldap.EscapeFilter(attr))
		if LDAPGroupSearchFilter() != "" {
			filter = fmt.Sprintf("(&%s%s)", LDAPGroupSearchFilter(), filter)
		}

		baseDN := LDAPBase()
		if LDAPGroupSearchDN() != "" {
			baseDN = LDAPGroupSearchDN()
		}

		// Search for the given users groups
		req := &ldap.SearchRequest{
			BaseDN:     baseDN,
			Filter:     filter,
			Scope:      ldap.ScopeWholeSubtree,
			Attributes: []string{LDAPGroupSearchNameAttr()},
		}

		fmt.Printf("Searching for LDAP groups:\n\tBaseDN: %s\n\tFilter: %s\n\tAttributes: %v\n", req.BaseDN, req.Filter, req.Attributes)

		sr, err := l.Search(req)
		if err != nil {
			fmt.Printf("Could not search active directory: %v", err)
			return nil, err
		}

		groups = append(groups, sr.Entries...)

		if len(groups) == 0 {
			fmt.Println("Warning: No groups found from LDAP search")
		} else {
			for _, group := range groups {
				name := getAttr(*group, LDAPGroupSearchNameAttr())
				if name != "" {
					groupNames = append(groupNames, name)
				}
			}

			fmt.Printf("Found %d groups for user\n", len(groupNames))
		}
	}

	return groupNames, nil
}

func main() {
	configFlag := flag.String("config-file", configPath, "The full path to the YAML config file")
	flag.Parse()

	if getResponse("This test application will read and output all LDAP data from your Flix config file.  This may include sensitive passwords.  Do you wish to continue? [y/N]", func(_ string) bool { return true }) != "y" {
		fmt.Println("Exiting test application due to user response.")
		os.Exit(0)
	}

	if configFlag != nil && *configFlag != "" {
		configPath = *configFlag
	}

	fmt.Printf("Loading config from %s\n", configPath)

	if err := processYML(); err != nil {
		fmt.Printf("Unable to load config from file\n")
		os.Exit(1)
	}

	if !LDAPUse() {
		fmt.Println("UseLDAP is set to false, cannot continue test")
		os.Exit(0)
	}

	fmt.Printf("Loaded config: %v\n", loadedConfig)

	// Use default port if not set
	ldapPort := LDAPPort()
	if ldapPort == 0 {
		fmt.Println("LDAP port not set in config, will use 389 as the port")
		ldapPort = 389
	}
	addr := net.JoinHostPort(LDAPHost(), strconv.Itoa(ldapPort))

	fmt.Printf("Attempting to dial the LDAP server at %s...\n", addr)
	l, err := ldap.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("Failed to dial the LDAP server: %v\n", err)
		os.Exit(1)
	}
	defer l.Close()

	fmt.Println("✓ Connected to LDAP Server")

	// Reconnect with TLS
	if LDAPUseSSL() {
		fmt.Println("Attempting to start TLS on the LDAP connection...")
		if err = l.StartTLS(&tls.Config{InsecureSkipVerify: true}); err != nil {
			fmt.Printf("Could not start TLS on the LDAP server connection: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✓ Start TLS Complete")
	}

	// decide to user the user search DN if its specified
	dn := LDAPBase()
	if LDAPUserSearchDN() != "" {
		dn = LDAPUserSearchDN()
	}

	username, password := getUsernamePass()

	fmt.Printf("Will log into the LDAP server as %s:%s\n", username, password)

	// Bind on the LDAP server. We only need to bind if the binduser and bind password have
	// been set in config. OR we have self auth set. We dont need to bind if these values are
	// not set.
	if LDAPBindUser() != "" || LDAPBindPassword() != "" || LDAPSelfAuth() {
		fmt.Println("Attempting to bind on LDAP server")

		// the bind users requires a full DN field for it to be able to search for the
		// user in the ldap. These values use the logging in users credentials to
		// set the user and password
		bindUser := fmt.Sprintf("%s=%s,%s", LDAPUserSearchUserAttr(), username, dn)
		bindPass := password

		// if we are not using selfAuth then we just want to use the values stored
		// in the config.
		if !LDAPSelfAuth() {
			// Do not self bind, use the provided readonly user credentials
			fmt.Println("Self auth is set to false, will attempt to bind with bind user set in config")
			bindUser = LDAPBindUser()
			bindPass = LDAPBindPassword()
		}
		fmt.Printf("✓ Established bind user - %s:%s\n", bindUser, bindPass)

		if err = l.Bind(bindUser, bindPass); err != nil {
			fmt.Printf("Could not bind to LDAP server: %v\n", err)
			os.Exit(1)
		}
	}
	fmt.Println("✓ Bind Complete")

	// create the search string
	filter := fmt.Sprintf("(%s=%s)", LDAPUserSearchUserAttr(), username)
	if LDAPUserSearchUserFilter() != "" {
		filter = fmt.Sprintf("(&%s%s)", LDAPUserSearchUserFilter(), filter)
	}

	// attributes
	attributes := []string{"dn", "cn", "uid", "displayName", "mail", "ou"}
	if LDAPUserSearchNameAttr() != "" {
		attributes = append(attributes, LDAPUserSearchNameAttr())
	}
	if LDAPUserSearchUserAttr() != "" {
		attributes = append(attributes, LDAPUserSearchUserAttr())
	}

	// Search for the given username
	req := ldap.NewSearchRequest(
		dn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		attributes,
		nil,
	)

	fmt.Printf("Searching for LDAP Users\n\tDN: %s\n\tFilter: %s\n\t Attributes %v\n", req.BaseDN, req.Filter, req.Attributes)

	sr, err := l.Search(req)
	if err != nil {
		fmt.Printf("Could not search active directory: %v\n", err)
		os.Exit(1)
	}

	if len(sr.Entries) != 1 {
		fmt.Printf("Expected exactly 1 entry returned, got %d entries\n", len(sr.Entries))
		os.Exit(1)
	}

	userdn := sr.Entries[0].DN

	// Bind as the user to verify their password
	fmt.Printf("Binding as returned user to verify credentials - %s:%s\n", userdn, password)
	if err := l.Bind(userdn, password); err != nil {
		fmt.Printf("Credentials mismatch: %v\n", err)
		os.Exit(1)
	}

	usernameAttr := LDAPUserSearchUserAttr()
	if usernameAttr == "" {
		usernameAttr = "uid"
	}

	fmt.Printf(
		"Got user details:\n\tUsername: %s\n\tCN: %s\n\tMail: %s\n\tDisplay name: %s\n",
		sr.Entries[0].GetAttributeValue(usernameAttr),
		sr.Entries[0].GetAttributeValue("cn"),
		sr.Entries[0].GetAttributeValue("mail"),
		sr.Entries[0].GetAttributeValue("displayName"),
	)

	// set the name as the specified search name attribute
	if LDAPUserSearchNameAttr() != "" {
		fmt.Printf("UserSearchNameAttr set, overriding display name to %s\n", sr.Entries[0].GetAttributeValue(LDAPUserSearchNameAttr()))
	}

	groupNames, err := groups(*sr.Entries[0], l)
	if err != nil {
		fmt.Printf("Failed to query groups: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Got group names: [%s]\n", strings.Join(groupNames, ", "))

	fmt.Println("Attempting to filter group names based on prefix and suffix")
	fmt.Printf("\tPrefix: %s\n", LDAPGroupPrefix())
	fmt.Printf("\tSuffix: %s\n", LDAPGroupSuffix())

	var matchedGroupNames []string

	for _, gn := range groupNames {
		fmt.Println("---")

		fmt.Printf("Testing group name '%s'\n", gn)

		if len(gn) < len(LDAPGroupPrefix()) {
			fmt.Println("Group name is too short to match prefix")
			continue
		}
		if len(gn) < len(LDAPGroupSuffix()) {
			fmt.Println("Group name is too short to match suffix")
			continue
		}

		if gn[:len(LDAPGroupPrefix())] != LDAPGroupPrefix() {
			fmt.Println("First part of group name does not match prefix")
			fmt.Printf("\t%s != %s\n", gn[:len(LDAPGroupPrefix())], LDAPGroupPrefix())
			continue
		}
		if gn[len(gn)-len(LDAPGroupSuffix()):] != LDAPGroupSuffix() {
			fmt.Println("Last part of group name does not match suffix")
			fmt.Printf("\t%s != %s\n", gn[:len(LDAPGroupSuffix())], LDAPGroupSuffix())
			continue
		}

		fmt.Println("Group name matches prefix and suffix")
		matchedGroupNames = append(matchedGroupNames, gn)
	}

	fmt.Printf("Matched group names: [%s]\n", strings.Join(matchedGroupNames, ", "))

	fmt.Println("Closing connection to LDAP server")
	l.Close()

	fmt.Println("✓ Finished.")
}
