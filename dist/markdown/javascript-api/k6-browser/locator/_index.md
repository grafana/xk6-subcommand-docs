
# Locator

The Locator API makes it easier to work with dynamically changing elements. Some of the benefits of using it over existing ways to locate an element (e.g. `Page.$()`) include:

- Helps with writing robust tests by finding an element even if the underlying frame navigates.
- Makes it easier to work with dynamic web pages and SPAs built with Svelte, React, Vue, etc.
- Enables the use of test abstractions like the Page Object Model (POM) pattern to simplify and organize tests.
- `strict` mode is enabled for all `locator` methods that are expected to target a single DOM element, meaning that if more than one element matches the given selector, an error will be thrown.

Locator can be created with the page.locator(selector[, options]) method.

| Method                                                                                                                                                 | Description                                                                                                             |
| ------------------------------------------------------------------------------------------------------------------------------------------------------ | ----------------------------------------------------------------------------------------------------------------------- |
| all()                                                                | When multiple elements match the selector, returns an array of `locator`.                                                                         |
| boundingBox([options])                                      | Returns the bounding box of the element.                                                                               |
| check([options])                        | Select the input checkbox.                                                                                              |
| clear([options])                                                   | Clears text boxes and input fields of any existing values.                                                              |
| click([options])                        | Mouse click on the chosen element.                                                                                      |
| contentFrame()                                                            | Returns a `FrameLocator` object pointing to the same `iframe` as this locator.   
| count()                                                            | Returns the number of elements matching the selector.                                                                   |
| dblclick([options])                  | Mouse double click on the chosen element.                                                                               |
| dispatchEvent(type, eventInit, [options])                  | Dispatches HTML DOM event types e.g. `'click'`.                                                                         |
| evaluate(pageFunction[, arg])                                              | Returns the value of the `pageFunction` invocation, called with the matching element as first argument.                                       |
| evaluateHandle(pageFunction[, arg])                                              | Returns the value of the `pageFunction` invocation as a JSHandle.                                       |
| fill(value, [options])                                              | Fill an `input`, `textarea` or `contenteditable` element with the provided value.                                       |
| filter(options)                                                   | Returns a new `locator` that matches only elements containing or excluding specified text.                               |
| first()                                                            | Returns a `locator` to the first matching element.                                                                        |
| focus([options])                                                   | Calls [focus](https://developer.mozilla.org/en-US/docs/Web/API/HTMLElement/focus) on the element, if it can be focused. |
| getAttribute(name, [options])                               | Returns the element attribute value for the given attribute name.                                                       |
| getByAltText(altText[, options])                                   | Returns a locator for elements with the specified `alt` attribute text.                                                                                    |
| getByLabel(text[, options])                                          | Returns a locator for form controls with the specified label text.                                                                                   |
| getByPlaceholder(placeholder[, options])                       | Returns a locator for input elements with the specified `placeholder` attribute text.                                                                                   |
| getByRole(role[, options])                                            | Returns a locator for elements with the specified ARIA role.                                                                                    |
| getByTestId(testId)                                                 | Returns a locator for elements with the specified `data-testid` attribute.                                                                              |
| getByText(text[, options])                                            | Returns a locator for elements containing the specified text.                                                                                  |
| getByTitle(title[, options])                                         | Returns a locator for elements with the specified `title` attribute.                                                                             |
| hover([options])                        | Hovers over the element.                                                                                                |
| innerHTML([options])                                           | Returns the `element.innerHTML`.                                                                                        |
| innerText([options])                                           | Returns the `element.innerText`.                                                                                        |
| inputValue([options])                                         | Returns `input.value` for the selected `input`, `textarea` or `select` element.                                         |
| isChecked([options])                                           | Checks if the `checkbox` `input` type is selected.                                                                      |
| isDisabled([options])                                         | Checks if the element is `disabled`.                                                                                    |
| isEditable([options])                                         | Checks if the element is `editable`.                                                                                    |
| isEnabled([options])                                           | Checks if the element is `enabled`.                                                                                     |
| isHidden()                                                      | Checks if the element is `hidden`.                                                                                      |
| isVisible()                                                    | Checks if the element is `visible`.                                                                                     |
| last()                                                              | Returns a `locator` to the last matching element.                                                                         |
| locator(selector[, options])                                     | Returns a new chained `locator` for the given `selector`.                                                                 |
| nth()                                                                | Returns a `locator` to the n-th matching element.                                                                         |
| press(key, [options])                                              | Press a single key on the keyboard or a combination of keys.                                                            |
| pressSequentially(text, [options])                     | Type text character by character, simulating real keyboard input.                                                       |
| selectOption(values, [options])  | Select one or more options which match the values.                                                                      |
| setChecked(checked[, options])                                | Sets the `checkbox` or `radio` input element's value to the specified checked or unchecked state.                       |
| tap([options])                            | Tap on the chosen element.                                                                                              |
| textContent([options])                                       | Returns the `element.textContent`.                                                                                      |
| type(text, [options])                                               | Type in the text into the input field.                                                                                  |
| uncheck([options])                    | Unselect the `input` checkbox.                                                                                          |
| waitFor([options])                    | Wait for the element to be in a particular state e.g. `visible`.                                                        |

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

  try {
    await page.goto('https://test.k6.io/flip_coin.php');

    /*
    In this example, we will use two locators, matching a
    different betting button on the page. If you were to query
    the buttons once and save them as below, you would see an
    error after the initial navigation. Try it!
  
      const heads = await page.$("input[value='Bet on heads!']");
      const tails = await page.$("input[value='Bet on tails!']");
  
    The Locator API allows you to get a fresh element handle each
    time you use one of the locator methods. And, you can carry a
    locator across frame navigations. Let's create two locators;
    each locates a button on the page.
    */
    const heads = page.locator("input[value='Bet on heads!']");
    const tails = page.locator("input[value='Bet on tails!']");

    const currentBet = page.locator("//p[starts-with(text(),'Your bet: ')]");

    // In the following Promise.all the tails locator clicks
    // on the tails button by using the locator's selector.
    // Since clicking on each button causes page navigation,
    // waitForNavigation is needed -- this is because the page
    // won't be ready until the navigation completes.
    // Setting up the waitForNavigation first before the click
    // is important to avoid race conditions.
    await Promise.all([page.waitForNavigation(), tails.click()]);
    console.log(await currentBet.innerText());
    // the heads locator clicks on the heads button
    // by using the locator's selector.
    await Promise.all([page.waitForNavigation(), heads.click()]);
    console.log(await currentBet.innerText());
    await Promise.all([page.waitForNavigation(), tails.click()]);
    console.log(await currentBet.innerText());
  } finally {
    await page.close();
  }
}
```

