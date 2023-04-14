const joinButton = document.getElementById('join');
const sendButton = document.getElementById('send');
const messageInput = document.getElementById('message');
const roomInput = document.getElementById('room');
const nicknameInput = document.getElementById('nickname');
const messagesDiv = document.getElementById('messages');
const chatDiv = document.getElementById('chat');
const joinForm = document.getElementById('join-form');

chatDiv.style.display = 'none';

let room;
let user;
let socket;

joinButton.addEventListener('click', () => {
  user = { nickname: nicknameInput.value };
  room = roomInput.value;

  if (!user.nickname || !room) {
    alert('Please enter a nickname and room name.');
    return;
  }

  socket = new WebSocket('ws://localhost:8080/ws');

  socket.addEventListener('open', () => {
    joinForm.style.display = 'none';
    chatDiv.style.display = 'block';
    nicknameInput.disabled = true;

    const joinMessage = {
      action: 'join',
      nickname: user.nickname,
      room: room,
    };
    socket.send(JSON.stringify(joinMessage));
  });

  socket.addEventListener('message', (event) => {
    const messageDiv = document.createElement('div');
    messageDiv.textContent = event.data;
    messagesDiv.appendChild(messageDiv);
  });

  socket.addEventListener('close', () => {
    alert('Connection closed.');
  });
});

sendButton.addEventListener('click', () => {
  sendMessage();
});

messageInput.addEventListener('keydown', (event) => {
  if (event.key === 'Enter') {
    event.preventDefault();
    sendMessage();
  }
});

function sendMessage() {
  const text = messageInput.value;
  if (!text) {
    return;
  }

  const message = {
    action: 'message',
    room: room,
    text: text,
  };
  socket.send(JSON.stringify(message));
  messageInput.value = '';
}