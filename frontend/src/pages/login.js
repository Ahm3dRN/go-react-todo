import React from "react";
import {Form, Button, Container,FloatingLabel } from "react-bootstrap";
import { Navigate } from "react-router-dom";

class Login extends React.Component {

    // Constructor 
    constructor(props) {
        super(props);
        this.state = {
            username: "",
            password: "",
            usernameError: "",
            passwordError: "",
            evalid: false,
            token: this.getToken()
        };
    }
    setUserName (username) {
        this.setState({username: username})
    }

    setPassword (password) {
        this.setState({password: password})
    }

    handleSubmit (event){
        event.preventDefault();
        event.stopPropagation();
        let username = this.state.username;
        let password = this.state.password;
        this.sendLoginRequest(username, password);
    }
    setToken(token){
        localStorage.setItem('token', JSON.stringify(token));
        this.setState({token: token});
    }
    getToken(){
        return localStorage.getItem('token')
    }

    sendLoginRequest(username, password){
        console.log(username, password);
        console.log(JSON.stringify({username:username, password:password}))
        fetch(`http://127.0.0.1:80/users/login/`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded'
            },
            body: `username=${username}&password=${password}`
        })
        .then(data => data.json()).then(data => {
            console.log(data)
            if (data.ok === 'false') {
                this.setState({evalid: false})
                if (data.username){
                    this.setState({usernameError: data.username})
                    console.log("yea")
                }
                if (data.password){
                    this.setState({passwordError: data.password})
                }
            } else{
                this.setToken(data.token)
                this.setState({evalid: true})
            }
        }
        )
    }
    render() {
        if (this.state.token) {
            return <Navigate to="/dashboard" />
        }
        return (
            <>
            <div className="row" style={{height:'20vh'}}></div>
            <Container className="justify-content-md-center col-6 vertical-center">
                <Form validated={this.state.evalid} onSubmit={this.handleSubmit.bind(this)} noValidate >
                    <Form.Group className="mb-3" controlId="formBasicUsername" hasValidation>
                    <FloatingLabel
                        controlId="floatingInput"
                        label="Username"
                        className="mb-3"
                    >
                        <Form.Control 
                            required 
                            isInvalid={!!this.state.usernameError}
                            type="text"
                            placeholder="Enter username" 
                            onChange={e => this.setUserName(e.target.value)}
                            feedback={this.state.usernameError}
                            feedbacktype="invalid"
                        />
                        </FloatingLabel>
                        <Form.Control.Feedback type="invalid">
                            {this.state.usernameError}
                            {/* Please choose a username. */}
                        </Form.Control.Feedback>
                    </Form.Group>
                    <Form.Group className="mb-3" controlId="formBasicPassword">
                    <FloatingLabel
                        controlId="floatingInput"
                        label="Password"
                        className="mb-3"
                    >
                        <Form.Control required isInvalid={!!this.state.passwordError}  type="password" placeholder="Password" onChange={e => this.setPassword(e.target.value)}/>
                    </FloatingLabel>
                    <Form.Control.Feedback type="invalid">
                            {this.state.passwordError}
                    </Form.Control.Feedback>
                    </Form.Group>
                    <Button variant="dark" type="submit">
                        Login
                    </Button>
                </Form>
                <p>{this.state.username}</p>
                <p>{this.state.password}</p>
            </Container>
            </>
        );
    }
}

export default Login;
