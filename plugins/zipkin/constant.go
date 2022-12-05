package zipkin

//For details, please see 'doc/middleware-span.md'
const (
	HTTP_TAG_ATTRIBUTE_ROUTE = "http.route"
	HTTP_TAG_METHOD          = "http.method"
	HTTP_TAG_PATH            = "http.path"
	HTTP_TAG_STATUS_CODE     = "http.status_code"
	HTTP_TAG_SCRIPT_FILENAME = "http.script.filename"
	HTTP_TAG_CLIENT_ADDRESS  = "Client Address"

	MYSQL_TAG_SQL = "sql"
	MYSQL_TAG_URL = "url"

	REDIS_TAG_METHOD = "redis.method"

	ELASTICSEARCH_TAG_INDEX     = "es.index"
	ELASTICSEARCH_TAG_OPERATION = "es.operation"
	ELASTICSEARCH_TAG_BODY      = "es.body"

	KAFKA_TAG_TOPIC  = "kafka.topic"
	KAFKA_TAG_KEY    = "kafka.key"
	KAFKA_TAG_BROKER = "kafka.broker"

	RABBIT_TAG_EXCHANGE    = "rabbit.exchange"
	RABBIT_TAG_ROUTING_KEY = "rabbit.routing_key"
	RABBIT_TAG_QUEUE       = "rabbit.queue"
	RABBIT_TAG_BROKER      = "rabbit.broker"

	MONGODB_TAG_COMMAND    = "mongodb.command"
	MONGODB_TAG_COLLECTION = "mongodb.collection"
	MONGODB_TAG_CLUSTER_ID = "mongodb.cluster_id"
)
