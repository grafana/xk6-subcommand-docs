
# browser

The browser module APIs are inspired by Playwright and other frontend testing frameworks.

You can find examples of using [the browser module API](#browser-module-api) in the getting started guide.

> **Note:** To work with the browser module, make sure you are using the latest [k6 version](https://github.com/grafana/k6/releases).

## Properties

The table below lists the properties you can import from the browser module (`'k6/browser'`).

| Property | Description                                                                                                                                                                          |
| -------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| browser  | The browser module API is the entry point for all your tests. See the [example](#example) and the [API](#browser-module-api) below.                                                  |
| devices  | Returns predefined emulation settings for many end-user devices that can be used to simulate browser behavior on a mobile device. See the [devices example](#devices-example) below. |

## Browser Module API

The browser module is the entry point for all your tests, and it is what interacts with the actual web browser via [Chrome DevTools Protocol](https://chromedevtools.github.io/devtools-protocol/) (CDP). It manages:

- BrowserContext which is where you can set a variety of attributes to control the behavior of pages;
- and Page which is where your rendered site is displayed.

| Method                                                                                                                                      | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| ------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| browser.closeContext()                                   | Closes the current BrowserContext.                                                                                                                                                                                                                                                                                                                                          |
| browser.context()                                             | Returns the current BrowserContext.                                                                                                                                                                                                                                                                                                                                         |
| browser.isConnected            | Indicates whether the [CDP](https://chromedevtools.github.io/devtools-protocol/) connection to the browser process is active or not.                                                                                                                                                                                                                                                                                                                             |
| browser.newContext([options])  | Creates and returns a new BrowserContext.                                                                                                                                                                                                                                                                                                                                   |
| browser.newPage([options])         | Creates a new Page in a new BrowserContext and returns the page. Pages that have been opened ought to be closed using `Page.close`. Pages left open could potentially distort the results of Web Vital metrics. |
| browser.version()                                             | Returns the browser application's version.                                                                                                                                                                                                                                                                                                                                                                                                                       |

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
  thresholds: {
    checks: ['rate==1.0'],
  },
};

export default async function () {
  const page = await browser.newPage();

  try {
    await page.goto('https://test.k6.io/');
  } finally {
    await page.close();
  }
}
```

Then, you can run the test with this command. Also, see the browser module options for customizing the browser module's behavior using environment variables.

```bash
k6 run script.js
```

```docker
# WARNING!
# The grafana/k6:master-with-browser image launches a Chrome browser by setting the
# 'no-sandbox' argument. Only use it with trustworthy websites.
#
# As an alternative, you can use a Docker SECCOMP profile instead, and overwrite the
# Chrome arguments to not use 'no-sandbox' such as:
# docker container run --rm -i -e K6_BROWSER_ARGS='' --security-opt seccomp=$(pwd)/chrome.json grafana/k6:master-with-browser run - <script.js
#
# You can find an example of a hardened SECCOMP profile in:
# https://raw.githubusercontent.com/jfrazelle/dotfiles/master/etc/docker/seccomp/chrome.json.
docker run --rm -i grafana/k6:master-with-browser run - <script.js
```

```windows
k6 run script.js
```

```windows-powershell
k6 run script.js
```

### Devices example

To emulate the browser behaviour on a mobile device and approximately measure the browser performance, you can import `devices` from `k6/browser`.

```javascript
import { browser, devices } from 'k6/browser';

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
  thresholds: {
    checks: ['rate==1.0'],
  },
};

export default async function () {
  const iphoneX = devices['iPhone X'];
  const context = await browser.newContext(iphoneX);
  const page = await context.newPage();

  try {
    await page.goto('https://test.k6.io/');
  } finally {
    page.close();
  }
}
```

## Browser-level APIs

| k6 Class                                                                                                               | Description                                                                                                                                              |
| ---------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------- |
| BrowserContext  | Enables independent browser sessions with separate Pages, cache, and cookies. |
| ElementHandle    | Represents an in-page DOM element.                                                                                                                       |
| Frame                    | Access and interact with the `Page`.'s `Frame`s.                              |
| JSHandle                                | Represents an in-page JavaScript object.                                                                                                                 |
| Keyboard                                | Used to simulate the keyboard interactions with the associated `Page`.        |
| Locator                                  | The Locator API makes it easier to work with dynamically changing elements.                                                                              |
| Mouse                                      | Used to simulate the mouse interactions with the associated `Page`.           |
| Page                      | Provides methods to interact with a single tab in a browser.                                                                                             |
| Request                | Used to keep track of the request the `Page` makes.                           |
| Response              | Represents the response received by the `Page`.                               |
| Touchscreen                          | Used to simulate touch interactions with the associated `Page`.               |
| Worker                                    | Represents a [WebWorker](https://developer.mozilla.org/en-US/docs/Web/API/Web_Workers_API).                                                              |

## Browser module options

You can customize the behavior of the browser module by providing browser options as environment variables.

| Environment Variable           | Description                                                                                                                                                                                                                                                                                                                                                              |
| ------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| K6_BROWSER_ARGS                | Extra command line arguments to include when launching browser process. See [this link](https://peter.sh/experiments/chromium-command-line-switches/) for a list of Chromium arguments. Note that arguments should not start with `--` (see the command example below).                                                                                                  |
| K6_BROWSER_DEBUG               | All CDP messages and internal fine grained logs will be logged if set to `true`.                                                                                                                                                                                                                                                                                         |
| K6_BROWSER_EXECUTABLE_PATH     | Override search for browser executable in favor of specified absolute path.                                                                                                                                                                                                                                                                                              |
| K6_BROWSER_HEADLESS            | Show browser GUI or not. `true` by default.                                                                                                                                                                                                                                                                                                                              |
| K6_BROWSER_IGNORE_DEFAULT_ARGS | Ignore any of the default arguments included when launching a browser process.                                                                                                                                                                                                   |
| K6_BROWSER_TIMEOUT             | Default timeout for initializing the connection to the browser instance. `'30s'` if not set.                                                                                                                                                                                                                                                                             |
| K6_BROWSER_TRACES_METADATA     | Sets additional _key-value_ metadata that is included as attributes in every span generated from browser module traces. Example: `K6_BROWSER_TRACES_METADATA=attr1=val1,attr2=val2`. This only applies if traces generation is enabled, refer to Traces output for more details. |

The following command passes the browser options as environment variables to launch a headful browser with custom arguments.

```bash
K6_BROWSER_HEADLESS=false K6_BROWSER_ARGS='show-property-changed-rects' k6 run script.js
```

```docker
# WARNING!
# The grafana/k6:master-with-browser image launches a Chrome browser by setting the
# 'no-sandbox' argument. Only use it with trustworthy websites.
#
# As an alternative, you can use a Docker SECCOMP profile instead, and overwrite the
# Chrome arguments to not use 'no-sandbox' such as:
# docker container run --rm -i -e K6_BROWSER_ARGS='' --security-opt seccomp=$(pwd)/chrome.json grafana/k6:master-with-browser run - <script.js
#
# You can find an example of a hardened SECCOMP profile in:
# https://raw.githubusercontent.com/jfrazelle/dotfiles/master/etc/docker/seccomp/chrome.json.
docker run --rm -i -e K6_BROWSER_HEADLESS=false -e K6_BROWSER_ARGS='show-property-changed-rects' grafana/k6:master-with-browser run - <script.js
```

```windows
set "K6_BROWSER_HEADLESS=false" && set "K6_BROWSER_ARGS='show-property-changed-rects' " && k6 run script.js
```

```windows-powershell
$env:K6_BROWSER_HEADLESS="false" ; $env:K6_BROWSER_ARGS='show-property-changed-rects' ; k6 run script.js
```

