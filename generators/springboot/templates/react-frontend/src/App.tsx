import { lazy, Suspense } from 'react';
import { Navigate, Route, Routes } from 'react-router';
import RequireAuth from './guards/RequireAuth';
import { useAuthNavigation } from './hooks/useAuthNavigation';

const Layout = lazy(() => import('./layout/Layout'));
const DashboardLayout = lazy(() => import('./layout/DashboardLayout'));
const Home = lazy(() => import('./pages/Home'));
const Login = lazy(() => import('./pages/Login'));
const Register = lazy(() => import('./pages/Register'));
const DashboardHome = lazy(() => import('./pages/dashboard/DashboardHome'));
const Profile = lazy(() => import('./pages/dashboard/Profile'));

function PageLoader() {
  return (
    <div className="flex justify-center items-center py-20" role="status">
      <div className="animate-spin rounded-full h-12 w-12 border-b-4 border-indigo-500" />
      <span className="sr-only">Loading...</span>
    </div>
  );
}

export default function App() {
  useAuthNavigation();

  return (
    <Suspense fallback={<PageLoader />}>
      <Routes>
        <Route element={<Layout />}>
          <Route index element={<Home />} />
          <Route path="login" element={<Login />} />
          <Route path="register" element={<Register />} />
        </Route>
        <Route
          path="dashboard"
          element={
            <RequireAuth>
              <DashboardLayout />
            </RequireAuth>
          }
        >
          <Route index element={<DashboardHome />} />
          <Route path="profile" element={<Profile />} />
        </Route>
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </Suspense>
  );
}
