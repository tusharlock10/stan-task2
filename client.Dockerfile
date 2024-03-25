# stage 1: build the client bundle
FROM node:latest as build
WORKDIR /app
COPY client/ ./
RUN npm install
RUN npm run build

# stage 2: servr the client using nginx
FROM nginx:alpine
COPY --from=build /app/build /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
