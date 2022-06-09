import React from "react";
import { Outlet, Link } from "react-router-dom";
import Header from "../components/header";

class Layout extends React.Component {

    // Constructor 
    constructor(props) {
        super(props);
        this.state = {};
    }

    render() {
        return (
            <>
                <Header />
                <Outlet />
            </>
        );
    }
}

export default Layout;
