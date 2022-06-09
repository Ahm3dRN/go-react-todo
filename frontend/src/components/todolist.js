import '../styles/todolist.css';
import React from "react";
import {Button, Card, Form, InputGroup, Modal} from "react-bootstrap";
import TodoItem from './todoitem';

class TodoList extends React.Component {

    constructor(props) {
        super(props);
        this.Handler = props.handler
        this.todosChange = props.todosChange
        this.state = {
            todo: props.todo,
            todos: [],
            DataisLoaded: false,
            current_todo_id:0,
            title: "",
            description: "",
            token: this.getToken()
        };
    }


    componentDidMount() {
        const task_list_id = this.state.todo.ID;
        // console.log(this.state.todo)
        // console.log("this")
        fetch(`http://127.0.0.1:80/tasks/`,{
            method: "POST",
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
                'Authentication': this.state.token
            },
            body: `task_list_id=${task_list_id}`
        })
            .then((res) => res.json())
            .then((json) => {
                this.setState({
                    todo: this.state.todo,
                    todos: json,
                    DataisLoaded: true,
                    token: this.getToken()
                });
                console.log(json)
            })
    }
    addTaskTodo = (e) => {
        this.setState({show: true})
    }
    handleTitleChange = (e) => {
        this.setState({title: e.target.value})
    } 
    handleDescriptionChange = (e) => {
        this.setState({description: e.target.value})
    }
    handleClose = () => {
        this.setState({show: false})
    }
    handleAdd = () => {
        let task_list_id = this.state.todo.ID
        let title = this.state.title
        let description = this.state.description
        console.log(task_list_id, title, description)
        console.log(this.state.token)
        this.sendAddTodo(task_list_id, title, description)
        this.setState({title:"", description: ""})
        this.todosChange()
    }
    sendAddTodo = (task_list_id, title, description) => {
        fetch("http://127.0.0.1:80/create-task/",{
            method: "POST",
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
                'Authentication': this.state.token
            },
            body: `task_list_id=${task_list_id}&title=${title}&description=${description}`
        }).then(data => data.json()).then( data => {
            console.log(data)
            if (data.ok) {  
                console.log("ok")
                this.setState({show: false, title:"", description: ""})
                this.todosChange()
            } else{
                console.log("error")
            }
        })
    }
    getToken = () => {
        return localStorage.getItem("token")
    }
    deleteTaskList = () => {
        const task_list_id = this.state.todo.ID;
        fetch(`http://127.0.0.1:80/delete-tasklist/`,{
            method: "POST",
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
                'Authentication': this.state.token
            },
            body: `task_list_id=${task_list_id}`
        }).then((res) => res.json()).then((json) => {
                if (json.ok) {
                    console.log("deleted")
                } else {
                    console.log("error")
                    console.log(json)
                }
            }
            )
    }
    render() {
        const { DataisLoaded, todos } = this.state;
        if (!DataisLoaded) return <div>
            <h1> Pleses wait some time.... </h1> </div> ;

        return (
            <>
                <Card
                    bg="dark"
                    key="dark"
                    text="white"
                    style={{ width: '100%' }}
                    className="mb-2"
                    >
                    <Card.Header>{this.state.todo.title}</Card.Header>
                    <Card.Body>
                            {todos.map((todoitem) => ( 
                                <>
                                <InputGroup style={{"paddingBottom": "1rem"}}>
                                    <TodoItem key={todoitem.id} todoitem={todoitem} task_list_id={this.state.todo.ID} handler={this.Handler} todosChange={this.todosChange}/>
                                </InputGroup>
                                </>
                                
                            ))}
                    </Card.Body>
                    <Card.Footer>
                        <Button
                            variant="secondary"
                            className="col-12"
                            onClick={this.addTaskTodo}
                        >
                            add new element
                        </Button>
                        <Button
                            className="col-12 custom-red"
                            onClick={this.deleteTaskList}
                        >
                            delete task list 
                        </Button>
                    </Card.Footer>
                </Card>
                <Modal show={this.state.show} onHide={this.handleClose}>
                    <Modal.Header closeButton>
                    <Modal.Title>Add todo item</Modal.Title>
                    </Modal.Header>
                    <Modal.Body>
                        <form>
                            <Form.Group className="mb-3">
                                <Form.Label>title</Form.Label>
                                <Form.Control type="text" name="title" onChange={this.handleTitleChange} value={this.state.title}/>
                            </Form.Group>
                            <Form.Group className="mb-3">
                                <Form.Label>description</Form.Label>
                                <Form.Control type="text" name="description" onChange={this.handleDescriptionChange} value={this.state.description}/>
                            </Form.Group>
                        </form>
                        </Modal.Body>
                    <Modal.Footer>
                    <Button variant="secondary" onClick={this.handleClose}>
                        Close
                    </Button>
                    <Button variant="primary" onClick={this.handleAdd}>
                        Save Changes
                    </Button>
                    </Modal.Footer>
                </Modal>
            </>
    );
}
}

export default TodoList;
