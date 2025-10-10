# 📌 Project Task Document — Authentication & Notification System (Go + MongoDB)

## 🎯 Goal
Build a **high-performance, reliable authentication system** with **Go** using **MongoDB**, **Redis**, and **NATS**. The system should support multiple signup methods, secure login/logout, token-based authentication, notification delivery, and be thoroughly tested.

---

## 🔧 Tech Stack
- **Language:** Go (1.21+)
- **Routing:** Fiber (we use fiber in soren but gin is good too)
- **Logging:** Zap
- **Database:** MongoDB (with indexes)
- **Cache / Rate-limiting / Blacklist:** Redis
- **Messaging:** NATS (with ack/nack to ensure message delivery)
- **Auth:** JWT with access + refresh tokens
- **Testing:** Go test, mockery for mocks, Postman for APIs
- **Version Control:** GitHub (branching, PR workflow, CI)

---

## 📂 Tasks

### 1. **Authentication Service**
- Implement **Signup, Login, Logout**.
- Use **access tokens** (short-lived) and **refresh tokens** (long-lived).
- Store JWTs in **secure HttpOnly cookies**.
- Passwords must be stored **only as salted hashes** (never plaintext).
- Use **MongoDB indexes** for performance.
- JWT should be signed with **private/public key** pair.
- Must support **two signup methods**:
  1. Email + 4-digit verification code (rate limit: 1 request/minute).
  2. Google OAuth2 login.
- Design signup flow in a way that allows adding **new methods easily**.
- Add **Brute-force protection**:
  - Track failed login attempts in Redis.
  - Lock account or apply exponential backoff after too many failures.

---

### 2. **MongoDB Library**
- Create a **separate library** for MongoDB connection and collection access.
- Handle indexes on startup (e.g., email unique index, TTL indexes).
- Expose DAO interfaces to be used in business logic.

---

### 3. **Redis Library**
- Create a **library for Redis** with reusable functions.
- Redis will be used for multiple purposes (implement these):
  - **Access token denylist** until expiry.
  - **Refresh token store** with TTL.
  - **Rate limiting middleware** (example: limit signup requests).
- Explore **different Redis patterns** (string keys, TTL, counters).

---

### 4. **Notification Service + Messaging**
- When a user signs up via email:
  - Auth service publishes a **NATS message** with `{email, code}`.
  - Notification service listens and sends the email.
- Requirements:
  - If **multiple notification services** are running, **only one** must handle the message.
  - Messages must be **acknowledged**; if not, they should be retried.
  - Develop a **library** to handle publishing and subscribing with ack/nack.

---

### 5. **Security & Middleware**
- **CORS Protection**:
  - Configure allowed origins, methods, and headers in Fiber.
- **CSRF Protection**:
  - Add CSRF middleware for state-changing requests.
- **Brute-force Protection**:
  - Lock accounts or IPs after repeated failed attempts.
- **Error Handling Standardization**:
  - Create a centralized error handler middleware.
  - All API responses should follow a consistent JSON format with `error_code`, `message`, and optional `details`.

---

### 6. **Testing**
- **Unit Tests**:
  - Must use `mockery` to mock DAOs (no DB dependency).
  - Cover meaningful cases (not just "happy path").
- **Integration Tests**:
  - Run against MongoDB, Redis, and NATS (use Docker).
- **Load Testing**:
  - Test system **with and without Redis** to compare performance.
- **Scenarios**:
  - Successful signup/login/logout.
  - Expired access token refresh.
  - Rate-limit exceeded.
  - Invalid/expired codes.
  - Concurrent notification service consumers.

---

### 7. **GitHub Workflow**
- **Branching:**  
  - No direct pushes to `main`.  
  - Use feature branches + PRs.  
- **CI/CD:**  
  - Pipeline should run `go test` on every push.  
- **Environment:**  
  - Use `.env.example` for tests (not committed to repo).  

---

## 📑 Deliverables
1. **Auth service** with:
   - Signup (Email + Google OAuth2)  
   - Login, Logout  
   - JWT auth in cookies  
2. **MongoDB library**  
3. **Redis library**  
4. **Notification service** with NATS (ack/nack handling)  
5. **Rate-limiting middleware**  
6. **Security middlewares** (Brute-force, CORS, CSRF)  
7. **Standardized error handling**  
8. **Unit + Integration + Load Tests**  
9. **Postman collection** for APIs  
10. **GitHub repo** with proper workflow  

---

## ✅ Acceptance Criteria
- Secure (no raw passwords stored, JWT signed with private key).
- Reliable messaging (only one notification service processes each event).
- Extensible signup methods.
- Brute-force protection implemented.
- Proper CORS and CSRF handling.
- Unified error response format.
- Tested (unit, integration, load).
- CI pipeline running automatically.
- Code organized into separate libraries for Mongo, Redis, NATS, and Auth.
