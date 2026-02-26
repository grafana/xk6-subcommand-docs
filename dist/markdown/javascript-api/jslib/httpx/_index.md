
# httpx

> **Note:** The source code for this library can be found in the [grafana/k6-jslib-httpx](https://github.com/k6io/k6-jslib-httpx) GitHub repository.

The `httpx` module is an external JavaScript library that wraps around the native k6/http module.
It's an HTTP client with features that are not yet available in the native module.

- Ability to set HTTP options globally (such as timeout).
- Ability to set default tags and headers that will be used for all requests.
- More user-friendly arguments to request functions (get, post, put take the same arguments).

The `httpx` module integrates well with the expect library.

> **Caution:** This library is in active development. It's stable enough to be useful, but you can watch the [GitHub repository](https://github.com/k6io/k6-jslib-httpx) to be notified when a new version is released.

### Methods

| Function                                                                                                                        | Description                                                        |
| ------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------ |
| asyncRequest(method, url, [body], [params]) | Generic method for making arbitrary, asynchronous HTTP requests.   |
| request(method, url, [body], [params])           | Generic method for making arbitrary HTTP requests.                 |
| get(url, [body], [params])                           | Makes GET request                                                  |
| post(url, [body], [params])                         | Makes POST request                                                 |
| put(url, [body], [params])                           | Makes PUT request                                                  |
| patch(url, [body], [params])                       | Makes PATCH request                                                |
| delete(url, [body], [params])                     | Makes DELETE request                                               |
| batch(requests)                                    | Batches multiple HTTP requests together to issue them in parallel. |
| setBaseUrl(url)                               | Sets the base URL for the session                                  |
| addHeader(key, value)                          | Adds a header to the session                                       |
| addHeaders(object)                            | Adds multiple headers to the session                               |
| clearHeader(name)                            | Removes header from the session                                    |
| addTag(key, value)                                | Adds a tag to the session                                          |
| addTags(object)                                  | Adds multiple tags to the session                                  |
| clearTag(name)                                  | Removes tag from the session                                       |

### Example

```javascript
import { fail } from 'k6';
import { Httpx } from 'https://jslib.k6.io/httpx/0.1.0/index.js';
import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

export const options = {
  vus: 1,
  iterations: 1,
};

const USERNAME = `user${randomIntBetween(1, 100000)}@example.com`; // random email address
const PASSWORD = 'secretpassword';

const session = new Httpx({
  baseURL: 'https://quickpizza.grafana.com',
  headers: {
    'User-Agent': 'My custom user agent',
  },
  timeout: 20000, // 20s timeout.
});

export default function testSuite() {
  const registrationResp = session.post(
    `/api/users`,
    JSON.stringify({
      username: USERNAME,
      password: PASSWORD,
    })
  );

  if (registrationResp.status !== 201) {
    fail('registration failed');
  }

  const loginResp = session.post(
    `/api/users/token/login`,
    JSON.stringify({
      username: USERNAME,
      password: PASSWORD,
    })
  );

  if (loginResp.status !== 200) {
    fail('Authentication failed');
  }

  const authToken = loginResp.json('token');

  // set the authorization header on the session for the subsequent requests.
  session.addHeader('Authorization', `Bearer ${authToken}`);

  const payload = {
    stars: 5,
    pizza_id: 1,
  };

  // this request uses the Authorization header set above.
  const respCreateRating = session.post(`/api/ratings`, JSON.stringify(payload));

  if (respCreateRating.status !== 201) {
    fail('Rating creation failed');
  } else {
    console.log('New rating created');
  }
}
```

