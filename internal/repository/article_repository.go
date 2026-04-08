// 未來 article_repository.go 裡，為了支援 tag filter 而寫的 SQL 查詢邏輯。

// 找 practical tag 對應到哪些 article」這件事，就是發生在 GetArticles(...)
// 然後裡面會有 SQL，概念上在做：
// article join article_tags
// article_tags join tags
// where tags.slug = 'practical'

// ListArticles(...)
// GetPublishedArticles(...)

// 有哪些 article_id 出現在 article_tags，且它對應到的 tag_id 是 practical 那個 tag
// 真正對應的是 repository 裡的 SQL。
package repository
