# Frontend - Dashboard Implementation 

## Dashboard Features
The dashboard was implemented with the following key features:  

### 1. User Authentication Display
- Shows logged in user details from [`AuthContext`](GatoRide/frontend/src/context/AuthContext.js)
- Displays email and username dynamically
- Handles loading states appropriately

### 2. Ride Map Integration
- Implemented interactive map using [`RideMap`](GatoRide/frontend/src/components/RideMap.js)
- Location search functionality with OpenStreetMap API
- Visual route plotting between pickup and destination

### 3. Navigation
- Integrated with [`NavBar`](GatoRide/frontend/src/components/NavBar.js) component
- Protected route handling
- Logout functionality

## Testing Implementation

### Cypress E2E Tests
1. **Login Tests** ([login_test.cy.js](cypress/e2e/login_test.cy.js))
```javascript
// filepath: cypress/e2e/login.cy.js
describe('Login Page Test', () => {
    it('should allow a user to log in', () => {
      cy.visit('http://localhost:3000/login'); // Visit the login page
      cy.get('input[name="email"]').type('testuser@example.com'); // Type email
      cy.get('input[name="password"]').type('password123'); // Type password
      cy.get('button[type="submit"]').click(); // Click login button
    });
  });
  
```
2. **SignUp Tests** ([signup_test.cy.js](cypress/e2e/signup_test.cy.js))
```javascript
// filepath: cypress/e2e/signup.cy.js
describe('SignUp Page Tests', () => {
    beforeEach(() => {
      // Visit the signup page before each test
      cy.visit('http://localhost:3000/signup');
    });
  
    it('should display signup form with all fields', () => {
      // Check if all form elements are present
      cy.get('form').should('exist');
      cy.get('h2').should('contain', 'Join GatoRides');
      cy.get('input[name="name"]').should('exist');
      cy.get('input[name="email"]').should('exist');
      cy.get('input[name="username"]').should('exist');
      cy.get('input[name="password"]').should('exist');
      cy.get('button[type="submit"]').should('exist');
    });
  
    it('should allow user to input data in all fields', () => {
      // Test input fields
      cy.get('input[name="name"]').type('Test User');
      cy.get('input[name="email"]').type('testuser@example.com');
      cy.get('input[name="username"]').type('testuser123');
      cy.get('input[name="password"]').type('password123');
  
      // Verify the input values
      cy.get('input[name="name"]').should('have.value', 'Test User');
      cy.get('input[name="email"]').should('have.value', 'testuser@example.com');
      cy.get('input[name="username"]').should('have.value', 'testuser123');
      cy.get('input[name="password"]').should('have.value', 'password123');
    });
  
    it('should submit the form with valid data', () => {
      // Intercept the signup API call
      cy.intercept('POST', '**/signup').as('signupRequest');
  
      // Fill out the form
      cy.get('input[name="name"]').type('Test User');
      cy.get('input[name="email"]').type('testuser@example.com');
      cy.get('input[name="username"]').type('testuser123');
      cy.get('input[name="password"]').type('password123');
  
      // Submit the form
      cy.get('button[type="submit"]').click();
  
      // Wait for the API call and verify alert
      cy.wait('@signupRequest').then(() => {
        cy.on('window:alert', (text) => {
          expect(text).to.contains('Sign up successful');
        });
      });
    });
  
    it('should show validation errors for empty fields', () => {
      // Click submit without filling any fields
      cy.get('button[type="submit"]').click();
  
      // Check if HTML5 validation is working
      cy.get('input[name="name"]').then($input => {
        expect($input[0].validationMessage).to.not.be.empty;
      });
    });
  });
```

### Unit Tests
1. **Login Component**(Login.test.js)
* Form rendering
* Input handling
* Authentication flow
* Error states
2. **Signup Component**(SignUp.test.js)
* Form validation
* Field updates
* Submission handling
* Error messaging
3. **RideMap Component**(RideMap.test.js)
* Map rendering
* Location search
* Route plotting
* Error handling

## Common Issues and Solutions

**Issue**: React Router Integration

**Solution**:
```javascript
ReactDOM.render(
  <BrowserRouter>
    <AuthProvider>
      <App />
    </AuthProvider>
  </BrowserRouter>,
  document.getElementById('root')
)
```

**Issue**: Map Marker icons not loading

**Solution**: Added proper icon imports and fixed paths in RideMap component
```javascript
import L from 'leaflet'
L.Icon.Default.mergeOptions({
    iconRetinaUrl: require('leaflet/dist/images/marker-icon-2x.png'),
    iconUrl: require('leaflet/dist/images/marker-icon.png'),
    shadowUrl: require('leaflet/dist/images/marker-shadow.png')
})
```  

**Issue**: Authentication State Persistence  

**Solution**: Implemented localStorage persistence in AuthContext.js
```javascript
useEffect(() => {
  const storedUser = localStorage.getItem('user')
  if (storedUser) {
    setUser(JSON.parse(storedUser))
  }
}, [])
```  
**Issue**: Protected Route Access  

**Solution**: Added route protection logic in routes.js  
```javascript
{
  path: '/dashboard',
  component: Dashboard,
  protected: true
}
```  

## Running Tests

```bash
# Run unit tests
npm test

# Run Cypress tests
npm run cypress:open

# Run test coverage
npm run test:coverage
```

## Future Enhancements

1. Real-time ride tracking
2. Payment integration
3. Chat functionality
4. Ride history
5. Advanced route optimization

### **Frontend Demo Video** - [Video](https://drive.google.com/file/d/1Qh8UhJgu-GSZ5W-R1sUsQejCtwFYOlLB/view?usp=drive_link)


# **Backend - Sprint 2  [video](https://drive.google.com/drive/folders/1gV2L_kqw48QFacPMA-2njNtdROZgyjAr?usp=drive_link)**

## **User Story 5: Update User Location**
**Summary:** Users should be able to update their last known location.
**Acceptance Criteria:**  
- User should send a `POST /users/location` request with latitude, longitude, and address.
- The system should validate and update the location in the database.
- Only authenticated users should be allowed to update their location.

### **Tasks:**
- [x] Implement `POST /users/location` endpoint.
- [x] Validate and extract location details from the request body.
- [x] Update user location in MongoDB.
- [x] Ensure the request is authenticated using JWT middleware.
- [x] Return appropriate success/error responses.

---

## **User Story 6: Provide a Ride**
**Summary:** Users should be able to provide a ride with pickup and dropoff locations.
**Acceptance Criteria:**  
- Users should send a `POST /user/provide-ride` request.
- System should validate the request and store ride details.
- Duplicate rides from the same driver should not be allowed.

### **Tasks:**
- [x] Implement `POST /user/provide-ride` endpoint.
- [x] Validate and extract ride details from request.
- [x] Ensure a driver cannot create a duplicate ride.
- [x] Store ride details in the database.
- [x] Return ride ID upon success.

---

## **User Story 7: Fetch Nearby Rides**
**Summary:** Users should see a list of available rides near their location.
**Acceptance Criteria:**  
- Users should send a `GET /home` request.
- System should fetch available rides near the user’s location.
- System should return ride details in JSON format.

### **Tasks:**
- [x] Implement `GET /home` endpoint.
- [x] Fetch rides from MongoDB using geospatial queries.
- [x] Return only rides that are currently open.
- [x] Ensure the request is authenticated.
- [x] Return appropriate success/error responses.

---

## **User Story 8: Unit Testing for Backend API**
**Summary:** Ensure API functionality with comprehensive unit tests.
**Acceptance Criteria:**  
- Each API function should have a corresponding unit test.
- Tests should run against a real database with proper cleanup.
- Middleware authentication should be handled in test cases.

### **Tasks:**
- [x] Write unit tests for `POST /users/location`.
- [x] Write unit tests for `POST /user/provide-ride`.
- [x] Write unit tests for `GET /home`.
- [x] Implement helper functions for database cleanup.
- [x] Ensure tests run successfully with database integration.

---

## **Backend API Documentation**

### **Implemented APIs:**
1. **`POST /signup`** – User registration with email verification.
2. **`GET /verify-email`** – Email verification via JWT token.
3. **`POST /login`** – User authentication and JWT issuance.
4. **`POST /users/location`** – Update user location.
5. **`POST /user/provide-ride`** – Allow users to offer rides.
6. **`GET /home`** – Fetch available nearby rides.

### **Unit Test Coverage:**
- A 1:1 test-to-function ratio was maintained.
- Test cases cover both **valid** and **invalid** scenarios.
- Tests include authentication checks.
- Database state is cleaned up after each test.

---

## **Sprint 2 Goal: Successfully Implemented**
✅ Implemented user location updates, ride creation, and ride fetching.  
✅ Developed and executed unit tests for backend API functions.  
✅ Ensured smooth JWT authentication in all API calls.  

---

## **Planned but Not Implemented**
- **Admin feature for listing rides per user** (Deferred to Sprint 3).  
- **Frontend and backend integration testing** (Postponed until frontend is fully developed).  

