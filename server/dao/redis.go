package dao

import (
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"time"
)

var (
	CommonCache        *cache.Cache
	WordToArticleCache *cache.Cache
	WordToAuthorCache  *cache.Cache
	ArticleCache       *cache.Cache
	AuthorCache        *cache.Cache
	BookCache          *cache.Cache
	JournalCache       *cache.Cache
	ArticleToAuthorRDB *redis.Client
	AuthorToArticleRDB *redis.Client
)

const (
	numCommonRDB = iota
	numWordToArticleRDB
	numWordToAuthorRDB
	numArticleCacheRDB
	numAuthorCacheRDB
	numBookCacheRDB
	numJournalCacheRDB
	numArticleToAuthorRDB
	numAuthorToArticleRDB
)

func init() {
	CommonCache = cache.New(&cache.Options{
		Redis: redis.NewFailoverClient(&redis.FailoverOptions{
			SentinelAddrs: []string{":17000", ":17001", ":17002"},
			DB:            numCommonRDB,
			Password:      "zxc05020519",
		}),
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	WordToArticleCache = cache.New(&cache.Options{
		Redis: redis.NewFailoverClient(&redis.FailoverOptions{
			SentinelAddrs: []string{":17000", ":17001", ":17002"},
			DB:            numWordToArticleRDB,
			Password:      "zxc05020519",
		}),
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	WordToAuthorCache = cache.New(&cache.Options{
		Redis: redis.NewFailoverClient(&redis.FailoverOptions{
			SentinelAddrs: []string{":17000", ":17001", ":17002"},
			DB:            numWordToAuthorRDB,
			Password:      "zxc05020519",
		}),
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	ArticleCache = cache.New(&cache.Options{
		Redis: redis.NewFailoverClient(&redis.FailoverOptions{
			SentinelAddrs: []string{":17000", ":17001", ":17002"},
			DB:            numArticleCacheRDB,
			Password:      "zxc05020519",
		}),
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	AuthorCache = cache.New(&cache.Options{
		Redis: redis.NewFailoverClient(&redis.FailoverOptions{
			SentinelAddrs: []string{":17000", ":17001", ":17002"},
			DB:            numAuthorCacheRDB,
			Password:      "zxc05020519",
		}),
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	BookCache = cache.New(&cache.Options{
		Redis: redis.NewFailoverClient(&redis.FailoverOptions{
			SentinelAddrs: []string{":17000", ":17001", ":17002"},
			DB:            numBookCacheRDB,
			Password:      "zxc05020519",
		}),
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	JournalCache = cache.New(&cache.Options{
		Redis: redis.NewFailoverClient(&redis.FailoverOptions{
			SentinelAddrs: []string{":17000", ":17001", ":17002"},
			DB:            numJournalCacheRDB,
			Password:      "zxc05020519",
		}),
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	ArticleToAuthorRDB = redis.NewFailoverClient(&redis.FailoverOptions{
		SentinelAddrs: []string{":17000", ":17001", ":17002"},
		DB:            numArticleToAuthorRDB,
		Password:      "zxc05020519",
	})

	AuthorToArticleRDB = redis.NewFailoverClient(&redis.FailoverOptions{
		SentinelAddrs: []string{":17000", ":17001", ":17002"},
		DB:            numAuthorToArticleRDB,
		Password:      "zxc05020519",
	})
}
