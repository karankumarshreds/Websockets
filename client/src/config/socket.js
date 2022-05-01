import { CHAT_SERVER_ENDPOINT, EVENT_NAMES } from './constants';
const events = require('events');

// let webSocketConnection = null;

export const eventEmitter = new events.EventEmitter();

/**
 * @param {number} userId
 */

export class Socket {
  webSocketConnection = null;

  constructor(userId) {
    if (!window || !window['WebSocket']) {
      return {
        error: 'The browser does not support websockets',
        webSocketConnection: null,
      };
    }
    if (!userId) {
      return {
        error: 'UserId is required to establish connection',
        webSocketConnection: null,
      };
    }

    this.webSocketConnection = new WebSocket(`ws://${CHAT_SERVER_ENDPOINT}/ws/${userId}`);
    return {
      error: null,
      webSocketConnection: this.webSocketConnection,
    };
  }
  sendDirectMessage = (payload) => {
    if (!this.webSocketConnection) return;
    this.webSocketConnection.send(
      JSON.stringify({
        eventName: EVENT_NAMES.DIRECT_MESSAGE,
        payload,
      })
    );
  };
}

// export const newSocketConnection = (userId) => {
//   if (!window['WebSocket']) {
//     return {
//       error: 'The browser does not support websockets.',
//       webSocketConnection: null,
//     };
//   }
//   if (!userId) {
//     return {
//       error: 'UserId is required to establish connection.',
//       webSocketConnection: null,
//     };
//   }
//   webSocketConnection = new WebSocket(`ws://${CHAT_SERVER_ENDPOINT}/ws/${userId}`);
//   return {
//     error: null,
//     webSocketConnection,
//   };
// };

// export const sendDirectMessage = (payload) => {
//   if (!webSocketConnection) return;
//   webSocketConnection.send(
//     JSON.stringify({
//       eventName: EVENT_NAMES.DIRECT_MESSAGE,
//       payload: payload,
//     })
//   );
// };

export const emitLogoutEvent = () => {};

// export const listenToSocketEvents = () => {
//   if (!webSocketConnection) return;
// };
