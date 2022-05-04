import { CHAT_SERVER_ENDPOINT, EVENT_NAMES } from './constants';
const events = require('events');

// let webSocketConnection = null;

export const eventEmitter = new events.EventEmitter();

/**
 * @param {number} userId
 */

export class Socket {
  webSocketConnection = null;
  error = null;

  constructor(userId, username) {
    if (!window || !window['WebSocket']) {
      this.error = 'The browser does not support websockets';
      this.webSocketConnection = null;
    }
    if (!userId) {
      this.error = 'UserId is required to establish connection';
      this.webSocketConnection = null;
    }
    // creating a new socket connection with the username and the password
    this.webSocketConnection = new WebSocket(`ws://${CHAT_SERVER_ENDPOINT}/ws/${userId}/${username}`);
    this.error = null;
  }

  listen = () => {
    if (!this.webSocketConnection) return;

    this.webSocketConnection.onclose = () => {
      eventEmitter.emit(EVENT_NAMES.DISCONNECT);
    };

    this.webSocketConnection.onmessage = (event) => {
      const { eventName, eventPayload } = JSON.parse(event.data);
      switch (eventName) {
        case EVENT_NAMES.NEW_USER:
          eventEmitter.emit(EVENT_NAMES.NEW_USER, eventPayload);
          console.log('NEW USER HAS JOINED', eventPayload);
        case EVENT_NAMES.DELETED_USER:
          console.log('USER HAS DISCONNECTED', eventPayload);
        case EVENT_NAMES.DIRECT_MESSAGE:
          console.log('A DIRECT MESSAGE RECEIVED', eventPayload);
        default:
          return;
      }
    };
  };

  sendDirectMessage = (payload) => {
    if (!this.webSocketConnection) return;
    this.webSocketConnection.send(
      JSON.stringify({
        eventName: EVENT_NAMES.DIRECT_MESSAGE,
        eventPayload: payload,
      })
    );
  };
}

export const emitLogoutEvent = () => {};
