
# Concepts

These topics explain the essential concepts of how scenarios and their executors work.

Different scenario configurations can affect many different aspects of your system,
including the generated load, utilized resources, and emitted metrics.
If you know a bit about how scenarios work, you'll design better tests and interpret test results with more understanding.

| On this page                                                                                                                  | Read about                                                                                                                            |
| ----------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------- |
| Open and closed models                 | Different ways k6 can schedule VUs, their affects on test results, and how k6 implements the open model in its arrival-rate executors |
| Graceful Stop                           | A configurable period for iterations to finish or ramp down after the test reaches its scheduled duration                             |
| Arrival-rate VU allocation | How k6 allocates VUs in arrival-rate executors                                                                                        |
| Dropped iterations                 | Possible reasons k6 might drop a scheduled iteration                                                                                  |
