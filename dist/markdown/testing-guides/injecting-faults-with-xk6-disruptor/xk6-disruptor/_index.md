
# xk6-disruptor API

The xk6-disruptor API is organized around _disruptors_ that affect specific targets such as pods or services. These disruptors can inject different types of faults on their targets.

| Class                                                                                                      | Description                                                      |
| ---------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------- |
| PodDisruptor         | Targets the Pods that match selection attributes such as labels. |
| ServiceDisruptor | Targets the Pods that back a Kubernetes Service                  |
