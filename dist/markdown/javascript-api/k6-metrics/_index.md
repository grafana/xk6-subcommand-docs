
# k6/metrics

The `k6/metrics` module provides functionality to create custom metrics of various types.

| Metric type                                                                           | Description                                                                                   |
| ------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------- |
| Counter | A metric that cumulatively sums added values.                                                 |
| Gauge     | A metric that stores the min, max and last values added to it.                                |
| Rate       | A metric that tracks the percentage of added values that are non-zero.                        |
| Trend     | A metric that calculates statistics on the added values (min, max, average, and percentiles). |

