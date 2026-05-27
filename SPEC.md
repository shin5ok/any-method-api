# SPEC: Security Hardening for any-method-api

## 1. Background & Goals

The `any-method-api` is a simple Go API server used for SLO monitoring demonstrations. To ensure the application adheres to maximum security standards ("secure the best"), we will perform security hardening in the following areas without compromising its core telemetry and demonstration capabilities:

*   **Prevent DoS Attacks (Slowloris)**: Configure explicit timeouts on the HTTP server instead of using default configurations (which have no timeouts).
*   **Enhance Response Headers**: Inject essential security headers (e.g., `X-Content-Type-Options`, `X-Frame-Options`) to mitigate cross-site scripting (XSS) and clickjacking risks.
*   **Prevent Information Leakage in Logs**: Mask or filter sensitive HTTP request headers (e.g., `Authorization`, `Cookie`) in logging sinks.
*   **Restrict HTTP Methods**: Restrict incoming HTTP methods to a safe allow-list, blocking potentially dangerous or unused methods (e.g., `TRACE`, `CONNECT`).
*   **Bind Address Configuration**: Allow binding to `127.0.0.1` for secure local testing, while keeping containerized deployment capability (binding to `0.0.0.0` via configuration/env).

## 2. Detailed Specifications

### 2.1. HTTP Server Configuration & Timeouts
The application currently starts the server with `http.ListenAndServe(":"+port, r)`. We will replace this with an explicit `http.Server` initialization:

*   **ReadHeaderTimeout**: `2 * time.Second` (Crucial mitigation for Slowloris attacks where headers are sent very slowly).
*   **ReadTimeout**: `5 * time.Second` (Total duration to read the entire request including the body).
*   **WriteTimeout**: `10 * time.Second` (Maximum duration before timing out writes of the response).
*   **IdleTimeout**: `30 * time.Second` (Keep-alive connection idle timeout).
*   **MaxHeaderBytes**: `1 << 20` (1 MB maximum request header size).

### 2.2. Security Headers Middleware
We will implement a standard HTTP middleware that injects the following security headers into all HTTP responses:

*   `X-Content-Type-Options: nosniff` (Prevents MIME sniffing)
*   `X-Frame-Options: DENY` (Mitigates clickjacking)
*   `X-XSS-Protection: 1; mode=block` (Enables XSS filtering in older browsers)
*   `Content-Security-Policy: default-src 'none'; frame-ancestors 'none';` (Restricts loading resources and embedding in iframes)
*   `Referrer-Policy: no-referrer` (Restricts leaking referrer data)

### 2.3. Log Masking of Sensitive Data
When logging requests in `commonHandler`, headers are serialized directly to the logs:
`Interface("headers", headers)`

We will modify this behavior to:
*   Filter or redact known sensitive headers: `Authorization`, `Cookie`, `Set-Cookie`, `Proxy-Authorization`, `X-Api-Key`.
*   Ensure that `r.Body` parsing and logging do not crash or consume excessive resources (limit body size to prevent memory-exhaustion DoS during body read).

### 2.4. Safe Binding and Listen Address
To align with safe testing guidelines:
*   Introduce `HOST` environment variable (defaults to `0.0.0.0` to maintain compatibility with Docker/Kubernetes, but allows overriding to `127.0.0.1` for local secure testing).

### 2.5. Restricting HTTP Methods
*   Configure the router to only allow the required methods for the SLO demo: `GET`, `POST`, `PUT`, `DELETE`, `PATCH`, `OPTIONS`, `HEAD`.
*   Explicitly block methods like `TRACE` and `CONNECT` by returning `405 Method Not Allowed`.

## 3. Verification Plan
A comprehensive set of tests will be performed to verify all changes:
1.  **Timeout Verification**: Unit tests checking that the `http.Server` config is correctly populated.
2.  **Security Headers Verification**: Making local test requests and asserting the presence and correctness of new security headers.
3.  **Log Masking Verification**: Sending a request with sensitive headers (e.g., `Authorization`) and asserting that the logs redact/mask them.
4.  **HTTP Method Restrictions Verification**: Sending `TRACE` or other unsupported methods and asserting a `405 Method Not Allowed` response.
5.  **Address Binding Verification**: Testing server startup bound to a specific IP address.

## 4. Cloud Run Deployment Fix

### 4.1. Background & Cause of Failure
When deploying to Cloud Run via `make deploy-cloudrun`, the build step using the local `Dockerfile` fails with the error:
`ERROR: build step 0 "gcr.io/cloud-builders/docker" failed: step exited with non-zero status: 1`

This occurs because:
* The `go.mod` file specifies `go 1.23`.
* The `Dockerfile` uses `golang:1.19` as the builder image base.
* Older Go toolchains (like 1.19) cannot compile modules defined with newer Go versions (like 1.23) and will fail to build.

### 4.2. Resolution
* Update the base image in `Dockerfile` from `golang:1.19` to `golang:1.23` to match the project's Go version.
* Ensure the runner image is also secure and compatible.

## 5. Verification Plan (Deployment Fix)
1. Verify that `Dockerfile` can build successfully with Go 1.23 toolchain.
2. Verify that `make deploy-cloudrun` can be executed.
