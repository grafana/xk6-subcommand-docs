
# Response

Response represents a response received by the page.

> **Caution:** This API is a work in progress. Some of the following functionalities might behave unexpectedly.

## Supported APIs

| Method                                                                                                                             | Description                                                                                                                                             |
| ---------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------- |
| allHeaders()  | Returns an object of headers associated to the response including headers added by the browser.                                                         |
| body()                                         | Returns the response body.                                                                                                                              |
| frame()                                       | The Frame that initiated the request which this response is associated to. |
| headers()                                   | Returns an object of headers associated to the response.                                                                                                |
| headersArray()                         | An array with all the response HTTP headers.                                                                                                            |
| headerValue(name)                       | Returns the value of the header matching the name. The name is case insensitive.                                                                        |
| headerValues(name)                     | Returns all values of the headers matching the name, for example `set-cookie`. The name is case insensitive.                                            |
| json()                                         | Returns the JSON representation of response body.                                                                                                       |
| ok()                                             | Returns a `boolean` stating whether the response was successful or not.                                                                                 |
| request()                                   | Returns the matching Request object.                                      |
| securityDetails()                   | Returns SSL and other security information.                                                                                                             |
| serverAddr()                             | Returns the IP address and port of the server for this response.                                                                                        |
| status()                                     | Contains the status code of the response (e.g., 200 for a success).                                                                                     |
| statusText()                             | Contains the status text of the response (e.g. usually an "OK" for a success).                                                                          |
| size()                                         | The size of the response body and the headers.                                                                                                          |
| text()                                         | Returns the response body as a string.                                                                                                                  |
| url()                                           | URL of the response.                                                                                                                                    |

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
    // Response returned once goto resolves.
    const res = await page.goto('https://test.k6.io/');
  } finally {
    await page.close();
  }
}
```

