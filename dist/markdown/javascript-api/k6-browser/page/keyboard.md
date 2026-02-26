
# keyboard

Returns the Keyboard instance to interact with a virtual keyboard on the page.

### Returns

| Type                                                                                                  | Description                                       |
| ----------------------------------------------------------------------------------------------------- | ------------------------------------------------- |
| Keyboard | The `Keyboard` instance associated with the page. |

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
  await page.keyboard.press('Tab');
}
```

