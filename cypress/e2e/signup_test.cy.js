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
    cy.get('input[name="location"]').should('exist'); // Address field check
    cy.get('button[type="submit"]').should('exist');
  });

  it('should allow user to input data in all fields including address', () => {
    // Test input fields
    cy.get('input[name="name"]').type('User 2');
    cy.get('input[name="email"]').type('user3@example.com');
    cy.get('input[name="username"]').type('user3');
    cy.get('input[name="password"]').type('password123');

    // Simulate typing address (OpenStreetMap dropdown)
    cy.get('input[name="location"]').type('Gainesville, Alachua County, Florida, United States');
    cy.wait(20000); // Wait for suggestions to appear

    // Select the first suggested address
    cy.get('.dropdown-item').first().click();

    // Verify input values
    cy.get('input[name="name"]').should('have.value', 'User 2');
    cy.get('input[name="email"]').should('have.value', 'user3@example.com');
    cy.get('input[name="username"]').should('have.value', 'user3');
    cy.get('input[name="password"]').should('have.value', 'password123');
    cy.get('input[name="location"]').should('have.value', 'Gainesville, Alachua County, Florida, United States');
  });

  it('should submit the form with valid data including address, latitude, and longitude', () => {
    // Intercept the signup API call
    cy.intercept('POST', '**/signup').as('signupRequest');

    // Fill out the form
    cy.get('input[name="name"]').type('User 2');
    cy.get('input[name="email"]').type('user3@example.com');
    cy.get('input[name="username"]').type('user3');
    cy.get('input[name="password"]').type('password123');

    // Type address and select a suggestion
    cy.get('input[name="location"]').type('Gainesville, Alachua County, Florida, United States');
    cy.wait(2000);
    cy.get('.dropdown-item').first().click();

    // Submit the form
    cy.get('button[type="submit"]').click();

    // Wait for the API call and verify request payload
    cy.wait('@signupRequest').its('request.body').should('deep.include', {
      name: 'User 2',
      email: 'user3@example.com',
      username: 'user3',
      password: 'password123',
      location: {
        address: 'Gainesville, Alachua County, Florida, United States',
        latitude: '29.6516', // Example latitude, replace with actual API response
        longitude: '-82.3248' // Example longitude, replace with actual API response
      }
    });

    // Verify success alert
    cy.on('window:alert', (text) => {
      expect(text).to.contains('Sign up successful');
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
