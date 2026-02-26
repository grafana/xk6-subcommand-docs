

# getByPlaceholder(placeholder[, options])

Returns a locator for input elements with the specified `placeholder` attribute. This method is useful for locating form fields that use `placeholder` attribute to provide hints or examples to users about the expected input format.

| Parameter       | Type             | Default | Description                                                                                                        |
| --------------- | ---------------- | ------- | ------------------------------------------------------------------------------------------------------------------ |
| `placeholder`   | string \| RegExp | -       | Required. The placeholder text to search for. Can be a string for exact match or a RegExp for pattern matching.    |
| `options`       | object           | `null`  |                                                                                                                    |
| `options.exact` | boolean          | `false` | Whether to match the placeholder text exactly with case sensitivity. When `true`, performs a case-sensitive match. |

## Returns

| Type                                                                                   | Description                                                                                                    |
| -------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------- |
| Locator | A locator object that can be used to interact with the input elements matching the specified placeholder text. |

## Example

Find and fill inputs by their placeholder text:

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
    const iframeHTML = `
      <input type="text" placeholder="First name">
      <input type="text" placeholder="Last name">
      <input type="text" placeholder="dd/mm/yyyy">
      <input type="text" placeholder="your.email@example.com">
      <input type="text" placeholder="+1 (555) 123-4567">
    `;

    await page.setContent(`
      <iframe id="my_frame" src="data:text/html,${encodeURIComponent(iframeHTML)}"></iframe>
    `);

    const frameLocator = page.locator("#my_frame").contentFrame();
    await frameLocator.getByPlaceholder('First name').fill('First');
    await frameLocator.getByPlaceholder('Last name').fill('Last');
    await frameLocator.getByPlaceholder('dd/mm/yyyy').fill('01/01/1990');

    await frameLocator.getByPlaceholder('your.email@example.com').fill('first.last@example.com');
    await frameLocator.getByPlaceholder('+1 (555) 123-4567').fill('+1 (555) 987-6543');
  } finally {
    await page.close();
  }
}
```

## Common use cases

- **Form field identification:**
  - Login and registration forms without explicit labels
  - Quick search boxes
  - Filter and input controls
  - Comment and feedback forms
- **E-commerce:**
  - Product search bars
  - Quantity input fields
  - Promo code entry
  - Address and payment information
- **Interactive applications:**
  - Chat input fields
  - Command input interfaces
  - Settings and configuration forms
  - Data entry applications

## Best practices

1. **Complement, don't replace labels**: Placeholder text should supplement, not replace proper form labels for accessibility.
1. **Use descriptive placeholders**: Ensure placeholder text clearly indicates the expected input format or content.
1. **Consider internationalization**: When testing multi-language applications, be aware that placeholder text may change.
1. **Accessibility considerations**: Remember that placeholder text alone may not be sufficient for users with disabilities.

## Related

- frameLocator.getByRole() - Locate by ARIA role
- frameLocator.getByAltText() - Locate by alt text
- frameLocator.getByLabel() - Locate by form labels (preferred for accessibility)
- frameLocator.getByTestId() - Locate by test ID
- frameLocator.getByTitle() - Locate by title attribute
- frameLocator.getByText() - Locate by visible text
