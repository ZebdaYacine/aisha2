# syntax=docker/dockerfile:1.7

FROM node:22-alpine AS deps
WORKDIR /app
COPY package.json package-lock.json ./
RUN --mount=type=cache,target=/root/.npm \
    npm ci --no-audit --no-fund

FROM node:22-alpine AS build
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
ENV NEXT_TELEMETRY_DISABLED=1
RUN npm run build
RUN cp -L /usr/lib/libstdc++.so.6 /tmp/libstdc++.so.6

FROM alpine:3.24 AS runtime
ENV NODE_ENV=production
ENV HOSTNAME=0.0.0.0
COPY --from=build /usr/local/bin/node /usr/local/bin/node
COPY --from=build /usr/lib/libgcc_s.so.1 /usr/lib/libgcc_s.so.1
COPY --from=build /tmp/libstdc++.so.6 /usr/lib/libstdc++.so.6
RUN addgroup -S aisha && adduser -S -G aisha aisha
WORKDIR /app
COPY --from=build --chown=aisha:aisha /app/.next/standalone ./
COPY --from=build --chown=aisha:aisha /app/.next/static ./.next/static
COPY --from=build --chown=aisha:aisha /app/public ./public
USER aisha
EXPOSE 3000
CMD ["node","server.js"]
