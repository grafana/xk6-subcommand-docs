
# xk6-disruptor first steps

[xk6-disruptor](https://github.com/grafana/xk6-disruptor) is an extension that adds fault injection capabilities to k6.

It provides a Javascript API to inject faults such as errors and delays into HTTP and gRPC requests served by selected Kubernetes Pods or Services.

```javascript
import { ServiceDisruptor } from 'k6/x/disruptor';

export default function () {
  // Create a new disruptor that targets a service
  const disruptor = new ServiceDisruptor('app-service', 'app-namespace');

  // Disrupt the targets by injecting delays and faults into HTTP request for 30 seconds
  const fault = {
    averageDelay: '500ms',
    errorRate: 0.1,
    errorCode: 500,
  };
  disruptor.injectHTTPFaults(fault, '30s');
}
```

## Next steps

Explore the fault injection API

See step-by-step examples.

Visit the [interactive demo environment](https://killercoda.com/grafana-xk6-disruptor/scenario/killercoda).

Learn the basics of using the disruptor in your test project:

- Requirements

- Installation
