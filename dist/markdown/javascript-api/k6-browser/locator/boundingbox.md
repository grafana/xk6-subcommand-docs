
# boundingBox([options])

Returns the bounding box of the element. The bounding box is calculated with respect to the position of the Frame of the current element, which is usually the Page's main frame.

| Parameter       | Type   | Default | Description                                                                                                                                                                                                                                                                                                                                   |
| --------------- | ------ | ------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| options         | object | `null`  |                                                                                                                                                                                                                                                                                                                                               |
| options.timeout | number | `30000` | Maximum time in milliseconds. Pass `0` to disable the timeout. Default is overridden by the `setDefaultTimeout` option on BrowserContext or Page. |

### Returns

| Type                    | Description                                                                                                                  |
| ----------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| `Promise<null \| Rect>` | A Promise that fulfills with the bounding box of the element as a [Rect](#rect). If the element is not visible, the Promise resolves to `null`. |

### Rect

The `Rect` object represents the bounding box of an element.

| Property | Type     | Description                                |
| -------- | -------- | ------------------------------------------ |
| x        | `number` | The x-coordinate of the element in pixels. |
| y        | `number` | The y-coordinate of the element in pixels. |
| width    | `number` | The width of the element in pixels.        |
| height   | `number` | The height of the element in pixels.       |

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

  const textbox = page.locator('#text1');
  const boundingBox = await textbox.boundingBox();
  console.log(`x: ${boundingBox.x}, y: ${boundingBox.y}, width: ${boundingBox.width}, height: ${boundingBox.height}`);

  await page.close();
}
```

