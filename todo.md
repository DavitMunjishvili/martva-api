# Project TODO

A prioritized list of tasks to improve the project's quality, security, and maintainability.

---

## 🔴 Critical

- [ ] **Fix Security Vulnerability: Disable `InsecureSkipVerify` in HTTP Client**
  - **Description**: The HTTP client is configured with `InsecureSkipVerify: true`, which disables SSL/TLS certificate validation. This makes the application vulnerable to man-in-the-middle (MITM) attacks.
  - **Definition of Done**:
    - [ ] Remove `InsecureSkipVerify: true` from the HTTP client's transport configuration.
    - [ ] Ensure the application communicates successfully with the external API using proper SSL/TLS validation.
    - [ ] If the external API uses a self-signed certificate, import it into the application's trust store.

- [ ] **Implement Proper Configuration Management**
  - **Description**: The application uses hardcoded values for API endpoints and other configuration settings, which is a security risk and makes configuration for different environments difficult.
  - **Definition of Done**:
    - [ ] Use a library like Viper or envconfig to load configuration from environment variables or a file.
    - [ ] Remove all hardcoded configuration values from the code.
    - [ ] Create a sample configuration file (e.g., `.env.example`) to document required variables.

---

## 🟠 High

- [ ] **Add Unit and Integration Tests**
  - **Description**: The project lacks tests, making it risky to refactor or add new features.
  - **Definition of Done**:
    - [ ] Write unit tests for the service layer.
    - [ ] Write integration tests for the API handlers.
    - [ ] Set up a CI/CD pipeline to run tests automatically.
    - [ ] Achieve at least 80% test coverage.

- [ ] **Implement Structured Logging**
  - **Description**: The current unstructured logging makes it hard to search, filter, and analyze logs in a production environment.
  - **Definition of Done**:
    - [ ] Use a structured logging library (e.g., `slog`, `zerolog`) to produce JSON logs.
    - [ ] Include contextual information (e.g., request ID) in logs.
    - [ ] Configure the log level based on the environment.

---

## 🟡 Medium

- [ ] **Use a Router with Middleware Support**
  - **Description**: The default `http.ServeMux` is not suitable for production and lacks features like middleware support.
  - **Definition of Done**:
    - [ ] Replace `http.ServeMux` with a more capable router (e.g., `chi`, `gorilla/mux`).
    - [ ] Implement middleware for logging, error handling, and request validation.

- [ ] **Implement Graceful Shutdown**
  - **Description**: The application may terminate abruptly, interrupting ongoing requests.
  - **Definition of Done**:
    - [ ] Implement a graceful shutdown mechanism to finish processing existing requests before shutting down.
    - [ ] Handle OS signals (`SIGINT`, `SIGTERM`) to trigger the shutdown.

- [ ] **Add Input Validation**
  - **Description**: The `AvailableHoursHandler` does not validate its query parameters, which could lead to errors or vulnerabilities.
  - **Definition of Done**:
    - [ ] Implement input validation for all API handlers.
    - [ ] Return a `400 Bad Request` error for invalid input.

---

## 🔵 Low

- [ ] **Use a Linter and Formatter**
  - **Description**: The code lacks a consistent style, making it harder to read and maintain.
  - **Definition of Done**:
    - [ ] Set up a linter (e.g., `golangci-lint`) and formatter (e.g., `gofmt`).
    - [ ] Integrate the linter and formatter into the CI/CD pipeline.

- [ ] **Improve Error Handling**
  - **Description**: Error handling is inconsistent. A centralized mechanism would be cleaner and more maintainable.
  - **Definition of Done**:
    - [ ] Implement a centralized error handling middleware.
    - [ ] Define a standard error response format.
