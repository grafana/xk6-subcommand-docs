
# k6/x/icmp

> **Note:** This module is implemented as an official extension and is not part of k6 by default. However, k6's [automatic extension resolution](/docs/k6/v1.5.0/extensions/run) makes it available seamlessly when you import it in your script. Refer to the [extensions documentation](/docs/k6/v1.5.0/extensions/explore) for available extensions and details.

The `k6/x/icmp` module adds support for sending ICMP echo requests (pings) to measure network latency and availability. This allows you to measure network latency and reachability of hosts directly within your load testing and synthetic monitoring scenarios.

The main use case for this extension is integration with Grafana's Synthetic Monitoring product to make a quick reachability check after a scripted check fails. It can also be used directly in k6 scripts, for example as a pre-check to verify network connectivity before running your main test.

## Key features

- Send ICMP echo requests (ping) to any host
- Synchronous and asynchronous ping functions
- Support for multiple ping attempts with configurable count and interval
- Configurable timeouts and deadlines
- Detailed ping results via callback function
- Support for IPv4 and IPv6
- Custom packet size and TTL settings
- Automatic metrics collection
- Success threshold configuration for multiple pings

### Use cases

- Network connectivity pre-checks before load tests
- Host reachability validation in synthetic monitoring
- Network latency measurement across different hosts
- Testing network infrastructure reliability
- Monitoring service availability at the network layer
- Validating network paths and routing

## API

| Function/Type | Description |
| ------------- | ----------- |
| ping( target, [options] ) | Send ICMP echo requests synchronously |
| pingAsync( target, [optsOrCallback], [callback] ) | Send ICMP echo requests asynchronously |
| PingOptions | Configuration options for ping operations |
| PingDetail | Detailed information about ping results |
| IP | IP protocol version enumeration |

## Metrics

The extension automatically generates metrics for ICMP operations:

- `icmp_pings`: Counter tracking the number of ping attempts
- `icmp_ping_duration`: Trend measuring round-trip time for successful pings
- `icmp_ping_failed`: Counter tracking failed ping attempts

## Examples

### Basic Usage

Here is a basic example showing how to use the `ping()` function:

```javascript
import { ping } from "k6/x/icmp"

export default function () {
  const host = "8.8.8.8"

  console.log(`Pinging ${host}:`);

  if (ping(host)) {
    console.log(`Host ${host} is reachable`);
  } else {
    console.error(`Host ${host} is unreachable`);
  }
}
```

### Advanced Usage

A more advanced example below demonstrates how to use the `PingCallback` to access detailed ping results for each request.

```javascript
import { pingAsync } from "k6/x/icmp"

export default async function () {
  const host = "8.8.8.8"

  console.log(`Pinging ${host} with callback:`);

  const opts = {
    timeout: 3000,
    count: 5
  };

  const result = await pingAsync(host, opts, (err, { target, sent_at, received_at, seq, ttl, size }) => {
    if (err) {
      console.error(`${target}: ${err}`);

      return
    }

    const rtt = received_at - sent_at;

    console.log(`${size} bytes from ${target}: icmp_seq=${seq} ttl=${ttl} time=${rtt} ms`);
  });

  if (result) {
    console.log(`Host ${host} is reachable`);
  } else {
    console.error(`Host ${host} is unreachable`);
  }
}
```
