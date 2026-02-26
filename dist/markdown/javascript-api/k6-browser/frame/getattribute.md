
# getAttribute(selector, name[, options])

> **Warning:** Use locator-based `locator.getAttribute()` instead.

Returns the element attribute value for the given attribute name.

| Parameter       | Type    | Default | Description                                                                                                                                                                                                                                                                                                                                   |
| --------------- | ------- | ------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| selector        | string  | `''`    | A selector to search for an element. If there are multiple elements satisfying the selector, the first will be used.                                                                                                                                                                                                                          |
| name            | string  | `''`    | Attribute name to get the value for.                                                                                                                                                                                                                                                                                                          |
| options         | object  | `null`  |                                                                                                                                                                                                                                                                                                                                               |
| options.strict  | boolean | `false` | When `true`, the call requires selector to resolve to a single element. If given selector resolves to more than one element, the call throws an exception.                                                                                                                                                                                    |
| options.timeout | number  | `30000` | Maximum time in milliseconds. Pass `0` to disable the timeout. Default is overridden by the `setDefaultTimeout` option on BrowserContext or Page. |

### Returns

| Type                      | Description                                                                       |
| ------------------------- | --------------------------------------------------------------------------------- |
| `Promise<string \| null>` | A Promise that fulfills with the value of the attribute. Else, it returns `null`. |
