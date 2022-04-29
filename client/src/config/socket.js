const events = require('events');

const CHAT_SERVER_ENDPOINT = '127.0.0.1:4444';
let webSocketConnection = null;

export const eventEmitter = new events.EventEmitter();

/**
 * @param {number} userId
 */
export const newSocketConnection = (userId) => {
  if (!window['WebSocket']) {
    return {
      error: 'The browser does not support websockets.',
      webSocketConnection: null,
    };
  }
  if (!userId) {
    return {
      error: 'UserId is required to establish connection.',
      webSocketConnection: null,
    };
  }
  webSocketConnection = new WebSocket(`ws://${CHAT_SERVER_ENDPOINT}/ws/${userId}`);
  return {
    error: null,
    webSocketConnection,
  };
};
