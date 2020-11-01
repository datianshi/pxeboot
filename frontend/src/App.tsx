import React from 'react';
import ServerForm from './components/serverlist'
import 'bootstrap/dist/css/bootstrap.min.css'
import Container from 'react-bootstrap/esm/Container';
import Card from 'react-bootstrap/Card'

// const nicService = new FakeNicService()

function App() {
  return (
    <div className="App">
      <Card>
        <Card.Header>PXE Boot Server UI</Card.Header>
      </Card>            
      <Container className="mt-5">
        <ServerForm/>
      </Container>
    </div>
  );
}

export default App;
