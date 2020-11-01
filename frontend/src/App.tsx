import React from 'react';
import ServerForm from './components/serverlist'
import 'bootstrap/dist/css/bootstrap.min.css'
import Container from 'react-bootstrap/esm/Container';
import Card from 'react-bootstrap/Card'

// const nicService = new FakeNicService()

function App() {
  return (
    <>
    <div className="App">
      <Card bg="dark" text="light">
        <Card.Header>PXE Boot Server UI</Card.Header>
      </Card>            
      <Container className="mt-5">
        <ServerForm/>
      </Container>
    </div>
    <footer className='footer mt-auto py-3 bg-dark text-white'>
      <div className='container'></div>
    </footer>
    </>    
  );
}

export default App;
