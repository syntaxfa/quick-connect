/**
 * Quick Connect Chat SDK - PREMIUM DARK UI (Fixed Bubbles)
 * Version: 2.3.0
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
            agentAvatar: null,
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

                    /* --- Dark Theme Palette --- */
                    --qc-bg-body: #0f172a;       /* Dark Slate Background */
                    --qc-bg-header: #1e293b;     /* Slightly Lighter Slate */
                    --qc-bg-input: #1e293b;
                    --qc-border: rgba(255, 255, 255, 0.1);

                    --qc-text-main: #f8fafc;
                    --qc-text-muted: #94a3b8;

                    /* Message Colors */
                    --qc-msg-own-bg: #6366f1;    /* Primary Color */
                    --qc-msg-own-text: #ffffff;
                    --qc-msg-other-bg: #1e293b;  /* Dark Bubble */
                    --qc-msg-other-text: #f1f5f9;

                    --qc-shadow: 0 10px 40px rgba(0,0,0,0.4);
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
                    width: 60px; height: 60px;
                    border-radius: 24px;
                    border: none;
                    background: var(--qc-primary);
                    box-shadow: 0 8px 24px rgba(99, 102, 241, 0.4);
                    cursor: pointer;
                    display: flex; align-items: center; justify-content: center;
                    transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
                    z-index: 100;
                    position: absolute;
                    bottom: 0;
                    ${this.config.position.includes('right') ? 'right' : 'left'}: 0;
                }
                #qc-btn:hover { transform: scale(1.1) rotate(-5deg); border-radius: 50%; }
                #qc-btn.hidden { transform: scale(0) rotate(90deg); opacity: 0; pointer-events: none; }

                #qc-badge {
                    position: absolute; top: -4px; right: -4px;
                    background: #ef4444; color: white;
                    min-width: 22px; height: 22px; border-radius: 11px;
                    font-size: 11px; font-weight: 700;
                    display: flex; align-items: center; justify-content: center;
                    border: 2px solid var(--qc-bg-body);
                    padding: 0 4px;
                }

                /* --- Main Window --- */
                #qc-window {
                    width: 380px; height: 650px;
                    max-height: calc(100vh - 100px);
                    background: var(--qc-bg-body);
                    border-radius: 20px;
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
                    visibility: hidden;
                    border: 1px solid var(--qc-border);
                }
                #qc-window.open {
                    opacity: 1; transform: scale(1) translateY(0);
                    visibility: visible;
                }

                /* --- Header --- */
                #qc-header {
                    padding: 16px 20px;
                    background: var(--qc-bg-header);
                    border-bottom: 1px solid var(--qc-border);
                    display: flex; align-items: center; justify-content: space-between;
                    z-index: 10;
                    height: 70px; flex-shrink: 0;
                }
                #qc-header-info { display: flex; align-items: center; gap: 12px; }
                #qc-avatar {
                    width: 40px; height: 40px;
                    background: rgba(255,255,255,0.1); border-radius: 14px;
                    display: flex; align-items: center; justify-content: center;
                    position: relative;
                }
                #qc-avatar img { width: 100%; height: 100%; object-fit: cover; border-radius: 14px; }
                #qc-status-dot {
                    position: absolute; bottom: -2px; right: -2px;
                    width: 10px; height: 10px; background: #10b981;
                    border: 2px solid var(--qc-bg-header); border-radius: 50%;
                }
                #qc-status-dot.offline { background: #64748b; }
                #qc-header-title { color: var(--qc-text-main); font-size: 15px; font-weight: 700; }
                #qc-status-text { color: var(--qc-text-muted); font-size: 11px; margin-top: 2px; }
                #qc-close-btn {
                    background: transparent; border: none; color: var(--qc-text-muted);
                    cursor: pointer; padding: 8px; border-radius: 8px; transition: 0.2s;
                }
                #qc-close-btn:hover { background: rgba(255,255,255,0.1); color: white; }

                /* --- Messages Container --- */
                #qc-messages {
                    flex: 1;
                    overflow-y: auto;
                    overflow-x: hidden;
                    padding: 16px;
                    display: flex;
                    flex-direction: column;
                    gap: 8px; /* Space between messages */
                    background-image: radial-gradient(circle at 50% 50%, #1e293b 0%, #0f172a 100%);
                    scroll-behavior: smooth;
                }

                /* --- Scrollbar --- */
                #qc-messages::-webkit-scrollbar { width: 5px; }
                #qc-messages::-webkit-scrollbar-track { background: transparent; }
                #qc-messages::-webkit-scrollbar-thumb { background: rgba(255,255,255,0.15); border-radius: 10px; }

                /* --- Date Separator --- */
                .qc-date-separator {
                    display: flex; justify-content: center;
                    margin: 20px 0 10px 0; position: sticky; top: 0; z-index: 5;
                }
                .qc-date-pill {
                    background: rgba(0, 0, 0, 0.4);
                    color: #94a3b8;
                    font-size: 11px; padding: 4px 12px;
                    border-radius: 12px; font-weight: 500;
                    border: 1px solid rgba(255,255,255,0.05);
                    backdrop-filter: blur(4px);
                }

                /* --- Message Bubbles (The Fix) --- */
                .qc-msg {
                    display: flex;
                    width: 100%;
                    margin-bottom: 2px;
                    animation: qc-fadeIn 0.25s ease-out;
                }

                .qc-msg.own { justify-content: flex-end; }
                .qc-msg.other { justify-content: flex-start; align-items: flex-end; }

                .qc-msg-avatar {
                    width: 30px; height: 30px;
                    margin-left: 8px; margin-bottom: 2px;
                    border-radius: 50%; overflow: hidden;
                    flex-shrink: 0; background: #334155;
                    border: 1px solid var(--qc-border);
                }
                .qc-msg-avatar img { width: 100%; height: 100%; object-fit: cover; }

                .qc-msg-bubble {
                    position: relative;
                    max-width: 78%; /* Prevent full width */
                    width: fit-content; /* Crucial for bubble size */
                    min-width: 60px;
                    padding: 8px 12px;
                    font-size: 14px;
                    line-height: 1.6;
                    word-wrap: break-word;
                    white-space: pre-wrap;
                    box-shadow: 0 2px 5px rgba(0,0,0,0.1);
                    display: flex;
                    flex-direction: column; /* Stack text and time */
                }

                /* Own Message Style */
                .qc-msg.own .qc-msg-bubble {
                    background: linear-gradient(135deg, #6366f1 0%, #4f46e5 100%);
                    color: white;
                    border-radius: 16px 16px 2px 16px; /* TopL, TopR, BotR(Sharp), BotL */
                }

                /* Other Message Style */
                .qc-msg.other .qc-msg-bubble {
                    background: var(--qc-msg-other-bg);
                    color: var(--qc-msg-other-text);
                    border: 1px solid var(--qc-border);
                    border-radius: 16px 16px 16px 2px; /* TopL, TopR, BotR, BotL(Sharp) */
                }

                /* System Message Style */
                .qc-msg.system { justify-content: center; margin: 12px 0; }
                .qc-msg.system .qc-msg-bubble {
                    background: rgba(255,255,255,0.05);
                    color: #94a3b8;
                    font-size: 12px;
                    padding: 4px 12px;
                    border-radius: 10px;
                    max-width: 90%;
                    text-align: center;
                    border: none;
                    width: auto;
                }

                .qc-msg-text {
                    margin-bottom: 2px;
                }

                .qc-msg-meta {
                    display: flex;
                    align-items: center;
                    justify-content: flex-end; /* Align time to right */
                    gap: 4px;
                    font-size: 10px;
                    opacity: 0.7;
                    margin-top: 2px;
                    user-select: none;
                    height: 14px;
                }

                .qc-msg.own .qc-msg-meta { color: rgba(255,255,255,0.9); }
                .qc-msg.other .qc-msg-meta { color: #94a3b8; }

                /* --- Input Area --- */
                #qc-input-area {
                    padding: 12px;
                    background: var(--qc-bg-header);
                    border-top: 1px solid var(--qc-border);
                    flex-shrink: 0;
                }

                #qc-emoji-picker {
                    max-height: 0; opacity: 0; overflow: hidden;
                    transition: all 0.2s ease; background: #334155;
                    border-radius: 12px; margin-bottom: 0;
                }
                #qc-emoji-picker.show {
                    max-height: 160px; opacity: 1; margin-bottom: 8px;
                    padding: 8px; border: 1px solid var(--qc-border); overflow-y: auto;
                }
                #qc-emoji-grid { display: grid; grid-template-columns: repeat(8, 1fr); gap: 4px; }
                .qc-emoji-btn {
                    font-size: 22px; padding: 4px; cursor: pointer;
                    background: none; border: none; border-radius: 6px; transition: transform 0.1s;
                }
                .qc-emoji-btn:hover { background: rgba(255,255,255,0.1); transform: scale(1.1); }

                #qc-input-container {
                    display: flex; align-items: flex-end; gap: 8px;
                    background: rgba(0,0,0,0.2);
                    border: 1px solid var(--qc-border);
                    border-radius: 20px; padding: 8px 12px;
                    transition: border-color 0.2s;
                }
                #qc-input-container:focus-within { border-color: var(--qc-primary); }

                #qc-input-wrapper { flex: 1; display: flex; align-items: center; min-height: 24px; }
                #qc-input {
                    width: 100%; background: transparent; border: none;
                    color: white; font-family: inherit; font-size: 14px;
                    resize: none; padding: 2px 0; max-height: 100px;
                    line-height: 20px;
                }
                #qc-input::placeholder { color: #64748b; }

                .qc-icon-btn {
                    width: 34px; height: 34px; border-radius: 50%; border: none;
                    background: transparent; color: #94a3b8;
                    display: flex; align-items: center; justify-content: center;
                    cursor: pointer; transition: 0.2s; flex-shrink: 0;
                }
                .qc-icon-btn:hover { color: var(--qc-primary); background: rgba(255,255,255,0.05); }

                #qc-send-btn {
                    color: white; background: var(--qc-primary);
                    width: 36px; height: 36px; border-radius: 12px;
                    box-shadow: 0 2px 10px rgba(99, 102, 241, 0.3);
                }
                #qc-send-btn:hover { transform: scale(1.05); background: var(--qc-primary-dark); }

                /* --- Typing Indicator --- */
                .qc-typing { display: flex; gap: 4px; padding: 4px 0; }
                .qc-typing-dot {
                    width: 6px; height: 6px; background: #94a3b8; border-radius: 50%;
                    animation: qc-bounce 1.4s infinite ease-in-out both;
                }
                .qc-typing-dot:nth-child(1) { animation-delay: -0.32s; }
                .qc-typing-dot:nth-child(2) { animation-delay: -0.16s; }

                /* --- Profile Banner --- */
                #qc-profile-banner {
                    background: rgba(249, 115, 22, 0.1);
                    border-bottom: 1px solid rgba(249, 115, 22, 0.2);
                    padding: 12px 16px; display: none;
                }
                #qc-profile-banner.show { display: block; animation: qc-fadeIn 0.3s; }
                #qc-profile-form input {
                    width: 100%; background: rgba(0,0,0,0.3);
                    border: 1px solid rgba(255,255,255,0.1);
                    padding: 8px 12px; border-radius: 8px; color: white;
                    margin-bottom: 6px; font-size: 13px;
                }
                #qc-profile-submit {
                    width: 100%; background: #f97316; color: white;
                    border: none; padding: 8px; border-radius: 8px; cursor: pointer; font-weight: bold;
                }

                @keyframes qc-fadeIn { from { opacity: 0; transform: translateY(5px); } to { opacity: 1; transform: translateY(0); } }
                @keyframes qc-bounce { 0%, 80%, 100% { transform: scale(0); } 40% { transform: scale(1); } }

                @media (max-width: 480px) {
                    #qc-widget { right: 0; left: 0; bottom: 0; width: 100%; }
                    #qc-btn { position: fixed; right: 24px; bottom: 24px; }
                    #qc-window {
                        width: 100%; height: 100%; border-radius: 0; border: none;
                        bottom: 0; right: 0; max-height: 100vh;
                    }
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
                                <svg id="qc-avatar-svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="#cbd5e1" stroke-width="2">
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
                        <div style="display:flex;justify-content:space-between;margin-bottom:8px;font-size:13px;color:#fdba74;font-weight:bold;">
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
                        this.state.agentAvatar = info.avatar;
                    } else {
                        img.style.display = 'none'; svg.style.display = 'block';
                        this.state.agentAvatar = null;
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
            input.style.height = 'auto';
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
            this.state.messages.push({ id: 'sys'+Date.now(), content: txt, isSystem: true, created_at: new Date().toISOString() });
            this.renderMessages();
        },

        handleTyping: function() {
            const now = Date.now();
            if (now - this.state.lastTypingSent > 5000) {
                this.sendTypingStarted(); this.state.lastTypingSent = now;
            }
            const input = document.getElementById('qc-input');
            input.style.height = 'auto';
            input.style.height = Math.min(input.scrollHeight, 120) + 'px';
        },

        sendTypingStarted: function() {
            if(this.state.ws) this.state.ws.send(JSON.stringify({type:'system', sub_type:'typing_started', conversation_id:this.state.conversationId}));
        },
        sendTypingStopped: function() {
            if(this.state.ws) this.state.ws.send(JSON.stringify({type:'system', sub_type:'typing_stopped', conversation_id:this.state.conversationId}));
            this.state.lastTypingSent = 0;
        },

        // --- Date Formatter Helper ---
        formatDateSeparator: function(dateString) {
            const date = new Date(dateString);
            const today = new Date();
            const isToday = date.getDate() === today.getDate() &&
                date.getMonth() === today.getMonth() &&
                date.getFullYear() === today.getFullYear();

            return isToday ? 'امروز' : date.toLocaleDateString('fa-IR', { weekday: 'long', day: 'numeric', month: 'long' });
        },

        renderMessages: function(preserve = false) {
            const c = document.getElementById('qc-messages');
            if (!c) return;
            c.innerHTML = '';

            let lastDateKey = null;

            this.state.messages.forEach(m => {
                // --- Date Separator Logic ---
                if (m.created_at) {
                    const dateObj = new Date(m.created_at);
                    // Create a simple key string YYYY-MM-DD for comparison
                    const dateKey = dateObj.toLocaleDateString('en-CA');

                    if (dateKey !== lastDateKey) {
                        const dateSep = document.createElement('div');
                        dateSep.className = 'qc-date-separator';
                        dateSep.innerHTML = `<span class="qc-date-pill">${this.formatDateSeparator(m.created_at)}</span>`;
                        c.appendChild(dateSep);
                        lastDateKey = dateKey;
                    }
                }

                const el = document.createElement('div');
                if (m.isSystem) {
                    el.className = 'qc-msg system';
                    el.innerHTML = `<div class="qc-msg-bubble">${this.escapeHtml(m.content)}</div>`;
                } else {
                    el.className = `qc-msg ${m.isOwn ? 'own' : 'other'}`;
                    const time = new Date(m.created_at).toLocaleTimeString('fa-IR', {hour:'2-digit', minute:'2-digit'});

                    // Tick Icon
                    let tick = '';
                    if (m.isOwn) {
                        if (m.status === 'pending') {
                            tick = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="14" height="14"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>`;
                        } else {
                            tick = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="14" height="14"><polyline points="20 6 9 17 4 12"></polyline></svg>`;
                        }
                    }

                    // Agent Avatar in Chat
                    let avatarHtml = '';
                    if (!m.isOwn) {
                        const avatarSrc = this.state.agentAvatar;
                        avatarHtml = `
                            <div class="qc-msg-avatar">
                                ${avatarSrc
                            ? `<img src="${avatarSrc}">`
                            : `<div style="width:100%;height:100%;display:flex;align-items:center;justify-content:center;">
                                         <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#94a3b8" stroke-width="2"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
                                       </div>`
                        }
                            </div>`;
                    }

                    el.innerHTML = `
                        ${avatarHtml}
                        <div class="qc-msg-bubble">
                            <div class="qc-msg-text">${this.escapeHtml(m.content)}</div>
                            <div class="qc-msg-meta">
                                <span>${time}</span>
                                ${m.isOwn ? tick : ''}
                            </div>
                        </div>`;
                }
                c.appendChild(el);
            });

            if (this.state.typingUsers.size > 0) {
                const t = document.createElement('div'); t.className = 'qc-msg other';
                const avatarSrc = this.state.agentAvatar;
                t.innerHTML = `
                    <div class="qc-msg-avatar">
                         ${avatarSrc
                    ? `<img src="${avatarSrc}">`
                    : `<div style="width:100%;height:100%;display:flex;align-items:center;justify-content:center;">...</div>`
                }
                    </div>
                    <div class="qc-msg-bubble" style="padding:8px 14px; min-width: 50px;">
                        <div class="qc-typing">
                            <div class="qc-typing-dot"></div>
                            <div class="qc-typing-dot"></div>
                            <div class="qc-typing-dot"></div>
                        </div>
                    </div>`;
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
            this.toggleEmojiPicker();
        },
        showProfileBanner: function() { document.getElementById('qc-profile-banner')?.classList.add('show'); },
        hideProfileBanner: function() { document.getElementById('qc-profile-banner')?.classList.remove('show'); },
    };

    if (document.readyState === 'loading') document.addEventListener('DOMContentLoaded', () => QuickConnect.init());
    else QuickConnect.init();
    window.QuickConnect = QuickConnect;

})(window, document);