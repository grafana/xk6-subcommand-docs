
# ReadableStreamDefaultReader

The `ReadableStreamDefaultReader` type represents a default reader that can be used to read stream data. It can be used to read from a ReadableStream object that has an underlying source of any type.

## Methods

| Name                                                                                                                                     | Description                                                                                                                                                       |
| ---------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| cancel(reason)     | Returns a `Promise` that resolves when the stream is canceled.                                                                                                    |
| read()               | Returns a `Promise` that resolves with an object containing the `done` and `value` properties, providing access to the next chunk in the stream's internal queue. |
| releaseLock() | Releases the lock on the stream.                                                                                                                                  |
