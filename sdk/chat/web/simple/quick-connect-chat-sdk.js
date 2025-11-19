/**
 * Quick Connect SDK - Client Side
 * Version: 2.2.0
 * Changes: Fixed echo duplication, Added JWT parsing, Reverted to Persian UI
 */

(function(window, document) {
    'use strict';

    // ==========================================
    // 1. CONFIGURATION & CONSTANTS
    // ==========================================
    const CONFIG = {
        API_BASE_URL: 'http://localhost:2531', // Manager Service
        CHAT_WS_URL: 'ws://localhost:2530/chats/clients',
        CHAT_API_URL: 'http://localhost:2530',
        TAILWIND_CDN: 'https://cdn.tailwindcss.com',

        // Storage Keys
        STORAGE_KEYS: {
            TOKEN: 'QC_TOKEN',
            USER_STATE: 'QC_USER_STATE', // 'guest' | 'client'
            CONVERSATION_ID: 'QC_CONVERSATION_ID'
        },

        // Timers
        TYPING_THROTTLE_MS: 5000,
        TYPING_DISPLAY_MS: 6000,
        RECONNECT_INTERVAL: 3000
    };

    // UI Strings (Persian)
    const UI_TEXT = {
        TITLE: 'پشتیبانی آنلاین',
        WELCOME: 'برای شروع گفتگو، پیام خود را تایپ کنید.',
        CONNECTING: 'در حال اتصال...',
        CONNECTED: 'متصل شد.',
        DISCONNECTED: 'اتصال قطع شد.',
        INPUT_PLACEHOLDER: 'پیام خود را بنویسید...',
        TYPING: 'پشتیبان در حال نوشتن...'
    };

    // ==========================================
    // 2. SDK LOGIC CLASS
    // ==========================================
    class QuickConnectSDK {
        constructor() {
            this.token = localStorage.getItem(CONFIG.STORAGE_KEYS.TOKEN);
            this.userState = localStorage.getItem(CONFIG.STORAGE_KEYS.USER_STATE);
            this.conversationId = localStorage.getItem(CONFIG.STORAGE_KEYS.CONVERSATION_ID) || null;
            this.userId = null; // We will extract this from token

            this.socket = null;
            this.lastTypingSentTime = 0;
            this.typingTimeout = null;

            // Event Callbacks
            this.onMessageReceived = null;
            this.onSystemMessage = null;
            this.onTypingStatusChange = null;
        }

        /**
         * Helper: Parse JWT to get User ID
         */
        parseJwt(token) {
            try {
                const base64Url = token.split('.')[1];
                const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
                const jsonPayload = decodeURIComponent(window.atob(base64).split('').map(function(c) {
                    return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
                }).join(''));
                return JSON.parse(jsonPayload);
            } catch (e) {
                return null;
            }
        }

        /**
         * Initialization Flow
         */
        async init(hostProvidedToken = null) {
            // 1. Authentication
            if (hostProvidedToken) {
                this.setSession(hostProvidedToken, 'client');
            } else if (!this.token) {
                await this.registerGuest();
            }

            // Extract User ID if token exists
            if (this.token) {
                const decoded = this.parseJwt(this.token);
                if (decoded && decoded.user_id) {
                    this.userId = decoded.user_id;
                    console.log("QC SDK: User ID detected:", this.userId);
                }
            }

            // 2. Connect
            if (this.token) {
                await this.fetchActiveConversation();
                this.connectWebSocket();
            }
        }

        setSession(token, type) {
            this.token = token;
            this.userState = type;
            localStorage.setItem(CONFIG.STORAGE_KEYS.TOKEN, token);
            localStorage.setItem(CONFIG.STORAGE_KEYS.USER_STATE, type);
        }

        /**
         * Guest Registration
         */
        async registerGuest() {
            try {
                const response = await fetch(`${CONFIG.API_BASE_URL}/users/guest/register`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ fullname: "Guest Visitor" })
                });

                if (!response.ok) throw new Error('Guest registration failed');
                const data = await response.json();

                if (data.qc_token) {
                    this.setSession(data.qc_token, 'guest');
                }
            } catch (error) {
                console.error("QC SDK: Auth Error", error);
            }
        }

        /**
         * Fetch Active Conversation
         */
        async fetchActiveConversation() {
            if (!this.token) return;

            try {
                const response = await fetch(`${CONFIG.CHAT_API_URL}/conversations/active`, {
                    method: 'GET',
                    headers: {
                        'Authorization': `Bearer ${this.token}`,
                        'Content-Type': 'application/json'
                    }
                });

                if (!response.ok) throw new Error('Failed to fetch conversation');

                const data = await response.json();

                if (data && data.id) {
                    this.conversationId = data.id;
                    localStorage.setItem(CONFIG.STORAGE_KEYS.CONVERSATION_ID, this.conversationId);
                }
            } catch (error) {
                console.error("QC SDK: Could not fetch active conversation", error);
            }
        }

        async updateGuestProfile(profileData) {
            if (this.userState !== 'guest' || !this.token) return;
            try {
                await fetch(`${CONFIG.API_BASE_URL}/users/guest/update`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${this.token}`
                    },
                    body: JSON.stringify(profileData)
                });
            } catch (e) { console.error(e); }
        }

        /**
         * WebSocket Connection
         */
        connectWebSocket() {
            if (this.socket && (this.socket.readyState === WebSocket.OPEN || this.socket.readyState === WebSocket.CONNECTING)) return;

            const wsUrl = CONFIG.CHAT_WS_URL;
            const protocols = [this.token];

            try {
                this.socket = new WebSocket(wsUrl, protocols);
            } catch (e) { return; }

            this.socket.onopen = () => {
                if(this.onSystemMessage) this.onSystemMessage(UI_TEXT.CONNECTED);
            };

            this.socket.onclose = (e) => {
                if(this.onSystemMessage) this.onSystemMessage(UI_TEXT.DISCONNECTED);
                setTimeout(() => this.connectWebSocket(), CONFIG.RECONNECT_INTERVAL);
            };

            this.socket.onmessage = (event) => {
                try {
                    const message = JSON.parse(event.data);
                    this.handleIncomingMessage(message);
                } catch (e) { }
            };
        }

        /**
         * Handle Incoming Messages
         */
        handleIncomingMessage(msg) {
            if (msg.type === 'text') {

                // 1. Reconciliation by Client Message ID (Preferred)
                if (msg.payload.client_message_id) {
                    const pendingMsgEl = document.getElementById(msg.payload.client_message_id);
                    if (pendingMsgEl) {
                        // Found pending message -> Confirm it
                        pendingMsgEl.classList.remove('message-pending');
                        pendingMsgEl.classList.add('message-sent-confirmed');
                        pendingMsgEl.id = msg.payload.id; // Update ID
                        return; // STOP HERE (It was our message)
                    }
                }

                // 2. Fallback Check: Is this message sent by ME?
                // If the server echoed the message but we missed the ID check above (e.g. page reload),
                // we still shouldn't show it as a "Received" message.
                if (this.userId && msg.payload.sender_id === this.userId) {
                    // It's my message, but I didn't find a pending element.
                    // Just ignore it to prevent duplicates on the left side.
                    return;
                }

                // 3. New Incoming Message (from Support)
                if (msg.payload && msg.payload.conversation_id && !this.conversationId) {
                    this.conversationId = msg.payload.conversation_id;
                    localStorage.setItem(CONFIG.STORAGE_KEYS.CONVERSATION_ID, this.conversationId);
                }

                if (this.onMessageReceived) this.onMessageReceived(msg.payload);

            } else if (msg.type === 'system') {
                this.handleSystemMessage(msg);
            }
        }

        handleSystemMessage(msg) {
            if (msg.sub_type === 'typing_started') {
                if (this.onTypingStatusChange) this.onTypingStatusChange(true);
                if (this.typingTimeout) clearTimeout(this.typingTimeout);
                this.typingTimeout = setTimeout(() => {
                    if (this.onTypingStatusChange) this.onTypingStatusChange(false);
                }, CONFIG.TYPING_DISPLAY_MS);
            } else if (msg.sub_type === 'typing_stopped') {
                if (this.onTypingStatusChange) this.onTypingStatusChange(false);
                if (this.typingTimeout) clearTimeout(this.typingTimeout);
            }
        }

        sendMessage(content) {
            if (!this.socket || this.socket.readyState !== WebSocket.OPEN) return false;

            const tempId = `temp_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
            const payload = {
                type: 'text',
                content: content,
                conversation_id: this.conversationId || "",
                client_message_id: tempId // Essential for reconciliation
            };

            this.socket.send(JSON.stringify(payload));
            return tempId;
        }

        sendTyping(isTyping) {
            if (!this.socket || this.socket.readyState !== WebSocket.OPEN || !this.conversationId) return;
            const now = Date.now();
            if (isTyping) {
                if (now - this.lastTypingSentTime > CONFIG.TYPING_THROTTLE_MS) {
                    this.socket.send(JSON.stringify({
                        type: 'system', sub_type: 'typing_started', conversation_id: this.conversationId
                    }));
                    this.lastTypingSentTime = now;
                }
            } else {
                this.socket.send(JSON.stringify({
                    type: 'system', sub_type: 'typing_stopped', conversation_id: this.conversationId
                }));
                this.lastTypingSentTime = 0;
            }
        }
    }

    // ==========================================
    // 3. UI RENDERER
    // ==========================================
    const sdk = new QuickConnectSDK();

    // CSS Styles (Telegram Theme - RTL)
    const customCSS = `
    :root { --tg-bg-primary: #212121; --tg-bg-secondary: #0F0F0F; --tg-text-primary: #FFFFFF; --tg-text-secondary: #AAAAAA; --tg-text-system: #A3B6C7; --tg-bubble-sent-bg: #766AC8; --tg-bubble-sent-text: #FFFFFF; --tg-bubble-received-bg: #182533; --tg-bubble-received-text: #FFFFFF; --tg-input-bg: #17212B; --tg-input-text: #D1D5DB; --tg-input-border: #303030; --tg-button-send-bg: #8774E1; --tg-emoji-picker-bg: #1F2D3A; }
    .chat-widget-container, .chat-widget-container * { direction: rtl; text-align: right; font-family: "Tahoma", "Vazir", sans-serif; }
    .message-bubble { max-width: 75%; padding: 8px 12px; border-radius: 15px; margin-bottom: 6px; font-size: 0.95rem; position: relative; word-wrap: break-word;}
    .message-sent { background-color: var(--tg-bubble-sent-bg); color: var(--tg-bubble-sent-text); margin-left: auto; border-bottom-right-radius: 4px; }
    .message-received { background-color: var(--tg-bubble-received-bg); color: var(--tg-bubble-received-text); margin-right: auto; border-bottom-left-radius: 4px; }
    .message-timestamp { font-size: 0.7rem; opacity: 0.7; text-align: left; margin-top: 2px; display: flex; align-items: center; justify-content: flex-end; gap: 4px;}

    /* Status Icons */
    .message-pending { opacity: 0.7; }
    .message-sent-confirmed { opacity: 1; }
    .msg-status-icon { font-size: 0.8em; }
    .message-pending .msg-status-icon::after { content: "🕒"; } /* Clock */
    .message-sent-confirmed .msg-status-icon::after { content: "✓"; } /* Check */

    .qc-chat-window { position: fixed; bottom: 90px; left: 20px; width: 350px; height: 500px; background-color: var(--tg-bg-primary); border: 1px solid var(--tg-input-border); border-radius: 12px; display: flex; flex-direction: column; box-shadow: 0 10px 25px rgba(0,0,0,0.3); z-index: 9999; transition: all 0.3s ease; opacity: 0; transform: translateY(20px); visibility: hidden; }
    .qc-chat-window.open { opacity: 1; transform: translateY(0); visibility: visible; }
    .qc-header { background-color: var(--tg-bg-secondary); padding: 12px; border-bottom: 1px solid #333; display: flex; justify-content: space-between; color: white; border-radius: 12px 12px 0 0; }
    .qc-messages { flex: 1; overflow-y: auto; padding: 10px; display: flex; flex-direction: column; }
    .qc-messages::-webkit-scrollbar { width: 5px; }
    .qc-messages::-webkit-scrollbar-thumb { background: #444; border-radius: 3px; }
    .qc-footer { padding: 10px; background-color: var(--tg-bg-secondary); display: flex; align-items: flex-end; gap: 8px; border-radius: 0 0 12px 12px; }
    .qc-input { flex: 1; background: var(--tg-input-bg); color: white; border: none; padding: 10px; border-radius: 20px; resize: none; max-height: 100px; min-height: 40px; outline: none; }
    .qc-btn { background: var(--tg-button-send-bg); border: none; width: 40px; height: 40px; border-radius: 50%; cursor: pointer; color: white; display: flex; align-items: center; justify-content: center; }
    .qc-fab { position: fixed; bottom: 20px; left: 20px; width: 60px; height: 60px; background: var(--tg-button-send-bg); border-radius: 50%; cursor: pointer; box-shadow: 0 4px 10px rgba(0,0,0,0.3); z-index: 10000; display: flex; align-items: center; justify-content: center; transition: transform 0.2s; }
    .qc-fab:hover { transform: scale(1.05); }
    .qc-fab svg { width: 30px; height: 30px; fill: white; }
    .qc-typing { font-size: 0.75rem; color: var(--tg-text-system); padding: 0 10px 5px 10px; display: none; }
    .qc-typing.active { display: block; }
    @media (max-width: 480px) { .qc-chat-window { width: calc(100% - 40px); height: 80%; bottom: 90px; right: 20px; left: 20px; } }
    `;

    let windowEl, messagesEl, inputEl, fabEl, typingEl;

    function injectStyles() { const style = document.createElement('style'); style.textContent = customCSS; document.head.appendChild(style); }
    function createDOM() {
        fabEl = document.createElement('div'); fabEl.className = 'qc-fab'; fabEl.innerHTML = `<svg viewBox="0 0 24 24"><path d="M20 2H4c-1.1 0-2 .9-2 2v18l4-4h14c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2z"/></svg>`; document.body.appendChild(fabEl);
        windowEl = document.createElement('div'); windowEl.className = 'qc-chat-window';
        windowEl.innerHTML = `<div class="qc-header"><span>${UI_TEXT.TITLE}</span><span style="cursor:pointer" id="qc-close">✕</span></div><div class="qc-messages" id="qc-messages-list"><div class="text-center text-xs py-2 system-message" style="color: #666; text-align: center; font-size: 0.8rem;">${UI_TEXT.WELCOME}</div></div><div class="qc-typing" id="qc-typing">${UI_TEXT.TYPING}</div><div class="qc-footer"><textarea class="qc-input" rows="1" placeholder="${UI_TEXT.INPUT_PLACEHOLDER}"></textarea><button class="qc-btn">➤</button></div>`;
        document.body.appendChild(windowEl);
        messagesEl = windowEl.querySelector('#qc-messages-list'); inputEl = windowEl.querySelector('.qc-input'); typingEl = windowEl.querySelector('#qc-typing');
        fabEl.addEventListener('click', toggleChat); windowEl.querySelector('#qc-close').addEventListener('click', toggleChat);
        windowEl.querySelector('.qc-btn').addEventListener('click', handleSend);
        inputEl.addEventListener('keydown', (e) => { if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); handleSend(); } handleTypingInput(); });
    }

    function toggleChat() {
        windowEl.classList.toggle('open');
        if (windowEl.classList.contains('open')) {
            inputEl.focus();
            sdk.init();
            scrollToBottom();
        }
    }

    function handleSend() {
        const text = inputEl.value.trim();
        if (!text) return;
        const tempId = sdk.sendMessage(text);
        if (tempId) { addMessageToUI(text, 'sent', tempId); inputEl.value = ''; sdk.sendTyping(false); } else { alert(UI_TEXT.DISCONNECTED); }
    }

    let typingDebounce = null;
    function handleTypingInput() {
        sdk.sendTyping(true);
        if (typingDebounce) clearTimeout(typingDebounce);
        typingDebounce = setTimeout(() => { sdk.sendTyping(false); }, 1000);
    }

    function addMessageToUI(text, type, messageId = null) {
        const msgDiv = document.createElement('div');
        if (messageId) { msgDiv.id = messageId; }

        let className = `message-bubble ${type === 'sent' ? 'message-sent' : 'message-received'}`;
        if (type === 'sent') { className += ' message-pending'; }
        msgDiv.className = className;

        const contentDiv = document.createElement('div'); contentDiv.className = 'message-content'; contentDiv.textContent = text;
        const timeDiv = document.createElement('div'); timeDiv.className = 'message-timestamp'; const now = new Date(); timeDiv.textContent = `${now.getHours()}:${now.getMinutes().toString().padStart(2, '0')}`;

        if (type === 'sent') {
            const statusIcon = document.createElement('span'); statusIcon.className = 'msg-status-icon'; timeDiv.appendChild(statusIcon);
        }

        msgDiv.appendChild(contentDiv); msgDiv.appendChild(timeDiv); messagesEl.appendChild(msgDiv); scrollToBottom();
        return msgDiv;
    }

    function scrollToBottom() { messagesEl.scrollTop = messagesEl.scrollHeight; }

    sdk.onMessageReceived = (payload) => {
        // Check for duplicates via ID
        if(payload.id && document.getElementById(payload.id)) return;
        addMessageToUI(payload.content, 'received', payload.id);
    };

    sdk.onSystemMessage = (text) => { const sysDiv = document.createElement('div'); sysDiv.style.textAlign = 'center'; sysDiv.style.fontSize = '0.8rem'; sysDiv.style.color = '#888'; sysDiv.style.padding = '5px'; sysDiv.textContent = text; messagesEl.appendChild(sysDiv); scrollToBottom(); };
    sdk.onTypingStatusChange = (isTyping) => { if (isTyping) typingEl.classList.add('active'); else typingEl.classList.remove('active'); };

    function loadTailwind() { if (!document.querySelector('script[src*="tailwindcss"]')) { const script = document.createElement('script'); script.src = CONFIG.TAILWIND_CDN; document.head.appendChild(script); } }

    if (document.readyState === 'loading') { document.addEventListener('DOMContentLoaded', () => { loadTailwind(); injectStyles(); createDOM(); }); } else { loadTailwind(); injectStyles(); createDOM(); }

    window.QuickConnect = { init: (token) => sdk.init(token), open: toggleChat, updateProfile: (data) => sdk.updateGuestProfile(data) };

})(window, document);