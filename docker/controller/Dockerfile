FROM scratch

ENV NATS_URI "nats://nats:4222"
ENV DJINN_PLUGIN_DIR_PATH "/app/plugins"

COPY ./controller /app/controller

CMD ["/app/controller"]