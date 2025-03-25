FROM node:22.12.0-alpine3.21 AS builder
ENV NODE_ENV production
WORKDIR /app
COPY ./package.json ./
RUN npm install
COPY . .
RUN npm run build-prod
FROM nginx:1.27.3-bookworm-perl
EXPOSE 3000
COPY --from=builder /app/dist /usr/share/nginx/html
COPY /nginx/nginx.conf /etc/nginx/conf.d/default.conf
