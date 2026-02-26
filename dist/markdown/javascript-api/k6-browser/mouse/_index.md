
# Mouse

`Mouse` provides a way to interact with a virtual mouse.

| Method                                                                                                         | Description                                                     |
| -------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------- |
| click(x, y[, options])       | Mouse clicks on the `x` and `y` coordinates.                    |
| dblclick(x, y[, options]) | Mouse double clicks on the `x` and `y` coordinates.             |
| down([options])               | Dispatches a `mousedown` event on the mouse's current position. |
| up([options])                   | Dispatches a `mouseup` event on the mouse's current position.   |
| move(x, y[, options])         | Dispatches a `mousemove` event on the mouse's current position. |

### Example

```javascript
import { browser } from 'k6/browser';

export const options = {
  scenarios: {
    ui: {
      executor: 'shared-iterations',
      options: {
        browser: {
            type: 'chromium',
        },
      },
    },
  }
}

export default async function () {
  const page = await browser.newPage();

  await page.goto('https://test.k6.io/', {
    waitUntil: 'networkidle'
});

  // Obtain ElementHandle for news link and navigate to it
  // by clicking in the 'a' element's bounding box
  const newsLinkBox = await page.$('a[href="/news.php"]');
  const boundingBox = await newsLinkBox.boundingBox();
  const x = newsLinkBox.x + newsLinkBox.width / 2; // center of the box
  const y = newsLinkBox.y;

  await page.mouse.click(x, y);

  await page.close();
}
```