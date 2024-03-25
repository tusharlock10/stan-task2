# Assignment Task 2
The project demonstrates the following -
- Golang chat server
- Using goroutines to manage concurrent connections in large quantities
- Reconnecting and disconnecting a client
- Saving messages in MongoDB and retrieving them
- Using prometheus and grafana for monitoring the no. of clients
- Using docker with and docker-compose to deploy and run

## Installations required
1. Docker

## Files explained
1. **server/main.go** used to initialize services like monitoring and database and create a http sever
2. **server/services/chat.go** used to create sockets and handle them concurrently
- accepts a socket connection from a client
- gracefully disconnects the clients and exits the goroutine
- reconnects the client if the client connects with the same username
- receives messages from a user broadcasts its to all users
- saves the message in mongodb
- updates prometheus if a client connects or disconnects for monitoring
3. **server/services/database.go** has utility functions for initializing database and creating and getting and inserting messages
4. **server/services/messages.go** an http route handle to respond will 100 latest messages in the database
5. **server/services/monitoring.go** has utility function for initializing prometheus monitoring
6. **server.Dockerfile** multi stage docker file for building the server and running it on port 8080
6. **client.Dockerfile** multi stage docker file for building the client and serving it on port 3000
7. **docker-compose.yml** builds the server and client and run them with mongodb, prometheus and grafana
8. **client/*** contains a react client to demonstrate the chat app

## Run
`docker compose up -d --build`

## Testing
1. The server is running at port `localhost:8080`.
- Get messages using `http://localhost:8080/messages`
- Connect to WS using `ws://localhost:8080/chat?username=<user_name>`
2. Use the Frontend UI on port `localhost:3000`. Run it on multiple tabs with different username
- Sending a message on one browser tab will make it receive on others
3. Monitor connected clients on Grafana
- Goto `http://localhost:9090`
- Login with username: `user` and password: `chatapp`
- Prometheus datasource is already configured, You need to create a dashboard
- To create a dashboard go to `Dashboard > Create Dashboard > Add Visualization > Select prometheus datasource`
- In the `Select metric` dropdown, choose `chatapp_connected_clients`. Click on apply
- In the top right set the refresh option to 5s for faster refreshes
- You should see the number of connections update every few seconds