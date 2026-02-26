
# k6/net/grpc

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

## gRPC metrics

k6 takes specific measurements for gRPC requests.
For the complete list, refer to the Metrics reference.

## Example

```javascript
import grpc from 'k6/net/grpc';
import { check, sleep } from 'k6';

const client = new grpc.Client();
client.load(null, 'quickpizza.proto');

export default () => {
  client.connect('grpc-quickpizza.grafana.com:443', {
    // plaintext: false
  });

  const data = { ingredients: ['Cheese'], dough: 'Thick' };
  const response = client.invoke('quickpizza.GRPC/RatePizza', data);

  check(response, {
    'status is OK': (r) => r && r.status === grpc.StatusOK,
  });

  console.log(JSON.stringify(response.message));

  client.close();
  sleep(1);
};
```

