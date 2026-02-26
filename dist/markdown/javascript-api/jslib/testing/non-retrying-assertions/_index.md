
# Non-Retrying Assertions

Non-retrying assertions are synchronous assertions that allow to test any conditions, but do not auto-retry. They are ideal for testing static values, API responses, and any scenario where the expected condition should be true at the moment of evaluation.

Non-retrying assertions differ from retrying assertions in that they:

- **Evaluate immediately** - They check the condition once and return the result
- **Are synchronous** - They don't need to be awaited and return results immediately
- **Have no timeout behavior** - They either pass or fail instantly
- **Are ideal for static content** - Perfect for testing values that don't change over time

```javascript
import http from 'k6/http';
import { expect } from 'https://jslib.k6.io/k6-testing//index.js';

export default function () {
  const response = http.get('https://quickpizza.grafana.com/');

  // Immediate assertions
  expect(response.status).toBe(200);
  expect(response.body).toBeDefined();
  expect(response.body).toBeTruthy();
  expect(typeof response.body).toBe('string');
  expect(response.body).toContain('Pizza');
}
```

## When to Use Non-Retrying Assertions

Non-retrying assertions are best suited for:

- **API response testing** - Checking status codes, response data, headers
- **Static values** - Testing constants, configuration values, or computed results
- **Data validation** - Verifying object properties, array contents, or data types
- **Known state verification** - Checking values that should be immediately available

## Non-retrying assertions methods

### Equality Assertions

| Method                                                                                                                | Description                      |
| --------------------------------------------------------------------------------------------------------------------- | -------------------------------- |
| toBe()       | Exact equality using Object.is() |
| toEqual() | Deep equality comparison         |

### Truthiness Assertions

| Method                                                                                                                            | Description            |
| --------------------------------------------------------------------------------------------------------------------------------- | ---------------------- |
| toBeTruthy()       | Value is truthy        |
| toBeFalsy()         | Value is falsy         |
| toBeDefined()     | Value is not undefined |
| toBeUndefined() | Value is undefined     |
| toBeNull()           | Value is null          |

### Comparison Assertions

| Method                                                                                                                                              | Description                   |
| --------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------- |
| toBeGreaterThan()               | Numeric greater than          |
| toBeGreaterThanOrEqual() | Numeric greater than or equal |
| toBeLessThan()                     | Numeric less than             |
| toBeLessThanOrEqual()       | Numeric less than or equal    |
| toBeCloseTo()                       | Floating point comparison     |

### Collection Assertions

| Method                                                                                                                              | Description                                 |
| ----------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------- |
| toContain()           | Array/string contains value                 |
| toContainEqual() | Array contains object with matching content |
| toHaveLength()     | Array/string has specific length            |

### Property Assertions

| Method                                                                                                                              | Description                  |
| ----------------------------------------------------------------------------------------------------------------------------------- | ---------------------------- |
| toHaveProperty() | Object has specific property |

### Type Assertions

| Method                                                                                                                              | Description                |
| ----------------------------------------------------------------------------------------------------------------------------------- | -------------------------- |
| toBeInstanceOf() | Value is instance of class |

## Common Patterns

### API Response Validation

```javascript
export default function () {
  const response = http.get('https://quickpizza.grafana.com/');

  // Response status and headers
  expect(response.status).toBe(200);
  expect(response.headers['content-type']).toContain('text/html');

  // Response data structure
  expect(response.body).toBeDefined();
  expect(response.body).toBeTruthy();
  expect(response.body).toContain('pizza');
  expect(response.body).toContain('Pizza');

  // Data types
  expect(typeof response.body).toBe('string');
  expect(response.body.length).toBeGreaterThan(0);
}
```

### Data Validation

```javascript
export default function () {
  const userData = {
    name: 'John Doe',
    email: 'john@example.com',
    age: 30,
    hobbies: ['reading', 'gaming'],
    address: {
      city: 'New York',
      country: 'USA',
    },
  };

  // Object structure validation
  expect(userData).toHaveProperty('name');
  expect(userData).toHaveProperty('address.city');
  expect(userData.hobbies).toHaveLength(2);

  // Value validation
  expect(userData.name).toBeDefined();
  expect(userData.email).toContain('@');
  expect(userData.age).toBeGreaterThan(0);
  expect(userData.hobbies).toContain('reading');
}
```

### Configuration Testing

```javascript
export default function () {
  const config = {
    apiUrl: __ENV.API_URL || 'https://api.example.com',
    timeout: parseInt(__ENV.TIMEOUT || '5000'),
    retries: parseInt(__ENV.RETRIES || '3'),
    features: (__ENV.FEATURES || 'feature1,feature2').split(','),
  };

  // Configuration validation
  expect(config.apiUrl).toContain('http');
  expect(config.timeout).toBeGreaterThan(0);
  expect(config.retries).toBeGreaterThanOrEqual(0);
  expect(config.features).toContain('feature1');
  expect(config.features).toHaveLength(2);
}
```

## Error Handling Patterns

### Graceful Error Handling

```javascript
export default function () {
  const response = http.get('https://quickpizza.grafana.com/');

  // Check response status first
  if (response.status === 200) {
    expect(response.body).toBeDefined();
    expect(response.body).toBeTruthy();
    expect(response.body).toContain('Pizza');
  } else {
    // Handle error response
    expect(response.status).toBeGreaterThanOrEqual(400);
    expect(response.body).toContain('error');
  }
}
```

