
# BrowserContext

`BrowserContext`s provide a way to operate multiple independent sessions, with separate pages, cache, and cookies. A default `BrowserContext` is created when a browser is launched.

The browser module API is used to create a new `BrowserContext`.

If a page opens another page, e.g. with a `window.open` call, the popup will belong to the parent page's `BrowserContext`.

| Method                                                                                                                                                      | Description                                                                                                                                                                                                                                          |
| ----------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| addCookies()                                                | Adds cookies into the `BrowserContext`.                                                                                                                  |
| addInitScript()                                          | Adds a script that will be evaluated on page creation, frame attached or a navigation occurs. |
| clearCookies()                                            | Clear the `BrowserContext`'s cookies.                                                                                                                    |
| clearPermissions()                                    | Clears all permission overrides for the `BrowserContext`.                                                                                                                                                                                            |
| cookies()                                                      | Returns a list of cookies from the `BrowserContext`.                                                                                                     |
| close()                                                          | Close the `BrowserContext` and all its pages.                                                                                                                             |
| grantPermissions(permissions[, options])              | Grants specified permissions to the `BrowserContext`.                                                                                                                                                                                                |
| newPage()                                                      | Uses the `BrowserContext` to create a new Page and returns it.                                                                                                            |
| pages()                               | Returns a list of pages that belongs to the `BrowserContext`.                                                                                                             |
| setDefaultNavigationTimeout(timeout)       | Sets the default navigation timeout in milliseconds.                                                                                                                                                                                                 |
| setDefaultTimeout(timeout)                           | Sets the default maximum timeout for all methods accepting a timeout option in milliseconds.                                                                                                                                                         |
| setGeolocation(geolocation)  | Sets the `BrowserContext`'s geolocation.                                                                                                                                                                                                             |
| setOffline(offline)                                         | Toggles the `BrowserContext`'s connectivity on/off.                                                                                                                                                                                                  |
| waitForEvent(event[, optionsOrPredicate])                 | Waits for the event to fire and passes its value into the predicate function.                                                                                                                                                                        |
