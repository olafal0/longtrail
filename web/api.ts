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

async function createPractice(event) {
  const response = await requests.post(`${config.apiUrl}/practices/new`, event);
  return response;
}

async function getPractice(id: string) {
  return await requests.get(`${config.apiUrl}/practice/${id}`);
}

async function getPractices(start, end) {
  return await requests.get(`${config.apiUrl}/practices`, { start, end });
}

async function setPractice(event) {
  const response = await requests.post(
    `${config.apiUrl}/practice/${event.id}`,
    event
  );
  return response;
}

async function deletePractice(id: string) {
  return await requests.delete(`${config.apiUrl}/practice/${id}`);
}

export default {
  echo,
  createPractice,
  getPractice,
  getPractices,
  setPractice,
  deletePractice,
};
