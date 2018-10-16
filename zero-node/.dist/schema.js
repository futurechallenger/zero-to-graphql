'use strict';Object.defineProperty(exports, "__esModule", { value: true });var _nodeFetch = require('node-fetch');var _nodeFetch2 = _interopRequireDefault(_nodeFetch);
var _graphql = require('graphql');







var _graphqlRelay = require('graphql-relay');function _interopRequireDefault(obj) {return obj && obj.__esModule ? obj : { default: obj };}var _nodeDefinitions =








(0, _graphqlRelay.nodeDefinitions)(
// A method that maps from a global id to an object
function (globalId, _ref) {var loaders = _ref.loaders;var _fromGlobalId =
  (0, _graphqlRelay.fromGlobalId)(globalId),id = _fromGlobalId.id,type = _fromGlobalId.type;
  if (type === 'Person') {
    return loaders.person.load(id);
  }
},
// A method that maps from an object to a type
function (obj) {
  if (obj.hasOwnProperty('username')) {
    return PersonType;
  }
}),nodeField = _nodeDefinitions.nodeField,nodeInterface = _nodeDefinitions.nodeInterface;


var PersonType = new _graphql.GraphQLObjectType({
  name: 'Person',
  description: 'Somebody that you used to know',
  fields: function fields() {return {
      id: (0, _graphqlRelay.globalIdField)('Person'),
      firstName: {
        type: _graphql.GraphQLString,
        description: 'What you yell at me',
        resolve: function resolve(obj) {return obj.first_name;} },

      lastName: {
        type: _graphql.GraphQLString,
        description: 'What you yell at me when I\'ve been bad',
        resolve: function resolve(obj) {return obj.last_name;} },

      fullName: {
        type: _graphql.GraphQLString,
        description: 'A name sandwich',
        resolve: function resolve(obj) {
          console.log(obj.first_name + ' ' + obj.last_name);
          return obj.first_name + ' ' + obj.last_name;
        } },

      email: {
        type: _graphql.GraphQLString,
        description: 'Where to send junk mail' },

      username: {
        type: _graphql.GraphQLString,
        description: 'Log in as this' },

      friends: {
        type: new _graphql.GraphQLList(PersonType),
        description: 'People who lent you money',
        resolve: function resolve(obj, args, _ref2) {var loaders = _ref2.loaders;return (
            loaders.person.loadManyByURL(obj.friends));} } };},


  interfaces: [nodeInterface] });


var QueryType = new _graphql.GraphQLObjectType({
  name: 'Query',
  description: 'The root of all... queries',
  fields: function fields() {return {
      allPeople: {
        type: new _graphql.GraphQLList(PersonType),
        description: 'Everyone, everywhere',
        resolve: function resolve(root, args, _ref3) {var loaders = _ref3.loaders;return loaders.person.loadAll();} },

      node: nodeField,
      person: {
        type: PersonType,
        args: {
          id: { type: new _graphql.GraphQLNonNull(_graphql.GraphQLID) } },

        resolve: function resolve(root, args, _ref4) {var loaders = _ref4.loaders;return loaders.person.load(args.id);} } };} });exports.default =




new _graphql.GraphQLSchema({
  query: QueryType });
//# sourceMappingURL=schema.js.map