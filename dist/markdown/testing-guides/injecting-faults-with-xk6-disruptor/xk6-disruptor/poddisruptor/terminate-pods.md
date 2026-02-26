
# terminatePods()

`terminatePods` terminates a number of the pods matching the selector configured in the PodDisruptor.

| Parameter | Type   | Description                                                                                                                              |
| --------- | ------ | ---------------------------------------------------------------------------------------------------------------------------------------- |
| fault     | object | description of the Pod Termination fault |

## Example

```javascript
const fault = {
  count: 2,
};
disruptor.terminatePods(fault);
```
