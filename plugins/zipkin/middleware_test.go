package zipkin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagValue(t *testing.T) {
	assert.Equal(t, "database", MySql.TagValue())
	assert.Equal(t, "redis", Redis.TagValue())
	assert.Equal(t, "elasticsearch", ElasticSearch.TagValue())
	assert.Equal(t, "kafka", Kafka.TagValue())
	assert.Equal(t, "rabbitmq", RabbitMQ.TagValue())
	assert.Equal(t, "mongodb", MongoDB.TagValue())
	assert.Equal(t, "unknow", MiddlewareType("aaa").TagValue())
}
