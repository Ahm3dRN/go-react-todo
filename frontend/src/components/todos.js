
import React from "react";
import { Col, Container, Row, Modal, Button, Form } from "react-bootstrap";
import TodoList from "./todolist";

class Todos extends React.Component {

    // Constructor 
    constructor(props) {
        super(props);
        this.Handler = this.Handler.bind(this);
        this.handleClose = this.handleClose.bind(this);
        this.handleTitleChange = this.handleTitleChange.bind(this);
        this.handleFieldChange = this.handleFieldChange.bind(this)
        this.fetchTaskList = this.fetchTaskList.bind(this)
        this.sendChangeTitle = this.sendChangeTitle.bind(this)
        this.state = {
            todoslist: [],
            DataisLoaded: false,
            show: false,
            task_list_id:0,
            task_id: 0,
            title: "",
            token: this.getToken(),

        };
    }
    getToken = () => {
        return localStorage.getItem("token")
    }
    Handler (todo_id, task_list_id) {
        this.setState({task_list_id:task_list_id, task_id: todo_id, show:true})
    }
    handleClose(){
        this.setState({show: false, title:""})
    }
    handleFieldChange(e){
        console.log(e.target.value)
        this.setState({title: e.target.value})
    }
    handleTitleChange (e){
        e.preventDefault();
        e.stopPropagation();
        console.log(e.target.value)
        let task_list_id = this.state.task_list_id;
        let task_id = this.state.task_id;
        this.sendChangeTitle(task_list_id, task_id, this.state.title);
        this.afterChange()
        this.freeTitle()
    }
    sendChangeTitle (task_list_id, task_id, title){
        fetch(`http://127.0.0.1:80/edit-task/`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
                'Authentication': this.state.token
            },
            body: `task_list_id=${task_list_id}&task_id=${task_id}&title=${title}`
        })
        .then(data => data.json()).then(data => {
            console.log(data)
            if (data.ok === 'false') {
                console.log('error')
                
            } else{
                this.setState({todoslist: []})
                this.fetchTaskList()
                this.setState({show: false, title:""})
            }
        }
        )
    }
    fetchTaskList = () => {
        fetch("http://127.0.0.1:80/tasklists",{
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
                'Authentication': this.state.token
            },
        })
            .then((res) => res.json())
            .then((json) => {
                this.setState({
                    todoslist: json,
                    DataisLoaded: true
                });
                console.log(json)
            })
    }
    componentDidMount() {
        this.fetchTaskList()
    }
    afterChange = () => {
        console.log("this called")
        this.setState({todoslist: []})
        this.fetchTaskList()
        this.forceUpdate();
    }
    freeTitle = () => {
        this.setState({title: ""})
    }
    render() {
        const { DataisLoaded, todoslist } = this.state;
        console.log(todoslist)
        if (!DataisLoaded) return <div>
            <h1> Pleses wait some time.... </h1> </div> ;

        return (
            <>
                <Modal show={this.state.show} onHide={this.handleClose}>
                    <Modal.Header closeButton>
                    <Modal.Title>Change title</Modal.Title>
                    </Modal.Header>
                    <Modal.Body>
                        Woohoo, you're reading {this.state.task_list_id} on {this.state.task_list_id}
                        <form>
                            <Form.Control type="title"  onChange={this.handleFieldChange} value={this.state.title}/>
                        </form>
                        </Modal.Body>
                    <Modal.Footer>
                    <Button variant="secondary" onClick={this.handleClose}>
                        Close
                    </Button>
                    <Button variant="primary" onClick={this.handleTitleChange}>
                        Save Changes
                    </Button>
                    </Modal.Footer>
                </Modal>
            <div className = "todos-lists">
                <Container fluid>
                <Row>
                    {
                        todoslist.map((todo) => ( 
                            <Col sm="12" lg="4" key={`col-${todo.ID}`}>
                                <TodoList key={todo.ID} todo={todo} handler={this.Handler} todosChange={this.afterChange}/>
                            </Col>
                        ))
                    }
                </Row>
                </Container>
            </div>
            </>
    );
}
}

export default Todos;
