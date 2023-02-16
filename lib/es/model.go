package es

type (
	// Search API Response
	ResponseSearch struct {
		Took     int64          `json:"took"`
		TimedOut bool           `json:"timed_out"`
		Shards   ResponseShards `json:"_shards"`
		Hits     ResponseHits   `json:"hits"`
	}
	ResponseShards struct{}
	ResponseHits   struct {
		Total    ResponseTotal    `json:"total"`
		MaxScore interface{}      `json:"max_score"`
		Hits     []ResponseRecord `json:"hits"`
	}
	ResponseTotal struct {
		Value    int64  `json:"value"`
		Relation string `json:"relation"`
	}
	ResponseRecord struct {
		Index  string      `json:"_index"`
		ID     string      `json:"_id"`
		Source interface{} `json:"_source"`
	}
)
