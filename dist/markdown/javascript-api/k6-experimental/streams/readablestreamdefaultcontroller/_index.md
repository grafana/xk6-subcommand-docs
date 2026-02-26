
# ReadableStreamDefaultController

The `ReadableStreamDefaultController` type allows you to control a `ReadableStream`'s state and internal queue.

## Methods

| Name                                                                                                                                      | Description                                          |
| ----------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------- |
| close()          | Closes the associated stream.                        |
| enqueue(chunk) | Enqueues a chunk of data into the associated stream. |
| error(reason)    | Causes the stream to become errored.                 |
