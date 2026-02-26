
# Single request

```javascript
import http from 'k6/http';

export const options = {
  iterations: 1,
};

export default function () {
  const response = http.get('https://quickpizza.grafana.com');
}
```

