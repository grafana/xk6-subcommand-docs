
# evaluateHandle(pageFunction[, arg])

Executes JavaScript code in the page and returns the value of the `pageFunction` invocation as a JSHandle.

The only difference between `page.evaluate()` and `page.evaluateHandle()` is that `page.evaluateHandle()` returns JSHandle.

| Parameter    | Type               | Defaults | Description                                                              |
| ------------ | ------------------ | -------- | ------------------------------------------------------------------------ |
| pageFunction | function or string |          | Function to be evaluated in the page context. This can also be a string. |
| arg          | string             | `''`     | Optional argument to pass to `pageFunction`                              |

### Returns

| Type                | Description                                                                                                                        |
| ------------------- | ---------------------------------------------------------------------------------------------------------------------------------- |
| `Promise<JSHandle>` | The JSHandle instance associated with the frame. |
