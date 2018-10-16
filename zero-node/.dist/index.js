'use strict';var _typeof = typeof Symbol === "function" && typeof Symbol.iterator === "symbol" ? function (obj) {return typeof obj;} : function (obj) {return obj && typeof Symbol === "function" && obj.constructor === Symbol && obj !== Symbol.prototype ? "symbol" : typeof obj;};var _dataloader = require('dataloader');var _dataloader2 = _interopRequireDefault(_dataloader);

var _express = require('express');var _express2 = _interopRequireDefault(_express);
var _nodeFetch = require('node-fetch');var _nodeFetch2 = _interopRequireDefault(_nodeFetch);
var _expressGraphql = require('express-graphql');var _expressGraphql2 = _interopRequireDefault(_expressGraphql);
var _schema = require('./schema');var _schema2 = _interopRequireDefault(_schema);function _interopRequireDefault(obj) {return obj && obj.__esModule ? obj : { default: obj };}

// const BASE_URL = 'http://localhost:8000';
var BASE_URL = 'http://localhost:1323';

function getJSONFromRelativeURL(relativeURL) {
  return (0, _nodeFetch2.default)('' + BASE_URL + relativeURL).
  then(function (res) {return res.json();});
}

function getPeople() {
  return getJSONFromRelativeURL('/people/all').
  then(function (json) {
    return json;
  });
}

function getPerson(id) {
  return getPersonByURL('/people/' + id + '/');
}

function getPersonByURL(relativeURL) {
  var requestURL = relativeURL;
  if (typeof relativeURL === 'number') {
    requestURL = '/people/' + requestURL;
  }
  return getJSONFromRelativeURL(requestURL).
  then(function (json) {
    console.log('===>getPersonByURL', json, requestURL);
    if ((typeof json === 'undefined' ? 'undefined' : _typeof(json)) === 'object' && json.length) {
      return json[0];
    }
    return json;});
}

var app = (0, _express2.default)();

app.use((0, _expressGraphql2.default)(function (req) {
  var cacheMap = new Map();
  var peopleLoader =
  new _dataloader2.default(function (keys) {return Promise.all(keys.map(getPeople));}, { cacheMap: cacheMap });
  var personLoader =
  new _dataloader2.default(function (keys) {return Promise.all(keys.map(getPerson));}, {
    cacheKeyFn: function cacheKeyFn(key) {return '/people/' + key + '/';},
    cacheMap: cacheMap });

  var personByURLLoader =
  new _dataloader2.default(function (keys) {return Promise.all(keys.map(getPersonByURL));}, { cacheMap: cacheMap });
  personLoader.loadAll = peopleLoader.load.bind(peopleLoader, '__all__');
  personLoader.loadByURL = personByURLLoader.load.bind(personByURLLoader);
  personLoader.loadManyByURL =
  personByURLLoader.loadMany.bind(personByURLLoader);
  var loaders = { person: personLoader };
  return {
    context: { loaders: loaders },
    graphiql: true,
    schema: _schema2.default };

}));

app.listen(
5000,
function () {return console.log('GraphQL Server running at http://localhost:5000');});
//# sourceMappingURL=index.js.map