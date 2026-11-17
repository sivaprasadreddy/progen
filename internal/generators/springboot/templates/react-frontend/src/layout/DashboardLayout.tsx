import { useState } from 'react';
import { Link, NavLink, Outlet, useNavigate } from 'react-router';
import { useAuth } from '../hooks/useAuth';
import { APP_TITLE } from '../config/app';
import { cn } from '../utils/classNames';
import UserAvatar from '../components/UserAvatar';

interface NavItem {
  label: string;
  icon: string;
  path: string;
}

const navItems: NavItem[] = [
  { label: 'Dashboard', icon: 'fa-home', path: '/dashboard' },
  { label: 'My Profile', icon: 'fa-user', path: '/dashboard/profile' },
];

export default function DashboardLayout() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const [isUserMenuOpen, setIsUserMenuOpen] = useState(false);

  const closeUserMenu = () => setIsUserMenuOpen(false);

  const handleLogout = () => {
    closeUserMenu();
    logout();
    navigate('/');
  };

  return (
    <div className="flex flex-col h-screen overflow-hidden">
      <header
        className="shrink-0 h-16 bg-white border-b border-gray-200 flex items-center justify-between px-6 z-50"
        role="banner"
      >
        <Link to="/dashboard" className="flex items-center gap-3 no-underline">
          <i className="fas fa-cube text-purple-600 text-xl" aria-hidden="true"></i>
          <span className="text-xl font-semibold text-gray-900">{APP_TITLE}</span>
        </Link>

        <div className="relative">
          <button
            type="button"
            onClick={() => setIsUserMenuOpen((open) => !open)}
            aria-expanded={isUserMenuOpen}
            aria-haspopup="true"
            aria-label="User menu"
            className="flex items-center gap-2 px-3 py-2 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors duration-200 focus-visible:outline-2 focus-visible:outline-purple-600 focus-visible:outline-offset-2"
          >
            <UserAvatar fullName={user?.name} size="sm" ariaLabel="User avatar" />
            <span className="text-sm font-medium hidden sm:block">{user?.name}</span>
            <i
              className={cn(
                'fas fa-chevron-down text-xs text-gray-400 transition-transform duration-200',
                isUserMenuOpen && 'rotate-180',
              )}
              aria-hidden="true"
            ></i>
          </button>

          {isUserMenuOpen && (
            <>
              <div className="fixed inset-0 z-10" onClick={closeUserMenu} aria-hidden="true"></div>
              <div
                role="menu"
                className="absolute right-0 top-full mt-1 w-52 bg-white border border-gray-200 rounded-lg shadow-lg z-20 py-1"
              >
                <div className="px-4 py-2 border-b border-gray-100">
                  <p className="text-sm font-semibold text-gray-900 m-0 truncate">{user?.name}</p>
                  <p className="text-xs text-gray-500 m-0 truncate">{user?.email}</p>
                </div>
                <Link
                  to="/dashboard/profile"
                  onClick={closeUserMenu}
                  role="menuitem"
                  className="flex items-center gap-3 px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 transition-colors duration-200 no-underline focus-visible:outline-2 focus-visible:outline-purple-600 focus-visible:outline-offset-2"
                >
                  <i className="fas fa-user w-4 text-center text-gray-400" aria-hidden="true"></i>
                  <span>My Profile</span>
                </Link>
                <button
                  type="button"
                  onClick={handleLogout}
                  role="menuitem"
                  className="w-full flex items-center gap-3 px-4 py-2 text-sm text-red-600 hover:bg-red-50 transition-colors duration-200 border-0 bg-transparent cursor-pointer focus-visible:outline-2 focus-visible:outline-red-500 focus-visible:outline-offset-2"
                >
                  <i className="fas fa-sign-out-alt w-4 text-center" aria-hidden="true"></i>
                  <span>Logout</span>
                </button>
              </div>
            </>
          )}
        </div>
      </header>

      <div className="flex flex-1 overflow-hidden">
        <aside
          className="w-64 shrink-0 bg-white border-r border-gray-200 overflow-y-auto"
          role="navigation"
          aria-label="Dashboard navigation"
        >
          <nav className="py-4" aria-label="Main menu">
            <ul className="space-y-1 px-3 list-none m-0">
              {navItems.map((item) => (
                <li key={item.path}>
                  <NavLink
                    to={item.path}
                    end={item.path === '/dashboard'}
                    className={({ isActive }) =>
                      cn(
                        'flex items-center gap-3 px-3 py-2 rounded-lg transition-colors duration-200 no-underline focus-visible:outline-2 focus-visible:outline-purple-600 focus-visible:outline-offset-2',
                        isActive
                          ? 'bg-purple-50 text-purple-600'
                          : 'text-gray-700 hover:bg-gray-50',
                      )
                    }
                  >
                    <i className={`fas ${item.icon} w-5 text-center`} aria-hidden="true"></i>
                    <span className="font-medium">{item.label}</span>
                  </NavLink>
                </li>
              ))}
            </ul>
          </nav>
        </aside>

        <main className="flex-1 overflow-y-auto bg-gray-50">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
