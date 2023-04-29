package valueobject

// UserAgent type describes a set of data parsed from the User-Agent header.
type UserAgent struct {
	SrcString string

	Bot bool

	Device   string
	Platform string
	Browser  string
	//Version  string
}
