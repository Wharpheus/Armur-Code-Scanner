# Use a base image with Python
FROM python:3.10-slim

# Set environment variables
ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1

# Set the working directory
WORKDIR /armur

# Install required dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends git curl build-essential gcc && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Install Go (use the amd64 version for x86-64 architecture)
RUN curl -OL https://go.dev/dl/go1.23.1.linux-amd64.tar.gz && \
    tar -C /usr/local -xvf go1.23.1.linux-amd64.tar.gz && \
    rm go1.23.1.linux-amd64.tar.gz

# Add Go to PATH
ENV PATH="/usr/local/go/bin:/root/go/bin:${PATH}"

# Install gosec
RUN go install github.com/securego/gosec/v2/cmd/gosec@v2.20.0

# Install golint
RUN go install golang.org/x/lint/golint@latest

# Install staticcheck
RUN go install honnef.co/go/tools/cmd/staticcheck@latest

# Install gocyclo
RUN go install github.com/fzipp/gocyclo/cmd/gocyclo@latest

# Install deadcode
RUN go install golang.org/x/tools/cmd/deadcode@latest

# Install osv-scanner
RUN go install github.com/google/osv-scanner/cmd/osv-scanner@latest

# Install Trivy
RUN curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b /usr/local/bin v0.55.2

# Install Node.js and npm using NodeSource
RUN curl -fsSL https://deb.nodesource.com/setup_22.x | bash - && \
    apt-get install -y nodejs && \
    npm install -g npm@latest

# Copy the current directory contents into the container at /app
COPY . /armur

# Install global npm packages
RUN npm install -g eslint

# Install jscpd
RUN npm install -g jscpd

# Install project-specific npm packages locally
RUN npm install @eslint/js eslint-plugin-jsdoc eslint-plugin-security

# Install Python dependencies
# RUN pip install --no-cache-dir -r requirements.txt

# Install Semgrep
RUN pip install semgrep

# Install Bandit
RUN pip install bandit

# Install pydocstyle
RUN pip install pydocstyle

# Install radon
RUN pip install radon

# Install pylint
RUN pip install pylint

# Install truffleHog
RUN pip install trufflehog3

# Install checkov
RUN pip install checkov

# Install vulture
RUN pip install vulture

# Copy the ESLint configuration files
COPY /rule_config/eslint/eslint.config.js /armur/eslint.config.js
COPY /rule_config/eslint/eslint_jsdoc.config.js /armur/eslint_jsdoc.config.js
COPY /rule_config/eslint/eslint_security.config.js /armur/eslint_security.config.js
COPY /rule_config/eslint/eslint_deadcode.config.js /armur/eslint_deadcode.config.js

# Expose the port that the app runs on
EXPOSE 4500

# Run the application
CMD ["go", "run", "./cmd/server/main.go"]