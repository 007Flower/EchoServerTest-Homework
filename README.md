# EchoServerTest-Homework
 
 You Run this Sever By  tpying the command go run server.go --port 5000"in a ternminal".
 Then open a new terminal an runn the command nc localhost 4000


 Which functionality was the most educationally enriching?

My most valuable learning experience came from implementing the inactivity timeout and working with time. Timer and managing it alongside goroutines helped me better understand how Go handles concurrency. Figuring out how to properly reset the timer and handle situations where a client goes silent without crashing the server gave me a clearer picture of how asynchronous operations and cleanup should be managed in a real-world scenario.

Which functionality required the most research?

The per-client logging. I had to dig into how to safely turn a client's IP and port into a valid filename, especially considering that characters like colons aren't allowed in filenames on some systems. I also needed to ensure that multiple clients could log into their files simultaneously without running into conflicts or race conditions and learning to use os.openfile properly.OpenFile with the correct flags and handling file writing line-by-line played a big part. I spent a reasonable amount of time researching best practices for logging from multiple goroutines while keeping things efficient and stable.
