
# Executors

**Executors** control how k6 schedules VUs and iterations.
The executor that you choose depends on the goals of your test and the type of traffic you want to model.

Define the executor in `executor` key of the scenario object.
The value is the executor name separated by hyphens.

```javascript
export const options = {
  scenarios: {
    arbitrary_scenario_name: {
      //Name of executor
      executor: 'ramping-vus',
      // more configuration here
    },
  },
};
```

## All executors

The following table lists all k6 executors and links to their documentation.

| Name                                                                                                                 | Value                   | Description                                                                                                                                                                                |
| -------------------------------------------------------------------------------------------------------------------- | ----------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| Shared iterations         | `shared-iterations`     | A fixed amount of iterations are shared between a number of VUs.                                                                                                                      |
| Per VU iterations         | `per-vu-iterations`     | Each VU executes an exact number of iterations.                                                                                                                                            |
| Constant VUs                   | `constant-vus`          | A fixed number of VUs execute as many iterations as possible for a specified amount of time.                                                                                          |
| Ramping VUs                     | `ramping-vus`           | A variable number of VUs execute as many iterations as possible for a specified amount of time.                                                                                       |
| Constant Arrival Rate | `constant-arrival-rate` | A fixed number of iterations are executed in a specified period of time.                                                                                                              |
| Ramping Arrival Rate   | `ramping-arrival-rate`  | A variable number of iterations are  executed in a specified period of time.                                                                                                          |
| Externally Controlled | `externally-controlled` | Control and scale execution at runtime via [k6's REST API](https://grafana.com/docs/k6/v1.5.0/misc/k6-rest-api) or the [CLI](https://k6.io/blog/how-to-control-a-live-k6-test). |

> **Note:** For any given scenario, you can't guarantee that a specific VU can run a specific iteration.
> With `SharedArray` and execution context variables, you can map a specific VU to a specific value in your test data.
> So the tenth VU could use the tenth item in your array (or the sixth iteration to the sixth item).
> But, you _cannot_ reliably map, for example, the tenth VU to the tenth iteration.

