import requests from './request';
import config from './config';

// GET/echo
// POST/practices/new
// GET/practice/{id}
// GET/practices
// POST/practice/{id}
// DELETE/practice/{id}

async function echo() {
  const response = await requests.get(`${config.apiUrl}/echo`);
}

async function createPractice({ start, end }) {
  const response = await requests.post(`${config.apiUrl}/practices/new`, {
    start: start,
    end: end,
  });
  return response;
}
async function getPractice(id: string) {
  return await requests.get(`${config.apiUrl}/practice/${id}`);
}
async function getPractices(start, end) {
  return await requests.get(`${config.apiUrl}/practices`, { start, end });
}
async function setPractice() {}
async function deletePractice() {}

export default {
  echo,
  createPractice,
  getPractice,
  getPractices,
  setPractice,
  deletePractice,
};
