package cache

type ICache interface {
	CacheTerms(terms []string)
	GetCacheValue(key string) (string, error)
}
