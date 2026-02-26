
# clearCookies()

Clears the browser context's cookies.

### Returns

| Type            | Description                                                                                                                                                                            |
| --------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `Promise<void>` | A Promise that fulfills when the cookies have been cleared from the browser context. |

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
  const page = await context.newPage();

  await page.goto('https://httpbin.org/cookies/set?testcookie=testcookievalue');
  let cookies = await context.cookies();
  console.log(cookies.length); // prints: 1

  await context.clearCookies();
  cookies = await context.cookies();
  console.log(cookies.length); // prints: 0
}
```

