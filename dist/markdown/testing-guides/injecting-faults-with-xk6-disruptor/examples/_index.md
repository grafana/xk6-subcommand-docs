
# xk6-disruptor examples

In this section, we present some examples of using the `xk6-disruptor` extension to introduce faults in `k6` tests.

- Injecting gRPC faults into a Service
- Injecting HTTP faults into a Pod
- [Interactive demo](https://killercoda.com/grafana-xk6-disruptor/scenario/killercoda) (Killercoda)

To follow the instructions of the examples, check first the system under test meets the requirements to receive faults, in particular:

- You have configured the credentials to access the Kubernetes cluster.
- This cluster exposes the service using an external IP.
