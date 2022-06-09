import { Component } from "react";
import { Container, Row, Col, Button, Modal, Form } from "react-bootstrap";

class CreateList extends Component {
    constructor(props) {
        super(props)
        this.changeRender = props.renderDashboard
        this.state = {
            token: this.getToken(),
            show: false,
            title: "",
            description: ""
        }
    }

    freeVars = () => {
        this.setState({
            title: "",
            description: ""
        })
    }
    addTask = (e) => {
        // load modal
    }
    handleClose = () => {
        this.setState({show: false})
        this.freeVars()
    }
    handleShow = () => {
        this.setState({show: true})
    }
    handleAdd = () => {
        console.log(this.state.title)
        console.log(this.state.description)
        let title = this.state.title
        let description = this.state.description
        this.createNewTaskList(title, description)
        this.freeVars()
        //after adding 
        this.changeRender()
    }
    handleTitleChange = (e) => {
        this.setState({title: e.target.value})
    }
    handleDescriptionChange = (e) => {
        this.setState({description: e.target.value})
    }
    createNewTaskList = (title, description) => {
        fetch("http://127.0.0.1:80/create-tasklist/",{
        method: "POST",
        headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                    'Authentication': this.state.token
                },
        body: `title=${title}&description=${description}&token=${this.state.token}`
        }).then(data => data.json()).then(data => {
            console.log(data)
            if (data.ok === 'false') {
                console.log('error')
                
            } else{
                console.log("added")
                this.setState({show: false})
                this.forceUpdate()
            }
        }
        )
    }
    getToken = () => {
        return localStorage.getItem("token")
    }
    render(){
        return (
        <>
            {/* <section className="containter-fluid">
                <Form >
                    <Form.Group className="mb-3">
                        <Form.Label className="col-6">Title</Form.Label>
                        <Form.Control className="col-6" type="text" placeholder="tile" />
                    </Form.Group>
                    <Form.Group className="mb-3">
                        <Form.Label>description </Form.Label>
                        <Form.Control type="text" placeholder="description" />
                    </Form.Group>
                    <Button variant="primary" type="submit">ADD</Button>
                </Form>
            </section> */}
            <div className="editor" style={{paddingTop:"50px", paddingBottom:"50px"}}>
                <Container>
                    <Row>
                        <Col>
                            <Button 
                                variant="dark"
                                style={{fontSize:"1.5rem"}}
                                onClick={this.handleShow}
                            >
                                <span style={{fontSize:"1.8rem", marginRight:"10px"}}>+</span>
                            Add New Todo-List</Button>
                        </Col>
                    </Row>
                </Container>
                </div>
                <Modal show={this.state.show} onHide={this.handleClose}>
                    <Modal.Header closeButton>
                    <Modal.Title>new task list</Modal.Title>
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

export default CreateList;