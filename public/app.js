new Vue({
  el: '#app',
  data: {
    ws: null, // websocket
    newMsg: '', // the new message to be sent to the server
    chatContent: '', // a running list of chat messages displayed on screen
    email: null, // email address and used for grabbing avatar
    username: null,
    joined: false // true if username and email != null
  },
  created: function() {
    var self = this;
    this.ws = new WebSocket('ws://' + window.location.host + '/ws');
    this.ws.addEventListener('message', function(e) {
      var msg = JSON.parse(e.data);
      self.chatContent += '<div class="chip">'
        + '<img src="' + self.gravatarURL(msg.email) + '">' // avatar
        + msg.username
        + '</div>'
        + emojione.toImage(msg.message) + '<br/>'; // parse emojis
      var element = document.getElementById('chat-messages');
      element.scrollTop = element.scrollHeight; // auto scroll to the bottom
    });
  },
  methods: {
    send: function() {
      if (this.newMsg != '') {
        this.ws.send(JSON.stringify({
          email: this.email,
          username: this.username,
          message: $('<p>').html(this.newMsg).text() // strip out html
        }));
        this.newMsg = ''; // reset newMsg to prepare for the next message
      }
    },
    join: function() {
      if (!this.email) {
        Materialize.toast('You must enter an email', 2000);
        return
      }
      if (!this.username) {
        Materialize.toast('You must enter a username', 2000);
        return
      }
      this.email = $('<p>').html(this.email).text();
      this.username = $('<p>').html(this.username).text();
      this.joined = true;
      console.log("joined!");
    },
    gravatarURL: function(email) {
      return 'http://www.gravatar.com/avatar/' + CryptoJS.MD5(email);
    }
  }
})