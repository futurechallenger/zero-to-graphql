import DataLoader from 'dataloader';

import express from 'express';
import graphqlHTTP from 'express-graphql';
import schema from './schema';
import {getPeople, getPerson, getPersonByURL} from './utils';

const app = express();

app.use(graphqlHTTP(req => {
  const cacheMap = new Map();
  const peopleLoader =
    new DataLoader(keys => Promise.all(keys.map(getPeople)), {cacheMap});
  const personLoader =
    new DataLoader(keys => Promise.all(keys.map(getPerson)), {
      cacheKeyFn: key => `/people/${key}`,
      cacheMap,
    });
  const personByURLLoader =
    new DataLoader(keys => Promise.all(keys.map(getPersonByURL)), {cacheMap});
  personLoader.loadAll = peopleLoader.load.bind(peopleLoader, '__all__');
  personLoader.loadByURL = personByURLLoader.load.bind(personByURLLoader);
  personLoader.loadManyByURL =
    personByURLLoader.loadMany.bind(personByURLLoader);
  const loaders = {person: personLoader};
  return {
    context: {loaders},
    graphiql: true,
    schema,
  };
}));

app.listen(
  5000,
  () => console.log('GraphQL Server running at http://localhost:5000')
);
