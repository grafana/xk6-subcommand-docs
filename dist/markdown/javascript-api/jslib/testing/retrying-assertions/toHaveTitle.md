
# toHaveTitle()

The `toHaveTitle()` method asserts that a page has a specific title. This is a retrying assertion that automatically waits for the page to have the expected title.

## Syntax

```javascript
await expect(page).toHaveTitle(expected);
await expect(page).not.toHaveTitle(expected);
await expect(page).toHaveTitle(expected, options);
```

## Parameters

| Parameter | Type                                                                                                                    | Description                    |
| --------- | ----------------------------------------------------------------------------------------------------------------------- | ------------------------------ |
| expected  | string \| RegExp                                                                                                        | The expected title             |
| options   | RetryConfig | Optional configuration options |

## Returns

| Type          | Description                                       |
| ------------- | ------------------------------------------------- |
| Promise<void> | A promise that resolves when the assertion passes |

## Description

The `toHaveTitle()` method checks if a page has a specific title. It retrives the title from the page's `<title>` element in the document's `<head>`.

When a string is provided, it performs an exact match. When a RegExp is provided, it tests the title against the regular expression pattern.

This is a retrying assertion that will automatically re-check the page's title until it matches the expected value or the timeout is reached.

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
  
  // Navigate and check exact title match
  await page.goto('https://quickpizza.grafana.com/');
  await expect(page).toHaveTitle('QuickPizza');

  // Check title with regex pattern
  await expect(page).toHaveTitle(/QuickPizza/);
  await expect(page).toHaveTitle(/^Quick/);
}
```
