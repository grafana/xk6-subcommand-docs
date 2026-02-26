
# HTML Forms

Scripting example on how to handle HTML forms.

In many cases using the Selection API (jQuery API clone) to interact with HTML data is enough, but for some use cases, like with forms, we can make things easier providing a higher-level API like the Response.submitForm( [params] ) API.

```javascript
import http from 'k6/http';
import { sleep } from 'k6';

export default function () {
  // Request page containing a form
  let res = http.get('https://quickpizza.grafana.com/admin');

  // Now, submit form setting/overriding some fields of the form
  res = res.submitForm({
    formSelector: 'form',
    fields: { username: 'admin', password: 'admin' },
  });
  sleep(3);
}
```

**Relevant k6 APIs**:

- Response.submitForm([params])
- Selection.find(selector)
  (the [jQuery Selector API](http://api.jquery.com/category/selectors/) docs are also a good
  resource on what possible selector queries can be made)
