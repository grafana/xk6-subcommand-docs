
# newPage()

Uses the browser context to create and return a new page.

### Returns

| Type            | Description                                                                                                                              |
| --------------- | ---------------------------------------------------------------------------------------------------------------------------------------- |
| `Promise` | A Promise that fulfills with a new page object. |

### Example

```javascript
import { browser } from 'k6/browser';

export default async function () {
  const context = await browser.newContext();
  const page = await context.newPage();

  try {
    await page.goto('https://test.k6.io/browser.php');
  } finally {
    await page.close();
  }
}
```

