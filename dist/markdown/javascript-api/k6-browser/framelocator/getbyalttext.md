

# getByAltText(altText[, options])

Returns a locator for elements with the specified alt text. This method is useful for locating images and other elements that have an `alt` attribute, making your tests more accessible and user-focused.

| Parameter       | Type             | Default | Description                                                                                                |
| --------------- | ---------------- | ------- | ---------------------------------------------------------------------------------------------------------- |
| `altText`       | string \| RegExp | -       | Required. The alt text to search for. Can be a string for exact match or a RegExp for pattern matching.    |
| `options`       | object           | `null`  |                                                                                                            |
| `options.exact` | boolean          | `false` | Whether to match the alt text exactly with case sensitivity. When `true`, performs a case-sensitive match. |

## Returns

| Type                                                                                   | Description                                                                                  |
| -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| Locator | A locator object that can be used to interact with elements matching the specified alt text. |

## Examples

### Basic alt text matching

Find and click an image by its alt text:

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
  },
};

export default async function () {
  const page = await browser.newPage();

  try {
    await page.setContent(`
      <iframe id="my_frame" src="https://quickpizza.grafana.com" width="50%" height="50%"></iframe>
    `);

    const logo = page.locator('#my_frame').contentFrame().getByAltText('LOGO');
    await logo.waitFor();

    await logo.click();
    await page.waitForLoadState();
  } finally {
    await page.close();
  }
}
```

### Exact alt text matching

Use exact matching for precise alt text:

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
  },
};

export default async function () {
  const page = await browser.newPage();

  try {
    await page.setContent(`
      <iframe id="my_frame" src="https://quickpizza.grafana.com" width="50%" height="50%"></iframe>
    `);

    const logo = page.locator('#my_frame').contentFrame().getByAltText('logo', { exact: true });
    await logo.waitFor();

    await logo.click();
    await page.waitForLoadState();
  } finally {
    await page.close();
  }
}
```

### Using regular expressions

Find images using pattern matching:

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
  },
};

export default async function () {
  const page = await browser.newPage();

  try {
    await page.setContent(`
      <iframe id="my_frame" src="https://quickpizza.grafana.com" width="50%" height="50%"></iframe>
    `);

    const logo = page.locator('#my_frame').contentFrame().getByAltText(/logo/s);
    await logo.waitFor();

    await logo.click();
    await page.waitForLoadState();
  } finally {
    await page.close();
  }
}
```

## Common use cases

- **Image testing:**
  - Testing image alt text for accessibility compliance
  - Interacting with clickable images/banners
- **Accessibility testing:**
  - Ensuring all images have meaningful alt text
  - Testing screen reader compatibility
  - Validating WCAG compliance
- **UI interaction:**
  - Clicking on logo images to navigate home
  - Selecting images in galleries or carousels
  - Interacting with image-based buttons

## Best practices

1. **Use descriptive alt text**: When creating tests, ensure your application uses meaningful alt text that describes the image content or function.
1. **Prefer exact matches**: Use `exact: true` when you need precise matching to avoid false positives.
1. **Accessibility-first**: Using `getByAltText()` encourages better accessibility practices by ensuring images have proper alt attributes.
1. **Regular expressions for patterns**: Use RegExp for flexible matching when dealing with dynamic or numbered content.
1. **Combine with assertions**: Always verify that the located elements behave as expected using assertions.

## Related

- frameLocator.getByRole() - Locate by ARIA role
- frameLocator.getByLabel() - Locate by form labels
- frameLocator.getByPlaceholder() - Locate by placeholder text
- frameLocator.getByTestId() - Locate by test ID
- frameLocator.getByTitle() - Locate by title attribute
- frameLocator.getByText() - Locate by visible text
