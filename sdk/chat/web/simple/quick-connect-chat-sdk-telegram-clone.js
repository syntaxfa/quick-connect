/**
 * Quick Connect Chat SDK - TELEGRAM WEB A CLONE (Perfected UI)
 * Version: 3.1.0 (Bubble Fixes)
 */

(function(window, document) {
    'use strict';

    const QuickConnect = {
        config: {
            managerUrl: 'http://localhost:2531',
            chatUrl: 'ws://localhost:2530/chats/clients',
            chatApiUrl: 'http://localhost:2530',
            position: 'bottom-right', // right or left
            theme: 'telegram-dark',
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
            fontLink.href = 'https://fonts.googleapis.com/css2?family=Vazirmatn:wght@100;300;400;500;600;700&display=swap';
            fontLink.rel = 'stylesheet';
            document.head.appendChild(fontLink);

            const style = document.createElement('style');
            style.textContent = `
                :root {
                    --qc-font: 'Vazirmatn', sans-serif;
                    /* Telegram Dark Theme Palette */
                    --tg-bg: #212121;
                    --tg-header-bg: #212121;
                    --tg-secondary-bg: #0f0f0f;
                    --tg-msg-own-bg: #8774e1;
                    --tg-msg-peer-bg: #212121;
                    --tg-text-white: #ffffff;
                    --tg-text-secondary: #aaaaaa;
                    --tg-text-meta: rgba(255, 255, 255, 0.6);
                    --tg-border: #0f0f0f;
                    --tg-input-bg: #181818;
                    --tg-hover: rgba(255, 255, 255, 0.08);
                    --tg-shadow: 0 1px 2px rgba(0,0,0,0.3);
                    --qc-shadow-lg: 0 10px 40px rgba(0,0,0,0.5);
                }

                #qc-widget * { box-sizing: border-box; margin: 0; padding: 0; outline: none; -webkit-tap-highlight-color: transparent; }

                #qc-widget {
                    position: fixed;
                    ${this.config.position.includes('right') ? 'right' : 'left'}: 20px;
                    bottom: 20px;
                    z-index: 999999;
                    font-family: var(--qc-font);
                    direction: rtl;
                    pointer-events: none;
                    display: flex;
                    flex-direction: column;
                    align-items: flex-end;
                }
                #qc-widget > * { pointer-events: auto; }

                /* --- Launcher Button --- */
                #qc-btn {
                    width: 56px; height: 56px;
                    border-radius: 50%;
                    border: none;
                    background: var(--tg-msg-own-bg);
                    color: white;
                    cursor: pointer;
                    display: flex; align-items: center; justify-content: center;
                    box-shadow: 0 4px 12px rgba(0,0,0,0.3);
                    transition: transform 0.2s cubic-bezier(0.33, 1, 0.68, 1);
                    z-index: 100;
                    position: relative;
                }
                #qc-btn:hover { transform: scale(1.05); }
                #qc-btn:active { transform: scale(0.95); }
                #qc-btn.hidden { transform: scale(0); opacity: 0; pointer-events: none; }

                #qc-badge {
                    position: absolute; top: 0; right: 0;
                    background: #ef4444; color: white;
                    min-width: 20px; height: 20px; border-radius: 10px;
                    font-size: 11px; font-weight: 700;
                    display: flex; align-items: center; justify-content: center;
                    border: 2px solid var(--tg-bg);
                    padding: 0 4px;
                }

                /* --- Main Window --- */
                #qc-window {
                    width: 380px; height: 600px;
                    max-height: calc(100vh - 100px);
                    background: #0f0f0f;
                    background-image: url("data:image/svg+xml,%3Csvg width='40' height='40' viewBox='0 0 40 40' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='%23181818' fill-opacity='0.4' fill-rule='evenodd'%3E%3Cpath d='M0 40L40 0H20L0 20M40 40V20L20 40'/%3E%3C/g%3E%3C/svg%3E");
                    border-radius: 16px;
                    box-shadow: var(--qc-shadow-lg);
                    display: flex; flex-direction: column;
                    overflow: hidden;
                    opacity: 0;
                    transform-origin: bottom ${this.config.position.includes('right') ? 'right' : 'left'};
                    transform: scale(0.9) translateY(20px);
                    transition: all 0.25s cubic-bezier(0.2, 0.8, 0.2, 1);
                    position: absolute;
                    bottom: 70px;
                    ${this.config.position.includes('right') ? 'right' : 'left'}: 0;
                    visibility: hidden;
                }
                #qc-window.open {
                    opacity: 1; transform: scale(1) translateY(0);
                    visibility: visible;
                }

                /* --- Header --- */
                #qc-header {
                    padding: 8px 16px;
                    background: var(--tg-header-bg);
                    display: flex; align-items: center; justify-content: space-between;
                    height: 56px; flex-shrink: 0;
                    cursor: pointer;
                    box-shadow: 0 1px 2px rgba(0,0,0,0.1);
                    z-index: 10;
                }
                #qc-header-info { display: flex; align-items: center; gap: 10px; flex: 1; overflow: hidden; }
                #qc-back-btn {
                    background: transparent; border: none; color: var(--tg-text-secondary);
                    cursor: pointer; padding: 8px; margin-right: -8px; border-radius: 50%;
                    display: flex; align-items: center; justify-content: center;
                }
                #qc-back-btn:hover { background: var(--tg-hover); color: var(--tg-text-white); }

                #qc-avatar {
                    width: 40px; height: 40px;
                    background: linear-gradient(135deg, #ff9a9e 0%, #fecfef 99%, #fecfef 100%);
                    border-radius: 50%;
                    display: flex; align-items: center; justify-content: center;
                    color: #555; font-weight: bold; font-size: 16px;
                    flex-shrink: 0;
                    position: relative;
                }
                #qc-avatar img { width: 100%; height: 100%; object-fit: cover; border-radius: 50%; }

                .qc-user-meta { display: flex; flex-direction: column; justify-content: center; overflow: hidden; }
                #qc-header-title { color: var(--tg-text-white); font-size: 16px; font-weight: 600; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
                #qc-status-text { color: var(--tg-text-secondary); font-size: 13px; margin-top: 2px; }
                #qc-status-text.online { color: #8774e1; }

                /* --- Messages Area (THE FIX) --- */
                #qc-messages {
                    flex: 1; overflow-y: auto; padding: 8px 10px;
                    display: flex; flex-direction: column; gap: 2px;
                    scroll-behavior: smooth;
                }

                .qc-msg {
                    display: flex;
                    width: 100%;
                    flex-direction: column; /* Important for alignment */
                }

                .qc-msg.own { align-items: flex-end; }
                .qc-msg.other { align-items: flex-start; }

                .qc-msg-bubble {
                    position: relative;
                    max-width: 75%;
                    min-width: 60px; /* Ensure space for time */
                    padding: 6px 8px 6px 10px;
                    border-radius: 12px;
                    font-size: 15px;
                    line-height: 1.5;
                    word-wrap: break-word;
                    word-break: break-word; /* Prevents overflow */
                    box-shadow: 0 1px 1px rgba(0,0,0,0.1);
                    text-align: right; /* Default for Persian */

                    /* The magic for tight wrapping */
                    display: inline-block;
                }

                /* Peer Message */
                .qc-msg.other .qc-msg-bubble {
                    background: var(--tg-msg-peer-bg);
                    color: var(--tg-text-white);
                    border-bottom-left-radius: 4px;
                }
                /* Tail for Peer */
                .qc-msg.other .qc-msg-bubble::before {
                    content: ""; position: absolute; bottom: 0; left: -6px;
                    width: 12px; height: 12px;
                    background: radial-gradient(circle at top left, transparent 10px, var(--tg-msg-peer-bg) 11px);
                    clip-path: polygon(50% 50%, 100% 50%, 100% 100%, 50% 100%);
                }

                /* Own Message */
                .qc-msg.own .qc-msg-bubble {
                    background: var(--tg-msg-own-bg);
                    color: var(--tg-text-white);
                    border-bottom-right-radius: 4px;
                }
                /* Tail for Own */
                .qc-msg.own .qc-msg-bubble::before {
                    content: ""; position: absolute; bottom: 0; right: -6px;
                    width: 12px; height: 12px;
                    background: radial-gradient(circle at top right, transparent 10px, var(--tg-msg-own-bg) 11px);
                    clip-path: polygon(0 50%, 50% 50%, 50% 100%, 0 100%);
                }

                /* Message Content & Meta Structure */
                .qc-msg-text {
                    display: inline; /* Lets text flow naturally */
                    white-space: pre-wrap;
                }

                .qc-msg-meta {
                    float: left; /* For RTL (Persian), float left. If LTR, should float right */
                    margin-top: 6px;
                    margin-right: 8px;
                    position: relative;
                    top: 3px;
                    display: inline-flex;
                    align-items: center;
                    gap: 3px;
                    font-size: 11px;
                    color: var(--tg-text-meta);
                    user-select: none;
                    vertical-align: bottom;
                    height: 14px;
                }

                /* Adjust float for LTR text if detected, but usually keep left for RTL bubbles */
                /* You can force right float if needed, but standard Telegram RTL puts time on left */

                .qc-msg.own .qc-msg-meta { color: rgba(255,255,255,0.75); }

                .qc-ticks { display: flex; }
                .qc-ticks svg { width: 14px; height: 14px; fill: currentColor; }

                /* System Messages */
                .qc-msg.system { align-items: center; margin: 10px 0; }
                .qc-msg.system .qc-msg-bubble {
                    background: rgba(0,0,0,0.3); color: #aaa;
                    font-size: 12px; padding: 4px 12px; border-radius: 16px;
                    box-shadow: none; max-width: 90%; text-align: center;
                }
                .qc-msg.system .qc-msg-bubble::before { display: none; }

                /* --- Input Area --- */
                #qc-input-area {
                    background: var(--tg-header-bg);
                    padding: 8px 12px;
                    display: flex; align-items: flex-end; gap: 8px;
                    border-top: 1px solid var(--tg-border);
                }

                #qc-input-wrapper {
                    flex: 1;
                    background: var(--tg-input-bg);
                    border-radius: 20px;
                    display: flex; align-items: center;
                    padding: 8px 12px;
                    min-height: 42px;
                }

                #qc-input {
                    width: 100%;
                    background: transparent; border: none;
                    color: var(--tg-text-white);
                    font-family: inherit; font-size: 15px;
                    resize: none;
                    max-height: 120px;
                    overflow-y: auto;
                    line-height: 20px;
                }
                #qc-input::placeholder { color: #707579; }

                .qc-icon-btn {
                    background: transparent; border: none;
                    color: #707579; cursor: pointer;
                    padding: 6px; border-radius: 50%;
                    display: flex; align-items: center; justify-content: center;
                    transition: color 0.2s;
                }
                .qc-icon-btn:hover { color: var(--tg-msg-own-bg); background: var(--tg-hover); }

                #qc-send-btn {
                    color: var(--tg-msg-own-bg);
                    width: 42px; height: 42px;
                    display: flex; align-items: center; justify-content: center;
                    border-radius: 50%;
                    transition: transform 0.1s;
                }
                #qc-send-btn:hover { background: rgba(135, 116, 225, 0.1); }
                #qc-send-btn:active { transform: scale(0.9); }

                /* --- Emoji Picker --- */
                #qc-emoji-picker {
                    height: 0; overflow: hidden; background: var(--tg-header-bg);
                    transition: height 0.25s cubic-bezier(0.2, 0.8, 0.2, 1);
                    border-top: 1px solid var(--tg-border);
                }
                #qc-emoji-picker.show { height: 250px; }
                #qc-emoji-grid {
                    display: grid; grid-template-columns: repeat(8, 1fr);
                    padding: 10px; gap: 5px; overflow-y: auto; height: 100%;
                }
                .qc-emoji-btn {
                    font-size: 24px; background: transparent; border: none; cursor: pointer;
                    padding: 5px; border-radius: 6px;
                }
                .qc-emoji-btn:hover { background: var(--tg-hover); }

                /* --- Profile Banner --- */
                #qc-profile-banner {
                    background: #2b2b2b; padding: 12px; border-bottom: 1px solid #000;
                    display: none;
                }
                #qc-profile-banner.show { display: block; }
                #qc-profile-form input {
                    width: 100%; background: #181818; border: 1px solid #333;
                    color: white; padding: 8px; margin: 4px 0; border-radius: 8px;
                }
                #qc-profile-submit {
                    width: 100%; background: var(--tg-msg-own-bg); color: white;
                    border: none; padding: 8px; border-radius: 8px; margin-top: 4px; cursor: pointer;
                }

                /* Scrollbar */
                ::-webkit-scrollbar { width: 5px; }
                ::-webkit-scrollbar-track { background: transparent; }
                ::-webkit-scrollbar-thumb { background: rgba(255,255,255,0.15); border-radius: 3px; }
                ::-webkit-scrollbar-thumb:hover { background: rgba(255,255,255,0.25); }

                @keyframes qc-fadeUp { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
                .qc-anim-in { animation: qc-fadeUp 0.2s ease-out; }

                /* Typing Animation */
                .qc-typing-dot {
                    width: 5px; height: 5px; background: #aaa; border-radius: 50%;
                    display: inline-block; margin: 0 1px;
                    animation: qc-wave 1.3s linear infinite;
                }
                .qc-typing-dot:nth-child(2) { animation-delay: -1.1s; }
                .qc-typing-dot:nth-child(3) { animation-delay: -0.9s; }
                @keyframes qc-wave { 0%, 60%, 100% { transform: translateY(0); } 30% { transform: translateY(-4px); } }

                @media (max-width: 480px) {
                    #qc-window {
                        width: 100%; height: 100%; max-height: 100%;
                        bottom: 0; right: 0; left: 0; border-radius: 0;
                    }
                    #qc-btn { bottom: 20px; right: 20px; }
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
                    <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path>
                    </svg>
                </button>

                <div id="qc-window">
                    <div id="qc-header">
                        <div id="qc-header-info">
                            <button id="qc-back-btn" style="display:none">
                                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M19 12H5M12 19l-7-7 7-7"/></svg>
                            </button>
                            <div id="qc-avatar">
                                <img id="qc-avatar-img" src="" style="display:none">
                                <span id="qc-avatar-text">?</span>
                            </div>
                            <div class="qc-user-meta">
                                <div id="qc-header-title">پشتیبانی</div>
                                <div id="qc-status-text">در انتظار شبکه...</div>
                            </div>
                        </div>
                        <div style="display:flex; gap:8px;">
                            <button class="qc-icon-btn" style="color:#aaa">
                                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="11" cy="11" r="8"/><path d="M21 21l-4.35-4.35"/></svg>
                            </button>
                            <button class="qc-icon-btn" style="color:#aaa" id="qc-close-btn">
                                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                            </button>
                        </div>
                    </div>

                    <div id="qc-profile-banner">
                        <div style="display:flex; justify-content:space-between; color:#ddd; margin-bottom:8px; font-size:12px;">
                            <span>شروع گفتگو</span>
                            <span id="qc-profile-close" style="cursor:pointer">✕</span>
                        </div>
                        <div id="qc-profile-form">
                            <input type="text" id="qc-profile-name" placeholder="نام و نام خانوادگی">
                            <input type="email" id="qc-profile-email" placeholder="آدرس ایمیل (اختیاری)">
                            <button id="qc-profile-submit">ثبت و شروع</button>
                        </div>
                    </div>

                    <div id="qc-messages"></div>

                    <div id="qc-input-area">
                        <button class="qc-icon-btn" style="padding:8px">
                            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#707579" stroke-width="1.5"><path d="M21.44 11.05l-9.19 9.19a6 6 0 0 1-8.49-8.49l9.19-9.19a4 4 0 0 1 5.66 5.66l-9.2 9.19a2 2 0 0 1-2.83-2.83l8.49-8.48"/></svg>
                        </button>
                        <div id="qc-input-wrapper">
                            <textarea id="qc-input" placeholder="پیام..." rows="1"></textarea>
                            <button id="qc-emoji-toggle" class="qc-icon-btn" style="margin-right:4px">
                                <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#707579" stroke-width="1.5"><circle cx="12" cy="12" r="10"/><path d="M8 14s1.5 2 4 2 4-2 4-2"/><line x1="9" y1="9" x2="9.01" y2="9"/><line x1="15" y1="9" x2="15.01" y2="9"/></svg>
                            </button>
                        </div>
                        <button id="qc-send-btn">
                            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right:2px; margin-top:2px">
                                <line x1="22" y1="2" x2="11" y2="13"></line>
                                <polygon points="22 2 15 22 11 13 2 9 22 2"></polygon>
                            </svg>
                        </button>
                    </div>
                    <div id="qc-emoji-picker">
                        <div id="qc-emoji-grid"></div>
                    </div>
                </div>
            `;
            document.body.appendChild(container);

            const emojis = ['👍', '👋', '😊', '😂', '🤔', '❤️', '🔥', '🎉', '🙏', '👀', '✅', '🚀', '💯', '✨', '💪', '😍', '😭', '🤣', '🥺', '😐', '👌', '😘', '😁', '😎'];
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
            document.getElementById('qc-back-btn').onclick = () => this.closeChat();
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
                    const txt = document.getElementById('qc-avatar-text');
                    const info = data.support_info || {};

                    title.textContent = (data.assigned_support_id && info.fullname) ? info.fullname : 'پشتیبانی آنلاین';

                    if (info.avatar) {
                        img.src = info.avatar; img.style.display = 'block'; txt.style.display = 'none';
                    } else {
                        img.style.display = 'none'; txt.style.display = 'block';
                        txt.textContent = info.fullname ? info.fullname.charAt(0).toUpperCase() : 'S';
                    }

                    if (info.last_online_at) {
                        const last = new Date(info.last_online_at);
                        const diff = (new Date() - last) / 60000;
                        if (diff < 5 && diff >= 0) {
                            this.state.agentStatusText = 'آنلاین'; this.state.agentIsOnline = true;
                        } else {
                            this.state.agentStatusText = 'آخرین بازدید به تازگی';
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
                        isOwn: m.sender_id === this.state.userId, status: 'read'
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
                this.addSystemMessage('اطلاعات شما با موفقیت ثبت شد');
            } catch (e) {}
        },

        connectWebSocket: function() {
            if (this.state.ws && this.state.ws.readyState === WebSocket.OPEN) return;
            try { this.state.ws = new WebSocket(this.config.chatUrl, [this.state.token]); } catch(e){ return; }

            this.state.ws.onopen = () => {
                this.state.isConnected = true;
                this.updateConnectionStatus();
                this.startHeartbeat();
                if(this.state.messages.length === 0 && this.state.historyLoaded) this.addSystemMessage('به پشتیبانی خوش آمدید');
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
            if(!txt) return;
            if (!this.state.isConnected) {
                txt.textContent = 'درحال اتصال...'; txt.className = '';
            } else if (this.state.typingUsers.size > 0) {
                txt.textContent = 'درحال نوشتن...'; txt.className = 'online';
            } else {
                txt.textContent = this.state.agentStatusText;
                txt.className = this.state.agentIsOnline ? 'online' : '';
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
            input.style.height = Math.min(input.scrollHeight, 120) + 'px';
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

            const fragment = document.createDocumentFragment();

            // SVG Icons
            const checkIcon = `<svg viewBox="0 0 24 24"><path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41L9 16.17z"/></svg>`;
            const doubleCheckIcon = `<svg viewBox="0 0 24 24"><path d="M18 7l-1.41-1.41-6.34 6.34 1.41 1.41L18 7zm4.24-1.41L11.66 16.17 7.48 12l-1.41 1.41L11.66 19l12-12-1.42-1.41zM.41 13.41L6 19l1.41-1.41L1.83 12 .41 13.41z"/></svg>`;
            const clockIcon = `<svg viewBox="0 0 24 24"><circle cx="12" cy="12" r="10" fill="none" stroke="currentColor" stroke-width="2"/><path d="M12 6v6l4 2" stroke="currentColor" stroke-width="2" stroke-linecap="round"/></svg>`;

            this.state.messages.forEach((m) => {
                const el = document.createElement('div');

                if (m.isSystem) {
                    el.className = 'qc-msg system qc-anim-in';
                    el.innerHTML = `<div class="qc-msg-bubble">${this.escapeHtml(m.content)}</div>`;
                } else {
                    const isOwn = m.isOwn;
                    const time = new Date(m.created_at).toLocaleTimeString('fa-IR', {hour:'2-digit', minute:'2-digit'});

                    let statusIcon = '';
                    if (isOwn) {
                        if (m.status === 'pending') statusIcon = clockIcon;
                        else if (m.status === 'read') statusIcon = doubleCheckIcon;
                        else statusIcon = checkIcon;
                    }

                    el.className = `qc-msg ${isOwn ? 'own' : 'other'} qc-anim-in`;

                    // IMPORTANT: dir="auto" handles LTR/RTL automatically.
                    // The meta span floats (float: left in Persian CSS).
                    // This causes text to wrap around the time stamp like Telegram.
                    el.innerHTML = `
                        <div class="qc-msg-bubble">
                            <div class="qc-msg-text" dir="auto">
                                ${this.escapeHtml(m.content)}
                                <span class="qc-msg-meta">
                                    <span>${time}</span>
                                    ${isOwn ? `<div class="qc-ticks">${statusIcon}</div>` : ''}
                                </span>
                            </div>
                        </div>
                    `;
                }
                fragment.appendChild(el);
            });

            if (this.state.typingUsers.size > 0) {
                const t = document.createElement('div');
                t.className = 'qc-msg other qc-anim-in';
                t.innerHTML = `
                    <div class="qc-msg-bubble" style="padding: 12px 16px; min-width: auto;">
                        <div class="qc-typing-dot"></div><div class="qc-typing-dot"></div><div class="qc-typing-dot"></div>
                    </div>`;
                fragment.appendChild(t);
            }

            c.innerHTML = '';
            c.appendChild(fragment);

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
        },
        showProfileBanner: function() { document.getElementById('qc-profile-banner')?.classList.add('show'); },
        hideProfileBanner: function() { document.getElementById('qc-profile-banner')?.classList.remove('show'); },
    };

    if (document.readyState === 'loading') document.addEventListener('DOMContentLoaded', () => QuickConnect.init());
    else QuickConnect.init();
    window.QuickConnect = QuickConnect;

})(window, document);