describe('RideMap Component', () => {
  beforeEach(() => {
    cy.viewport(2160, 1400); // Set the viewport size for consistent testing
    // Visit the login page
    cy.visit('http://localhost:3000/login');

    // Perform login
    cy.get('input[name="email"]').type('r60y739794@tidissajiiu.com'); // Type email
    cy.get('input[name="password"]').type('asdf'); // Type password
    cy.get('button[type="submit"]').click(); // Click login button

    // Wait for the login to complete and navigate to the home page
    cy.url().should('include', '/'); // Ensure the user is redirected to the home page
  });

  it('should display suggestions when typing in the "From" input', () => {
    // Type into the "From" input
    cy.get('#from').type('New York');

    // Mock the API response for location suggestions
    cy.intercept('GET', '**/search?format=json&q=New%20York', {
      statusCode: 200,
      body: [
        {
          lat: '40.712776',
          lon: '-74.005974',
          display_name: 'New York',
        },
        {
          lat: '40.73061',
          lon: '-73.935242',
          display_name: 'Brooklyn',
        },
      ],
    }).as('getFromSuggestions');

    // Wait for the API response and dropdown to appear
    cy.wait('@getFromSuggestions');
    cy.get('.dropdown-menu').should('be.visible');
    cy.get('.dropdown-item').should('have.length', 2);

    // Verify the suggestions
    cy.get('.dropdown-item').first().should('contain.text', 'New York, NY, USA');
    cy.get('.dropdown-item').last().should('contain.text', 'Brooklyn, NY, USA');
  });

  it('should allow selecting a suggestion for "From" location', () => {
    // Type into the "From" input
    cy.get('#from').type('New York');

    // Mock the API response for location suggestions
    cy.intercept('GET', '**/search?format=json&q=New%20York', {
      statusCode: 200,
      body: [
        {
          lat: '40.712776',
          lon: '-74.005974',
          display_name: 'New York, NY, USA',
        },
      ],
    }).as('getFromSuggestions');

    // Wait for the API response and dropdown to appear
    cy.wait('@getFromSuggestions');
    cy.get('.dropdown-item').first().click();

    // Verify the input value is updated
    cy.get('#from').should('have.value', 'New York, NY, USA');
  });

  it('should display suggestions when typing in the "To" input', () => {
    // Type into the "To" input
    cy.get('#to').type('Los Angeles');

    // Mock the API response for location suggestions
    cy.intercept('GET', '**/search?format=json&q=Los%20Angeles', {
      statusCode: 200,
      body: [
        {
          lat: '34.052235',
          lon: '-118.243683',
          display_name: 'Los Angeles, CA, USA',
        },
      ],
    }).as('getToSuggestions');

    // Wait for the API response and dropdown to appear
    cy.wait('@getToSuggestions');
    cy.get('.dropdown-menu').should('be.visible');
    cy.get('.dropdown-item').should('have.length', 1);

    // Verify the suggestion
    cy.get('.dropdown-item').first().should('contain.text', 'Los Angeles, CA, USA');
  });

  it('should allow selecting a suggestion for "To" location', () => {
    // Type into the "To" input
    cy.get('#to').type('Los Angeles');

    // Mock the API response for location suggestions
    cy.intercept('GET', '**/search?format=json&q=Los%20Angeles', {
      statusCode: 200,
      body: [
        {
          lat: '34.052235',
          lon: '-118.243683',
          display_name: 'Los Angeles, CA, USA',
        },
      ],
    }).as('getToSuggestions');

    // Wait for the API response and dropdown to appear
    cy.wait('@getToSuggestions');
    cy.get('.dropdown-item').first().click();

    // Verify the input value is updated
    cy.get('#to').should('have.value', 'Los Angeles, CA, USA');
  });

  it('should display an error if "From" or "To" is not selected', () => {
    // Submit the form without selecting locations
    cy.get('form').submit();

    // Verify the error message
    cy.get('.error-message').should('be.visible').and('contain.text', 'Please select valid locations.');
  });

  it('should submit the form with valid "From" and "To" locations', () => {
    // Type and select "From" location
    cy.get('#from').type('New York');
    cy.intercept('GET', '**/search?format=json&q=New%20York', {
      statusCode: 200,
      body: [
        {
          lat: '40.712776',
          lon: '-74.005974',
          display_name: 'New York, NY, USA',
        },
      ],
    });
    cy.get('.dropdown-item').first().click();

    // Type and select "To" location
    cy.get('#to').type('Los Angeles');
    cy.intercept('GET', '**/search?format=json&q=Los%20Angeles', {
      statusCode: 200,
      body: [
        {
          lat: '34.052235',
          lon: '-118.243683',
          display_name: 'Los Angeles, CA, USA',
        },
      ],
    });
    cy.get('.dropdown-item').first().click();

    // Submit the form
    cy.get('form').submit();

    // Verify the console logs (mocked for testing)
    cy.window().then((win) => {
      cy.stub(win.console, 'log').as('consoleLog');
    });
    cy.get('@consoleLog').should('be.calledWith', 'From:', {
      lat: 40.712776,
      lon: -74.005974,
      display_name: 'New York, NY, USA',
    });
    cy.get('@consoleLog').should('be.calledWith', 'To:', {
      lat: 34.052235,
      lon: -118.243683,
      display_name: 'Los Angeles, CA, USA',
    });
  });
});