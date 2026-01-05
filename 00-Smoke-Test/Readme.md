## Protohackers #0 - Smoke Test

### **Introduction**

The goal of the first challenge, **Smoke Test**, is to implement an **Echo Server**. According to the specification, the server must accept TCP connections and, whenever it receives data, it must transmit that same data back to the client. The connection should remain open until the client initiates a close.

### **The Challenge**

While the logic sounds straightforward, the main challenge lies in handling the **streaming nature of TCP**.

* **Reliability:** Every byte received must be echoed back.
* **Concurrency:** The server must handle multiple simultaneous connections without blocking.
* **Termination:** The server should only close the connection once it reaches the end of the input (EOF).

### **System Design**

To achieve high concurrency, I leveraged Go's concurrency model, specifically using goroutines managed by the GPM scheduler. This allows the server to handle thousands of simultaneous connections with minimal memory overhead.

### **Implementation Details**

The implementation can be broken down into a simple loop:

1. **Listen & Accept:** The server binds to a port and waits for incoming TCP segments.
2. **The Read-Write Loop:** * We read data into a **buffer**.
* We immediately write the contents of that buffer back to the socket.
* We repeat this until the client closes the connection (detected by a zero-byte read).



**Key Code Logic (Conceptual):**

```go
listen, err := net.Listen("tcp", ":8080")
if err != nil {
    return
}
defer listen.Close()

for {
    conn, err := listen.Accept()
    if err != nil {
        continue
    }

    go handle(conn)
}
```

### **Lessons Learned**

The most important takeaway from this "Smoke Test" is understanding that **TCP is a stream, not a packet**. There are no "message boundaries" here; we simply pipe the input stream directly to the output stream. This sets the foundation for more complex framing challenges in later problems.