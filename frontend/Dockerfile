FROM node:22-alpine AS deps
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
FROM node:22-alpine AS build
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN npm run build
FROM node:22-alpine AS runtime
ENV NODE_ENV=production
ENV HOSTNAME=0.0.0.0
RUN addgroup -S aisha && adduser -S -G aisha aisha
WORKDIR /app
COPY --from=build --chown=aisha:aisha /app/.next/standalone ./
COPY --from=build --chown=aisha:aisha /app/.next/static ./.next/static
COPY --from=build --chown=aisha:aisha /app/public ./public
USER aisha
EXPOSE 3000
CMD ["node","server.js"]
