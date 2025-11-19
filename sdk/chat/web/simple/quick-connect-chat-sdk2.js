/**
 * Quick Connect Chat SDK
 * نسخه: 1.0.0
 *
 * نحوه استفاده:
 * <script src="quick-connect-sdk.js"></script>
 * <script>
 *   QuickConnect.init({
 *     managerUrl: 'https://your-manager-url.com',
 *     chatUrl: 'wss://your-chat-url.com'
 *   });
 * </script>
 */

(function(window, document) {
    'use strict';

    const QuickConnect = {
        config: {
            managerUrl: 'http://localhost:2531',
            chatUrl: 'ws://localhost:2530',
            position: 'bottom-left',
            theme: 'purple',
            lang: 'fa'
        },

        state: {
            isOpen: false,
            isConnected: false,
            token: null,
            userState: null,
            conversationId: '01KAAWW5V95WFFDERJZCTDE2QQ',
            messages: [],
            typingUsers: new Set(),
            ws: null,
            lastTypingSent: 0,
            typingTimeouts: {},
            reconnectAttempts: 0
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
            this.injectStyles();
            this.injectHTML();
            this.attachEventListeners();
            this.initializeAuth();
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
        }

        #qc-btn:hover {
          transform: scale(1.1);
          box-shadow: 0 12px 32px rgba(102, 126, 234, 0.6);
        }

        #qc-btn.hidden { display: none; }

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
          display: flex;
          align-items: center;
          justify-content: center;
          animation: qc-bounce 1s infinite;
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
          border-radius: 16px;
          box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
          overflow: hidden;
          flex-direction: column;
          border: 1px solid rgba(168, 85, 247, 0.3);
          animation: qc-slideIn 0.3s ease;
        }

        #qc-window.open { display: flex; }

        #qc-header {
          background: linear-gradient(135deg, #7c3aed 0%, #2563eb 100%);
          padding: 16px;
          display: flex;
          align-items: center;
          justify-content: space-between;
        }

        #qc-header-info {
          display: flex;
          align-items: center;
          gap: 12px;
        }

        #qc-avatar {
          width: 40px;
          height: 40px;
          background: white;
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          position: relative;
        }

        #qc-status-dot {
          position: absolute;
          bottom: 0;
          right: 0;
          width: 12px;
          height: 12px;
          background: #10b981;
          border: 2px solid white;
          border-radius: 50%;
        }

        #qc-status-dot.offline { background: #6b7280; }

        #qc-header-text h3 {
          color: white;
          font-size: 18px;
          font-weight: bold;
          margin-bottom: 2px;
        }

        #qc-header-text p {
          color: rgba(255, 255, 255, 0.8);
          font-size: 12px;
        }

        #qc-close-btn {
          background: rgba(255, 255, 255, 0.2);
          border: none;
          color: white;
          width: 32px;
          height: 32px;
          border-radius: 8px;
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
          padding: 12px 16px;
        }

        #qc-profile-banner.show { display: block; }

        #qc-profile-header {
          display: flex;
          align-items: center;
          justify-content: space-between;
          margin-bottom: 8px;
        }

        #qc-profile-header p {
          color: white;
          font-size: 14px;
          font-weight: 500;
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
          gap: 8px;
        }

        #qc-profile-form input {
          width: 100%;
          padding: 10px 12px;
          border-radius: 8px;
          border: none;
          background: rgba(255, 255, 255, 0.9);
          font-size: 13px;
        }

        #qc-profile-actions {
          display: flex;
          gap: 8px;
        }

        #qc-profile-submit {
          flex: 1;
          padding: 10px;
          background: white;
          color: #ea580c;
          border: none;
          border-radius: 8px;
          font-weight: 600;
          font-size: 13px;
          cursor: pointer;
          transition: background 0.2s;
        }

        #qc-profile-submit:hover { background: rgba(255, 255, 255, 0.9); }

        #qc-messages {
          flex: 1;
          overflow-y: auto;
          padding: 16px;
          background: rgba(15, 23, 42, 0.5);
          display: flex;
          flex-direction: column;
          gap: 12px;
        }

        #qc-messages::-webkit-scrollbar { width: 6px; }
        #qc-messages::-webkit-scrollbar-track { background: transparent; }
        #qc-messages::-webkit-scrollbar-thumb { background: rgba(148, 163, 184, 0.3); border-radius: 3px; }

        .qc-msg {
          display: flex;
          animation: qc-slideIn 0.2s ease;
        }

        .qc-msg.own { justify-content: flex-end; }
        .qc-msg.other { justify-content: flex-start; }
        .qc-msg.system { justify-content: center; }

        .qc-msg-bubble {
          max-width: 75%;
          padding: 12px 16px;
          border-radius: 16px;
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
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
        }

        .qc-msg.system .qc-msg-bubble {
          background: rgba(51, 65, 85, 0.5);
          color: #cbd5e1;
          font-size: 12px;
          padding: 6px 12px;
          border-radius: 12px;
        }

        .qc-msg-content {
          font-size: 14px;
          line-height: 1.5;
          word-wrap: break-word;
        }

        .qc-msg-footer {
          display: flex;
          align-items: center;
          justify-content: flex-end;
          gap: 4px;
          margin-top: 4px;
        }

        .qc-msg-time {
          font-size: 11px;
          opacity: 0.7;
        }

        .qc-msg-status {
          width: 16px;
          height: 16px;
          opacity: 0.7;
        }

        .qc-typing {
          display: flex;
          gap: 4px;
          padding: 12px 16px;
          background: #1e293b;
          border-radius: 16px;
          border-bottom-left-radius: 4px;
          width: fit-content;
        }

        .qc-typing-dot {
          width: 8px;
          height: 8px;
          background: #64748b;
          border-radius: 50%;
          animation: qc-typing 1.4s infinite;
        }

        .qc-typing-dot:nth-child(2) { animation-delay: 0.2s; }
        .qc-typing-dot:nth-child(3) { animation-delay: 0.4s; }

        @keyframes qc-typing {
          0%, 60%, 100% { transform: translateY(0); }
          30% { transform: translateY(-8px); }
        }

        #qc-emoji-picker {
          display: none;
          background: #1e293b;
          border-top: 1px solid #334155;
          padding: 12px;
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
          padding: 4px;
          border-radius: 8px;
          transition: background 0.2s;
        }

        .qc-emoji-btn:hover { background: #334155; }

        #qc-input-area {
          background: #0f172a;
          border-top: 1px solid #334155;
          padding: 16px;
        }

        #qc-input-container {
          display: flex;
          align-items: flex-end;
          gap: 8px;
        }

        #qc-emoji-toggle, #qc-send-btn {
          background: rgba(71, 85, 105, 0.5);
          border: none;
          color: #94a3b8;
          width: 40px;
          height: 40px;
          border-radius: 12px;
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
          border-radius: 12px;
          padding: 10px 14px;
          border: 1px solid #334155;
          transition: border-color 0.2s;
        }

        #qc-input-wrapper:focus-within { border-color: #7c3aed; }

        #qc-input {
          width: 100%;
          background: transparent;
          border: none;
          color: white;
          font-size: 14px;
          outline: none;
          resize: none;
          max-height: 100px;
          font-family: inherit;
        }

        #qc-input::placeholder { color: #64748b; }
      `;
            document.head.appendChild(style);
        },

        // Inject HTML Structure
        injectHTML: function() {
            const container = document.createElement('div');
            container.id = 'qc-widget';
            container.innerHTML = `
        <button id="qc-btn" aria-label="Open chat">
          <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
          </svg>
          <span id="qc-badge">3</span>
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
            document.getElementById('qc-btn').onclick = () => this.openChat();
            document.getElementById('qc-close-btn').onclick = () => this.closeChat();
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

        // Authentication
        initializeAuth: function() {
            const token = localStorage.getItem('QC_TOKEN');
            const userState = localStorage.getItem('QC_USER_STATE');

            if (token && userState) {
                this.state.token = token;
                this.state.userState = userState;
                this.connectWebSocket();
            } else {
                this.registerGuest();
            }
        },

        registerGuest: function() {
            // Simulate API call
            setTimeout(() => {
                const token = 'guest_token_' + Date.now();
                localStorage.setItem('QC_TOKEN', token);
                localStorage.setItem('QC_USER_STATE', 'guest');

                this.state.token = token;
                this.state.userState = 'guest';

                this.connectWebSocket();
                this.showProfileBanner();
            }, 500);
        },

        updateProfile: function() {
            const name = document.getElementById('qc-profile-name').value;
            const email = document.getElementById('qc-profile-email').value;

            if (!name || !email) {
                alert('لطفا تمام فیلدها را پر کنید');
                return;
            }

            // Simulate API call
            console.log('Updating profile:', { name, email });
            this.hideProfileBanner();
            this.addSystemMessage('پروفایل شما با موفقیت ثبت شد ✅');
        },

        // WebSocket
        connectWebSocket: function() {
            // Simulate WebSocket connection
            setTimeout(() => {
                this.state.isConnected = true;
                this.updateConnectionStatus();
                this.addSystemMessage('به چت پشتیبانی خوش آمدید! 👋');

                // Simulate support message
                setTimeout(() => {
                    this.receiveMessage({
                        type: 'text',
                        payload: {
                            id: 'msg_' + Date.now(),
                            content: 'سلام! چطور میتونم کمکتون کنم؟ 😊',
                            created_at: new Date().toISOString(),
                            isSupport: true
                        }
                    });
                }, 1500);
            }, 1000);
        },

        updateConnectionStatus: function() {
            const statusText = document.getElementById('qc-status-text');
            const statusDot = document.getElementById('qc-status-dot');

            if (this.state.isConnected) {
                statusText.textContent = 'آنلاین';
                statusDot.classList.remove('offline');
            } else {
                statusText.textContent = 'آفلاین';
                statusDot.classList.add('offline');
            }
        },

        // UI Control
        openChat: function() {
            this.state.isOpen = true;
            document.getElementById('qc-btn').classList.add('hidden');
            document.getElementById('qc-window').classList.add('open');
            document.getElementById('qc-input').focus();
        },

        closeChat: function() {
            this.state.isOpen = false;
            document.getElementById('qc-btn').classList.remove('hidden');
            document.getElementById('qc-window').classList.remove('open');
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
            document.getElementById('qc-profile-banner').classList.add('show');
        },

        hideProfileBanner: function() {
            document.getElementById('qc-profile-banner').classList.remove('show');
        },

        // Messaging
        sendMessage: function() {
            const input = document.getElementById('qc-input');
            const content = input.value.trim();

            if (!content || !this.state.isConnected) return;

            const message = {
                id: 'local_' + Date.now(),
                content: content,
                created_at: new Date().toISOString(),
                isOwn: true,
                status: 'sent'
            };

            this.state.messages.push(message);
            this.renderMessages();

            input.value = '';
            input.style.height = 'auto';
            this.sendTypingStopped();

            // Simulate response
            setTimeout(() => {
                this.receiveMessage({
                    type: 'text',
                    payload: {
                        id: 'msg_' + Date.now(),
                        content: 'پیام شما دریافت شد. یک لحظه صبر کنید...',
                        created_at: new Date().toISOString(),
                        isSupport: true
                    }
                });
            }, 1000);
        },

        receiveMessage: function(data) {
            if (data.type === 'text') {
                this.state.messages.push({
                    ...data.payload,
                    isOwn: false
                });
                this.renderMessages();
            } else if (data.type === 'system') {
                this.handleSystemMessage(data);
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
            // Handle typing indicators
            if (data.sub_type === 'typing_started') {
                this.state.typingUsers.add(data.payload.sender_id);
                this.renderMessages();

                setTimeout(() => {
                    this.state.typingUsers.delete(data.payload.sender_id);
                    this.renderMessages();
                }, 6000);
            } else if (data.sub_type === 'typing_stopped') {
                this.state.typingUsers.delete(data.payload.sender_id);
                this.renderMessages();
            }
        },

        handleTyping: function() {
            const now = Date.now();
            if (now - this.state.lastTypingSent > 5000) {
                this.sendTypingStarted();
                this.state.lastTypingSent = now;
            }

            // Auto-resize textarea
            const input = document.getElementById('qc-input');
            input.style.height = 'auto';
            input.style.height = input.scrollHeight + 'px';
        },

        sendTypingStarted: function() {
            // Send typing_started via WebSocket
            console.log('Typing started');
        },

        sendTypingStopped: function() {
            // Send typing_stopped via WebSocket
            console.log('Typing stopped');
        },

        // Render
        renderMessages: function() {
            const container = document.getElementById('qc-messages');
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
                    msgEl.innerHTML = `
            <div class="qc-msg-bubble">
              <div class="qc-msg-content">${this.escapeHtml(msg.content)}</div>
              <div class="qc-msg-footer">
                <span class="qc-msg-time">${this.formatTime(msg.created_at)}</span>
                ${msg.isOwn ? `
                  <svg class="qc-msg-status" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="20 6 9 17 4 12"></polyline>
                  </svg>
                ` : ''}
              </div>
            </div>
          `;
                }

                container.appendChild(msgEl);
            });

            // Add typing indicator
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

            // Scroll to bottom
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

        // Public API
        login: function(token) {
            localStorage.setItem('QC_TOKEN', token);
            localStorage.setItem('QC_USER_STATE', 'client');
            this.state.token = token;
            this.state.userState = 'client';
            this.connectWebSocket();
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

    // Auto-initialize on load
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', function() {
            QuickConnect.init();
        });
    } else {
        QuickConnect.init();
    }

    // Expose to global scope
    window.QuickConnect = QuickConnect;

})(window, document);