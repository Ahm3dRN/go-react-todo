import React from "react";
import '../styles/navbar.css';
import { Navbar, Container, Nav} from 'react-bootstrap';
import { NavLink } from 'react-router-dom';

class NavBar extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            token: this.getToken()
        };
    }
    getToken = () => {
      return localStorage.getItem("token")
    }
    render() {
        return (
            <Navbar bg="dark" variant="dark" expand="lg">
            <Container>
              <Navbar.Brand href="#home">let's Go React</Navbar.Brand>
              <Navbar.Toggle aria-controls="basic-navbar-nav" />
              <Navbar.Collapse id="basic-navbar-nav">
                <Nav  className="me-auto">
                  <Nav.Link as={NavLink} to="/">Home</Nav.Link>
                  <Nav.Link as={NavLink} to="/dashboard">Dashboard</Nav.Link>
                  {(() => {
                  if (this.state.token === null){
                      return (
                        <>
                          <Nav.Link as={NavLink} to="/login">Login</Nav.Link>
                          <Nav.Link as={NavLink} to="/register">Register</Nav.Link>
                        </>
                      )
                  } else{
                    return (<Nav.Link as={NavLink} to="/logout">Logout</Nav.Link>);
                  }
                  })()}
                  
                </Nav>
              </Navbar.Collapse>
            </Container>
          </Navbar>
      );
  }
}

export default NavBar;
