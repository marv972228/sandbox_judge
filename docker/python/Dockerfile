# Python runner for Sandbox Judge
# Minimal image for secure, fast code execution

FROM python:3.11-slim

# Create non-root user for security
RUN useradd --create-home --shell /bin/bash runner

# Set working directory
WORKDIR /sandbox

# Install common packages that problems might need
# Keep minimal for security - can expand later
RUN pip install --no-cache-dir \
    sortedcontainers \
    && rm -rf /var/lib/apt/lists/*

# Switch to non-root user
USER runner

# Default command - will be overridden
CMD ["python3", "-c", "print('Ready')"]
