
# toBeGreaterThanOrEqual()

The `toBeGreaterThanOrEqual()` method asserts that a numeric value is greater than or equal to another value.

## Syntax

```javascript
expect(actual).toBeGreaterThanOrEqual(expected);
expect(actual).not.toBeGreaterThanOrEqual(expected);
```

## Parameters

| Parameter | Type   | Description                  |
| --------- | ------ | ---------------------------- |
| expected  | number | The value to compare against |

## Returns

| Type | Description     |
| ---- | --------------- |
| void | No return value |

## Description

The `toBeGreaterThanOrEqual()` method performs a numeric comparison using the `>=` operator. Both values must be numbers, or the assertion will fail.

## Usage

```javascript
import { expect } from 'https://jslib.k6.io/k6-testing//index.js';

export default function () {
  expect(5).toBeGreaterThanOrEqual(3);
  expect(5).toBeGreaterThanOrEqual(5); // Equal values pass
  expect(10.5).toBeGreaterThanOrEqual(10);
  expect(0).toBeGreaterThanOrEqual(-1);
  expect(-1).toBeGreaterThanOrEqual(-5);
}
```

