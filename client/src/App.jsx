import React, { useState, useEffect, useRef } from 'react';
import axios from 'axios';
import { TextField, Button, Box, Container, Typography, List, ListItem, ListItemText, Divider, Paper } from '@mui/material';
import SendIcon from '@mui/icons-material/Send';
import { MESSAGES_URL, CHAT_WS_URL } from './constants';

const App = () => {
  const [enteredUsername, setEnteredUsername] = useState('');
  const [confirmedUsername, setConfirmedUsername] = useState('');
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState('');
  const [isLoading, setIsLoading] = useState(true);
  const ws = useRef(null);
  const messagesEndRef = useRef(null);

  const handleUsernameSubmit = (e) => {
    e.preventDefault();
    setConfirmedUsername(enteredUsername);
  };

  const handleSendMessage = (e) => {
    e.preventDefault();
    if (input.trim() && ws.current) {
      ws.current.send(JSON.stringify({ text: input }));
      setInput('');
    }
  };

  const getMessages = async () => {
    try {
      const response = await axios.get(MESSAGES_URL);
      setMessages(response.data);
    } catch (error) {
      console.error('Error fetching messages:', error);
    }
    setIsLoading(false);
  };

  useEffect(() => {
    getMessages();
  }, []);

  useEffect(() => {
    if (confirmedUsername) {
      ws.current = new WebSocket(`${CHAT_WS_URL}${confirmedUsername}`);
      ws.current.onmessage = (event) => {
        const message = JSON.parse(event.data);
        setMessages((prevMessages) => [...prevMessages, message]);
      };

      ws.current.onopen = () => console.log('WebSocket Connected');
      ws.current.onerror = (error) => console.log('WebSocket Error: ', error);

      return () => ws.current && ws.current.close();
    }
  }, [confirmedUsername]);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  if (!confirmedUsername || isLoading) {
    return (
      <Container maxWidth="xs" style={{ marginTop: '20vh' }}>
        <Typography variant="h4" gutterBottom>
          Enter Your Username
        </Typography>
        <form onSubmit={handleUsernameSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
          <TextField
            label="Username"
            variant="outlined"
            value={enteredUsername}
            onChange={(e) => setEnteredUsername(e.target.value)}
            fullWidth
          />
          <Button variant="contained" color="primary" type="submit">
            Confirm Username
          </Button>
        </form>
      </Container>
    );
  }

  return (
    <Box sx={{ p: 3 }}>
      <Typography variant="h5" gutterBottom>
        Chat
      </Typography>
      <Paper style={{ maxHeight: '70vh', overflow: 'auto', marginBottom: '20px', padding: '10px' }}>
        <List>
          {messages.map((message, index) => (
            <React.Fragment key={index}>
              <ListItem alignItems="flex-start">
                <ListItemText
                  primary={message.username}
                  secondary={
                    <>
                      <Typography
                        sx={{ display: 'inline' }}
                        component="span"
                        variant="body2"
                        color="text.primary"
                      >
                        {message.text}
                      </Typography>
                      {" — " + new Date(message.createdAt).toLocaleTimeString()}
                    </>
                  }
                />
              </ListItem>
              {index < messages.length - 1 && <Divider variant="inset" component="li" />}
            </React.Fragment>
          ))}
          <div ref={messagesEndRef} />
        </List>
      </Paper>
      <form onSubmit={handleSendMessage} style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
        <TextField
          label="Type a message"
          variant="outlined"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          fullWidth
        />
        <Button variant="contained" color="primary" type="submit">
          <SendIcon />
        </Button>
      </form>
    </Box>
  );
};

export default App;
