import { useCallback, useRef, useState } from 'react';
import { Link, NavLink, useNavigate } from 'react-router';
import { useAuth } from '../hooks/useAuth';
import { useClickOutside } from '../hooks/useClickOutside';
import { APP_TITLE } from '../config/app';
import { cn } from '../utils/classNames';
import UserAvatar from './UserAvatar';

const navLinkClass = ({ isActive }: { isActive: boolean }) =>
  cn(
    'text-gray-300 no-underline font-medium px-4 py-2 rounded-md transition-all duration-300 hover:text-white hover:bg-white/10 focus-visible:outline-2 focus-visible:outline-blue-400 focus-visible:outline-offset-2',
    isActive && 'bg-blue-400/20 text-white',
  );

export default function Navbar() {
  const { user, isAuthenticated, logout } = useAuth();
  const navigate = useNavigate();
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const navRef = useRef<HTMLElement>(null);

  const closeDropdown = useCallback(() => setIsDropdownOpen(false), []);
  useClickOutside(navRef, closeDropdown, isDropdownOpen);

  const handleLogout = () => {
    logout();
    closeDropdown();
    navigate('/');
  };

  return (
    <nav
      ref={navRef}
      className="bg-gray-900 text-white shadow-md sticky top-0 z-50"
      role="navigation"
      aria-label="Main navigation"
    >
      <div className="max-w-7xl mx-auto px-4 flex flex-col md:flex-row justify-between items-center h-auto md:h-16 py-4 md:py-0 gap-4 md:gap-0">
        <div className="flex items-center">
          <Link
            to="/"
            className="text-2xl font-bold text-white hover:text-blue-400 transition-colors duration-300 focus-visible:outline-2 focus-visible:outline-blue-400 focus-visible:outline-offset-4 rounded no-underline"
            aria-label="Home"
          >
            <span className="bg-linear-to-r from-blue-400 to-purple-400 bg-clip-text text-transparent">
              {APP_TITLE}
            </span>
          </Link>
        </div>

        <ul
          className="flex flex-wrap justify-center items-center gap-4 md:gap-8 list-none m-0 p-0"
          role="menubar"
        >
          <li role="none">
            <NavLink to="/" end className={navLinkClass} role="menuitem">
              Home
            </NavLink>
          </li>

          {isAuthenticated ? (
            <>
              <li role="none">
                <NavLink to="/dashboard" className={navLinkClass} role="menuitem">
                  Dashboard
                </NavLink>
              </li>
              <li role="none" className="relative">
                <button
                  type="button"
                  onClick={() => setIsDropdownOpen((open) => !open)}
                  className="flex items-center gap-2 text-gray-300 font-medium px-2 py-2 rounded-md transition-all duration-300 hover:text-white hover:bg-white/10 focus-visible:outline-2 focus-visible:outline-blue-400 focus-visible:outline-offset-2"
                  aria-expanded={isDropdownOpen}
                  aria-haspopup="true"
                  aria-label="User menu"
                >
                  <UserAvatar fullName={user?.name} size="sm" ariaLabel="User menu avatar" />
                  <span className="text-sm">{user?.name}</span>
                  <i className="fas fa-chevron-down text-xs" aria-hidden="true"></i>
                </button>

                {isDropdownOpen && (
                  <div className="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg py-2 z-50">
                    <div className="px-4 py-2 border-b border-gray-200">
                      <p className="text-sm font-semibold text-gray-900 m-0">{user?.name}</p>
                      <p className="text-xs text-gray-500 m-0">{user?.email}</p>
                    </div>
                    <Link
                      to="/dashboard/profile"
                      onClick={closeDropdown}
                      className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 no-underline transition-colors duration-200"
                    >
                      <i className="fas fa-user mr-2" aria-hidden="true"></i>
                      My Profile
                    </Link>
                    <button
                      type="button"
                      onClick={handleLogout}
                      className="w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-gray-100 transition-colors duration-200 border-0 bg-transparent cursor-pointer"
                    >
                      <i className="fas fa-sign-out-alt mr-2" aria-hidden="true"></i>
                      Logout
                    </button>
                  </div>
                )}
              </li>
            </>
          ) : (
            <>
              <li role="none">
                <NavLink to="/login" className={navLinkClass} role="menuitem">
                  Login
                </NavLink>
              </li>
              <li role="none">
                <NavLink to="/register" className={navLinkClass} role="menuitem">
                  Register
                </NavLink>
              </li>
            </>
          )}
        </ul>
      </div>
    </nav>
  );
}
