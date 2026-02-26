
# toBeTruthy()

The `toBeTruthy()` method asserts that a value is truthy in JavaScript. A value is truthy if it converts to `true` when evaluated in a boolean context.

## Syntax

```javascript
expect(actual).toBeTruthy();
expect(actual).not.toBeTruthy();
```

## Returns

| Type | Description     |
| ---- | --------------- |
| void | No return value |

## Description

The `toBeTruthy()` method checks if a value is truthy. In JavaScript, all values are truthy except for the following falsy values:

- `false`
- `0`
- `-0`
- `0n` (BigInt)
- `""` (empty string)
- `null`
- `undefined`
- `NaN`

## Usage

```javascript
import { expect } from 'https://jslib.k6.io/k6-testing//index.js';

export default function () {
  expect(true).toBeTruthy();
  expect(1).toBeTruthy();
  expect(-1).toBeTruthy();
  expect('hello').toBeTruthy();
  expect(' ').toBeTruthy(); // Non-empty string
  expect(42).toBeTruthy();
  expect(3.14).toBeTruthy();
  expect(Infinity).toBeTruthy();
  expect(-Infinity).toBeTruthy();
}
```

