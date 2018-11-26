// import React from 'react';
import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Router, Route, Switch } from 'react-router';
import Header from './header.jsx';

class App extends Component {
  render() {
    return (
      <div>
        <Header />
      </div>
    );
  }
}

// ReactDOM.render(
// 	<App />,
// 	document.getElementById('appDiv')
// );
// <Router history = {browserHistory}>
// ReactRouter
var browserHistory = {};
ReactDOM.render(
  <Router history = {browserHistory}>
    <Route path = "/" component = {App}>
    </Route>
  </Router>,
	document.getElementById('appDiv')
);