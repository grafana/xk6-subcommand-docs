
# ElementHandle

> **Caution:** This API is a work in progress. Some of the following functionalities might behave unexpectedly.

| Method                                                                                                                                        | Description                                                                                                                                                                         |
| --------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| $(selector)                                         | Queries the element for the given selector.                                                                                                                                         |
| $$(selector)                                     | Queries the elements for the given selector.                                                                                                                                        |
| boundingBox()                                 | Returns the bounding box of the element.                                                                                                                                            |
| check([options])                                    | Checks the element if it is a `checkbox` or `radio` button.                                                                                                                         |
| click([options])                                    | Clicks on the element.                                                                                                                                                              |
| contentFrame()                               | Returns the Frame of the element.                                                                       |
| dblclick([options])                              | Double clicks on the element.                                                                                                                                                       |
| dispatchEvent(type[, eventInit])            | Dispatches an event to the element.                                                                                                                                                 |
| fill(value[, options])                               | Fills the specified value into the element.                                                                                                                                         |
| focus()                                             | Focuses on the element.                                                                                                                                                             |
| getAttribute(name)                           | Returns the specified attribute of the element.                                                                                                                                     |
| hover([options])                                    | Hovers over the element.                                                                                                                                                            |
| innerHTML()                                     | Returns the inner HTML of the element.                                                                                                                                              |
| innerText()                                     | Returns the inner text of the element.                                                                                                                                              |
| inputValue([options])                          | Returns the value of the input element.                                                                                                                                             |
| isChecked()                                     | Checks if the `checkbox` input type is selected.                                                                                                                                    |
| isDisabled()                                   | Checks if the element is `disabled`.                                                                                                                                                |
| isEditable()                                   | Checks if the element is `editable`.                                                                                                                                                |
| isEnabled()                                     | Checks if the element is `enabled`.                                                                                                                                                 |
| isHidden()                                      | Checks if the element is `hidden`.                                                                                                                                                  |
| isVisible()                                    | Checks if the element is `visible`.                                                                                                                                                 |
| ownerFrame()                                   | Returns the Frame of the element.                                                                       |
| press(key[, options])                               | Focuses on the element and presses a single key or a combination of keys using the virtual keyboard. |
| screenshot([options])                          | Takes a screenshot of the element.                                                                                                                                                  |
| scrollIntoViewIfNeeded([options])  | Scrolls the element into view if needed.                                                                                                                                            |
| selectOption(values[, options])              | Selects the `select` element's one or more options which match the values.                                                                                                          |
| selectText([options])                          | Selects the text of the element.                                                                                                                                                    |
| setChecked(checked[, options])                 | Sets the `checkbox` or `radio` input element's value to the specified checked or unchecked state.                                                                                   |
| setInputFiles(file[, options])              | Sets the file input element's value to the specified files.                                                                                                                         |
| tap(options)                                          | Taps the element.                                                                                                                                                                   |
| textContent()                                 | Returns the text content of the element.                                                                                                                                            |
| type(text[, options])                                | Focuses on the element and types the specified text into the element using the virtual keyboard.     |
| uncheck([options])                                | Unchecks the element if it is a `checkbox` or `radio` button.                                                                                                                       |
| waitForElementState(state[, options]) | Waits for the element to reach the specified state.                                                                                                                                 |
| waitForSelector(selector[, options])      | Waits for the element to be present in the DOM and to be visible.                                                                                                                   |

## Examples

```javascript
import { browser } from 'k6/browser';
import { check } from "https://jslib.k6.io/k6-utils/1.5.0/index.js";

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

  // Goto front page, find login link and click it
  try {
    await page.goto('https://test.k6.io/');
    const messagesLink = await page.$('a[href="/my_messages.php"]');

    await Promise.all([page.waitForNavigation(), messagesLink.click()]);
    // Enter login credentials and login
    const login = await page.$('input[name="login"]');
    await login.type('admin');
    const password = await page.$('input[name="password"]');
    await password.type('123');

    const submitButton = await page.$('input[type="submit"]');

    await Promise.all([page.waitForNavigation(), submitButton.click()]);

    await check(page, {
      'header': async p => {
        const h2 = await p.$('h2');
        return await h2.textContent() == 'Welcome, admin!';
      },
    });
  } finally {
    await page.close();
  }
}
```

```javascript
import { browser } from 'k6/browser';
import { check } from 'https://jslib.k6.io/k6-utils/1.5.0/index.js';

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
    // Inject page content
    await page.setContent(`
          <div class="visible">Hello world</div>
          <div style="display:none" class="hidden"></div>
          <div class="editable" editable>Edit me</div>
          <input type="checkbox" enabled class="enabled">
          <input type="checkbox" disabled class="disabled">
          <input type="checkbox" checked class="checked">
          <input type="checkbox" class="unchecked">
    `);

    // Check state
    await check(page, {
      'is visible': async p => {
        const e = await p.$('.visible');
        return e.isVisible();
      },
      'is hidden': async p => {
        const e = await p.$('.hidden');
        return e.isHidden();
      },
      'is editable': async p => {
        const e = await p.$('.editable');
        return e.isEditable();
      },
      'is enabled': async p => {
        const e = await p.$('.enabled');
        return e.isEnabled();
      },
      'is disabled': async p => {
        const e = await p.$('.disabled');
        return e.isDisabled();
      },
      'is checked': async p => {
        const e = await p.$('.checked');
        return e.isChecked();
      },
      'is unchecked': async p => {
        const e = await p.$('.unchecked');
        return await e.isChecked() === false;
      },
    });
  } finally {
    await page.close();
  }
}
```
