

# getByTitle(title[, options])

Returns a locator for elements with the specified `title` attribute. This method is useful for locating elements that use the `title` attribute to provide additional information, tooltips, or accessibility context.

| Parameter       | Type             | Default | Description                                                                                                  |
| --------------- | ---------------- | ------- | ------------------------------------------------------------------------------------------------------------ |
| `title`         | string \| RegExp | -       | Required. The title text to search for. Can be a string for exact match or a RegExp for pattern matching.    |
| `options`       | object           | `null`  |                                                                                                              |
| `options.exact` | boolean          | `false` | Whether to match the title text exactly with case sensitivity. When `true`, performs a case-sensitive match. |

## Returns

| Type                                                                                   | Description                                                                                             |
| -------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------- |
| Locator | A locator object that can be used to interact with the elements matching the specified title attribute. |

## Examples

Find and interact with elements by their title attribute:

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
      <button title="Hello World">Hi</button>
      <select title="Favorite Color">
        <option value="Red">Red</option>
        <option value="Green">Green</option>
        <option value="Blue">Blue</option>
      </select>
      <input type="checkbox" title="Check me">
    `);

    const frame = page.mainFrame();

    await frame.getByTitle('Hello World').click();

    await frame.getByTitle('Favorite Color').selectOption('Red');

    await frame.getByTitle('Check me').check();
  } finally {
    await page.close();
  }
}
```

## Common use cases

- **User interface controls:**
  - Toolbar buttons and action items
  - Navigation controls (previous/next, pagination)
  - Media player controls
  - Menu and drop-down triggers
- **Informational elements:**
  - Help icons and tooltips
  - Status indicators and badges
  - Progress indicators
  - Warning and error messages
- **Accessibility support:**
  - Screen reader descriptions
  - Alternative text for complex elements
  - Context-sensitive help
  - Form field explanations

## Best practices

1. **Meaningful titles**: Ensure title attributes provide clear, helpful information about the element's purpose or content.
1. **Accessibility compliance**: Use titles to enhance accessibility, especially for elements that might not have clear visual labels.
1. **Avoid redundancy**: Don't duplicate visible text in the title attribute unless providing additional context.
1. **Dynamic content**: When testing applications with changing title content, use flexible matching patterns or regular expressions.
1. **Tooltip testing**: Remember that title attributes often create tooltips on hover, which can be useful for UI testing.
1. **Internationalization**: Consider that title text may change in different language versions of your application.

## Related

- frame.getByRole() - Locate by ARIA role
- frame.getByAltText() - Locate by alt text
- frame.getByLabel() - Locate by form labels
- frame.getByPlaceholder() - Locate by placeholder text
- frame.getByTestId() - Locate by test ID
- frame.getByText() - Locate by visible text
