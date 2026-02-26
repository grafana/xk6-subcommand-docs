
# Init context

Before the k6 starts the test logic, code in the _init context_ prepares the script.
A few functions are available only in init context.
For details about the runtime, refer to the Test lifecycle.

| Function                                                                                              | Description                                          |
| ----------------------------------------------------------------------------------------------------- | ---------------------------------------------------- |
| open( filePath, [mode] ) | Opens a file and reads all the contents into memory. |

