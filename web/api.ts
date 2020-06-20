import requests from './request';
import config from './config';

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
  createPractice,
  getPractice,
  getPractices,
  setPractice,
  deletePractice,
};
