FROM scratch

COPY tismd /
COPY config.yaml /

ENTRYPOINT ["/tismd"]
