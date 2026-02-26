
# selectText(values, [options])

Selects the text of the element.

| Parameter           | Type                         | Default | Description                                                                                                                                                                                                                                                                                                                                   |
| ------------------- | ---------------------------- | ------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| options             | object                       | `null`  |                                                                                                                                                                                                                                                                                                                                               |
| options.force       | boolean                      | `false` | Setting this to `true` will bypass the actionability checks (`visible`, `stable`, `enabled`).                                                                                                                                                                                                                                                 |
| options.timeout     | number                       | `30000` | Maximum time in milliseconds. Pass `0` to disable the timeout. Default is overridden by the `setDefaultTimeout` option on BrowserContext or Page. |

### Returns

| Type            | Description                                        |
| --------------- | -------------------------------------------------- |
| `Promise<void>` | A Promise that fulfills when the text is selected. |

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
  await page.goto('https://test.k6.io');

  const element = await page.$('h1');
  await element.selectText();

  await page.close();
}
```

