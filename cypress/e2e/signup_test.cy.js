describe('Signup Form Tests', () => {
  beforeEach(() => {
    cy.visit('http://localhost:3000/signup'); // Adjust URL as needed
  });

  it('should display the signup form with all fields', () => {
    cy.get('form').should('exist');
    cy.get('input[name="name"]').should('exist');
    cy.get('input[name="email"]').should('exist');
    cy.get('input[name="username"]').should('exist');
    cy.get('input[name="password"]').should('exist');
    cy.get('input[name="location"]').should('exist');
    cy.get('button[type="submit"]').should('exist');
  });

  it('should allow the user to enter text into the form fields', () => {
    cy.get('input[name="name"]').type('John Doe').should('have.value', 'John Doe');
    cy.get('input[name="email"]').type('john@example.com').should('have.value', 'john@example.com');
    cy.get('input[name="username"]').type('johndoe').should('have.value', 'johndoe');
    cy.get('input[name="password"]').type('password123').should('have.value', 'password123');
  });

  it('should show location suggestions when typing a location', () => {
    cy.get('input[name="location"]').type('New York, United States');
    cy.wait(1000); // Wait for API response
    cy.get('.dropdown-menu').should('be.visible');
  });

  it('should allow the user to select a location from suggestions', () => {
    cy.get('input[name="location"]').type('New York, United States');
    cy.wait(100000);
    cy.get('.dropdown-item').first().click();
    cy.get('input[name="location"]').should('not.have.value', '');
  });

  it('should submit the form successfully', () => {
    cy.get('input[name="name"]').type('John Doe');
    cy.get('input[name="email"]').type('john@example.com');
    cy.get('input[name="username"]').type('johndoe');
    cy.get('input[name="password"]').type('password123');
    cy.get('input[name="location"]').type('New York, United States');
    cy.wait(1000);
    cy.get('.dropdown-item').first().click();

    cy.intercept('POST', 'http://localhost:5001/signup', {
      statusCode: 200,
      body: { token: 'mock-token' }
    }).as('signupRequest');

    cy.get('button[type="submit"]').click();
    cy.wait('@signupRequest');
    cy.window().its('localStorage.token').should('eq', 'mock-token');
  });
});
