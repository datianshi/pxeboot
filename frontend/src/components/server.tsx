import React, {FunctionComponent} from 'react';
import 'bootstrap/dist/css/bootstrap.min.css'
import Form from 'react-bootstrap/Form'

export interface ServerProp {
    edit?: boolean;
    hostname?: string;
    ip?: string;
    mac_address?: string;
    gateway?: string;
    netmask?: string;        
}

export interface InputProp {
    serverProp?: ServerProp;
    edit: boolean;
}


const Server: FunctionComponent<InputProp> = (prop) =>

        <>
        {prop.edit? 
        <Form.Group controlId="formHostName" >
            <Form.Label>Host Name</Form.Label>
            <Form.Control placeholder="server1.example.org" />
            <Form.Text className="text-muted">
            ESX Host Name
            </Form.Text>
        </Form.Group> : null
        }
        <Form.Group controlId="formIP">
            <Form.Label>Static IP</Form.Label>
            {prop.edit? <Form.Control placeholder="example: 10.65.62.25"/>: <Form.Control placeholder="IP" value={prop.serverProp === undefined ? "" : prop.serverProp.ip}/>}            
            <Form.Text className="text-muted">
            ESX Host Static IP Address
            </Form.Text>
        </Form.Group>
        <Form.Group controlId="formMAC">
            <Form.Label>Mac Address</Form.Label>
            {prop.edit? <Form.Control placeholder="example: 00:50:A6:83:75:98"/> : <Form.Control placeholder="example: 00:50:A6:83:75:98" value={prop.serverProp === undefined ? "" : prop.serverProp.mac_address}/> }            
            <Form.Text className="text-muted">
            ESX First NIC Mac Address
            </Form.Text>
        </Form.Group>                
        <Form.Group controlId="formGateway">
            <Form.Label>Gateway</Form.Label>
            {prop.edit? <Form.Control placeholder="example: 10.65.62.1"/>: <Form.Control placeholder="example: 10.65.62.21" value={prop.serverProp === undefined ? "" :  prop.serverProp.gateway}/>}
            <Form.Text className="text-muted">
            ESX Host Management Gateway
            </Form.Text>
        </Form.Group>
        <Form.Group controlId="formNetmask">
            <Form.Label>Gateway</Form.Label>
            {prop.edit? <Form.Control placeholder="example: 255.255.255.0"/>: <Form.Control placeholder="example: 255.255.255.0" value={prop.serverProp === undefined ? "" :  prop.serverProp.netmask}/>}
            <Form.Text className="text-muted">
            ESX Host net mask
            </Form.Text>
        </Form.Group>
        </>                
export default Server    