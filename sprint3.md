# Sprint 3 Documentation - [Video](#)

## Overview
This document outlines the functionalities implemented during Sprint 3 of the GatoRide application backend.

## New Features

### 1. User Profile Management

#### 1.1 Update Profile Functionality
- **Endpoint**: `POST /user/update-profile`
- **Description**: Implemented a comprehensive profile update system that allows users to modify their personal information
- **Features**:
  - Update user's name
  - Update username (with validation to prevent duplicates)
  - Update location information (latitude, longitude, address)
- **Implementation**: Created in `controllers/update_profile_controller.go`
- **Testing**: Comprehensive tests in `tests/controllers/update_profile_test.go`
  - Tests for successful updates of various fields
  - Tests for duplicate username validation
  - Tests for empty update validation
  - Tests for invalid data format handling

#### 1.2 Profile Controller
- **Endpoints**: 
  - `GET /user/profile` - Retrieve user profile information
  - `GET /user/rides` - Retrieve user's ride history
- **Description**: Allows users to view their profile information and ride history
- **Features**:
  - Get user profile details with sensitive data (password, verification tokens) removed
  - Get list of rides offered as a driver
  - Get list of rides taken as a passenger
- **Implementation**: Created in `controllers/profile_controller.go` 
- **Testing**: Implemented in profile controller tests

### 2. Authentication Improvements

#### 2.1 Logout Functionality
- **Endpoint**: `POST /user/logout`
- **Description**: Securely invalidates user sessions
- **Features**:
  - Extracts JWT token from authorization header
  - Removes the session from database
  - Prevents token reuse after logout
- **Implementation**: Created in `controllers/logout_controller.go`
- **Testing**: Implemented tests that verify:
  - Successful logout operations
  - Handling of invalid tokens
  - Handling of missing tokens

### 3. Ride Management

#### 3.1 Search Ride Functionality
- **Endpoint**: `POST /rides/search-ride`
- **Description**: Allows users to search for available rides based on location, destination, date, and seats required
- **Features**:
  - Search by geographical proximity (origin and destination)
  - Filter by date of travel
  - Filter by minimum available seats
  - Returns only "open" (available) rides
- **Implementation**: Created in `controllers/search_ride_controller.go`
- **Testing**: Comprehensive tests verifying:
  - Successful searches with multiple criteria
  - Empty result handling
  - Invalid input handling

#### 3.2 Book Ride Functionality (Partial Implementation)
- **Endpoint**: `POST /rides/book-ride_id`
- **Description**: Allows users to request to book a ride
- **Features**:
  - Creates a booking alert for the ride driver
  - Validates seat availability
  - Status tracking (pending, confirmed)
- **Implementation**: Created in `controllers/book_ride_controller.go`
- **Note**: This feature is partially implemented and does not have test coverage yet

## Testing

All new functionality (except the partially implemented Book Ride feature) is covered by unit tests. Tests verify:

- Successful operations with valid data
- Validation of input data
- Error handling for edge cases
- Database state verification after operations
- Authentication and authorization checks

## API Documentation - [Link](#)

## Technical Improvements
- Enhanced error handling across all endpoints
- Improved validation for all user inputs
- Optimized database queries for better performance
- Consistent response format across all endpoints
- Added detailed documentation for all new functionality

## Future Enhancements
- Complete the Book Ride functionality with tests
- Add ride cancellation functionality
- Implement ride ratings and reviews
- Add payment integration
- Develop notification system for ride updates