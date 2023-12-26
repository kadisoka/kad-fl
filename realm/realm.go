// package realm provides realm-related functionalities.
//
// A realm is a kind of a domain or a site. In general, it has the same
// definition as HTTP's definition of realm.
//
// The concept is closely related to https://schema.org/OnlineBusiness
//
// A realm can be comprised of multiple services.
package realm

import (
	"github.com/rez-go/stev"

	"github.com/alloyzeus/go-azfl/v2/errors"
)

const EnvVarsPrefixDefault = "REALM_"

const (
	NameDefault                                   = "UNKNOWN Realm"
	URLDefault                                    = "https://github.com/kadisoka"
	ContactEmailAddressDefault                    = "nop@example.com"
	ServiceNotificationEmailsSenderAddressDefault = "no-reply@example.com"
	DeveloperNameDefault                          = "Team Kadisoka"
	MaintainerNameDefault                         = DeveloperNameDefault
)

func DefaultInfo() Info {
	return Info{
		Name:                                   NameDefault,
		URL:                                    URLDefault,
		Contact:                                ContactInfo{EmailAddress: ContactEmailAddressDefault},
		ServiceNotificationEmailsSenderAddress: ServiceNotificationEmailsSenderAddressDefault,
		System: SystemInfo{
			Developer:  DeveloperInfo{Name: DeveloperNameDefault},
			Maintainer: MaintainerInfo{Name: MaintainerNameDefault},
		},
	}
}

// Info holds information about a realm.
type Info struct {
	// Name of the realm
	Name string
	// Description contains a brief textual information about the realm.
	Description string
	// Canonical URL of the realm. This URL usually refers to the website
	// of the realm.
	URL string

	TermsOfServiceURL string
	PrivacyPolicyURL  string

	Contact ContactInfo

	ServiceNotificationEmailsSenderAddress string

	System SystemInfo
}

// InfoZero returns a zero-valued Info.
func InfoZero() Info { return Info{} }

func InfoFromEnvOrDefault(envVarsPrefix string) (Info, error) {
	defaultInfo := DefaultInfo()
	return InfoFromEnv(envVarsPrefix, &defaultInfo)
}

// InfoFromEnv creates an instance of Info with values looked up from their
// respective environment variables which names are prefixed with the value
// of envVarsPrefix. The argument defaultInfo serves as the the default, that's
// it, a value in defaultInfo will be preserved unless it's overriden
// by the value specified in environment variables.
func InfoFromEnv(envVarsPrefix string, defaultInfo *Info) (Info, error) {
	if envVarsPrefix == "" {
		envVarsPrefix = EnvVarsPrefixDefault
	}
	if defaultInfo == nil {
		info := InfoZero()
		defaultInfo = &info
	}
	err := stev.LoadFromEnv(envVarsPrefix, defaultInfo)
	if err != nil {
		return InfoZero(), errors.Op("info loading from environment variables", err)
	}
	return *defaultInfo, nil
}

// TODO: person or organization.
// https://schema.org/ContactPoint
type ContactInfo struct {
	EmailAddress string
}

type SystemInfo struct {
	Developer  DeveloperInfo
	Maintainer MaintainerInfo
}

// TODO: person or organization.
type DeveloperInfo struct {
	Name string
}

// TODO: person or organization.
type MaintainerInfo struct {
	Name string
}
