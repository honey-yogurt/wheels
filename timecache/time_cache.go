package timecache

import "time"

// Strategy is the TimeCache expiration strategy to use.
type Strategy uint8

const (
	// StrategyFirstSeen expires an entry from the time it was added.
	StrategyFirstSeen Strategy = iota
	// StrategyLastSeen expires an entry from the last time it was touched by an Add or Has.
	StrategyLastSeen
)

// TimeCache is a cache of recently seen messages (by id).
type TimeCache interface {
	// Add adds an id into the cache, if it is not already there.
	// Returns true if the id was newly added to the cache.
	// Depending on the implementation strategy, it may or may not update the expiry of
	// an existing entry.
	Add(string) bool
	// Has checks the cache for the presence of an id.
	// Depending on the implementation strategy, it may or may not update the expiry of
	// an existing entry.
	Has(string) bool
	// Done signals that the user is done with this cache, which it may stop background threads
	// and relinquish resources.
	Done()
}

// NewTimeCache defaults to the original ("first seen") cache implementation
func NewTimeCache(ttl time.Duration) TimeCache {
	return NewTimeCacheWithStrategy(StrategyFirstSeen, ttl)
}

func NewTimeCacheWithStrategy(strategy Strategy, ttl time.Duration) TimeCache {
	switch strategy {
	case StrategyFirstSeen:
		return newFirstSeenCache(ttl)
	case StrategyLastSeen:
		return newLastSeenCache(ttl)
	default:
		// Default to the original time cache implementation
		return newFirstSeenCache(ttl)
	}
}
