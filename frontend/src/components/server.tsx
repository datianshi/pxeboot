import React, {FunctionComponent, useState} from 'react';
import 'bootstrap/dist/css/bootstrap.min.css'
import Form from 'react-bootstrap/Form'
import Button from 'react-bootstrap/Button'

export interface ServerProp {
    hostname?: string;
    ip?: string;
    mac_address?: string;
    gateway?: string;
    netmask?: string; 
}

export interface InputProp {
    serverProp?: ServerProp;
    edit: boolean;
    refresh: Function;
}

const Server: FunctionComponent<InputProp> = (prop) => {
        const [ip, setIp] = useState("");
        const [mac, setMac] = useState("");
        const [hostname, setHostname] = useState("");
        const [gateway, setGateway] = useState("");
        const [mask, setNetmask] = useState("");

        function post() {
            fetch("/api/conf/nic", {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ 
                    ip: ip,
                    hostname: hostname,
                    mac_address: mac,
                    netmask: mask,
                    gateway: gateway
                })})
              .then(
                (resp) => {
                    prop.refresh()
                },
                (error) => {
                  alert(error)
                }
            )
        }
        return <>
        {prop.edit? 
        <Form.Group controlId="formHostName" >
            <Form.Label>Host Name</Form.Label>
            <Form.Control placeholder="server1.example.org" onChange={e => setHostname(e.target.value)}/>
            <Form.Text className="text-muted">
            ESX Host Name
            </Form.Text>
        </Form.Group> : null
        }
        <Form.Group controlId="formIP">
            <Form.Label>Static IP</Form.Label>
            {prop.edit? <Form.Control placeholder="example: 10.65.62.25" onChange={e => setIp(e.target.value)} />: <Form.Control readOnly placeholder="IP" value={prop.serverProp === undefined ? "" : prop.serverProp.ip}/>}            
            <Form.Text className="text-muted">
            ESX Host Static IP Address
            </Form.Text>
        </Form.Group>
        <Form.Group controlId="formMAC">
            <Form.Label>Mac Address</Form.Label>
            {prop.edit? <Form.Control placeholder="example: 00:50:A6:83:75:98" onChange={e => setMac(e.target.value)}/> : <Form.Control readOnly placeholder="example: 00:50:A6:83:75:98" value={prop.serverProp === undefined ? "" : prop.serverProp.mac_address}/> }            
            <Form.Text className="text-muted">
            ESX First NIC Mac Address
            </Form.Text>
        </Form.Group>                
        <Form.Group controlId="formGateway">
            <Form.Label>Gateway</Form.Label>
            {prop.edit? <Form.Control placeholder="example: 10.65.62.1" onChange={e => setGateway(e.target.value)}/>: <Form.Control readOnly placeholder="example: 10.65.62.21" value={prop.serverProp === undefined ? "" :  prop.serverProp.gateway}/>}
            <Form.Text className="text-muted">
            ESX Host Management Gateway
            </Form.Text>
        </Form.Group>
        <Form.Group controlId="formNetmask">
            <Form.Label>Netmask</Form.Label>
            {prop.edit? <Form.Control placeholder="example: 255.255.255.0" onChange={e => setNetmask(e.target.value)}/>: <Form.Control readOnly placeholder="example: 255.255.255.0" value={prop.serverProp === undefined ? "" :  prop.serverProp.netmask}/>}
            <Form.Text className="text-muted">
            ESX Host net mask
            </Form.Text>
        </Form.Group>
        {prop.edit? <Button variant="primary" onClick={post}>Add</Button> : null}
        </>
}
export default Server    