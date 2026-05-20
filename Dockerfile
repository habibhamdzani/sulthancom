FROM node:22-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm install --legacy-peer-deps

ARG CACHEBUST=1
COPY . .
RUN npm run build

FROM node:22-alpine AS runner

WORKDIR /app

COPY package*.json ./
RUN npm install --omit=dev --legacy-peer-deps

COPY --from=builder /app/build ./build
COPY --from=builder /app/.svelte-kit ./.svelte-kit

ENV NODE_ENV=production
ENV PORT=3000

EXPOSE 3000

CMD ["node", "build/index.js"]
