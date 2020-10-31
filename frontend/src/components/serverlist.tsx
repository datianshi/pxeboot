import React, {ChangeEvent, Component, MouseEvent} from 'react';
import Server, {ServerProp} from './server'
import Form from 'react-bootstrap/Form'
import Button from 'react-bootstrap/Button'
import Row from 'react-bootstrap/esm/Row';
import Col from 'react-bootstrap/esm/Col';
import FormGroup from 'react-bootstrap/esm/FormGroup';


interface NicService {
    getServers(): ServerProp[];
    deleteServer(hostname: string): void;
}

class FakeNicService implements NicService {
    constructor(){
        alert('I am called')
    }
    mutableProps: ServerProp[] =             [
        {
            ip: "192.168.10.5",
            gateway: "192.168.10.1",
            netmask: "255.255.255.0",
            hostname: "test-1.example.org",
            mac_address: "00:50:A6:83:75:98"
        },
        {
            ip: "192.168.10.9",
            gateway: "192.168.10.1",
            netmask: "255.255.255.0",
            hostname: "test-2.example.org",
            mac_address: "00:50:A6:83:75:97"
        },
        {
            ip: "192.168.10.10",
            gateway: "192.168.10.1",
            netmask: "255.255.255.0",
            hostname: "test-3.example.org",
            mac_address: "00:50:A6:83:75:95"
        }               
    ]

    deleteServer(hostname: string) {
        this.mutableProps = this.mutableProps.filter(s => s.hostname !== hostname)
    }
    getServers(): ServerProp[] {
        return (
            this.mutableProps
        )
    }

}

export {FakeNicService}

interface MyProp {
    nicService: NicService;
}

interface MyState {
    availableServers: ServerProp[];
    currentServer?: ServerProp
}


class ServerForm extends Component<MyProp, MyState>{
    state: MyState;
    constructor(props: MyProp) {
        super(props)
        let response: ServerProp[] = props.nicService.getServers()
        this.state = {
            availableServers: response,
            currentServer: response[0]
        }
        this.selectHost = this.selectHost.bind(this);
        this.deleteHost = this.deleteHost.bind(this);
    }

    selectHost(e: ChangeEvent<HTMLSelectElement>) {        
        var s: MyState = {
            availableServers: this.props.nicService.getServers() 
        }
        s.currentServer = this.props.nicService.getServers().find(server => server.hostname === e.target.value)
        this.setState(s)
    }

    deleteHost(e: MouseEvent<HTMLButtonElement>) {        
        if (this.state.currentServer != null && this.state.currentServer.hostname != null) {
            alert('delete, Hello')
            this.props.nicService.deleteServer(this.state.currentServer.hostname)
        }
        var s: MyState = {
            availableServers: this.props.nicService.getServers() 
        }
        alert(s.availableServers.length)
        s.currentServer = s.availableServers.length === 0 ? undefined : s.availableServers[0]
        this.setState(s)
    }

    render() {
        return (
            <>
            <Row>
                <Col>
                    <Form>
                        <FormGroup>
                            <Form.Label>Host Name</Form.Label>
                            <Form.Control as="select" onChange={this.selectHost}>
                                {this.state.availableServers.map( server => 
                                <option>{server.hostname}</option>
                                )}
                            </Form.Control>
                            <Form.Text className="text-muted">
                            ESX Host Name
                            </Form.Text>                                                 
                        </FormGroup>
                        <Server serverProp={this.state.currentServer} edit={false}/>
                        <Button variant="primary" type="submit" disabled={this.state.currentServer === undefined} onClick={this.deleteHost}> 
                            Delete
                        </Button>                
                    </Form>                
                </Col>
                <Col>
                    <Form>
                        <Server edit={true}/>
                        <Button variant="primary" type="submit"> 
                            Add
                        </Button>                
                    </Form>
                </Col>
            </Row>
            </>
        );
    }
}

export default ServerForm