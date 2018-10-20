import fetch from 'node-fetch';

// const BASE_URL = 'http://localhost:8000';
const BASE_URL = 'http://localhost:1323';

function getJSONFromRelativeURL(relativeURL) {
  return fetch(`${BASE_URL}${relativeURL}`)
    .then(res => res.json());
}

export function getPeople() {
  return getJSONFromRelativeURL('/people/all')
    .then(json => {
      return json
    });
}

export function getPerson(id) {
  return getPersonByURL(`/people/${id}`);
}

export function getPersonByURL(relativeURL) {
  let requestURL = relativeURL;
  if(typeof relativeURL === 'number') {
    requestURL = '/people/' + requestURL;
  }
  return getJSONFromRelativeURL(requestURL)
    .then(json => {
      console.log('===>getPersonByURL', json, requestURL);
      if(typeof json === 'object' && json.length) {
        return json[0];
      }
      return json});
}