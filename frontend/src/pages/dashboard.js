import React from "react";
import Todos from "../components/todos"
// import Login from "../pages/login"
import {Navigate } from "react-router-dom"
import CeateList from '../components/createlist'
class Dashboard extends React.Component {

    // Constructor 
    constructor(props) {
        super(props);
        this.state = {
            token: this.getToken(),
            rerender: false
        };
    }
    changeRender = () => {
        console.log("change render")
        this.forceUpdate();
    }
    getToken = () => {
        return localStorage.getItem("token")
    }
    render() {
        if (!this.state.token) {
            return <Navigate to="/login" />
        }
        return (
            <>
                <CeateList renderDashboard={this.changeRender}/>
                <Todos />i
            </>
        );
    }
}

export default Dashboard;
