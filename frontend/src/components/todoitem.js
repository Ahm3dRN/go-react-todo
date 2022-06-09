import React from "react";
import '../styles/todoitem.css';
import {Form, Button} from 'react-bootstrap';

class TodoItem extends React.Component {

    // Constructor 
    constructor(props) {
        super(props);
        this.Handler = props.handler
        this.todosChange = props.todosChange
        this.state = {
            todo: props.todoitem,
            task_list_id: props.task_list_id,
            token: this.getToken()
        };
    }
    HandleChange = (e) => {
        console.log("handle change")
        let todo_id = e.target.id
        fetch(`http://127.0.0.1:80/check-task/`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
                "Authentication": this.state.token
            },
            body: `task_list_id=${this.state.task_list_id}&task_id=${todo_id}`
        })
            .then((res) => res.json())
            .then((json) => {
                console.log(json)
                if (json.ok) {
                    this.setState({
                        todo:json.todo
                    });
                    console.log(json)
                } else{
                    console.log(json)
                    console.log("error")
                }
            })
        console.log("here")
    }
    HandleDelete = (e) => {
        console.log("delete task")
        let todo_id = e.target.id
        fetch(`http://127.0.0.1:80/delete-task/`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
                "Authentication": this.state.token
            },
            body: `task_list_id=${this.state.task_list_id}&task_id=${todo_id}`
        })
            .then((res) => res.json())
            .then((json) => {
                console.log(json)
                if (json.ok) {
                    console.log(json)
                    this.todosChange()
                } else{
                    console.log(json)
                    console.log("error")
                }
            })
    }
    HandleEdit = (e) => {
        console.log("in edit")
        let todo_id = e.target.id
        this.Handler(todo_id, this.state.task_list_id)
    }
    getToken = () => {
        return localStorage.getItem("token")
    }
        render() {
        const todo = this.state.todo;
        return (
            <>
                {/* <Card.Title> */}
                    <Form.Check 
                        // inline
                        defaultChecked={todo.is_complete}
                        onChange={this.HandleChange}
                        className="custom-checkbox h5 col-10" 
                        type="checkbox"
                        id={`${todo.ID}`}
                        label={todo.title}
                    />
                    <Button variant="outline-light col-1" onClick={this.HandleEdit} id={`${todo.ID}`} size="sm">E</Button>
                    <Button variant="outline-light col-1" onClick={this.HandleDelete} id={`${todo.ID}`} size="sm">D</Button>
                    {/* <InputGroup.Checkbox className="custom-checkbox" label="Checkbox for following text input" /> */}
            {/* </Card.Title> */}
            </>

    );
}
}

export default TodoItem;
