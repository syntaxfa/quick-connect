# Quick Connect Chat WebSocket API Documentation

## 1. Overview

The Quick Connect Chat Service uses **WebSockets** to provide real-time, bidirectional communication between Clients and Support Agents. The system is built on a **Pub/Sub architecture**, ensuring scalability and persistent connection handling throughout the application lifecycle.

## 2. Connection & Authentication

### Endpoints

There are two distinct endpoints depending on the user role. Both require a valid **WebSocket (ws:// or wss://)** connection upgrade.

| Role | Method | Endpoint | Access Level |
| :--- | :--- | :--- | :--- |
| **Client** | GET | `/chats/clients` | `RoleGuest`, `RoleClient` |
| **Support** | GET | `/chats/supports` | `RoleSupport` |

*Base URL Example:* `ws://localhost:2530/chats/clients`

### Authentication

The connection is secured via **JWT (JSON Web Token)**. The token must be present in the request (typically via cookies or Authorization headers, depending on your HTTP client configuration) to identify the `user_id`.

* **Note:** The server extracts the `user_id` from the token. You do not need to send `sender_id` in the message payload; the server handles this securely.

---

## 3. Message Structures

All messages sent and received are formatted as **JSON**.

### 3.1. Sending Messages (Client $\rightarrow$ Server)

This structure is used when you send data to the websocket.

```json
{
  "type": "string",           // Required: "text" | "system"
  "sub_type": "string",       // Required for "system" messages
  "conversation_id": "string", // Required: ULID of the active conversation
  "content": "string"         // Required for "text" messages
}
````

### 3.2. Receiving Messages (Server $\rightarrow$ Client)

This structure is used when you receive data from the websocket.

```json
{
  "type": "string",           // "text" | "system"
  "sub_type": "string",       // e.g., "typing_started"
  "timestamp": "string",      // ISO 8601 Date Time
  "payload": { ... }          // Dynamic object (Message Object or System Object)
}
```

-----

## 4\. Message Types & Workflows

### 4.1. Text Messages (`text`)

Used for standard chat messages. These are saved to the database and broadcast to all participants.

**Sending a Text Message:**

```json
{
  "type": "text",
  "conversation_id": "01KAAWW5V95WFFDERJZCTDE2QQ",
  "content": "Hello, I need help with my order."
}
```

**Receiving a Text Message:**
The `payload` will contain the full Database Message Object.

```json
{
  "type": "text",
  "sub_type": "",
  "timestamp": "2025-11-18T10:30:00Z",
  "payload": {
    "id": "01KAB...",
    "conversation_id": "01KAAWW5V95WFFDERJZCTDE2QQ",
    "sender_id": "01K9RDKY84GC0Q7F9F1NY54B5C",
    "message_type": "text",
    "content": "Hello, I need help with my order.",
    "created_at": "2025-11-18T10:30:00Z",
    "read_at": null
  }
}
```

-----

### 4.2. System Messages (`system`)

Used for ephemeral states like typing indicators. These are **not saved** to the database and are only broadcast to *other* participants in the conversation.

#### Supported Sub-Types:

1.  `typing_started`
2.  `typing_stopped`

**Sending "Typing Started":**

```json
{
  "type": "system",
  "sub_type": "typing_started",
  "conversation_id": "01KAAWW5V95WFFDERJZCTDE2QQ"
}
```

**Sending "Typing Stopped":**

```json
{
  "type": "system",
  "sub_type": "typing_stopped",
  "conversation_id": "01KAAWW5V95WFFDERJZCTDE2QQ"
}
```

**Receiving System Messages:**
The `payload` contains identification data.

```json
{
  "type": "system",
  "sub_type": "typing_started",
  "timestamp": "2025-11-18T10:30:05Z",
  "payload": {
    "conversation_id": "01KAAWW5V95WFFDERJZCTDE2QQ",
    "sender_id": "01K9RDKY84GC0Q7F9F1NY54B5C"
  }
}
```

-----

## 5\. Frontend Implementation Guidelines

To ensure a smooth user experience and reduce server load, follow these logic rules for Typing Indicators.

### 5.1. Sending Logic (The 5-Second Rule)

Do not send a `typing_started` event for every keystroke. Use a throttling mechanism.

**Algorithm:**

1.  **User presses a key.**
2.  Check `lastTypingSentTime`.
3.  **IF** `(CurrentTime - lastTypingSentTime) > 5000ms` (5 seconds):
    * Send `{"type": "system", "sub_type": "typing_started", ...}`
    * Update `lastTypingSentTime` to now.
4.  **IF** User stops typing (e.g., `keyup` event with a debounce of 1000ms):
    * Send `{"type": "system", "sub_type": "typing_stopped", ...}`
    * Reset `lastTypingSentTime`.

### 5.2. Receiving Logic (Display Handling)

The receiver should treat the `typing_started` signal as valid for a limited time window.

**Algorithm:**

1.  **Receive** `typing_started` from `Sender X`.
2.  **Show** UI: "Sender X is typing..."
3.  **Clear existing timeout** (if any) for Sender X.
4.  **Set new timeout** (e.g., 6 seconds):
    * *Action:* Hide the typing indicator automatically if no new signal arrives.
5.  **IF Receive** `typing_stopped` from `Sender X`:
    * **Hide** UI immediately.
    * Clear any pending timeouts.

-----

## 6\. Error Handling

If an error occurs (e.g., invalid JSON, closed conversation, or permission denied), the server logs the error and may close the connection depending on severity.

* **Validation Errors:** If the payload is missing `conversation_id` or `type`, the message is dropped, and an error is logged on the server.
* **Closed Conversations:** If you attempt to send a message to a `closed` conversation, the server will ignore the message. The UI should disable the input box when the conversation status is closed.