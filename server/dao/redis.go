package dao

import (
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"time"
)

var (
	WordToArticleCache  *cache.Cache
	WordToAuthorCache   *cache.Cache
	ArticleCache        *cache.Cache
	TitleCache          *cache.Cache
	AuthorCache         *cache.Cache
	NameCache           *cache.Cache
	BookCache           *cache.Cache
	JournalCache        *cache.Cache
	ArticleWordCntCache *cache.Cache
	AuthorWordCntCache  *cache.Cache

	ArticleResRDB *redis.Client

	ArticleToAuthorRDB *redis.Client
	AuthorToArticleRDB *redis.Client
)

const (
	numWordToArticleRDB = iota
	numWordToAuthorRDB
	numArticleCacheRDB
	numTitleCacheRDB
	numAuthorCacheRDB
	numNameCacheRDB
	numBookCacheRDB
	numJournalCacheRDB
	numArticleWordCntRDB
	numAuthorWordCntRDB

	numArticleResRDB

	numArticleToAuthorRDB
	numAuthorToArticleRDB
)

func init() {
	WordToArticleCache = cache.New(&cache.Options{
		Redis: redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    "mymaster",
			SentinelAddrs: []string{":17000", ":17001", ":17002"},
			DB:            numWordToArticleRDB,
			Password:      "zxc05020519",
		}),
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	WordToAuthorCache = cache.New(&cache.Options{
		Redis: redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    "mymaster",
			SentinelAddrs: []string{":17000", ":17001", ":17002"},
			DB:            numWordToAuthorRDB,
			Password:      "zxc05020519",
		}),
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	ArticleCache = cache.New(&cache.Options{
		Redis: redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    "mymaster",
			SentinelAddrs: []string{":17000", ":17001", ":17002"},
			DB:            numArticleCacheRDB,
			Password:      "zxc05020519",
		}),
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	TitleCache = cache.New(&cache.Options{
		Redis: redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    "mymaster",
			SentinelAddrs: []string{":17000", ":17001", ":17002"},
			DB:            numTitleCacheRDB,
			Password:      "zxc05020519",
		}),
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	AuthorCache = cache.New(&cache.Options{
		Redis: redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    "mymaster",
			SentinelAddrs: []string{":17000", ":17001", ":17002"},
			DB:            numAuthorCacheRDB,
			Password:      "zxc05020519",
		}),
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	NameCache = cache.New(&cache.Options{
		Redis: redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    "mymaster",
			SentinelAddrs: []string{":17000", ":17001", ":17002"},
			DB:            numNameCacheRDB,
			Password:      "zxc05020519",
		}),
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	BookCache = cache.New(&cache.Options{
		Redis: redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    "mymaster",
			SentinelAddrs: []string{":17000", ":17001", ":17002"},
			DB:            numBookCacheRDB,
			Password:      "zxc05020519",
		}),
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	JournalCache = cache.New(&cache.Options{
		Redis: redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    "mymaster",
			SentinelAddrs: []string{":17000", ":17001", ":17002"},
			DB:            numJournalCacheRDB,
			Password:      "zxc05020519",
		}),
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	ArticleWordCntCache = cache.New(&cache.Options{
		Redis: redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    "mymaster",
			SentinelAddrs: []string{":17000", ":17001", ":17002"},
			DB:            numArticleWordCntRDB,
			Password:      "zxc05020519",
		}),
		LocalCache: cache.NewTinyLFU(10000, time.Hour),
	})

	AuthorWordCntCache = cache.New(&cache.Options{
		Redis: redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    "mymaster",
			SentinelAddrs: []string{":17000", ":17001", ":17002"},
			DB:            numAuthorWordCntRDB,
			Password:      "zxc05020519",
		}),
		LocalCache: cache.NewTinyLFU(10000, time.Hour),
	})

	ArticleResRDB = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "mymaster",
		SentinelAddrs: []string{":17000", ":17001", ":17002"},
		DB:            numArticleResRDB,
		Password:      "zxc05020519",
	})

	ArticleToAuthorRDB = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "mymaster",
		SentinelAddrs: []string{":17000", ":17001", ":17002"},
		DB:            numArticleToAuthorRDB,
		Password:      "zxc05020519",
	})

	AuthorToArticleRDB = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    "mymaster",
		SentinelAddrs: []string{":17000", ":17001", ":17002"},
		DB:            numAuthorToArticleRDB,
		Password:      "zxc05020519",
	})
}
