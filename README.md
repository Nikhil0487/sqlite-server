# sqlite-server
Simple HTTP based SQLite server which serialises writes.

For an application with mutiple processes wanting to write to a single DB, we can

1. Enable WAL mode to handle concurrent reads within the application,
2. Serialize writes by sending SQL writes to this this SQLite server,

Internally, 

1. We setup a HTTP server exposing an endpoint,
2. Get multipe SQL write requests, which is serialized by a go routine,


Potential Issues:

1. Yet to bind to local host,
2. Some form of authnetication to HTTP requests to prevent attacks leading to DB corruption,
3. Executable size. Its ~13 MB in macOS Big Sur. We can compress using tools like UPX.


TODO:

- [ ] Can we handle SQL READs and formulate a way to respond data needed precisely for the query
- [ ] Authenticate DB requests
- [ ] Bind HTTP server to local host 
