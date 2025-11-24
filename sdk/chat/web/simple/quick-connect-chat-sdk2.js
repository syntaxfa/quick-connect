/**
 * Quick Connect Chat SDK - PREMIUM ULTRA (Telegram Style)
 * Version: 2.1.0
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
            isLoadingHistory: false,
            // Agent State
            agentStatusText: 'آنلاین',
            agentIsOnline: true,
            onlineTimeout: null,
            heartbeatInterval: null
        },

        init: function(options) {
            Object.assign(this.config, options);
            if (document.readyState === 'loading') {
                document.addEventListener('DOMContentLoaded', () => this.setup());
            } else {
                this.setup();
            }
        },

        setup: function() {
            if (document.getElementById('qc-widget')) return;
            this.injectStyles();
            this.injectHTML();
            this.attachEventListeners();
            this.initializeAuth();
        },

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

        injectStyles: function() {
            const fontLink = document.createElement('link');
            fontLink.href = 'https://fonts.googleapis.com/css2?family=Vazirmatn:wght@100;300;400;500;700;900&display=swap';
            fontLink.rel = 'stylesheet';
            document.head.appendChild(fontLink);

            const style = document.createElement('style');
            style.textContent = `
                :root {
                    --qc-font: 'Vazirmatn', sans-serif;
                    --qc-primary: #6366f1;
                    --qc-primary-dark: #4f46e5;
                    --qc-grad-bg: linear-gradient(135deg, #1e293b 0%, #0f172a 100%);
                    --qc-glass: rgba(30, 41, 59, 0.75);
                    --qc-border: rgba(255, 255, 255, 0.08);
                    --qc-own-msg: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
                    --qc-other-msg: rgba(255, 255, 255, 0.08);
                    --qc-text-main: #f1f5f9;
                    --qc-text-muted: #94a3b8;
                    --qc-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
                }

                #qc-widget * { box-sizing: border-box; margin: 0; padding: 0; outline: none; -webkit-tap-highlight-color: transparent; }

                #qc-widget {
                    position: fixed;
                    ${this.config.position.includes('right') ? 'right' : 'left'}: 24px;
                    bottom: 24px;
                    z-index: 999999;
                    font-family: var(--qc-font);
                    direction: rtl;
                    pointer-events: none;
                    width: auto;
                    height: auto;
                }
                #qc-widget > * { pointer-events: auto; }

                /* --- Launcher Button --- */
                #qc-btn {
                    width: 64px; height: 64px;
                    border-radius: 26px;
                    border: none;
                    background: var(--qc-own-msg);
                    box-shadow: 0 10px 25px rgba(99, 102, 241, 0.4);
                    cursor: pointer;
                    display: flex; align-items: center; justify-content: center;
                    transition: all 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
                    z-index: 100;
                    position: absolute;
                    bottom: 0;
                    ${this.config.position.includes('right') ? 'right' : 'left'}: 0;
                }
                #qc-btn:hover { transform: scale(1.1) rotate(-3deg); border-radius: 50%; }
                #qc-btn.hidden { transform: scale(0) rotate(90deg); opacity: 0; pointer-events: none; }

                #qc-badge {
                    position: absolute; top: -5px; right: -5px;
                    background: #ef4444; color: white;
                    min-width: 24px; height: 24px; border-radius: 12px;
                    font-size: 11px; font-weight: 800;
                    display: flex; align-items: center; justify-content: center;
                    border: 3px solid #0f172a;
                }

                /* --- Main Window --- */
                #qc-window {
                    width: 380px; height: 680px;
                    max-height: calc(100vh - 100px);
                    background: #0f172a;
                    border-radius: 24px;
                    box-shadow: var(--qc-shadow);
                    position: absolute;
                    bottom: 80px;
                    ${this.config.position.includes('right') ? 'right' : 'left'}: 0;
                    display: flex; flex-direction: column;
                    overflow: hidden;
                    opacity: 0;
                    transform-origin: bottom ${this.config.position.includes('right') ? 'right' : 'left'};
                    transform: scale(0.9) translateY(20px);
                    transition: all 0.35s cubic-bezier(0.4, 0, 0.2, 1);
                    border: 1px solid var(--qc-border);
                    visibility: hidden;
                }
                #qc-window.open {
                    opacity: 1; transform: scale(1) translateY(0);
                    visibility: visible;
                }

                /* --- Header --- */
                #qc-header {
                    padding: 16px 20px;
                    background: rgba(15, 23, 42, 0.95);
                    border-bottom: 1px solid var(--qc-border);
                    display: flex; align-items: center; justify-content: space-between;
                    z-index: 10;
                    height: 70px; flex-shrink: 0;
                }
                #qc-header-info { display: flex; align-items: center; gap: 12px; }
                #qc-avatar {
                    width: 42px; height: 42px;
                    background: #1e293b; border-radius: 16px;
                    display: flex; align-items: center; justify-content: center;
                    position: relative;
                    border: 1px solid rgba(255,255,255,0.1);
                }
                #qc-avatar img { width: 100%; height: 100%; object-fit: cover; border-radius: 16px; }
                #qc-status-dot {
                    position: absolute; bottom: -2px; right: -2px;
                    width: 12px; height: 12px; background: #10b981;
                    border: 2px solid #0f172a; border-radius: 50%;
                }
                #qc-status-dot.offline { background: #64748b; }
                #qc-header-title { color: var(--qc-text-main); font-size: 15px; font-weight: 700; }
                #qc-status-text { color: var(--qc-text-muted); font-size: 11px; margin-top: 2px; }
                #qc-close-btn {
                    background: transparent; border: none; color: var(--qc-text-muted);
                    cursor: pointer; padding: 8px; border-radius: 8px; transition: 0.2s;
                }
                #qc-close-btn:hover { background: rgba(255,255,255,0.1); color: white; }

                #qc-messages, #qc-emoji-picker, #qc-input {
                    scrollbar-width: thin;
                    scrollbar-color: rgba(255, 255, 255, 0.15) transparent;
                }

                #qc-messages::-webkit-scrollbar,
                #qc-emoji-picker::-webkit-scrollbar,
                #qc-input::-webkit-scrollbar {
                    width: 6px;
                }

                #qc-messages::-webkit-scrollbar-track,
                #qc-emoji-picker::-webkit-scrollbar-track,
                #qc-input::-webkit-scrollbar-track {
                    background: transparent;
                }

                #qc-messages::-webkit-scrollbar-thumb,
                #qc-emoji-picker::-webkit-scrollbar-thumb,
                #qc-input::-webkit-scrollbar-thumb {
                    background-color: rgba(255, 255, 255, 0.15);
                    border-radius: 10px;
                }

                #qc-messages::-webkit-scrollbar-thumb:hover,
                #qc-emoji-picker::-webkit-scrollbar-thumb:hover,
                #qc-input::-webkit-scrollbar-thumb:hover {
                    background-color: rgba(255, 255, 255, 0.3);
                }

                /* Corner fix for input area if both scrolls appear */
                #qc-input::-webkit-scrollbar-corner { background: transparent; }


                /* --- Messages --- */
                #qc-messages {
                    flex: 1; overflow-y: auto; padding: 16px;
                    display: flex; flex-direction: column; gap: 10px;
                    background-image: radial-gradient(circle at center, #1e293b 0%, #0f172a 100%);
                    scroll-behavior: smooth;
                }

                .qc-msg { display: flex; width: 100%; animation: qc-fadeUp 0.3s ease; }
                .qc-msg.own { justify-content: flex-end; }

                .qc-msg-bubble {
                    max-width: 85%;
                    padding: 12px 18px;
                    border-radius: 20px;
                    position: relative;
                    font-size: 14px;
                    line-height: 1.6;
                    box-shadow: 0 2px 5px rgba(0,0,0,0.1);
                    word-wrap: break-word;
                    overflow-wrap: break-word;
                }
                .qc-msg.own .qc-msg-bubble {
                    background: var(--qc-own-msg); color: white;
                    border-bottom-right-radius: 4px;
                }
                .qc-msg.other .qc-msg-bubble {
                    background: #1e293b; color: var(--qc-text-main);
                    border-bottom-left-radius: 4px; border: 1px solid var(--qc-border);
                }
                .qc-msg.system { justify-content: center; margin: 10px 0; }
                .qc-msg.system .qc-msg-bubble {
                    background: rgba(0,0,0,0.3); color: var(--qc-text-muted);
                    font-size: 11px; padding: 4px 12px; border-radius: 100px;
                }
                .qc-msg-footer {
                    display: flex; justify-content: flex-end; align-items: center;
                    gap: 4px; margin-top: 4px; opacity: 0.7; font-size: 10px;
                }

                /* --- Input Area (Floating) --- */
                #qc-input-area {
                    padding: 0 12px 12px 12px;
                    background: transparent;
                    flex-shrink: 0;
                    position: relative;
                }

                #qc-emoji-picker {
                    max-height: 0;
                    opacity: 0;
                    overflow: hidden;
                    margin-bottom: 0;
                    padding: 0;
                    border: 0 solid var(--qc-border);

                    background: #1e293b;
                    border-radius: 16px;
                    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
                }

                #qc-emoji-picker.show {
                    max-height: 200px;
                    opacity: 1;
                    margin-bottom: 8px;
                    padding: 8px;
                    border-width: 1px;
                    overflow-y: auto;
                }

                #qc-emoji-grid { display: grid; grid-template-columns: repeat(8, 1fr); gap: 4px; }
                .qc-emoji-btn {
                    background: transparent; border: none; font-size: 22px;
                    padding: 6px; border-radius: 8px; cursor: pointer;
                    transition: transform 0.1s;
                }
                .qc-emoji-btn:hover { background: rgba(255,255,255,0.1); transform: scale(1.1); }

                #qc-input-container {
                    background: #1e293b;
                    border: 1px solid var(--qc-border);
                    border-radius: 20px;
                    padding: 8px;
                    display: flex; align-items: flex-end;
                    gap: 8px;
                    box-shadow: 0 4px 20px rgba(0,0,0,0.2);
                    transition: border-color 0.2s;
                }
                #qc-input-container:focus-within { border-color: var(--qc-primary); }

                #qc-input-wrapper { flex: 1; padding-bottom: 2px; }

                #qc-input {
                    width: 100%; background: transparent; border: none;
                    color: white; font-family: inherit; font-size: 14px;
                    line-height: 20px;
                    max-height: 120px;
                    resize: none; padding: 6px 0;
                    display: block;
                    /* مطمئن می‌شویم اسکرول بار فضای کافی دارد */
                    overflow-y: auto;
                }
                #qc-input::placeholder { color: #64748b; }

                .qc-icon-btn {
                    width: 36px; height: 36px;
                    border-radius: 50%; border: none;
                    background: transparent; color: #94a3b8;
                    display: flex; align-items: center; justify-content: center;
                    cursor: pointer; flex-shrink: 0; transition: 0.2s;
                }
                .qc-icon-btn:hover { background: rgba(255,255,255,0.1); color: var(--qc-primary); }

                #qc-send-btn {
                    background: var(--qc-primary); color: white;
                    border-radius: 12px;
                    width: 36px; height: 36px;
                    box-shadow: 0 2px 10px rgba(99, 102, 241, 0.3);
                }
                #qc-send-btn:hover { transform: scale(1.05); }
                #qc-send-btn:active { transform: scale(0.95); }

                /* --- Typing Indicator --- */
                .qc-typing { display: flex; gap: 4px; padding: 12px 16px; background: #1e293b; border-radius: 18px; border-bottom-left-radius: 4px; width: fit-content; border: 1px solid var(--qc-border); }
                .qc-typing-dot { width: 6px; height: 6px; background: #94a3b8; border-radius: 50%; animation: qc-bounce 1.4s infinite ease-in-out both; }
                .qc-typing-dot:nth-child(1) { animation-delay: -0.32s; }
                .qc-typing-dot:nth-child(2) { animation-delay: -0.16s; }

                /* --- Profile Banner --- */
                #qc-profile-banner {
                    background: rgba(249, 115, 22, 0.1); border-bottom: 1px solid rgba(249, 115, 22, 0.2);
                    padding: 12px 16px; display: none;
                }
                #qc-profile-banner.show { display: block; animation: qc-fadeUp 0.3s; }
                #qc-profile-header { display: flex; justify-content: space-between; margin-bottom: 8px; font-size: 13px; color: #fdba74; font-weight: bold; }
                #qc-profile-form input {
                    width: 100%; background: rgba(0,0,0,0.3); border: 1px solid rgba(255,255,255,0.1);
                    padding: 8px 12px; border-radius: 8px; color: white; margin-bottom: 6px; font-size: 13px;
                }
                #qc-profile-submit { width: 100%; background: #f97316; color: white; border: none; padding: 8px; border-radius: 8px; cursor: pointer; font-weight: bold; }

                @keyframes qc-fadeUp { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
                @keyframes qc-bounce { 0%, 80%, 100% { transform: scale(0); } 40% { transform: scale(1); } }

                @media (max-width: 480px) {
                    #qc-widget {
                        top: 0; left: 0; right: 0; bottom: 0;
                        width: 100%; height: 100%;
                    }
                    #qc-btn {
                        position: absolute;
                        bottom: 24px;
                        ${this.config.position.includes('right') ? 'right' : 'left'}: 24px;
                    }
                    #qc-window {
                        width: 100% !important;
                        height: 100% !important;
                        max-height: 100vh !important;
                        border-radius: 0 !important;
                        top: 0; left: 0; right: 0; bottom: 0;
                        transform: translateY(100%);
                        border: none;
                    }
                    #qc-window.open {
                        transform: translateY(0);
                    }
                    #qc-header { padding-top: 10px; height: 60px; }
                    #qc-messages { padding-bottom: 80px; }
                }
            `;
            document.head.appendChild(style);
        },

        injectHTML: function() {
            const container = document.createElement('div');
            container.id = 'qc-widget';
            container.innerHTML = `
                <button id="qc-btn">
                    <span id="qc-badge" style="display: none">0</span>
                    <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="white" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
                    </svg>
                </button>

                <div id="qc-window">
                    <div id="qc-header">
                        <div id="qc-header-info">
                            <div id="qc-avatar">
                                <img id="qc-avatar-img" src="" style="display:none">
                                <svg id="qc-avatar-svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#64748b" stroke-width="2">
                                    <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                                    <circle cx="12" cy="7" r="4"></circle>
                                </svg>
                                <span id="qc-status-dot"></span>
                            </div>
                            <div>
                                <div id="qc-header-title">پشتیبانی</div>
                                <div id="qc-status-text">اتصال...</div>
                            </div>
                        </div>
                        <button id="qc-close-btn">
                            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                <line x1="18" y1="6" x2="6" y2="18"></line>
                                <line x1="6" y1="6" x2="18" y2="18"></line>
                            </svg>
                        </button>
                    </div>

                    <div id="qc-profile-banner">
                        <div id="qc-profile-header">
                            <span>اطلاعات تماس</span>
                            <span id="qc-profile-close" style="cursor:pointer">✕</span>
                        </div>
                        <div id="qc-profile-form">
                            <input type="text" id="qc-profile-name" placeholder="نام شما">
                            <input type="email" id="qc-profile-email" placeholder="ایمیل">
                            <button id="qc-profile-submit">ثبت و ادامه</button>
                        </div>
                    </div>

                    <div id="qc-messages"></div>

                    <div id="qc-input-area">
                        <div id="qc-emoji-picker">
                            <div id="qc-emoji-grid"></div>
                        </div>
                        <div id="qc-input-container">
                            <button id="qc-emoji-toggle" class="qc-icon-btn">
                                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                    <circle cx="12" cy="12" r="10"></circle>
                                    <path d="M8 14s1.5 2 4 2 4-2 4-2"></path>
                                    <line x1="9" y1="9" x2="9.01" y2="9"></line>
                                    <line x1="15" y1="9" x2="15.01" y2="9"></line>
                                </svg>
                            </button>
                            <div id="qc-input-wrapper">
                                <textarea id="qc-input" placeholder="پیام..." rows="1"></textarea>
                            </div>
                            <button id="qc-send-btn" class="qc-icon-btn">
                                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                    <line x1="22" y1="2" x2="11" y2="13"></line>
                                    <polygon points="22 2 15 22 11 13 2 9 22 2"></polygon>
                                </svg>
                            </button>
                        </div>
                    </div>
                </div>
            `;
            document.body.appendChild(container);

            const emojis = ['👍', '👋', '😊', '😂', '🤔', '❤️', '🔥', '🎉', '🙏', '👀', '✅', '🚀', '💯', '✨', '💪', '😍'];
            const emojiGrid = document.getElementById('qc-emoji-grid');
            emojis.forEach(emoji => {
                const btn = document.createElement('button');
                btn.className = 'qc-emoji-btn';
                btn.textContent = emoji;
                btn.onclick = () => this.insertEmoji(emoji);
                emojiGrid.appendChild(btn);
            });
        },

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

            const msgs = document.getElementById('qc-messages');
            msgs.addEventListener('scroll', () => {
                if (msgs.scrollTop === 0 && this.state.hasMore && !this.state.isLoadingHistory) {
                    this.fetchChatHistory();
                }
            });
        },

        initializeAuth: async function() {
            const token = localStorage.getItem('QC_TOKEN');
            const userState = localStorage.getItem('QC_USER_STATE');
            if (token && userState) {
                this.state.token = token;
                this.state.userState = userState;
            } else {
                await this.registerGuest();
            }
            if (this.state.token) {
                const decoded = this.parseJwt(this.state.token);
                if (decoded && decoded.user_id) this.state.userId = decoded.user_id;
                await this.fetchActiveConversation();
                this.connectWebSocket();
            }
        },

        registerGuest: async function() {
            try {
                const res = await fetch(`${this.config.managerUrl}/users/guest/register`, {
                    method: 'POST', headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ fullname: "Guest Visitor" })
                });
                const data = await res.json();
                if (data.qc_token) {
                    localStorage.setItem('QC_TOKEN', data.qc_token);
                    localStorage.setItem('QC_USER_STATE', 'guest');
                    this.state.token = data.qc_token;
                    this.state.userState = 'guest';
                    this.showProfileBanner();
                }
            } catch (e) { console.error(e); }
        },

        fetchActiveConversation: async function() {
            if (!this.state.token) return;
            try {
                const res = await fetch(`${this.config.chatApiUrl}/conversations/active`, {
                    headers: { 'Authorization': `Bearer ${this.state.token}` }
                });
                const data = await res.json();
                if (data && data.id) {
                    this.state.conversationId = data.id;
                    localStorage.setItem('QC_CONVERSATION_ID', data.id);

                    const title = document.getElementById('qc-header-title');
                    const img = document.getElementById('qc-avatar-img');
                    const svg = document.getElementById('qc-avatar-svg');
                    const info = data.support_info || {};

                    title.textContent = (data.assigned_support_id && info.fullname) ? info.fullname : 'پشتیبانی آنلاین';

                    if (info.avatar) {
                        img.src = info.avatar; img.style.display = 'block'; svg.style.display = 'none';
                    } else {
                        img.style.display = 'none'; svg.style.display = 'block';
                    }

                    if (info.last_online_at) {
                        const last = new Date(info.last_online_at);
                        const diff = (new Date() - last) / 60000;
                        if (diff < 5 && diff >= 0) {
                            this.state.agentStatusText = 'آنلاین'; this.state.agentIsOnline = true;
                        } else {
                            this.state.agentStatusText = last.toLocaleTimeString('fa-IR', {hour:'2-digit', minute:'2-digit'});
                            this.state.agentIsOnline = false;
                        }
                    }
                    this.updateConnectionStatus();
                    if (this.state.isOpen && !this.state.historyLoaded) this.fetchChatHistory();
                }
            } catch (e) { console.error(e); }
        },

        fetchChatHistory: async function() {
            if (this.state.isLoadingHistory) return;
            this.state.isLoadingHistory = true;
            const container = document.getElementById('qc-messages');
            const prevH = container ? container.scrollHeight : 0;

            try {
                const res = await fetch(`${this.config.chatApiUrl}/chats`, {
                    method: 'POST',
                    headers: { 'Authorization': `Bearer ${this.state.token}`, 'Content-Type': 'application/json' },
                    body: JSON.stringify({ conversation_id: this.state.conversationId, pagination: { cursor: this.state.nextCursor, limit: 20 } })
                });
                const data = await res.json();
                if (data.results) {
                    const msgs = data.results.map(m => ({
                        id: m.id, content: m.content, created_at: m.created_at,
                        isOwn: m.sender_id === this.state.userId, status: 'sent'
                    })).reverse();
                    this.state.messages = [...msgs, ...this.state.messages];
                    if (data.paginate) {
                        this.state.nextCursor = data.paginate.next_cursor;
                        this.state.hasMore = data.paginate.has_more;
                    }
                    this.state.historyLoaded = true;
                    this.renderMessages(prevH > 0);
                    if (prevH > 0) container.scrollTop = container.scrollHeight - prevH;
                }
            } catch (e) { console.error(e); }
            finally { this.state.isLoadingHistory = false; }
        },

        updateProfile: async function() {
            const name = document.getElementById('qc-profile-name').value;
            const email = document.getElementById('qc-profile-email').value;
            if (!name) return;
            try {
                await fetch(`${this.config.managerUrl}/users/guest/update`, {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${this.state.token}` },
                    body: JSON.stringify({ fullname: name, email: email })
                });
                this.hideProfileBanner();
                this.addSystemMessage('اطلاعات ثبت شد ✅');
            } catch (e) {}
        },

        connectWebSocket: function() {
            if (this.state.ws && this.state.ws.readyState === WebSocket.OPEN) return;
            try { this.state.ws = new WebSocket(this.config.chatUrl, [this.state.token]); } catch(e){ return; }

            this.state.ws.onopen = () => {
                this.state.isConnected = true;
                this.updateConnectionStatus();
                this.startHeartbeat();
                if(this.state.messages.length === 0 && this.state.historyLoaded) this.addSystemMessage('خوش آمدید! 👋');
            };
            this.state.ws.onclose = () => {
                this.state.isConnected = false; this.updateConnectionStatus(); this.stopHeartbeat();
                setTimeout(() => this.connectWebSocket(), 3000);
            };
            this.state.ws.onmessage = (e) => this.handleIncomingMessage(JSON.parse(e.data));
        },

        startHeartbeat: function() {
            this.stopHeartbeat();
            const ping = () => {
                if(this.state.ws?.readyState === WebSocket.OPEN && this.state.conversationId) {
                    this.state.ws.send(JSON.stringify({ type: 'system', sub_type: 'online', conversation_id: this.state.conversationId }));
                }
            };
            ping();
            this.state.heartbeatInterval = setInterval(ping, 50000);
        },
        stopHeartbeat: function() { clearInterval(this.state.heartbeatInterval); },

        updateConnectionStatus: function() {
            const txt = document.getElementById('qc-status-text');
            const dot = document.getElementById('qc-status-dot');
            if(!txt) return;
            if (!this.state.isConnected) {
                txt.textContent = 'درحال اتصال...'; dot.classList.add('offline');
            } else if (this.state.typingUsers.size > 0) {
                txt.textContent = 'درحال نوشتن...'; dot.classList.remove('offline');
            } else {
                txt.textContent = this.state.agentStatusText;
                this.state.agentIsOnline ? dot.classList.remove('offline') : dot.classList.add('offline');
            }
        },

        updateBadgeUI: function() {
            const b = document.getElementById('qc-badge');
            if(b) {
                b.style.display = this.state.unreadCount > 0 ? 'flex' : 'none';
                b.textContent = this.state.unreadCount;
            }
        },

        sendMessage: function() {
            const input = document.getElementById('qc-input');
            const content = input.value.trim();
            if (!content || !this.state.ws) return;

            const tempId = `loc_${Date.now()}`;
            this.state.ws.send(JSON.stringify({
                type: 'text', content: content,
                conversation_id: this.state.conversationId, client_message_id: tempId
            }));

            this.state.messages.push({
                id: tempId, content: content, created_at: new Date().toISOString(),
                isOwn: true, status: 'pending'
            });
            this.renderMessages();

            input.value = '';
            input.style.height = 'auto'; // Reset height
            this.sendTypingStopped();
        },

        handleIncomingMessage: function(msg) {
            if (msg.type === 'text') {
                if (msg.client_message_id) {
                    const idx = this.state.messages.findIndex(m => m.id === msg.client_message_id);
                    if (idx !== -1) {
                        this.state.messages[idx].status = 'sent';
                        if(msg.payload?.id) this.state.messages[idx].id = msg.payload.id;
                        this.renderMessages(); return;
                    }
                }
                if (msg.payload.sender_id === this.state.userId) return;

                if (!this.state.isOpen) { this.state.unreadCount++; this.updateBadgeUI(); }
                if (msg.payload.conversation_id && !this.state.conversationId) this.state.conversationId = msg.payload.conversation_id;

                this.state.messages.push({ ...msg.payload, isOwn: false });
                this.renderMessages();
            } else if (msg.type === 'system') {
                this.handleSystemMessage(msg);
            }
        },

        handleSystemMessage: function(data) {
            const sid = data.payload?.sender_id;
            if (!sid) return;

            if (data.sub_type === 'typing_started') {
                this.state.typingUsers.add(sid);
                this.updateConnectionStatus(); this.renderMessages(true);
                clearTimeout(this.state.typingTimeouts[sid]);
                this.state.typingTimeouts[sid] = setTimeout(() => {
                    this.state.typingUsers.delete(sid);
                    this.updateConnectionStatus(); this.renderMessages(true);
                }, 6000);
            } else if (data.sub_type === 'typing_stopped') {
                this.state.typingUsers.delete(sid);
                this.updateConnectionStatus(); this.renderMessages(true);
            } else if (data.sub_type === 'online') {
                this.state.agentIsOnline = true; this.state.agentStatusText = 'آنلاین';
                this.updateConnectionStatus();
                clearTimeout(this.state.onlineTimeout);
                this.state.onlineTimeout = setTimeout(() => this.fetchActiveConversation(), 60000);
            }
        },

        addSystemMessage: function(txt) {
            this.state.messages.push({ id: 'sys'+Date.now(), content: txt, isSystem: true });
            this.renderMessages();
        },

        handleTyping: function() {
            const now = Date.now();
            if (now - this.state.lastTypingSent > 5000) {
                this.sendTypingStarted(); this.state.lastTypingSent = now;
            }
            const input = document.getElementById('qc-input');
            input.style.height = 'auto';
            input.style.height = Math.min(input.scrollHeight, 120) + 'px'; // Telegram grow effect
        },

        sendTypingStarted: function() {
            if(this.state.ws) this.state.ws.send(JSON.stringify({type:'system', sub_type:'typing_started', conversation_id:this.state.conversationId}));
        },
        sendTypingStopped: function() {
            if(this.state.ws) this.state.ws.send(JSON.stringify({type:'system', sub_type:'typing_stopped', conversation_id:this.state.conversationId}));
            this.state.lastTypingSent = 0;
        },

        renderMessages: function(preserve = false) {
            const c = document.getElementById('qc-messages');
            if (!c) return;
            c.innerHTML = '';

            this.state.messages.forEach(m => {
                const el = document.createElement('div');
                if (m.isSystem) {
                    el.className = 'qc-msg system';
                    el.innerHTML = `<div class="qc-msg-bubble">${this.escapeHtml(m.content)}</div>`;
                } else {
                    el.className = `qc-msg ${m.isOwn ? 'own' : 'other'}`;
                    const time = new Date(m.created_at).toLocaleTimeString('fa-IR', {hour:'2-digit', minute:'2-digit'});
                    const tick = m.isOwn ? (m.status === 'pending' ? '🕒' : '✓') : '';
                    el.innerHTML = `
                        <div class="qc-msg-bubble">
                            ${this.escapeHtml(m.content)}
                            <div class="qc-msg-footer">
                                <span>${time}</span> <span>${tick}</span>
                            </div>
                        </div>`;
                }
                c.appendChild(el);
            });

            if (this.state.typingUsers.size > 0) {
                const t = document.createElement('div'); t.className = 'qc-msg other';
                t.innerHTML = `<div class="qc-typing"><div class="qc-typing-dot"></div><div class="qc-typing-dot"></div><div class="qc-typing-dot"></div></div>`;
                c.appendChild(t);
            }
            if (!preserve) setTimeout(() => c.scrollTop = c.scrollHeight, 10);
        },

        escapeHtml: function(t) {
            const d = document.createElement('div'); d.textContent = t;
            return d.innerHTML.replace(/\n/g, '<br>');
        },

        openChat: function() {
            this.state.isOpen = true; this.state.unreadCount = 0; this.updateBadgeUI();
            this.fetchActiveConversation();
            document.getElementById('qc-btn').classList.add('hidden');
            document.getElementById('qc-window').classList.add('open');
            setTimeout(() => document.getElementById('qc-input').focus(), 300);
        },
        closeChat: function() {
            this.state.isOpen = false;
            document.getElementById('qc-btn').classList.remove('hidden');
            document.getElementById('qc-window').classList.remove('open');
        },
        toggleEmojiPicker: function() {
            document.getElementById('qc-emoji-picker').classList.toggle('show');
        },
        insertEmoji: function(emoji) {
            const input = document.getElementById('qc-input');
            input.value += emoji;
            input.focus();
            this.toggleEmojiPicker(); // Auto close
        },
        showProfileBanner: function() { document.getElementById('qc-profile-banner')?.classList.add('show'); },
        hideProfileBanner: function() { document.getElementById('qc-profile-banner')?.classList.remove('show'); },
    };

    if (document.readyState === 'loading') document.addEventListener('DOMContentLoaded', () => QuickConnect.init());
    else QuickConnect.init();
    window.QuickConnect = QuickConnect;

})(window, document);