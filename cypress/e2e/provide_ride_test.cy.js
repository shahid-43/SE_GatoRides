describe('ProvideRide Component', () => {
    beforeEach(() => {
      cy.visit('/ride-request'); // Update with the correct route
      cy.intercept('GET', 'https://nominatim.openstreetmap.org/search**', { fixture: 'locationSuggestions.json' }).as('getLocationSuggestions');
      cy.intercept('POST', 'http://localhost:5001/user/ride-request', { statusCode: 200, body: { message: 'Ride provided successfully!' } }).as('provideRide');
    });
  
    it('should display the form elements correctly', () => {
      cy.get('h2').should('contain', 'Request a Ride');
      cy.get('input[name="pickup"]').should('exist');
      cy.get('input[name="dropoff"]').should('exist');
      cy.get('input[name="price"]').should('exist');
      cy.get('input[name="date"]').should('exist');
      cy.get('button[type="submit"]').should('contain', 'Submit Ride Request');
    });
  
    it('should allow users to enter pickup and dropoff locations and fetch suggestions', () => {
      cy.get('input[name="pickup"]').type('123 Main St');
      cy.wait('@getLocationSuggestions');
      cy.get('.dropdown-menu').should('be.visible');
      cy.get('.dropdown-item').first().click();
      cy.get('input[name="pickup"]').should('not.have.value', '');
  
      cy.get('input[name="dropoff"]').type('456 Elm St');
      cy.wait('@getLocationSuggestions');
      cy.get('.dropdown-menu').should('be.visible');
      cy.get('.dropdown-item').first().click();
      cy.get('input[name="dropoff"]').should('not.have.value', '');
    });
  
    it('should allow users to enter price and date', () => {
      cy.get('input[name="price"]').type('20');
      cy.get('input[name="price"]').should('have.value', '20');
  
      const today = new Date().toISOString().split('T')[0];
      cy.get('input[name="date"]').type(today);
      cy.get('input[name="date"]').should('have.value', today);
    });
  
    it('should submit the ride request successfully', () => {
      cy.get('input[name="pickup"]').type('123 Main St');
      cy.wait('@getLocationSuggestions');
      cy.get('.dropdown-item').first().click();
  
      cy.get('input[name="dropoff"]').type('456 Elm St');
      cy.wait('@getLocationSuggestions');
      cy.get('.dropdown-item').first().click();
  
      cy.get('input[name="price"]').type('20');
      const today = new Date().toISOString().split('T')[0];
      cy.get('input[name="date"]').type(today);
  
      cy.get('button[type="submit"]').click();
      cy.wait('@provideRide');
      cy.on('window:alert', (txt) => {
        expect(txt).to.contains('Ride provided successfully!');
      });
    });
  });
  