
# $(selector)

> **Warning:** Use locator-based `frame.locator(selector)` instead.

The method finds an element matching the specified selector within the frame. If no elements match the selector, the return value resolves to `null`. To wait for an element on the frame, use locator.waitFor([options]).

### Returns

| Type                             | Description                                                                                                                                                                                                   |
| -------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `Promise` | A Promise that fulfills with the ElementHandle  of the selector when a matching element is found or `null`. |
