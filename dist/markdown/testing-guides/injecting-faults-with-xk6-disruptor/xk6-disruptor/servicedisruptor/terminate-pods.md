
# terminatePods()

`terminatePods` terminates a number of pods that belong to the service specified in the ServiceDisruptor.

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
