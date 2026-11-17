import { fireEvent, render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router';
import Login from './Login';

describe('Login', () => {
  it('shows validation errors when submitting an empty form', async () => {
    render(
      <MemoryRouter>
        <Login />
      </MemoryRouter>,
    );

    fireEvent.click(screen.getByRole('button', { name: /sign in/i }));

    expect(await screen.findByText('Email is required.')).toBeInTheDocument();
    expect(screen.getByText('Password is required.')).toBeInTheDocument();
  });

  it('shows an email format error', async () => {
    render(
      <MemoryRouter>
        <Login />
      </MemoryRouter>,
    );

    const email = screen.getByLabelText('Email Address');
    fireEvent.change(email, { target: { value: 'not-an-email' } });
    fireEvent.blur(email);

    expect(await screen.findByText('Please enter a valid email address.')).toBeInTheDocument();
  });
});
