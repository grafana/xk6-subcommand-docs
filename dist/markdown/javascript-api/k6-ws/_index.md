
# k6/ws

> **Note:** A module with a better and standard API exists.
> 
> 
> The new [k6/experimental/websockets API](/docs/k6/v1.5.0/javascript-api/k6-experimental/websockets/) partially implements the [WebSockets API living standard](https://websockets.spec.whatwg.org/).
> 
> 
> When possible, we recommend using the new API. It uses a global event loop for consistency with other k6 APIs and better performance.

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

