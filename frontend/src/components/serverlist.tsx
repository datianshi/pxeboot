import React, {ChangeEvent, Component, MouseEvent} from 'react';
import Server, {ServerProp} from './server'
import Form from 'react-bootstrap/Form'
import Button from 'react-bootstrap/Button'
import Row from 'react-bootstrap/esm/Row';
import Col from 'react-bootstrap/esm/Col';
import FormGroup from 'react-bootstrap/esm/FormGroup';
import Upload from './upload'


// interface NicService {
//     getServers(): ServerProp[];
//     deleteServer(hostname: string): void;
// }

// class FakeNicService implements NicService {
//     constructor(){
//         alert('I am the real nic service')
//     }

//     getServers(): ServerProp[] {
//         throw new Error('Method not implemented.');
//     }
//     deleteServer(hostname: string): void {
//         throw new Error('Method not implemented.');
//     }
// }

// class FakeNicService implements NicService {
//     constructor(){
//         alert('I am called')
//     }
//     mutableProps: ServerProp[] =             [
//         {
//             ip: "192.168.10.5",
//             gateway: "192.168.10.1",
//             netmask: "255.255.255.0",
//             hostname: "test-1.example.org",
//             mac_address: "00:50:A6:83:75:98"
//         },
//         {
//             ip: "192.168.10.9",
//             gateway: "192.168.10.1",
//             netmask: "255.255.255.0",
//             hostname: "test-2.example.org",
//             mac_address: "00:50:A6:83:75:97"
//         },
//         {
//             ip: "192.168.10.10",
//             gateway: "192.168.10.1",
//             netmask: "255.255.255.0",
//             hostname: "test-3.example.org",
//             mac_address: "00:50:A6:83:75:95"
//         }               
//     ]

//     deleteServer(hostname: string) {
//         this.mutableProps = this.mutableProps.filter(s => s.hostname !== hostname)
//     }
//     getServers(): ServerProp[] {
//         return (
//             this.mutableProps
//         )
//     }

// }

// export {FakeNicService}

interface MyProp {
    // nicService: NicService;
}

interface MyState {
    availableServers: ServerProp[];
    currentServer?: ServerProp
}


class ServerForm extends Component<MyProp, MyState>{
    state: MyState;
    constructor(props: MyProp) {
        super(props)
        this.state = {
            availableServers: [],
        }
        this.selectHost = this.selectHost.bind(this);
        this.deleteHost = this.deleteHost.bind(this);
        this.getServers = this.getServers.bind(this);
    }

    getServers() {
        fetch("/api/conf/nics")
        .then(res => res.json())
        .then(
          (result) => {
            this.setState({
              availableServers: result,
              currentServer: result[0]
            });
          },
          (error) => {
            alert(error)
            this.setState({
            });
          }
        )
    }

    componentDidMount() {
        this.getServers()
      }    

    selectHost(e: ChangeEvent<HTMLSelectElement>) {        
        var s: MyState = {
            availableServers: this.state.availableServers
        }
        s.currentServer = this.state.availableServers.find(server => server.hostname === e.target.value)
        this.setState(s)
    }

    deleteHost(e: MouseEvent<HTMLButtonElement>) {
        if (this.state.currentServer == null) {
            return
        }
        const url = "api/conf/nic/" + this.state.currentServer.mac_address
        alert(url)
        fetch(url, {method: 'DELETE'})
          .then(
            (resp) => {
                this.getServers()
            },
            (error) => {
              alert(error)
            }
        )                
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
                                <option key={server.hostname}>{server.hostname}</option>
                                )}
                            </Form.Control>
                            <Form.Text className="text-muted">
                            ESX Host Name
                            </Form.Text>                                                 
                        </FormGroup>
                        <Server serverProp={this.state.currentServer} edit={false} refresh={this.getServers}/>
                        <Button variant="primary" disabled={this.state.currentServer === undefined} onClick={this.deleteHost}> 
                            Delete
                        </Button>                
                    </Form>                
                </Col>
                <Col>
                    <Form>
                        <Server edit={true} refresh={this.getServers}/>
                    </Form>
                </Col>
            </Row>
            <Row>
                <Col>
                    <Upload/>    
                </Col>
                <Col/>
            </Row>
            </>
        );
    }
}

export default ServerForm