import Home from '../pages/Home';
import SignUp from '../components/SignUp';
import Login from '../components/Login';
import Dashboard from '../components/Dashboard';
import VerifyEmail from '../components/VerifyEmail';
import RideRequest from '../components/ProvideRide';

const RouteConfig = [
  {
    path: '/',
    component: Home,
    exact: true,
  },
  {
    path: '/signup',
    component: SignUp,
  },
  {
    path: '/login',
    component: Login,
  },
  {
    path: '/dashboard',
    component: Dashboard,
  },
  {
    path: '/ride-request',
    component: RideRequest,
  },
  {
    path: '/verify-email/:token',
    component: VerifyEmail,
  },
];

export default RouteConfig;
