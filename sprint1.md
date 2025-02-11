# **Backend Authentication System - User Stories**


## **User Story 1: User Registration**
**Summary:** Implement user signup with email and password.  
**Acceptance Criteria:**  
- User should be able to send a `POST /signup` request with name, email, username, and password.
- Email and username should be unique.
- Password should be **hashed** before storing.
- User should receive an email verification token after successful registration.

### **Tasks:**
- [x] Create `User` model with `name`, `email`, `username`, `password`, `is_verified`, `verification_token`.
- [x] Implement `POST /signup` endpoint in `auth_controller.go`.
- [x] Hash password using `bcrypt` before storing in the database.
- [x] Ensure unique constraints on `email` and `username`.
- [x] Send verification email with JWT token.

---

## **User Story 2: Email Verification**
**Summary:** Users must verify their email before logging in.  
**Acceptance Criteria:**  
- User receives an email containing a verification link with a JWT token.
- Clicking the link should validate the token and mark the user as verified in the database.
- If the token is expired or invalid, return an appropriate error.

### **Tasks:**
- [x] Generate JWT token for email verification (`GenerateVerificationToken`).
- [x] Implement `GET /verify-email?token=<JWT_TOKEN>` endpoint.
- [x] Update `is_verified` flag to `true` in the database.
- [x] Handle token expiration and invalid cases.
- [x] Log errors for debugging.

---

## **User Story 3: Secure User Login**
**Summary:** Implement user login with JWT authentication.  
**Acceptance Criteria:**  
- User should be able to send a `POST /login` request with email and password.
- System should validate the credentials and return a JWT token.
- If the user is not verified, they should not be allowed to log in.

### **Tasks:**
- [x] Implement `POST /login` endpoint.
- [x] Verify the hashed password using `bcrypt.CompareHashAndPassword`.
- [x] Ensure only verified users can log in (`is_verified: true`).
- [x] Generate a JWT token for authenticated users (`GenerateJWT`).
- [x] Return appropriate error messages if credentials are incorrect or the user is not verified.

---

## **Sprint Goal:**
Deliver a **secure, production-ready authentication system** for GatoRides with **signup, email verification, login, and JWT-based authentication**.


