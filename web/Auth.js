import Amplify from '@aws-amplify/core';
import Auth from '@aws-amplify/auth';

import authConfig from './config';

if (process.env.NODE_ENV === 'development') {
  Amplify.configure({
    Auth: authConfig.development,
  });
}
if (process.env.NODE_ENV === 'production') {
  Amplify.configure({
    Auth: authConfig.production,
  });
}

export default Auth;
