# Quick Connect SDK & WebSocket API Documentation

## 1. Introduction

This document serves as the comprehensive guide for the **Quick Connect Client SDK**. It covers the client-side architecture, authentication flows, storage requirements, and the full WebSocket API reference needed to implement a real-time chat system with **Optimistic UI** capabilities.

The SDK handles authentication transparently, supporting both **Guest Users** (anonymous) and **Identified Clients** (users already logged into the host application).

---

## 2. Storage & State Management

The SDK must maintain persistent state across browser sessions and application restarts. It should utilize long-term storage (e.g., `localStorage` for Web, `AsyncStorage` for Mobile).

### Required Storage Keys

1.  **`QC_TOKEN`**: The JWT authentication token used for API calls and WebSocket connections.
2.  **`QC_USER_STATE`**: Indicates the type of the current user.
    * Values: `'guest'` | `'client'`
3.  **`QC_CONVERSATION_ID`**: Stores the ID of the active conversation to ensure proper message routing.

---

## 3. Initialization Logic

When the SDK initializes, it must execute the following decision tree. A critical step is fetching the active conversation context **before** establishing the WebSocket connection.

### Logic Diagram

```mermaid
graph TD
    A[SDK Initialize] --> B{Check Storage for QC_TOKEN}
    B -- Token Exists --> C[Validate/Use Existing Token]

    B -- No Token --> E{Host App Provided User Info?}

    E -- Yes (Scenario A) --> F[Client Auth Flow]
    F --> G[Set QC_TOKEN & USER_STATE='client']

    E -- No (Scenario B) --> H[Guest Auth Flow]
    H --> I[Set QC_TOKEN & USER_STATE='guest']

    C --> J[Fetch Active Conversation API]
    G --> J
    I --> J

    J --> K[Connect WebSocket]
````

-----

## 4\. Authentication Scenarios

### Scenario A: Identified Client (Host App Integration)

Used when the user is already logged into the main application (e.g., an Online Shop).

1.  **Host App Action:** Detects a logged-in user.
2.  **Backend-to-Backend Call:** The host backend calls the Quick Connect Manager API.
    * **Endpoint:** `POST /auth/identify-client?APIKey={SECURE_KEY}`
    * **Body:** User details.
    * **Response:** Returns a `token` (QC\_TOKEN).
3.  **SDK Action:** The host frontend passes this token to the SDK (e.g., `QuickConnectSDK.login(token)`).
4.  **State Update:** Save `QC_TOKEN` and set state to `'client'`.

### Scenario B: Guest User (Anonymous)

Used when a visitor opens the site but is not logged in.

1.  **SDK Action:** Detects no existing token.
2.  **API Call:** Calls the public Guest Register endpoint.
    * **Endpoint:** `POST {MANAGER_URL}/users/guest/register`
    * **Body:**
      ```json
      {
        "fullname": "Guest User"
      }
      ```
3.  **Response:** Server returns the `qc_token`.
4.  **State Update:** Save `qc_token` as `QC_TOKEN` and set state to `'guest'`.

-----

## 5\. Pre-Connection Logic (Critical)

Before connecting to the WebSocket, the SDK **MUST** fetch the active conversation context. This ensures that `conversation_id` is available for the first message sent by the user.

### Fetch Active Conversation

* **Endpoint:** `GET {CHAT_API_URL}/conversations/active`
* **Headers:** `Authorization: Bearer <QC_TOKEN>`
* **Response:**
  ```json
  {
    "id": "01K...", // Store this value as QC_CONVERSATION_ID
    "status": "new",
    "created_at": "..."
  }
  ```

-----

## 6\. WebSocket Connection

### Connection Details

* **URL:** `ws://{CHAT_URL}/chats/clients`

### Browser Constraints & Authentication

Standard browser `WebSocket` API **does not support** custom headers (like `Authorization`). To bypass this limitation securely:

1.  **Solution:** Send the token in the **`Sec-WebSocket-Protocol`** header.
2.  **Client Implementation:**
    ```javascript
    const protocols = [token]; // Token must be the first element
    const socket = new WebSocket(wsUrl, protocols);
    ```
3.  **Server Validation:** The server extracts the token from the protocol header, validates it, and echoes it back in the response handshake to confirm the connection.

### Connection Events

* **onOpen:** Connection established. Set UI to "Online".
* **onMessage:** Parse JSON message. Trigger **Reconciliation Logic** (see Section 9).
* **onClose:** Implement Reconnection Strategy (Exponential backoff: 1s, 2s, 5s...).
* **onError:** Handle connection errors.

-----

## 7\. API Endpoint Reference

### Manager Service (Authentication)

| Method | Endpoint | Auth Required | Description |
| :--- | :--- | :--- | :--- |
| **POST** | `/users/guest/register` | No | Registers a new guest. Returns `qc_token`. |
| **PUT** | `/users/guest/update` | **Yes** (JWT) | Updates the guest's profile information. |

### Chat Service (Real-time)

| Method | Endpoint | Auth Required | Description |
| :--- | :--- | :--- | :--- |
| **GET** | `/conversations/active` | **Yes** (JWT) | **Required step.** Returns current conversation context/ID. |
| **GET** | `/chats/clients` | **Yes** (Protocol) | WebSocket endpoint for clients. |

-----

## 8\. Message Structures

All WebSocket messages are formatted as **JSON**.

### 8.1. Sending Messages (Client $\rightarrow$ Server)

The client **MUST** generate a unique `client_message_id` for every text message to enable optimistic UI updates.

```json
{
  "type": "text",               // Required: "text" | "system"
  "conversation_id": "string",   // Required: From /conversations/active API
  "content": "Hello Support",    // Required for text
  "client_message_id": "temp_123" // Required: Unique ID generated by client
}
```

### 8.2. Receiving Messages (Server $\rightarrow$ Client)

The server broadcasts the message back to **all** participants (including the sender). The payload contains the full database object.

```json
{
  "type": "text",
  "timestamp": "2025-11-18T10:30:00Z",
  "payload": {
    "id": "01K...",               // Real Database ID
    "conversation_id": "01K...",
    "sender_id": "01K...",
    "content": "Hello Support",
    "client_message_id": "temp_123" // Echoed back for reconciliation
  }
}
```

-----

## 9\. Optimistic UI & Reconciliation Strategy

To ensure a smooth, responsive user experience, the SDK should implement the following flow:

### Step 1: Sending (Client Side)

1.  User sends a message.
2.  SDK generates a temporary ID (e.g., `temp_xyz`).
3.  SDK adds the message to the UI immediately with a **Pending** status (e.g., Clock icon, 0.7 opacity).
4.  SDK sends the JSON payload including `"client_message_id": "temp_xyz"`.

### Step 2: Receiving & Reconciliation

The server processes the message and broadcasts it via WebSocket.

1.  **SDK receives a message.**
2.  **Check 1 (Reconciliation):** Does the payload contain a `client_message_id` that matches a pending message in the DOM?
    * **YES:**
        * Find the DOM element with ID `temp_xyz`.
        * Change status to **Sent/Confirmed** (Check icon, 1.0 opacity).
        * Update the DOM element ID to the real `id` from the server payload.
        * **Stop processing.** (Do not render a duplicate message).
3.  **Check 2 (Fallback/Echo Prevention):**
    * If reconciliation failed (e.g., page refresh), check: `payload.sender_id === current_user_id`.
    * **YES:** Ignore the message (it's an echo of our own message).
    * **NO:** Display as a new incoming message from the support agent.

-----

## 10\. System Messages

System messages are ephemeral (not saved to DB) and used for typing indicators.

**Sending Typing Started:**

```json
{
  "type": "system",
  "sub_type": "typing_started",
  "conversation_id": "01K..."
}
```

**Receiving Typing Started:**

```json
{
  "type": "system",
  "sub_type": "typing_started",
  "payload": {
    "conversation_id": "01K...",
    "sender_id": "01K..."
  }
}
```

**Frontend Logic:**

* **Throttling:** Only send `typing_started` once every 5 seconds.
* **Display:** When receiving `typing_started`, show "Support is typing..." for 6 seconds, then auto-hide if no further updates or `typing_stopped` events are received.

-----

## 11\. Implementation Checklist

1.  [ ] **Storage Adapter:** Implement wrapper for `localStorage` / `AsyncStorage`.
2.  [ ] **Auth Manager:** Handle Guest Registration (save `qc_token`) vs Client Login.
3.  [ ] **Context Fetch:** Call `/conversations/active` before WS connection to get `conversation_id`.
4.  [ ] **WebSocket Client:**
    * Use `Sec-WebSocket-Protocol` for authentication.
    * Implement **Optimistic UI** (add message immediately as pending).
    * Implement **Reconciliation** (match `client_message_id` to confirm delivery).
    * Implement **Typing Indicators** with throttling and timeouts.