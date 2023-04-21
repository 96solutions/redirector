package valueobject

const LinkRedirectType = "link"
const SlugRedirectType = "slug"
const SmartSlugRedirectType = "smart"
const NoRedirectType = "block"

// RedirectRules type describes a set of options used to handle traffic redirect correctly
//in case traffic doesn't satisfy campaign requirements.
type RedirectRules struct {
	RedirectType string

	RedirectURL       string
	RedirectSlug      string
	RedirectSmartSlug []string
}
