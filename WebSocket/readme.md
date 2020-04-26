# Websocket Example

see https://github.com/golang-samples/websocket/blob/master/cli/src/server.go



client wählt sich bei server mit connection point localhost:8080/echo ein


|        |---------------->|        |
| client |                 | server |  
|        |<----------------|        |


A cell has a server and a client combined together


      cell A                          cell B
+----------------+              +----------------+
|                |              |                |
o server  client |------------->o /echo          |-------------->
|                |              | server  client |
+----------------+              +----------------+
  localhost:8001                  localhost:8002
