/**
 * Quick Connect Chat SDK
 * Version: 1.6.0 (Added Chat History Feature)
 */

(function(window, document) {
    'use strict';

    const QuickConnect = {
        config: {
            managerUrl: 'http://localhost:2531',
            chatUrl: 'ws://localhost:2530/chats/clients',
            chatApiUrl: 'http://localhost:2530',
            position: 'bottom-right',
            theme: 'purple',
            lang: 'fa'
        },

        state: {
            isOpen: false,
            isConnected: false,
            token: null,
            userId: null,
            userState: null,
            conversationId: null,
            messages: [],
            typingUsers: new Set(),
            ws: null,
            lastTypingSent: 0,
            typingTimeouts: {},
            reconnectAttempts: 0,
            reconnectInterval: 3000,
            unreadCount: 0,
            // History State
            historyLoaded: false,
            nextCursor: null,
            hasMore: false,
            isLoadingHistory: false
        },

        // Initialize SDK
        init: function(options) {
            Object.assign(this.config, options);

            if (document.readyState === 'loading') {
                document.addEventListener('DOMContentLoaded', () => this.setup());
            } else {
                this.setup();
            }
        },

        // Setup SDK
        setup: function() {
            if (document.getElementById('qc-widget')) return;

            this.injectStyles();
            this.injectHTML();
            this.attachEventListeners();
            this.initializeAuth();
        },

        // Helper: Parse JWT
        parseJwt: function(token) {
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
        },

        // Inject CSS Styles
        injectStyles: function() {
            const style = document.createElement('style');
            style.textContent = `
        #qc-widget * { box-sizing: border-box; margin: 0; padding: 0; }

        #qc-widget {
          position: fixed;
          ${this.config.position.includes('right') ? 'right' : 'left'}: 24px;
          bottom: 24px;
          z-index: 999999;
          font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
          direction: rtl;
          pointer-events: none;
        }

        #qc-widget > * {
          pointer-events: auto;
        }

        #qc-btn {
          width: 64px;
          height: 64px;
          border-radius: 50%;
          border: none;
          background: linear-gradient(135deg, #667eea 0%, #764ba2 50%, #f093fb 100%);
          box-shadow: 0 8px 24px rgba(102, 126, 234, 0.4);
          cursor: pointer;
          position: relative;
          transition: all 0.3s ease;
          animation: qc-pulse 2s infinite;
          display: flex;
          align-items: center;
          justify-content: center;
          z-index: 1000000;
        }

        #qc-btn:hover {
          transform: scale(1.1);
          box-shadow: 0 12px 32px rgba(102, 126, 234, 0.6);
        }

        #qc-btn.hidden {
          display: none !important;
          opacity: 0;
          pointer-events: none;
        }

        #qc-badge {
          position: absolute;
          top: -4px;
          right: -4px;
          background: #ef4444;
          color: white;
          width: 24px;
          height: 24px;
          border-radius: 50%;
          font-size: 12px;
          font-weight: bold;
          display: flex; /* Managed by JS logic now */
          align-items: center;
          justify-content: center;
          animation: qc-bounce 1s infinite;
          box-shadow: 0 2px 5px rgba(0,0,0,0.2);
        }

        @keyframes qc-pulse {
          0%, 100% { box-shadow: 0 8px 24px rgba(102, 126, 234, 0.4); }
          50% { box-shadow: 0 8px 40px rgba(102, 126, 234, 0.6); }
        }

        @keyframes qc-bounce {
          0%, 100% { transform: translateY(0); }
          50% { transform: translateY(-4px); }
        }

        @keyframes qc-slideIn {
          from { transform: translateY(20px); opacity: 0; }
          to { transform: translateY(0); opacity: 1; }
        }

        #qc-window {
          display: none;
          width: 380px;
          height: 600px;
          background: linear-gradient(135deg, #1e1b4b 0%, #581c87 50%, #1e1b4b 100%);
          border-radius: 20px;
          box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
          overflow: hidden;
          flex-direction: column;
          border: 1px solid rgba(168, 85, 247, 0.3);
          animation: qc-slideIn 0.3s ease;
          position: absolute;
          bottom: 80px;
          ${this.config.position.includes('right') ? 'right' : 'left'}: 0;
        }

        #qc-window.open { display: flex; }

        #qc-header {
          background: linear-gradient(135deg, #7c3aed 0%, #2563eb 100%);
          padding: 18px;
          display: flex;
          align-items: center;
          justify-content: space-between;
          flex-shrink: 0;
        }

        #qc-header-info {
          display: flex;
          align-items: center;
          gap: 12px;
        }

        #qc-avatar {
          width: 44px;
          height: 44px;
          background: white;
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          position: relative;
          box-shadow: 0 4px 12px rgba(0,0,0,0.1);
        }

        #qc-status-dot {
          position: absolute;
          bottom: 0;
          right: 0;
          width: 13px;
          height: 13px;
          background: #10b981;
          border: 2px solid white;
          border-radius: 50%;
        }

        #qc-status-dot.offline { background: #6b7280; }

        #qc-header-text h3 {
          color: white;
          font-size: 18px;
          font-weight: 800;
          margin-bottom: 2px;
        }

        #qc-header-text p {
          color: rgba(255, 255, 255, 0.85);
          font-size: 13px;
        }

        #qc-close-btn {
          background: rgba(255, 255, 255, 0.2);
          border: none;
          color: white;
          width: 36px;
          height: 36px;
          border-radius: 10px;
          cursor: pointer;
          transition: background 0.2s;
          display: flex;
          align-items: center;
          justify-content: center;
        }

        #qc-close-btn:hover { background: rgba(255, 255, 255, 0.3); }

        #qc-profile-banner {
          display: none;
          background: linear-gradient(135deg, #f59e0b 0%, #ea580c 100%);
          padding: 14px 18px;
          flex-shrink: 0;
        }

        #qc-profile-banner.show { display: block; }

        #qc-profile-header {
          display: flex;
          align-items: center;
          justify-content: space-between;
          margin-bottom: 10px;
        }

        #qc-profile-header p {
          color: white;
          font-size: 14px;
          font-weight: 600;
        }

        #qc-profile-close {
          background: transparent;
          border: none;
          color: rgba(255, 255, 255, 0.8);
          cursor: pointer;
          padding: 0;
          width: 20px;
          height: 20px;
        }

        #qc-profile-form {
          display: flex;
          flex-direction: column;
          gap: 10px;
        }

        #qc-profile-form input {
          width: 100%;
          padding: 12px 14px;
          border-radius: 10px;
          border: none;
          background: rgba(255, 255, 255, 0.9);
          font-size: 14px;
        }

        #qc-profile-actions {
          display: flex;
          gap: 10px;
        }

        #qc-profile-submit {
          flex: 1;
          padding: 12px;
          background: white;
          color: #ea580c;
          border: none;
          border-radius: 10px;
          font-weight: 700;
          font-size: 14px;
          cursor: pointer;
          transition: background 0.2s;
        }

        #qc-profile-submit:hover { background: rgba(255, 255, 255, 0.9); }

        #qc-messages {
          flex: 1;
          overflow-y: auto;
          padding: 20px;
          background: rgba(15, 23, 42, 0.5);
          display: flex;
          flex-direction: column;
          gap: 16px;
          width: 100%;
        }

        #qc-messages::-webkit-scrollbar { width: 6px; }
        #qc-messages::-webkit-scrollbar-track { background: transparent; }
        #qc-messages::-webkit-scrollbar-thumb { background: rgba(148, 163, 184, 0.3); border-radius: 3px; }

        .qc-msg {
          display: flex;
          width: 100%;
          animation: qc-slideIn 0.2s ease;
        }

        .qc-msg.own { justify-content: flex-end; }
        .qc-msg.other { justify-content: flex-start; }
        .qc-msg.system { justify-content: center; }

        .qc-msg-bubble {
          max-width: 75%;
          padding: 14px 18px;
          border-radius: 18px;
          box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
          position: relative;
          min-width: 60px;
        }

        .qc-msg.own .qc-msg-bubble {
          background: linear-gradient(135deg, #7c3aed 0%, #2563eb 100%);
          color: white;
          border-bottom-right-radius: 4px;
        }

        .qc-msg.other .qc-msg-bubble {
          background: #1e293b;
          color: white;
          border-bottom-left-radius: 4px;
          border: 1px solid #334155;
        }

        .qc-msg.system .qc-msg-bubble {
          background: rgba(51, 65, 85, 0.6);
          color: #cbd5e1;
          font-size: 12px;
          padding: 8px 16px;
          border-radius: 20px;
          border: 1px solid rgba(255,255,255,0.1);
          max-width: 90%;
          text-align: center;
        }

        .qc-msg-content {
          font-size: 15px;
          line-height: 1.6;
          word-break: break-word;
          overflow-wrap: anywhere;
          white-space: pre-wrap;
        }

        .qc-msg-footer {
          display: flex;
          align-items: center;
          justify-content: flex-end;
          gap: 6px;
          margin-top: 8px;
        }

        .qc-msg-time {
          font-size: 11px;
          opacity: 0.7;
          font-weight: 300;
        }

        .qc-msg-status {
          width: 14px;
          height: 14px;
          opacity: 0.8;
        }

        .qc-typing {
          display: flex;
          gap: 5px;
          padding: 14px 18px;
          background: #1e293b;
          border-radius: 18px;
          border-bottom-left-radius: 4px;
          width: fit-content;
          border: 1px solid #334155;
        }

        .qc-typing-dot {
          width: 8px;
          height: 8px;
          background: #94a3b8;
          border-radius: 50%;
          animation: qc-typing 1.4s infinite;
        }

        .qc-typing-dot:nth-child(2) { animation-delay: 0.2s; }
        .qc-typing-dot:nth-child(3) { animation-delay: 0.4s; }

        #qc-emoji-picker {
          display: none;
          background: #1e293b;
          border-top: 1px solid #334155;
          padding: 12px;
          flex-shrink: 0;
        }

        #qc-emoji-picker.show { display: block; }

        #qc-emoji-grid {
          display: grid;
          grid-template-columns: repeat(8, 1fr);
          gap: 8px;
        }

        .qc-emoji-btn {
          background: transparent;
          border: none;
          font-size: 24px;
          cursor: pointer;
          padding: 6px;
          border-radius: 8px;
          transition: background 0.2s;
        }

        .qc-emoji-btn:hover { background: #334155; }

        #qc-input-area {
          background: #0f172a;
          border-top: 1px solid #334155;
          padding: 18px;
          flex-shrink: 0;
        }

        #qc-input-container {
          display: flex;
          align-items: flex-end;
          gap: 10px;
        }

        #qc-emoji-toggle, #qc-send-btn {
          background: rgba(71, 85, 105, 0.4);
          border: none;
          color: #94a3b8;
          width: 44px;
          height: 44px;
          border-radius: 14px;
          cursor: pointer;
          transition: all 0.2s;
          display: flex;
          align-items: center;
          justify-content: center;
          flex-shrink: 0;
        }

        #qc-emoji-toggle:hover { background: #475569; color: #a78bfa; }

        #qc-send-btn {
          background: linear-gradient(135deg, #7c3aed 0%, #2563eb 100%);
          color: white;
        }

        #qc-send-btn:hover {
          transform: scale(1.05);
          box-shadow: 0 4px 12px rgba(124, 58, 237, 0.4);
        }

        #qc-send-btn:disabled {
          opacity: 0.5;
          cursor: not-allowed;
          transform: scale(1);
        }

        #qc-input-wrapper {
          flex: 1;
          background: #1e293b;
          border-radius: 14px;
          padding: 12px 16px;
          border: 1px solid #334155;
          transition: border-color 0.2s;
        }

        #qc-input-wrapper:focus-within { border-color: #7c3aed; }

        #qc-input {
          width: 100%;
          background: transparent;
          border: none;
          color: white;
          font-size: 15px;
          outline: none;
          resize: none;
          max-height: 100px;
          font-family: inherit;
          line-height: 1.5;
        }

        #qc-input::placeholder { color: #64748b; }

        /* Mobile Responsive Styles */
        @media screen and (max-width: 600px) {
            #qc-widget {
                bottom: 0 !important;
                left: 0 !important;
                right: 0 !important;
                top: 0 !important;
                width: 100%;
                height: 100%;
            }

            #qc-window {
                width: 100%;
                height: 100%;
                border-radius: 0;
                border: none;
                bottom: 0;
                right: 0;
                left: 0;
            }

            #qc-messages {
                padding-bottom: 20px;
            }

            #qc-btn {
                position: absolute;
                bottom: 20px;
                right: 20px;
                left: auto;
            }
        }
      `;
            document.head.appendChild(style);
        },

        // Inject HTML Structure
        injectHTML: function() {
            const container = document.createElement('div');
            container.id = 'qc-widget';
            // Badge is hidden by default using style="display: none"
            container.innerHTML = `
        <button id="qc-btn" aria-label="Open chat">
          <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
          </svg>
          <span id="qc-badge" style="display: none">0</span>
        </button>

        <div id="qc-window">
          <div id="qc-header">
            <div id="qc-header-info">
              <div id="qc-avatar">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#7c3aed" stroke-width="2">
                  <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
                </svg>
                <span id="qc-status-dot"></span>
              </div>
              <div id="qc-header-text">
                <h3>پشتیبانی آنلاین</h3>
                <p id="qc-status-text">در حال اتصال...</p>
              </div>
            </div>
            <button id="qc-close-btn" aria-label="Close chat">
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"></line>
                <line x1="6" y1="6" x2="18" y2="18"></line>
              </svg>
            </button>
          </div>

          <div id="qc-profile-banner">
            <div id="qc-profile-header">
              <p>اطلاعات خود را تکمیل کنید 📝</p>
              <button id="qc-profile-close">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="18" y1="6" x2="6" y2="18"></line>
                  <line x1="6" y1="6" x2="18" y2="18"></line>
                </svg>
              </button>
            </div>
            <div id="qc-profile-form">
              <input type="text" id="qc-profile-name" placeholder="نام و نام خانوادگی">
              <div id="qc-profile-actions">
                <input type="email" id="qc-profile-email" placeholder="ایمیل">
                <button id="qc-profile-submit">ثبت</button>
              </div>
            </div>
          </div>

          <div id="qc-messages"></div>

          <div id="qc-emoji-picker">
            <div id="qc-emoji-grid"></div>
          </div>

          <div id="qc-input-area">
            <div id="qc-input-container">
              <button id="qc-emoji-toggle" aria-label="Toggle emoji picker">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10"></circle>
                  <path d="M8 14s1.5 2 4 2 4-2 4-2"></path>
                  <line x1="9" y1="9" x2="9.01" y2="9"></line>
                  <line x1="15" y1="9" x2="15.01" y2="9"></line>
                </svg>
              </button>
              <div id="qc-input-wrapper">
                <textarea id="qc-input" placeholder="پیام خود را بنویسید..." rows="1"></textarea>
              </div>
              <button id="qc-send-btn" aria-label="Send message">
                <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="22" y1="2" x2="11" y2="13"></line>
                  <polygon points="22 2 15 22 11 13 2 9 22 2"></polygon>
                </svg>
              </button>
            </div>
          </div>
        </div>
      `;
            document.body.appendChild(container);

            // Add emojis
            const emojis = ['😊', '😂', '❤️', '👍', '🙏', '🎉', '🔥', '✨', '💯', '🚀', '💪', '😍', '🤔', '👌', '🎯', '⭐'];
            const emojiGrid = document.getElementById('qc-emoji-grid');
            emojis.forEach(emoji => {
                const btn = document.createElement('button');
                btn.className = 'qc-emoji-btn';
                btn.textContent = emoji;
                btn.onclick = () => this.insertEmoji(emoji);
                emojiGrid.appendChild(btn);
            });
        },

        // Attach Event Listeners
        attachEventListeners: function() {
            const btn = document.getElementById('qc-btn');
            const closeBtn = document.getElementById('qc-close-btn');

            if(btn) btn.onclick = () => this.openChat();
            if(closeBtn) closeBtn.onclick = () => this.closeChat();

            document.getElementById('qc-emoji-toggle').onclick = () => this.toggleEmojiPicker();
            document.getElementById('qc-send-btn').onclick = () => this.sendMessage();
            document.getElementById('qc-profile-close').onclick = () => this.hideProfileBanner();
            document.getElementById('qc-profile-submit').onclick = () => this.updateProfile();

            const input = document.getElementById('qc-input');
            input.oninput = () => this.handleTyping();
            input.onkeydown = (e) => {
                if (e.key === 'Enter' && !e.shiftKey) {
                    e.preventDefault();
                    this.sendMessage();
                }
            };
        },

        // Authentication Logic
        initializeAuth: async function() {
            const token = localStorage.getItem('QC_TOKEN');
            const userState = localStorage.getItem('QC_USER_STATE');

            if (token && userState) {
                this.state.token = token;
                this.state.userState = userState;
            } else {
                await this.registerGuest();
            }

            // Extract user ID from token
            if (this.state.token) {
                const decoded = this.parseJwt(this.state.token);
                if (decoded && decoded.user_id) {
                    this.state.userId = decoded.user_id;
                }
                // Connect flow
                await this.fetchActiveConversation();
                this.connectWebSocket();
            }
        },

        registerGuest: async function() {
            try {
                const response = await fetch(`${this.config.managerUrl}/users/guest/register`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        fullname: "Guest Visitor"
                    })
                });

                if (!response.ok) throw new Error('Guest registration failed');

                const data = await response.json();
                if (data.qc_token) {
                    localStorage.setItem('QC_TOKEN', data.qc_token);
                    localStorage.setItem('QC_USER_STATE', 'guest');
                    this.state.token = data.qc_token;
                    this.state.userState = 'guest';
                    this.showProfileBanner();
                }
            } catch (error) {
                console.error("QC SDK: Auth Error", error);
            }
        },

        fetchActiveConversation: async function() {
            if (!this.state.token) return;

            try {
                const response = await fetch(`${this.config.chatApiUrl}/conversations/active`, {
                    method: 'GET',
                    headers: {
                        'Authorization': `Bearer ${this.state.token}`,
                        'Content-Type': 'application/json'
                    }
                });

                if (!response.ok) throw new Error('Failed to fetch conversation');

                const data = await response.json();
                if (data && data.id) {
                    this.state.conversationId = data.id;
                    localStorage.setItem('QC_CONVERSATION_ID', data.id);

                    // If chat window is already open but history not loaded (rare case), load it now
                    if (this.state.isOpen && !this.state.historyLoaded) {
                        this.fetchChatHistory();
                    }
                }
            } catch (error) {
                console.error("QC SDK: Active conversation fetch error", error);
            }
        },

        // NEW: Fetch Chat History
        fetchChatHistory: async function() {
            if (!this.state.token || !this.state.conversationId || this.state.historyLoaded || this.state.isLoadingHistory) return;

            this.state.isLoadingHistory = true;

            try {
                const payload = {
                    conversation_id: this.state.conversationId,
                    pagination: {
                        cursor: this.state.nextCursor, // null for first page
                        limit: 20
                    }
                };

                const response = await fetch(`${this.config.chatApiUrl}/chats`, {
                    method: 'POST',
                    headers: {
                        'Authorization': `Bearer ${this.state.token}`,
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(payload)
                });

                if (!response.ok) throw new Error('Failed to fetch history');

                const data = await response.json();

                if (data && data.results) {
                    // API returns messages Newest to Oldest. We reverse them for UI (Oldest to Newest).
                    const historyMessages = data.results.map(msg => ({
                        id: msg.id,
                        content: msg.content,
                        created_at: msg.created_at,
                        // Check if sender is current user
                        isOwn: msg.sender_id === this.state.userId,
                        status: 'sent' // Historical messages are always sent
                    })).reverse();

                    // Prepend history to existing messages (if any)
                    this.state.messages = [...historyMessages, ...this.state.messages];

                    // Update pagination state
                    if (data.paginate) {
                        this.state.nextCursor = data.paginate.next_cursor;
                        this.state.hasMore = data.paginate.has_more;
                    }

                    this.state.historyLoaded = true;
                    this.renderMessages();
                }

            } catch (error) {
                console.error("QC SDK: History Fetch Error", error);
            } finally {
                this.state.isLoadingHistory = false;
            }
        },

        updateProfile: async function() {
            const name = document.getElementById('qc-profile-name').value;
            const email = document.getElementById('qc-profile-email').value;

            if (!name || !email) {
                alert('لطفا تمام فیلدها را پر کنید');
                return;
            }

            if (this.state.userState !== 'guest' || !this.state.token) return;

            try {
                await fetch(`${this.config.managerUrl}/users/guest/update`, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${this.state.token}`
                    },
                    body: JSON.stringify({
                        fullname: name,
                        email: email
                    })
                });
                this.hideProfileBanner();
                this.addSystemMessage('پروفایل شما با موفقیت ثبت شد ✅');
            } catch (e) {
                console.error("Profile Update Error", e);
            }
        },

        // WebSocket Logic
        connectWebSocket: function() {
            if (this.state.ws && (this.state.ws.readyState === WebSocket.OPEN || this.state.ws.readyState === WebSocket.CONNECTING)) return;

            const protocols = [this.state.token];
            try {
                this.state.ws = new WebSocket(this.config.chatUrl, protocols);
            } catch (e) {
                console.error("WebSocket creation failed", e);
                return;
            }

            this.state.ws.onopen = () => {
                this.state.isConnected = true;
                this.updateConnectionStatus();
                if(this.state.messages.length === 0 && this.state.historyLoaded) {
                    this.addSystemMessage('به چت پشتیبانی خوش آمدید! 👋');
                } else if (!this.state.historyLoaded) {
                    // Wait for history to potentially load message, if not loaded yet
                    // Just a welcome log
                }
            };

            this.state.ws.onclose = () => {
                this.state.isConnected = false;
                this.updateConnectionStatus();
                setTimeout(() => this.connectWebSocket(), this.state.reconnectInterval);
            };

            this.state.ws.onerror = (err) => {
                console.error("QC WS Error", err);
            };

            this.state.ws.onmessage = (event) => {
                try {
                    const message = JSON.parse(event.data);
                    this.handleIncomingMessage(message);
                } catch (e) {
                    console.error("Parse Error", e);
                }
            };
        },

        updateConnectionStatus: function() {
            const statusText = document.getElementById('qc-status-text');
            const statusDot = document.getElementById('qc-status-dot');

            if(!statusText || !statusDot) return;

            // Check connectivity first
            if (!this.state.isConnected) {
                statusText.textContent = 'آفلاین';
                statusDot.classList.add('offline');
                return;
            }

            // If connected, check typing status (Prioritize Typing over "Online")
            statusDot.classList.remove('offline');

            if (this.state.typingUsers.size > 0) {
                statusText.textContent = 'در حال نوشتن...';
            } else {
                statusText.textContent = 'آنلاین';
            }
        },

        updateBadgeUI: function() {
            const badge = document.getElementById('qc-badge');
            if (!badge) return;

            if (this.state.unreadCount > 0) {
                badge.style.display = 'flex';
                badge.textContent = this.state.unreadCount;
            } else {
                badge.style.display = 'none';
            }
        },

        // Messaging Logic
        sendMessage: function() {
            const input = document.getElementById('qc-input');
            const content = input.value.trim();

            if (!content || !this.state.isConnected || !this.state.ws) return;

            const tempId = `local_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;

            const payload = {
                type: 'text',
                content: content,
                conversation_id: this.state.conversationId || "",
                client_message_id: tempId
            };

            this.state.ws.send(JSON.stringify(payload));

            // Optimistic UI Add
            const message = {
                id: tempId,
                content: content,
                created_at: new Date().toISOString(),
                isOwn: true,
                status: 'pending'
            };

            this.state.messages.push(message);
            this.renderMessages();

            input.value = '';
            input.style.height = 'auto';
            this.sendTypingStopped();
        },

        handleIncomingMessage: function(msg) {
            if (msg.type === 'text') {
                const incomingClientId = msg.client_message_id;

                if (incomingClientId) {
                    const pendingIndex = this.state.messages.findIndex(m => m.id === incomingClientId);

                    if (pendingIndex !== -1) {
                        const msgRef = this.state.messages[pendingIndex];
                        msgRef.status = 'sent';
                        if(msg.payload && msg.payload.id) {
                            msgRef.id = msg.payload.id;
                        }
                        this.renderMessages();
                        return;
                    }
                }

                if (this.state.userId && msg.payload.sender_id === this.state.userId) {
                    return;
                }

                // Check if chat is closed to increment badge
                if (!this.state.isOpen) {
                    this.state.unreadCount++;
                    this.updateBadgeUI();
                }

                if (msg.payload && msg.payload.conversation_id && !this.state.conversationId) {
                    this.state.conversationId = msg.payload.conversation_id;
                    localStorage.setItem('QC_CONVERSATION_ID', this.state.conversationId);
                }

                this.receiveMessage({
                    type: 'text',
                    payload: {
                        ...msg.payload,
                        isOwn: false
                    }
                });

            } else if (msg.type === 'system') {
                this.handleSystemMessage(msg);
            }
        },

        receiveMessage: function(data) {
            if (data.type === 'text') {
                this.state.messages.push({
                    ...data.payload,
                    isOwn: false
                });
                this.renderMessages();
            }
        },

        addSystemMessage: function(content) {
            this.state.messages.push({
                id: 'system_' + Date.now(),
                content: content,
                created_at: new Date().toISOString(),
                isSystem: true
            });
            this.renderMessages();
        },

        handleSystemMessage: function(data) {
            if (!data.payload || !data.payload.sender_id) return;

            const senderId = data.payload.sender_id;

            if (data.sub_type === 'typing_started') {
                this.state.typingUsers.add(senderId);

                // Update Header immediately
                this.updateConnectionStatus();
                this.renderMessages();

                if (this.state.typingTimeouts[senderId]) {
                    clearTimeout(this.state.typingTimeouts[senderId]);
                }

                // Auto-stop after 6 seconds to handle packet loss
                this.state.typingTimeouts[senderId] = setTimeout(() => {
                    this.state.typingUsers.delete(senderId);
                    this.updateConnectionStatus(); // Revert header
                    this.renderMessages();
                }, 6000);

            } else if (data.sub_type === 'typing_stopped') {
                this.state.typingUsers.delete(senderId);

                if (this.state.typingTimeouts[senderId]) {
                    clearTimeout(this.state.typingTimeouts[senderId]);
                }

                // Update Header immediately
                this.updateConnectionStatus();
                this.renderMessages();
            }
        },

        handleTyping: function() {
            const now = Date.now();
            if (now - this.state.lastTypingSent > 5000) {
                this.sendTypingStarted();
                this.state.lastTypingSent = now;
            }

            const input = document.getElementById('qc-input');
            input.style.height = 'auto';
            input.style.height = input.scrollHeight + 'px';
        },

        sendTypingStarted: function() {
            if (!this.state.ws || !this.state.conversationId) return;
            this.state.ws.send(JSON.stringify({
                type: 'system',
                sub_type: 'typing_started',
                conversation_id: this.state.conversationId
            }));
        },

        sendTypingStopped: function() {
            if (!this.state.ws || !this.state.conversationId) return;
            this.state.ws.send(JSON.stringify({
                type: 'system',
                sub_type: 'typing_stopped',
                conversation_id: this.state.conversationId
            }));
            this.state.lastTypingSent = 0;
        },

        // Render Logic
        renderMessages: function() {
            const container = document.getElementById('qc-messages');
            if (!container) return;

            container.innerHTML = '';

            this.state.messages.forEach(msg => {
                const msgEl = document.createElement('div');

                if (msg.isSystem) {
                    msgEl.className = 'qc-msg system';
                    msgEl.innerHTML = `
            <div class="qc-msg-bubble">
              <div class="qc-msg-content">${this.escapeHtml(msg.content)}</div>
            </div>
          `;
                } else {
                    msgEl.className = `qc-msg ${msg.isOwn ? 'own' : 'other'}`;

                    let statusIcon = '';
                    if (msg.isOwn) {
                        if (msg.status === 'pending') {
                            statusIcon = `
                            <svg class="qc-msg-status" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <circle cx="12" cy="12" r="10"></circle>
                                <polyline points="12 6 12 12 16 14"></polyline>
                            </svg>`;
                        } else {
                            statusIcon = `
                            <svg class="qc-msg-status" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <polyline points="20 6 9 17 4 12"></polyline>
                            </svg>`;
                        }
                    }

                    msgEl.innerHTML = `
            <div class="qc-msg-bubble">
              <div class="qc-msg-content">${this.escapeHtml(msg.content)}</div>
              <div class="qc-msg-footer">
                <span class="qc-msg-time">${this.formatTime(msg.created_at)}</span>
                ${statusIcon}
              </div>
            </div>
          `;
                }

                container.appendChild(msgEl);
            });

            if (this.state.typingUsers.size > 0) {
                const typingEl = document.createElement('div');
                typingEl.className = 'qc-msg other';
                typingEl.innerHTML = `
          <div class="qc-typing">
            <div class="qc-typing-dot"></div>
            <div class="qc-typing-dot"></div>
            <div class="qc-typing-dot"></div>
          </div>
        `;
                container.appendChild(typingEl);
            }

            container.scrollTop = container.scrollHeight;
        },

        // Utilities
        formatTime: function(timestamp) {
            const date = new Date(timestamp);
            return date.toLocaleTimeString('fa-IR', { hour: '2-digit', minute: '2-digit' });
        },

        escapeHtml: function(text) {
            const div = document.createElement('div');
            div.textContent = text;
            return div.innerHTML;
        },

        // UI Control
        openChat: function() {
            this.state.isOpen = true;
            this.state.unreadCount = 0; // Reset count
            this.updateBadgeUI(); // Update UI

            // Try to fetch history if not loaded yet and conversation ID is available
            if (!this.state.historyLoaded && this.state.conversationId) {
                this.fetchChatHistory();
            }

            const btn = document.getElementById('qc-btn');
            const win = document.getElementById('qc-window');
            const input = document.getElementById('qc-input');

            if (btn) btn.classList.add('hidden');
            if (win) win.classList.add('open');

            if (input) {
                setTimeout(() => input.focus(), 100);
            }
        },

        closeChat: function() {
            this.state.isOpen = false;
            const btn = document.getElementById('qc-btn');
            const win = document.getElementById('qc-window');

            if (btn) btn.classList.remove('hidden');
            if (win) win.classList.remove('open');
        },

        toggleEmojiPicker: function() {
            const picker = document.getElementById('qc-emoji-picker');
            picker.classList.toggle('show');
        },

        insertEmoji: function(emoji) {
            const input = document.getElementById('qc-input');
            input.value += emoji;
            input.focus();
            document.getElementById('qc-emoji-picker').classList.remove('show');
        },

        showProfileBanner: function() {
            const banner = document.getElementById('qc-profile-banner');
            if(banner) banner.classList.add('show');
        },

        hideProfileBanner: function() {
            const banner = document.getElementById('qc-profile-banner');
            if(banner) banner.classList.remove('show');
        },

        // Public API
        login: function(token) {
            localStorage.setItem('QC_TOKEN', token);
            localStorage.setItem('QC_USER_STATE', 'client');
            this.state.token = token;
            this.state.userState = 'client';
            if (this.state.ws) this.state.ws.close();
            this.initializeAuth();
        },

        logout: function() {
            localStorage.removeItem('QC_TOKEN');
            localStorage.removeItem('QC_USER_STATE');
            this.state.token = null;
            this.state.userState = null;
            if (this.state.ws) {
                this.state.ws.close();
            }
        },

        open: function() {
            this.openChat();
        },

        close: function() {
            this.closeChat();
        }
    };

    // Auto-initialize
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', function() {
            QuickConnect.init();
        });
    } else {
        QuickConnect.init();
    }

    window.QuickConnect = QuickConnect;

})(window, document);