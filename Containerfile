FROM registry.redhat.io/ubi9/go-toolset AS builder

COPY . .
RUN go build -o ../better-sr-rss -buildvcs=false .

FROM registry.redhat.io/ubi9-micro
RUN mkdir /app && chown 1001 /app && chmod a+rwx /app
USER 1001
WORKDIR /app
COPY --from=builder /etc/pki/ca-trust /etc/pki/ca-trust
COPY --from=builder /opt/app-root/better-sr-rss .
COPY --from=builder /opt/app-root/src/templates ./templates

EXPOSE 8080
CMD /app/better-sr-rss
