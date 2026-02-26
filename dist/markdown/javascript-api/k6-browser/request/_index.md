
# Request

The request that the browser performs can be retrieved from the Response when a navigation occurs.

> **Caution:** This API is a work in progress. Some of the following functionalities might behave unexpectedly.

## Supported APIs

| Method                                                                                                                            | Description                                                                                                          |
| --------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------- |
| allHeaders()  | Returns an object of headers associated to the request including headers added by the browser.                       |
| frame()                                       | The Frame that initiated the request.   |
| headers()                                   | Returns an object of headers associated to the request.                                                              |
| headersArray()                         | An array with all the request HTTP headers.                                                                          |
| headerValue(name)                       | Returns the value of the header matching the name. The name is case insensitive.                                     |
| isNavigationRequest()           | Returns a boolean stating whether the request is for a navigation.                                                   |
| method()                                     | Request's method (GET, POST, etc.).                                                                                  |
| postData()                                 | Contains the request's post body, if any.                                                                            |
| postDataBuffer()                     | Request's post body in a binary form, if any.                                                                        |
| resourceType()                         | Contains the request's resource type as it was perceived by the rendering engine.                                    |
| response()                                 | Returns the matching Response object. |
| size()                                         | Returns an object containing the size of the request headers and body.                   |
| timing()                                     | Returns resource timing information for given request.                                                               |
| url()                                           | URL of the request.                                                                                                  |

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
  },
};

export default async function () {
  const page = await browser.newPage();

  try {
    const res = await page.goto('https://test.k6.io/');
    const req = res.request();
  } finally {
    await page.close();
  }
}
```

