import React from "react";
import {
    Link,
} from 'react-router-dom';

class NavBar extends React.Component {
    render() {
      return <div>
        <ul>
            <li>
                <Link to="/">Home</Link>
            </li>
            <li>
                <Link to="/test">Test</Link>
            </li>
        </ul>
      </div>
    }
  }

  export default NavBar;