
# toBeEnabled()

The `toBeEnabled()` method asserts that an element is enabled and interactive. This is a retrying assertion that automatically waits for the element to become enabled.

## Syntax

```javascript
await expect(locator).toBeEnabled();
await expect(locator).not.toBeEnabled();
await expect(locator).toBeEnabled(options);
```

## Parameters

| Parameter | Type                                                                                                                    | Description                    |
| --------- | ----------------------------------------------------------------------------------------------------------------------- | ------------------------------ |
| options   | RetryConfig | Optional configuration options |

## Returns

| Type          | Description                                       |
| ------------- | ------------------------------------------------- |
| Promise<void> | A promise that resolves when the assertion passes |

## Description

The `toBeEnabled()` method checks if an element is enabled and interactive. An element is considered enabled if:

- It exists in the DOM
- It does not have the `disabled` attribute
- It is not disabled through CSS or JavaScript
- It can receive user interactions

This is a retrying assertion that will automatically re-check the element's enabled state until it becomes enabled or the timeout is reached.

## Usage

```javascript
import { browser } from 'k6/browser';
import { expect } from 'https://jslib.k6.io/k6-testing//index.js';

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
};

export default async function () {
  const page = await browser.newPage();
  await page.goto('https://quickpizza.grafana.com/');

  // Check that the pizza button is enabled
  await expect(page.locator('button[name="pizza-please"]')).toBeEnabled();

  // Check that buttons are generally enabled on the page
  await expect(page.locator('button')).toBeEnabled();
}
```

