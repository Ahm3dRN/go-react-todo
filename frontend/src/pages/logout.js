import React from "react";
import { Navigate } from "react-router-dom";

class Logout extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            token: this.getToken()
        };
    }
    getToken(){
        return localStorage.getItem('token')
    }
    removeToken(){
        localStorage.removeItem('token');
    }
    render() {
        if (this.state.token) {
            this.removeToken()
        }
        return <Navigate to="/login/" />
         
    }
}

export default Logout;
