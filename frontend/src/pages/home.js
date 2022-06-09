import { Component } from 'react';
import { Navigate } from 'react-router-dom';
import CeateList from '../components/createlist' 
import Todos from '../components/todos' 

class Home extends Component {
    constructor(props) {
        super(props)
        this.state = {
            shit: false,
            token: this.getToken()
        }
    }

    changeRender = () => {
        this.setState({shit: true})
    }

    getToken = () => {
        return localStorage.getItem("token")
    }

    render () {
        if (this.state.token) {
            return <Navigate to="/dashboard/" />
        }
        return (
            <>
                {/* <CeateList renderDashboard={this.changeRender}/>
                <Todos /> */}
                <a>login</a>
            </>
        
        );
    }
}
export default Home;
