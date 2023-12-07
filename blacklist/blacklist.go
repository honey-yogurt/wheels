package blacklist

// Blacklist is an interface for something blacklisting.
type Blacklist interface {
	Add(string) bool
	Contains(string) bool
}
