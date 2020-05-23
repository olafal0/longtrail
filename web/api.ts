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
  console.log(response);
}

async function createPractice({ start, end }) {
  const response = await requests.post(`${config.apiUrl}/practices/new`, {
    startTime: start,
    endTime: end,
  });
  return response;
}
async function getPractice() {}
async function getPractices() {}
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
