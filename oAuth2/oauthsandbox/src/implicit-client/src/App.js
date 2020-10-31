import './App.css';
import React from 'react';

import {
  BrowserRouter as Router,
  Switch,
  Route,
  Link
} from "react-router-dom";

// function App() {
//   return (
//     <div className="App">
//         <h1>Implicit Grant Type</h1>
//         <div><a href="http://192.168.2.10:8080/auth/realms/learningApp/protocol/openid-connect/auth?client_id=implicitClient&response_type=token">Login</a></div>
//         <div><a href="http://localhost:8002/billing/v1/services">Service</a></div>
//         <div><a>Logout</a></div>
//     </div>
//   );
// }

class App extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      access_token: "",
      expires_in: "",
      session_state: "",
      token_type: ""
    }
  }

  setStateValue = (k, v) => {
    if (this.state[k] !== v) {
      this.setState({ [k]: v })
    }
  }

  onCheckStateClick = () => {
    console.log(this.state);
  }

  render() {
    return (
      <Router>
        <div>
          <div className="App">
            <h1>Implicit Grant Type</h1>
          </div>
          <button onClick={this.onCheckStateClick}>Check state</button>
          <nav>
            <ul>
              <li>
                <Link to="/">Home</Link>
              </li>
              <li>
                <Link to="/login">Login</Link>
              </li>
              <li>
                <Link to="/service">Service</Link>
              </li>
              <li>
                <Link to="/logout ">Logout</Link>
              </li>
            </ul>
          </nav>

          {/* A <Switch> looks through its children <Route>s and
              renders the first one that matches the current URL. */}
          <Switch>
            <Route path="/login">
              <Login />
            </Route>
            <Route path="/callback">
              <Callback setStateValue={this.setStateValue} />
            </Route>
            <Route path="/service">
              <Service accessToken={this.state.access_token} />
            </Route>
            <Route path="/">
              <Home />
            </Route>
          </Switch>
        </div>
      </Router>
    );
  }
}

function Home() {
  return <h2>Home</h2>;
}

function Login() {
  window.location = "http://192.168.2.10:8080/auth/realms/learningApp/protocol/openid-connect/auth?client_id=implicitClient&response_type=token&redirect_uri=http://localhost:3000/callback&scope=getBillingService";
  return null;
}

function Callback(props) {
  // get access token
  const hashStr = window.location.hash;
  const hashMap = hashStr.substr(1).split('&').reduce((accum, item) => {
    // add item to accumulator
    const kv = item.split('=');
    accum[kv[0]] = kv[1];
    // return accumulator
    return accum;
  }, {});

  // setState...
  const { setStateValue } = props;
  setStateValue("access_token", hashMap.access_token);
  setStateValue("expires_in", hashMap.expires_in);
  setStateValue("session_state", hashMap.session_state);
  setStateValue("token_type", hashMap.token_type);
  return <h2>Callback</h2>;
}

function Service(props) {
  const { accessToken } = props;
  // access protected resources
  // POST + form
  const formData = new FormData();
  formData.append("access_token", accessToken);
  fetch('http://localhost:8082/billing/v1/services', {
    method: 'POST',
    body: formData
  })
    .then(response => response.json())
    .then(data => {
      console.log(data);
    });

  // parse response 
  const services = [];
  services.push(<div key="a">billingA</div>);
  services.push(<div key="b">billingb</div>);
  return services;
}

export default App;
