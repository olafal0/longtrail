import Auth from './Auth';

export default {
  get: async (url) => {
    const session = await Auth.currentSession();
    const response = await fetch(url, {
      method: 'GET',
      mode: 'cors',
      headers: {
        Authorization: session.getIdToken().getJwtToken(),
      },
    });
    return response.json();
  },

  post: async (url, data) => {
    const session = await Auth.currentSession();
    console.log(session);
    const response = await fetch(url, {
      method: 'POST',
      mode: 'cors',
      headers: {
        Authorization: session.getIdToken().getJwtToken(),
      },
      body: JSON.stringify(data),
    });
    return response.json();
  },
};
