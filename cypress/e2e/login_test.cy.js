describe('Login Page Test', () => {
    it('should allow a user to log in', () => {
      cy.visit('http://localhost:3000/login'); // Visit the login page
      cy.get('input[name="email"]').type('testuser@example.com'); // Type email
      cy.get('input[name="password"]').type('password123'); // Type password
      cy.get('button[type="submit"]').click(); // Click login button
    });
  });
  