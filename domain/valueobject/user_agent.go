// Package valueobject contains immutable value objects that represent business concepts.
// These objects are defined by their attributes and are considered equal when all their attributes match.
package valueobject

// UserAgent type describes a set of data parsed from the User-Agent header.
type UserAgent struct {
	// SrcString contains the original User-Agent header string.
	SrcString string

	// Bot indicates whether the request comes from a bot/crawler.
	Bot bool

	// Device represents the type of device (desktop, mobile, tablet, etc.).
	Device string
	// Platform represents the operating system (Windows, iOS, Android, etc.).
	Platform string
	// Browser represents the web browser name (Chrome, Firefox, Safari, etc.).
	Browser string
	// Version represents the browser version (currently unused).
	// Version  string
}
