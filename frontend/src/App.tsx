import React from 'react';
import ServerForm, {FakeNicService} from './components/serverlist'
import 'bootstrap/dist/css/bootstrap.min.css'
import Container from 'react-bootstrap/esm/Container';
import Card from 'react-bootstrap/Card'

const nicService = new FakeNicService()

function App() {
  return (
    <div className="App">
      <Card>
        <Card.Header>PXE Boot Server UI</Card.Header>
      </Card>            
      <Container className="mt-5">
        <ServerForm nicService={nicService}/>
      </Container>
    </div>
  );
}

export default App;
