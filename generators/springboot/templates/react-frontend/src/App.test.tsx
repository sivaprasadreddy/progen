import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router';
import App from './App';
import { authStore } from './services/authStore';

function renderAt(path: string) {
  return render(
    <MemoryRouter initialEntries={[path]}>
      <App />
    </MemoryRouter>,
  );
}

describe('App', () => {
  beforeEach(() => {
    authStore.clear();
  });

  it('renders the home page for guests', async () => {
    renderAt('/');
    expect(await screen.findByRole('heading', { name: /welcome guest/i })).toBeInTheDocument();
    expect(screen.getByRole('menuitem', { name: 'Login' })).toBeInTheDocument();
    expect(screen.getByRole('menuitem', { name: 'Register' })).toBeInTheDocument();
  });

  it('redirects unauthenticated users from the dashboard to login', async () => {
    renderAt('/dashboard');
    expect(await screen.findByRole('heading', { name: /welcome back/i })).toBeInTheDocument();
  });

  it('shows the dashboard for authenticated users', async () => {
    authStore.setAuthData(
      { name: 'John Doe', email: 'john@example.com', role: 'ROLE_USER' },
      'token',
    );
    renderAt('/');
    expect(await screen.findByRole('heading', { name: 'Dashboard' })).toBeInTheDocument();
    expect(screen.getAllByText('John Doe').length).toBeGreaterThan(0);
  });
});
