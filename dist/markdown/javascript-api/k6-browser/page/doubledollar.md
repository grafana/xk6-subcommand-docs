
# page.$$(selector)

> **Warning:** When possible, use locator-based `page.locator(selector)` instead.
> However, using `locator`s may not always work when selecting an element from a list or table, especially if there isn't a reliable way to consistently identify a single element (for example, due to changing or non-unique attributes). In such cases, `$$` remains useful.

The method finds all elements matching the specified selector within the page. If no elements match the selector, the return value resolves to `[]`. The results are returned in DOM order. This is particularly useful when you want to retrieve a list of elements, and iterate through them to find the one that you need for your test case.

### Returns

| Type                       | Description                                                                                                                                                                                 |
| -------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `Promise` | A Promise that fulfills with the ElementHandle array of the selector when matching elements are found. |

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

  await page.goto('https://test.k6.io/');

  // Retrieve all the td elements.
  const cells = await page.$$('td');
  for (let i = 0; i < cells.length; i++) {
    if ((await cells[i].innerText()) == '/pi.php?decimals=3') {
      // When the element is found, click on it and
      // wait for the navigation.
      await Promise.all([page.waitForNavigation(), cells[i].click()]);
      break;
    }
  }

  // Wait for an important element to load.
  await page.locator('//pre[text()="3.141"]').waitFor();
}
```

