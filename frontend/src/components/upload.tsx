import React, {FunctionComponent, useState } from 'react';
import Form from 'react-bootstrap/Form';
import Button from 'react-bootstrap/Button';
import ProgressBar from 'react-bootstrap/esm/ProgressBar';
import Axios from 'axios';
import { Col, Container, Row } from 'react-bootstrap';

export interface UploadProp {
}

const Upload: FunctionComponent<UploadProp> = (prop) => {
    const [file, setFile] = useState<File | undefined >(undefined)
    const [progress, setProgress] = useState(0)
    function uploadFile() {
        if (file?.type === undefined) {
            alert('Please select a file')
            return
        }
        const formData = new FormData()
        formData.append('image', file)
        Axios.post('api/image', formData, {
            onUploadProgress: (event: ProgressEvent) => {
                setProgress(Math.round(event.loaded * 100/event.total))
            }
        })
    }
    return <>
    <Form className="pt-3">
        <Form.Label>ESXi ISO Image</Form.Label>
            <Form.Control type="file" placeholder="ESXi ISO Image" accept="*.iso" onChange={(e: React.ChangeEvent<HTMLInputElement>) : void => setFile(e.target.files![0])} ></Form.Control>
        <Container className="my-2">
        <Row>            
            <Col><Button variant="primary" onClick={uploadFile}>Upload</Button></Col>        
            <Col xs={5}>
                <ProgressBar now={progress} label={`${progress}%`}/>
            </Col>        
        </Row>
        </Container>
    </Form>
    </>
}

export default Upload