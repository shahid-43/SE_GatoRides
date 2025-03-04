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