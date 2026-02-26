
# focus()

> **Warning:** Use `locator.focus([options])` instead.

Calls [focus](https://developer.mozilla.org/en-US/docs/Web/API/HTMLElement/focus) on the element, if it can be focused on.

### Returns

| Type            | Description                                                |
| --------------- | ---------------------------------------------------------- |
| `Promise<void>` | A Promise that fulfills when the focus action is finished. |

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
  const page = await browser.newPage();
  await page.goto('https://test.k6.io/browser.php');

  const textbox = await page.$('#text1');
  await textbox.focus();

  await page.close();
}
```

