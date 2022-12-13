package zipkin

// MiddlewareType A special type about middleware
type MiddlewareType string

// decorate a special type Span by tag
const (
	MiddlewareTag string         = "component.type"
	MySQL         MiddlewareType = "mysql"
	Redis         MiddlewareType = "redis"
	ElasticSearch MiddlewareType = "elasticsearch"
	Kafka         MiddlewareType = "kafka"
	RabbitMQ      MiddlewareType = "rabbitmq"
	MongoDB       MiddlewareType = "mongodb"
)

// TagValue return the middleware tag value for decorate Span
func (t MiddlewareType) TagValue() string {
	switch t {
	case MySQL:
		return "database" //for datastore ui
	case Redis, ElasticSearch, Kafka, RabbitMQ, MongoDB:
		return string(t)
	default:
		return "unknow"
	}
}
