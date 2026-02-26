
# k6/http

The `k6/http` module contains functionality for performing HTTP transactions.

| Function                                                                                                                       | Description                                                                                                                               |
| ------------------------------------------------------------------------------------------------------------------------------ | ----------------------------------------------------------------------------------------------------------------------------------------- |
| batch( requests )                                     | Issue multiple HTTP requests in parallel (like e.g. browsers tend to do).                                                                 |
| cookieJar()                                | Get active HTTP Cookie jar.                                                                                                               |
| del( url, [body], [params] )                            | Issue an HTTP DELETE request.                                                                                                             |
| file( data, [filename], [contentType] )                | Create a file object that is used for building multi-part requests.                                                                       |
| get( url, [params] )                                    | Issue an HTTP GET request.                                                                                                                |
| head( url, [params] )                                  | Issue an HTTP HEAD request.                                                                                                               |
| options( url, [body], [params] )                    | Issue an HTTP OPTIONS request.                                                                                                            |
| patch( url, [body], [params] )                        | Issue an HTTP PATCH request.                                                                                                              |
| post( url, [body], [params] )                          | Issue an HTTP POST request.                                                                                                               |
| put( url, [body], [params] )                            | Issue an HTTP PUT request.                                                                                                                |
| request( method, url, [body], [params] )            | Issue any type of HTTP request.                                                                                                           |
| asyncRequest( method, url, [body], [params] )  | Issue any type of HTTP request asynchronously.                                                                                            |
| setResponseCallback(expectedStatuses) | Sets a response callback to mark responses as expected.                                                                                   |
| url\`url\`                                              | Creates a URL with a name tag. Read more on URL Grouping. |
| expectedStatuses( statusCodes )           | Create a callback for setResponseCallback that checks status codes.                                                                       |

| Class                                                                                  | Description                                                                              |
| -------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------- |
| CookieJar | Used for storing cookies, set by the server and/or added by the client.                  |
| FileData   | Used for wrapping data representing a file when doing multipart requests (file uploads). |
| Params       | Used for setting various HTTP request-specific parameters such as headers, cookies, etc. |
| Response   | Returned by the http.\* methods that generate HTTP requests.                             |

