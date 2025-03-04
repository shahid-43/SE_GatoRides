# Frontend - Dashboard Implementation 

## Dashboard Features
The dashboard was implemented with the following key features:  

### 1. User Authentication Display
- Shows logged in user details from [`AuthContext`](src/context/AuthContext.js)
- Displays email and username dynamically
- Handles loading states appropriately

### 2. Ride Map Integration
- Implemented interactive map using [`RideMap`](src/components/RideMap.js)
- Location search functionality with OpenStreetMap API
- Visual route plotting between pickup and destination

### 3. Navigation
- Integrated with [`NavBar`](src/components/NavBar.js) component
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
