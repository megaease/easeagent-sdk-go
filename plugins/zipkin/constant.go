package zipkin

//For details, please see 'doc/middleware-span.md'
const (
	HTTPTagAttributeRoute = "http.route"
	HTTPTagMethod         = "http.method"
	HTTPTagPath           = "http.path"
	HTTPTagStatusCode     = "http.status_code"
	HTTPTagScriptFilename = "http.script.filename"
	HTTPTagClientAddress  = "Client Address"

	MysqlTagSQL = "sql"
	MysqlTagURL = "url"

	RedisTagMethod = "redis.method"

	ElasticsearchTagIndex     = "es.index"
	ElasticsearchTagOperation = "es.operation"
	ElasticsearchTagBody      = "es.body"

	KafkaTagTopic  = "kafka.topic"
	KafkaTagKey    = "kafka.key"
	KafkaTagOroker = "kafka.broker"

	RabbitTagExchange   = "rabbit.exchange"
	RabbitTagRoutingKey = "rabbit.routing_key"
	RabbitTagQueue      = "rabbit.queue"
	RabbitTagBroker     = "rabbit.broker"

	MongodbTagCommand    = "mongodb.command"
	MongodbTagCollection = "mongodb.collection"
	MongodbTagClusterID  = "mongodb.cluster_id"
)
