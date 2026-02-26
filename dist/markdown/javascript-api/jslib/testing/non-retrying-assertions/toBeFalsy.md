
# toBeFalsy()

The `toBeFalsy()` method asserts that a value is falsy in JavaScript. A value is falsy if it converts to `false` when evaluated in a boolean context.

## Syntax

```javascript
expect(actual).toBeFalsy();
expect(actual).not.toBeFalsy();
```

## Returns

| Type | Description     |
| ---- | --------------- |
| void | No return value |

## Description

The `toBeFalsy()` method checks if a value is falsy. In JavaScript, the following values are falsy:

- `false`
- `0`
- `-0`
- `0n` (BigInt)
- `""` (empty string)
- `null`
- `undefined`
- `NaN`

All other values are truthy.

## Usage

```javascript
import { expect } from 'https://jslib.k6.io/k6-testing//index.js';

export default function () {
  expect(false).toBeFalsy();
  expect(0).toBeFalsy();
  expect(-0).toBeFalsy();
  expect(0n).toBeFalsy();
  expect('').toBeFalsy();
  expect(null).toBeFalsy();
  expect(undefined).toBeFalsy();
  expect(NaN).toBeFalsy();
}
```

