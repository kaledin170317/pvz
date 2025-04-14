FROM ubuntu:latest
LABEL authors="kaled"

ENTRYPOINT ["top", "-b"]