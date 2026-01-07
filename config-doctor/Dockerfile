FROM bash:5.2

WORKDIR /app

# Copy validation scripts and config
COPY scripts/validate-config.sh /app/validate-config.sh
COPY config/deprecations.json /app/config/deprecations.json

# Make script executable
RUN chmod +x /app/validate-config.sh

# Default command runs validation
CMD ["/app/validate-config.sh"]

