
# mouse

Returns the Mouse instance to interact with a virtual mouse on the page.

### Returns

| Type                                                                                            | Description                                    |
| ----------------------------------------------------------------------------------------------- | ---------------------------------------------- |
| Mouse | The `Mouse` instance associated with the page. |

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
  await page.mouse.down();
}
```

ß
