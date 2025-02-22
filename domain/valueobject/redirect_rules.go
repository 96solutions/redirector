// Package valueobject contains immutable value objects that represent business concepts.
// These objects are defined by their attributes and are considered equal when all their attributes match.
package valueobject

// Predefined redirect types that determine how traffic should be handled.
const (
	// LinkRedirectType indicates a direct URL redirect.
	LinkRedirectType = "link"
	// SlugRedirectType indicates a redirect to another tracking link by slug.
	SlugRedirectType = "slug"
	// SmartSlugRedirectType indicates a redirect to one of multiple tracking links.
	SmartSlugRedirectType = "smart"
	// NoRedirectType indicates that traffic should be blocked.
	NoRedirectType = "block"
	// NoClickType indicates that no click should be recorded.
	NoClickType = "no-click"
)

// RedirectRules type describes a set of options used to handle traffic redirect correctly
// in case traffic doesn't satisfy campaign requirements.
type RedirectRules struct {
	// RedirectType determines how the redirect should be handled.
	// Valid values are defined in the constants above.
	RedirectType string

	// RedirectURL is the target URL for LinkRedirectType.
	RedirectURL string
	// RedirectSlug is the target tracking link slug for SlugRedirectType.
	RedirectSlug string
	// RedirectSmartSlug is a list of possible target slugs for SmartSlugRedirectType.
	RedirectSmartSlug []string
}
