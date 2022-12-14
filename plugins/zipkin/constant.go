/**
 * Copyright 2022 MegaEase
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package zipkin

// For details, please see https://github.com/megaease/easeagent-sdk-go/blob/main/doc/middleware-span.md
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
