# MegaCloud Configuration

Modify the [agent.yml](./agent.yml) file to configure your information.

## 1. Name

You'll service_name to find your data later. It's important to use a unique and meaningful name.

The service_name of megacloud consists of three parts: zone, domain, name. They are joined by `.` into `ServiceName`

```yaml
service_name: zone.domain.service
```

## 2. Reporter

MegaCloud uses HTTP to receive data, so you need to change the configuration to HTTP and MegaCloud's address.
```yaml
reporter.output.server: {MEGA_CLOUD_URL}/application-tracing-log
```
## 3. MTLS

MTLS is a secure authentication protocol for EaseAgent to connect to MegaCloud.

Config: Get TLS
```yaml
reporter.output.server.tls.enable: true
reporter.output.server.tls.key: YOUR_TLS_KEY
reporter.output.server.tls.cert: YOUR_TLS_CERT
```

## 4. Tracing
we have some
```yaml
tracing_type: log-tracing
tracing.enable: true
tracing.sample.rate: 1.0
```

## Third: About `MEGA_CLOUD_URL` And `TLS`

When you download the `agent.yaml` file through our megacloud, `MEGA_CLOUD_URL` and `TLS` will be filled in for you automatically.

If you need it separately, please download the `agent.yaml` and get it by yourself.