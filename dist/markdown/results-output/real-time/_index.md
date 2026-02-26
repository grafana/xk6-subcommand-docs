
# Real time

Besides the end-of-test summary, you can also view metrics as granular data points.
k6 can stream the metrics in real-time and either:

- Write output to a file
- Send output to an external service.

## Write to file {#write}

Currently, k6 supports writing to the following file formats:

- CSV
- JSON

## Stream to service {#service}

You can also stream real-time metrics to:

- Grafana Cloud k6

As well as the following third-party services:

- Amazon CloudWatch
- Apache Kafka
- Datadog
- Dynatrace
- Elasticsearch
- Grafana Cloud Prometheus
- InfluxDB
- Netdata
- New Relic
- OpenTelemetry
- Prometheus remote write
- TimescaleDB
- StatsD
- [Other alternative with a custom output extension](https://grafana.com/docs/k6/v1.5.0/extensions/create/output-extensions)

> **Note:** This list applies to local tests, not to [cloud tests](https://grafana.com/docs/grafana-cloud/testing/k6/).

## Read more

- [Ways to visualize k6 results](https://k6.io/blog/ways-to-visualize-k6-results/)
- [k6 data collection pipeline](https://grafana.com/blog/2023/08/10/understanding-grafana-k6-a-simple-guide-to-the-load-testing-tool/)
