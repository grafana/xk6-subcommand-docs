
# isConnected()

> **Caution:** This feature has **known issues**.
> For details, refer to [#453](https://github.com/grafana/xk6-browser/issues/453).

Indicates whether the [CDP](https://chromedevtools.github.io/devtools-protocol/) connection to the browser process is active or not.

### Returns

| Type    | Description                                                                                                                                                                                  |
| ------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| boolean | Returns `true` if the browser module is connected to the browser application. Otherwise, returns `false`. |

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

export default function () {
  const isConnected = browser.isConnected();
  console.log(isConnected); // true
}
```

