import React from "react";
import ReactDOM from "react-dom";
// import { Router, Route, Switch } from 'react-router';
import {
    BrowserRouter as Router,
    Route,
    Switch,
    Link,
} from 'react-router-dom';
import Header from './js/header.jsx';
import NavBar from './js/navbar.jsx';
// import Header from './js/header.js';
import TestTag from './js/test.js';
// https://reacttraining.com/react-router/web/example/basic
// https://medium.freecodecamp.org/part-1-react-app-from-scratch-using-webpack-4-562b1d231e75
// const Index = () => {
//   return <div>Hello React!</div>;
// };

// ReactDOM.render(<Index />, document.getElementById("index"));

class App extends React.Component {
    render() {
      return (
        <div>
          <Header />
          <NavBar />
        </div>
      );
    }
  }
  
//   ReactDOM.render(
//   	<App />,
//   	document.getElementById('appDiv')
//   );
  // <Router history = {browserHistory}>
  // ReactRouter
  // var browserHistory = {};
  ReactDOM.render(
    /*<Router history = {browserHistory}>
      <Route path = "/" component = {App}>
      </Route>
    </Router>,*/
    // <Router history = {browserHistory}>
    <Router>
      <Switch>
          <Route exact path = "/" component = {App}></Route>
          <Route path = "/test" component = {TestTag}></Route>
      </Switch>
    </Router>,
    document.getElementById('appDiv')
  );
  // https://medium.com/@thejasonfile/basic-intro-to-react-router-v4-a08ae1ba5c42

  // https://medium.com/@Preda/getting-started-on-building-a-personal-website-with-react-b44ee93b1710