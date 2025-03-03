import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { BrowserRouter } from 'react-router-dom';
import Login from '../Login';
import AuthContext from '../../context/AuthContext';
import '@testing-library/jest-dom';

const mockNavigate = jest.fn();

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useNavigate: () => mockNavigate
}));

// Mock AuthService
jest.mock('../../services/AuthService', () => ({
  login: jest.fn(),
  signup: jest.fn(),
  verifyEmail: jest.fn(),
}));

const mockHandleLogin = jest.fn();
const mockContextValue = {
  handleLogin: mockHandleLogin,
  user: null
};

const renderWithContext = (component) => {
  return render(
    <BrowserRouter>
      <AuthContext.Provider value={mockContextValue}>
        {component}
      </AuthContext.Provider>
    </BrowserRouter>
  );
};



describe('Login Component', () => {
  beforeEach(() => {
    mockHandleLogin.mockClear();
  });

  it('renders login form with all elements', () => {
    renderWithContext(<Login />);
    
    expect(screen.getByText('Login to GatoRides')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('Email')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('Password')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /login/i })).toBeInTheDocument();
  });

  it('updates form data when user types', async () => {
    renderWithContext(<Login />);
    
    const emailInput = screen.getByPlaceholderText('Email');
    const passwordInput = screen.getByPlaceholderText('Password');

    await userEvent.type(emailInput, 'test@example.com');
    await userEvent.type(passwordInput, 'password123');

    expect(emailInput).toHaveValue('test@example.com');
    expect(passwordInput).toHaveValue('password123');
  });

  it('calls handleLogin with correct data on form submission', async () => {
    renderWithContext(<Login />);
    
    const emailInput = screen.getByPlaceholderText('Email');
    const passwordInput = screen.getByPlaceholderText('Password');
    const submitButton = screen.getByRole('button', { name: /login/i });

    await userEvent.type(emailInput, 'test@example.com');
    await userEvent.type(passwordInput, 'password123');
    await userEvent.click(submitButton);

    expect(mockHandleLogin).toHaveBeenCalledWith('test@example.com', 'password123');
  });

  it('shows success alert on successful login', async () => {
    const mockAlert = jest.spyOn(window, 'alert').mockImplementation(() => {});
    mockHandleLogin.mockResolvedValueOnce();
    
    renderWithContext(<Login />);
    
    const emailInput = screen.getByPlaceholderText('Email');
    const passwordInput = screen.getByPlaceholderText('Password');
    const submitButton = screen.getByRole('button', { name: /login/i });

    await userEvent.type(emailInput, 'test@example.com');
    await userEvent.type(passwordInput, 'password123');
    await userEvent.click(submitButton);

    await waitFor(() => {
      expect(mockAlert).toHaveBeenCalledWith('Login successful');
    });

    mockAlert.mockRestore();
  });

  it('shows error alert on login failure', async () => {
    const mockAlert = jest.spyOn(window, 'alert').mockImplementation(() => {});
    mockHandleLogin.mockRejectedValueOnce(new Error('Login failed'));
    
    renderWithContext(<Login />);
    
    const submitButton = screen.getByRole('button', { name: /login/i });
    await userEvent.click(submitButton);

    await waitFor(() => {
      expect(mockAlert).toHaveBeenCalledWith('Error logging in');
    });

    mockAlert.mockRestore();
  });
});