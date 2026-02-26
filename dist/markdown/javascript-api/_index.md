
# JavaScript API

The list of k6 modules natively supported in your k6 scripts.

## Init context

Before the k6 starts the test logic, code in the _init context_ prepares the script.
A few functions are available only in init context.
For details about the runtime, refer to the Test lifecycle.

| Function                                                                                              | Description                                          |
| ----------------------------------------------------------------------------------------------------- | ---------------------------------------------------- |
| open( filePath, [mode] ) | Opens a file and reads all the contents into memory. |

## import.meta

`import.meta` is only available in ECMAScript modules, but not CommonJS ones.

| Function                                                                                           | Description                                               |
| -------------------------------------------------------------------------------------------------- | --------------------------------------------------------- |
| import.meta.resolve | resolve path to URL the same way that an ESM import would |

## k6

The `k6` module contains k6-specific functionality.

| Function                                                                                     | Description                                                                                                                                  |
| -------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------- |
| check(val, sets, [tags]) | Runs one or more checks on a value and generates a pass/fail result but does not throw errors or otherwise interrupt execution upon failure. |
| fail([err])               | Throws an error, failing and aborting the current VU script iteration immediately.                                                           |
| group(name, fn)          | Runs code inside a group. Used to organize results in a test.                                                                                |
| randomSeed(int)    | Set seed to get a reproducible pseudo-random number using `Math.random`.                                                                     |
| sleep(t)                 | Suspends VU execution for the specified duration.                                                                                            |

## k6/browser

The `k6/browser` module provides browser-level APIs to interact with browsers and collect frontend performance metrics as part of your k6 tests.

| Method                                                                                                                                      | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| ------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| browser.closeContext()                                   | Closes the current BrowserContext.                                                                                                                                                                                                                                                                                                                                          |
| browser.context()                                             | Returns the current BrowserContext.                                                                                                                                                                                                                                                                                                                                         |
| browser.isConnected            | Indicates whether the [CDP](https://chromedevtools.github.io/devtools-protocol/) connection to the browser process is active or not.                                                                                                                                                                                                                                                                                                                             |
| browser.newContext([options])  | Creates and returns a new BrowserContext.                                                                                                                                                                                                                                                                                                                                   |
| browser.newPage([options])         | Creates a new Page in a new BrowserContext and returns the page. Pages that have been opened ought to be closed using `Page.close`. Pages left open could potentially distort the results of Web Vital metrics. |
| browser.version()                                             | Returns the browser application's version.                                                                                                                                                                                                                                                                                                                                                                                                                       |

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

## k6/crypto

The `k6/crypto` module provides common hashing functionality available in the GoLang [crypto](https://golang.org/pkg/crypto/) package.

| Function                                                                                                                | Description                                                                                                                  |
| ----------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| createHash(algorithm)                   | Create a Hasher object, allowing the user to add data to hash multiple times, and extract hash digests along the way.        |
| createHMAC(algorithm, secret)           | Create an HMAC hashing object, allowing the user to add data to hash multiple times, and extract hash digests along the way. |
| hmac(algorithm, secret, data, outputEncoding) | Use HMAC to sign an input string.                                                                                            |
| md4(input, outputEncoding)                     | Use MD4 to hash an input string.                                                                                             |
| md5(input, outputEncoding)                     | Use MD5 to hash an input string.                                                                                             |
| randomBytes(int)                       | Return an array with a number of cryptographically random bytes.                                                             |
| ripemd160(input, outputEncoding)         | Use RIPEMD-160 to hash an input string.                                                                                      |
| sha1(input, outputEncoding)                   | Use SHA-1 to hash an input string.                                                                                           |
| sha256(input, outputEncoding)               | Use SHA-256 to hash an input string.                                                                                         |
| sha384(input, outputEncoding)               | Use SHA-384 to hash an input string.                                                                                         |
| sha512(input, outputEncoding)               | Use SHA-512 to hash an input string.                                                                                         |
| sha512_224(input, outputEncoding)       | Use SHA-512/224 to hash an input string.                                                                                     |
| sha512_256(input, outputEncoding)       | Use SHA-512/256 to hash an input string.                                                                                     |

| Class                                                                              | Description                                                                                                                                                                                           |
| ---------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Hasher | Object returned by crypto.createHash(). It allows adding more data to be hashed and to extract digests along the way. |

## k6/data

The `k6/data` module provides helpers to work with data.

| Class/Method                                                                               | Description                                                   |
| ------------------------------------------------------------------------------------------ | ------------------------------------------------------------- |
| SharedArray | read-only array like structure that shares memory between VUs |

## k6/encoding

The `k6/encoding` module provides [base64](https://en.wikipedia.org/wiki/Base64)
encoding/decoding as defined by [RFC4648](https://tools.ietf.org/html/rfc4648).

| Function                                                                                                                 | Description             |
| ------------------------------------------------------------------------------------------------------------------------ | ----------------------- |
| b64decode(input, [encoding], [format]) | Base64 decode a string. |
| b64encode(input, [encoding])           | Base64 encode a string. |

## k6/execution

The `k6/execution` module provides the capability to get information about the current test execution state inside the test script. You can read in your script the execution state during the test execution and change your script logic based on the current state.

`k6/execution` provides the test execution information with the following properties:

- instance
- scenario
- test
- vu

## k6/experimental

`k6/experimental` modules are stable modules that may introduce breaking changes. Once they become fully stable, they may graduate to become k6 core modules.

| Modules                                                                                          | Description                                                                                                                |
| ------------------------------------------------------------------------------------------------ | -------------------------------------------------------------------------------------------------------------------------- |
| csv               | Provides support for efficient and convenient parsing of CSV files.                                                        |
| fs                 | Provides a memory-efficient way to handle file interactions within your test scripts.                                      |
| streams       | Provides an implementation of the Streams API specification, offering support for defining and consuming readable streams. |
| websockets | Implements the browser's [WebSocket API](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket).                      |

## k6/html

The `k6/html` module contains functionality for HTML parsing.

| Function                                                                                    | Description                                                                                                                        |
| ------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------- |
| parseHTML(src) | Parse an HTML string and populate a Selection object. |

| Class                                                                                  | Description                                                                                                                        |
| -------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------- |
| Element     | An HTML DOM element as returned by the Selection API. |
| Selection | A jQuery-like API for accessing HTML DOM elements.                                                                                 |

## k6/http

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

## k6/metrics

The `k6/metrics` module provides functionality to create custom metrics of various types.

| Metric type                                                                           | Description                                                                                   |
| ------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------- |
| Counter | A metric that cumulatively sums added values.                                                 |
| Gauge     | A metric that stores the min, max and last values added to it.                                |
| Rate       | A metric that tracks the percentage of added values that are non-zero.                        |
| Trend     | A metric that calculates statistics on the added values (min, max, average, and percentiles). |

## k6/net/grpc

The `k6/net/grpc` module provides a [gRPC](https://grpc.io/) client for Remote Procedure Calls (RPC) over HTTP/2.

| Class/Method                                                                                                                                 | Description                                                                                                                                                                         |
| -------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Client                                                         | gRPC client used for making RPC calls to a gRPC Server.                                                                                                                             |
| Client.load(importPaths, ...protoFiles)            | Loads and parses the given protocol buffer definitions to be made available for RPC requests.                                                                                       |
| Client.connect(address [,params])               | Connects to a given gRPC service.                                                                                                                                                   |
| Client.invoke(url, request [,params])            | Makes an unary RPC for the given service/method and returns a Response.                             |
| Client.asyncInvoke(url, request [,params]) | Asynchronously makes an unary RPC for the given service/method and returns a Promise with Response. |
| Client.close()                                    | Close the connection to the gRPC service.                                                                                                                                           |
| Params                                                         | RPC Request specific options.                                                                                                                                                       |
| Response                                                     | Returned by RPC requests.                                                                                                                                                           |
| Constants                                                   | Define constants to distinguish between gRPC Response statuses.                                     |
| Stream(client, url, [,params])                                 | Creates a new gRPC stream.                                                                                                                                                          |
| Stream.on(event, handler)                            | Adds a new listener to one of the possible stream events.                                                                                                                           |
| Stream.write(message)                             | Writes a message to the stream.                                                                                                                                                     |
| Stream.end()                                        | Signals to the server that the client has finished sending.                                                                                                                         |
| EventHandler                                     | The function to call for various events on the gRPC stream.                                                                                                                         |
| Metadata                                      | The metadata of a gRPC stream’s message.                                                                                                                                            |

## k6/secrets

The `k6/secrets` module gives access to secrets provided by configured secret sources.

| Property                                                                                      | Description                                                                                         |
| --------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------- |
| get([String])       | asynchrounsly get a secret from the default secret source                                           |
| source([String]) | returns a source for the provided name that can than be used to get a secret from a concrete source |

## k6/timers

The `k6/timers` module implements timers to work with k6's event loop. They mimic the functionality found in browsers and other JavaScript runtimes.

| Function                                                                      | Description                                          |
| :---------------------------------------------------------------------------- | :--------------------------------------------------- |
| [setTimeout](https://developer.mozilla.org/en-US/docs/Web/API/setTimeout)     | Sets a function to be run after a given timeout.     |
| [clearTimeout](https://developer.mozilla.org/en-US/docs/Web/API/clearTimeout) | Clears a previously set timeout with `setTimeout`.   |
| [setInterval](https://developer.mozilla.org/en-US/docs/Web/API/setInterval)   | Sets a function to be run on a given interval.       |
| [clearInterval](https://developer.mozilla.org/en-US/docs/Web/API/setInterval) | Clears a previously set interval with `setInterval`. |

> **Note:** The timer methods are available globally, so you can use them in your script without including an import statement.

## k6/ws

The `k6/ws` module provides a [WebSocket](https://en.wikipedia.org/wiki/WebSocket) client implementing the [WebSocket protocol](http://www.rfc-editor.org/rfc/rfc6455.txt).

| Function                                                                                                  | Description                                                                                                                                                                                                                               |
| --------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| connect( url, params, callback ) | Create a WebSocket connection, and provides a Socket client to interact with the service. The method blocks the test finalization until the connection is closed. |

| Class/Method                                                                                                                      | Description                                                                                                                                                                    |
| --------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| Params                                                    | Used for setting various WebSocket connection parameters such as headers, cookie jar, compression, etc.                                                                        |
| Socket                                                    | WebSocket client used to interact with a WS connection.                                                                                                                        |
| Socket.close()                               | Close the WebSocket connection.                                                                                                                                                |
| Socket.on(event, callback)                      | Set up an event listener on the connection for any of the following events:- open- binaryMessage- message- ping- pong- close- error. |
| Socket.ping()                                 | Send a ping.                                                                                                                                                                   |
| Socket.send(data)                             | Send string data.                                                                                                                                                              |
| Socket.sendBinary(data)                 | Send binary data.                                                                                                                                                              |
| Socket.setInterval(callback, interval) | Call a function repeatedly at certain intervals, while the connection is open.                                                                                                 |
| Socket.setTimeout(callback, period)     | Call a function with a delay, if the connection is open.                                                                                                                       |

## crypto

The `crypto` module provides a WebCrypto API implementation.

| Class/Method                                                                                      | Description                                                                                                                                                                                                        |
| ------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| getRandomValues | Fills the passed `TypedArray` with cryptographically sound random values.                                                                                                                                          |
| randomUUID           | Returns a randomly generated, 36 character long v4 UUID.                                                                                                                                                           |
| subtle             | The SubtleCrypto interface provides access to common cryptographic primitives, such as hashing, signing, encryption, or decryption. |

> **Note:** The `crypto` object is available globally, so you can use it in your script without including an import statement.

## Error codes

The following specific error codes are currently defined:

- 1000: A generic error that isn't any of the ones listed below.
- 1010: A non-TCP network error - this is a place holder there is no error currently known to trigger it.
- 1020: An invalid URL was specified.
- 1050: The HTTP request has timed out.
- 1100: A generic DNS error that isn't any of the ones listed below.
- 1101: No IP for the provided host was found.
- 1110: Blacklisted IP was resolved or a connection to such was tried to be established.
- 1111: Blacklisted hostname using The Block Hostnames option.
- 1200: A generic TCP error that isn't any of the ones listed below.
- 1201: A "broken pipe" on write - the other side has likely closed the connection.
- 1202: An unknown TCP error - We got an error that we don't recognize but it is from the operating system and has `errno` set on it. The message in `error` includes the operation(write,read) and the errno, the OS, and the original message of the error.
- 1210: General TCP dial error.
- 1211: Dial timeout error - the timeout for the dial was reached.
- 1212: Dial connection refused - the connection was refused by the other party on dial.
- 1213: Dial unknown error.
- 1220: Reset by peer - the connection was reset by the other party, most likely a server.
- 1300: General TLS error
- 1310: Unknown authority - the certificate issuer is unknown.
- 1311: The certificate doesn't match the hostname.
- 1400 to 1499: error codes that correspond to the [HTTP 4xx status codes for client errors](https://en.wikipedia.org/wiki/List_of_HTTP_status_codes#4xx_Client_errors)
- 1500 to 1599: error codes that correspond to the [HTTP 5xx status codes for server errors](https://en.wikipedia.org/wiki/List_of_HTTP_status_codes#5xx_Server_errors)
- 1600: A generic HTTP/2 error that isn't any of the ones listed below.
- 1610: A general HTTP/2 GoAway error.
- 1611 to 1629: HTTP/2 GoAway errors with the value of the specific [HTTP/2 error code](https://tools.ietf.org/html/rfc7540#section-7) added to 1611.
- 1630: A general HTTP/2 stream error.
- 1631 to 1649: HTTP/2 stream errors with the value of the specific [HTTP/2 error code](https://tools.ietf.org/html/rfc7540#section-7) added to 1631.
- 1650: A general HTTP/2 connection error.
- 1651 to 1669: HTTP/2 connection errors with the value of the specific [HTTP/2 error code](https://tools.ietf.org/html/rfc7540#section-7) added to 1651.
- 1701: Decompression error.

Read more about Error codes.

## jslib

jslib is a collection of JavaScript libraries maintained by the k6 team that can be used in k6 scripts.

| Library                                                                                                                        | Description                                                                                                            |
| ------------------------------------------------------------------------------------------------------------------------------ | ---------------------------------------------------------------------------------------------------------------------- |
| aws                                                       | Library allowing to interact with Amazon AWS services                                                                  |
| httpx                                                   | Wrapper around k6/http to simplify session handling |
| k6chaijs                                             | BDD assertion style                                                                                                    |
| http-instrumentation-pyroscope | Library to instrument k6/http to send baggage headers for pyroscope to read back                                       |
| http-instrumentation-tempo         | Library to instrument k6/http to send tracing data                                                                     |
| testing                                               | Advanced assertion library with Playwright-inspired API for protocol and browser testing                             |
| totp                                                     | TOTP (Time-based One-Time Password) generation and verification                                                        |
| utils                                                   | Small utility functions useful in every day load testing                                                               |

