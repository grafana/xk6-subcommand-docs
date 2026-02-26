
# $$(selector)

> **Warning:** Use locator-based `frame.locator(selector)` instead.

The method finds all elements matching the specified selector within the page. If no elements match the selector, the return value resolves to `[]`. The results are returned in DOM order.

### Returns

| Type                       | Description                                                                                                                                                                                              |
| -------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `Promise` | A Promise that fulfills with the ElementHandle array of the selector when matching elements are found. |
