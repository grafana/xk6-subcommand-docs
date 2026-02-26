
# pages()

> **Caution:** This feature has **known issues**. For details, refer to
> [#444](https://github.com/grafana/xk6-browser/issues/444).

Returns all open Pages in the `BrowserContext`.

### Returns

| Type          | Description                                                                                           |
| ------------- | ----------------------------------------------------------------------------------------------------- |
| `Array` | An array of page objects. |

### Example

```javascript
import { browser } from 'k6/browser';

export const options = {
  scenarios: {
    browser: {
      executor: 'shared-iterations',
      options: {
        browser: {
          type: 'chromium',
        },
      },
    },
  },
};

export default async function () {
  const context = await browser.newContext();
  await context.newPage();
  const pages = context.pages();
  console.log(pages.length); // 1
  await context.close();
}
```

