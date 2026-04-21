import { Navigate, Route, Routes } from 'react-router-dom';

import { useAppSelector } from '../hooks'

import { selectIsAuthenticated } from '../../features/auth/auth.selector'

import Landing from '../../features/static/pages/Landing/Landing';
import Login from '../../features/auth/pages/LogIn/Login-new';
import GoogleCallback from '../../features/auth/pages/GoogleCallback/GoogleCallback';
import SignUp from '../../features/auth/pages/SignUp/SignUp';

const RoleRoutes = (): React.ReactElement | null => {
  const isAuthenticated = useAppSelector(selectIsAuthenticated);
  return (
    !isAuthenticated ? <Routes>
      <Route path="*" element={<Navigate to="/login" />} />
      <Route path="/login" element={<Login />} />
      <Route path="/signup" element={<SignUp />} />
      <Route path="/oauth/google/callback" element={<GoogleCallback />} />
      <Route path="/" element={<Landing />} />
    </Routes> : <Routes>
      <Route path="*" element={<Navigate to="/" />} />
      <Route path="/oauth/google/callback" element={<GoogleCallback />} />
      <Route path="/" element={<Landing />} />
    </Routes>
  )
}

export default RoleRoutes