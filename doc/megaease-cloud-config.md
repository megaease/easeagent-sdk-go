# MegaEase Cloud Configuration

Modify the [agent.yml](./agent.yml) file to configure your information.

## 1. Name

You'll serviceName to find your data later. It's important to use a unique and meaningful name.

The serviceName of MegaEase Cloud consists of three parts: zone, domain, name. They are joined by `.` into `ServiceName`

```yaml
serviceName: zone.domain.service
```

## 2. Reporter

MegaEase Cloud uses HTTP to receive data, so you need to change the configuration to HTTP and MegaEase Cloud's address.
```yaml
reporter.output.server: {MEGA_CLOUD_URL}/application-tracing-log
```
## 3. MTLS

MTLS is a secure authentication protocol for EaseAgent to connect to MegaEase Cloud.

Config: Get TLS
```yaml
reporter.output.server.tls.enable: true
reporter.output.server.tls.key: YOUR_TLS_KEY
reporter.output.server.tls.cert: YOUR_TLS_CERT
```

## 4. Tracing
we have some
```yaml
tracing.type: log-tracing
tracing.enable: true
tracing.sample.rate: 1.0
```

## Third: About `MEGA_CLOUD_URL` And `TLS`

When you download the `agent.yaml` file through our MegaEase Cloud, `MEGA_CLOUD_URL` and `TLS` will be filled in for you automatically.

If you need it separately, please download the `agent.yaml` and get it by yourself.