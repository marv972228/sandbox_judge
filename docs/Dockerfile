# MkDocs documentation server
FROM python:3.11-slim

# Install MkDocs with Material theme and plugins
RUN pip install --no-cache-dir \
    mkdocs==1.5.3 \
    mkdocs-material==9.5.3 \
    mkdocs-minify-plugin==0.7.2 \
    pymdown-extensions==10.7

WORKDIR /docs

# Copy documentation source
COPY content/ /docs/

# Expose MkDocs dev server port
EXPOSE 8000

# Default command: serve docs with live reload
CMD ["mkdocs", "serve", "--dev-addr", "0.0.0.0:8000"]
