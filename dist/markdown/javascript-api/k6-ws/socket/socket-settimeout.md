
# Socket.setTimeout(callback, delay)

> **Note:** A module with a better and standard API exists.
> 
> 
> The new [k6/experimental/websockets API](/docs/k6/v1.5.0/javascript-api/k6-experimental/websockets/) partially implements the [WebSockets API living standard](https://websockets.spec.whatwg.org/).
> 
> 
> When possible, we recommend using the new API. It uses a global event loop for consistency with other k6 APIs and better performance.

Call a function at a later time, if the WebSocket connection is still open then.

| Parameter | Type     | Description                                    |
| --------- | -------- | ---------------------------------------------- |
| callback  | function | The function to call when `delay` has expired. |
| delay     | number   | The delay time, in milliseconds.               |

### Example

```javascript
import ws from 'k6/ws';
import { sleep } from 'k6';

export default function () {
  console.log('T0: Script started');
  const url = 'wss://echo.websocket.org';
  const response = ws.connect(url, null, function (socket) {
    console.log('T0: Entered WebSockets run loop');
    socket.setTimeout(function () {
      console.log('T0+1: This is printed');
    }, 1000);
    socket.setTimeout(function () {
      console.log('T0+2: Closing socket');
      socket.close();
    }, 2000);
    socket.setTimeout(function () {
      console.log('T0+3: This is not printed, because socket is closed');
    }, 3000);
  });
  console.log('T0+2: Exited WebSockets run loop');
  sleep(2);
  console.log('T0+4: Script finished');
}
```

