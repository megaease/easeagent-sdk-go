package zipkin

type MiddlewareType string

const (
	MIDDLEWARE_TAG string         = "component.type"
	MySql          MiddlewareType = "mysql"
	Redis          MiddlewareType = "redis"
	ElasticSearch  MiddlewareType = "elasticsearch"
	Kafka          MiddlewareType = "kafka"
	RabbitMQ       MiddlewareType = "rabbitmq"
	MongoDB        MiddlewareType = "mongodb"
)

func (t MiddlewareType) TagValue() string {
	switch t {
	case MySql:
		return "database"
	case Redis, ElasticSearch, Kafka, RabbitMQ, MongoDB:
		return string(t)
	default:
		return "unknow"
	}
}
