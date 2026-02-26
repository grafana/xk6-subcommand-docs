
# contentFrame()

This method returns a FrameLocator object pointing to the same `iframe` as this locator. This is useful when you have a Locator object obtained somewhere, and later on would like to interact with the content inside the frame.

## Returns

| Type                                                                                   | Description                                              |
| -------------------------------------------------------------------------------------- | -------------------------------------------------------- |
| FrameLocator | The `FrameLocator` pointing to the same`iframe` as this locator. |

## Example

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

  try {
    await page.setContent(`
      <iframe id="my_frame" src="https://quickpizza.grafana.com" width="50%" height="50%"></iframe>
    `);

    const frameLocator = page.locator('#my_frame').contentFrame();
    await frameLocator.getByText('Pizza, Please!').click();

    const noThanksBtn = frameLocator.getByText('No thanks');
    await noThanksBtn.click();
  } finally {
    await page.close();
  }
}
```
