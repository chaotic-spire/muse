FROM oven/bun:1 AS builder

WORKDIR /app

COPY ../backend-search/package.json bun.lock ./
RUN bun install --frozen-lockfile

COPY ../backend-search/main.ts main.ts

RUN bun build main.ts --compile --outfile server

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/server ./server

EXPOSE 3000
CMD ["./server"]