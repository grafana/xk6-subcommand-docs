
# textContent([options])

Returns the `element.textContent`.

| Parameter       | Type   | Default | Description                                                                                                                                                                                                                                                                                                         |
| --------------- | ------ | ------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| options         | object | `null`  |                                                                                                                                                                                                                                                                                                                     |
| options.timeout | number | `30000` | Maximum time in milliseconds. Pass `0` to disable the timeout. Default is overridden by the `setDefaultTimeout` option on BrowserContext or Page. |

### Returns

| Type                      | Description                                                              |
| ------------------------- | ------------------------------------------------------------------------ |
| `Promise<null \| string>` | A Promise that fulfills with the text content of the selector or `null`. |

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
  const options = page.locator('#checkbox1');
  console.log(await options.textContent());
}
```

