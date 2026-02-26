
# inputValue(selector[, options])

> **Warning:** Use locator-based `locator.inputValue([options])` instead.

Returns `input.value` for the selected `input`, `textarea` or `select` element.

| Parameter       | Type    | Default | Description                                                                                                                                                                                                                                                                                                                                   |
| --------------- | ------- | ------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| selector        | string  | `''`    | A selector to search for an element. If there are multiple elements satisfying the selector, the first will be used.                                                                                                                                                                                                                          |
| options         | object  | `null`  |                                                                                                                                                                                                                                                                                                                                               |
| options.strict  | boolean | `false` | When `true`, the call requires selector to resolve to a single element. If given selector resolves to more than one element, the call throws an exception.                                                                                                                                                                                    |
| options.timeout | number  | `30000` | Maximum time in milliseconds. Pass `0` to disable the timeout. Default is overridden by the `setDefaultTimeout` option on BrowserContext or Page. |

### Returns

| Type              | Description                                                  |
| ----------------- | ------------------------------------------------------------ |
| `Promise<string>` | A Promise that fullfils with the input value of the element. |
