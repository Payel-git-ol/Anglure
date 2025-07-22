const params = new URLSearchParams(window.location.search);
const username = params.get('username') || '';
document.getElementById('username').value = username;

