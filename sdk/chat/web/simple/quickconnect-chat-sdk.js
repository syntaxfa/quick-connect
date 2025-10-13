// QuickConnect Chat Widget SDK
// Developer: Your Name/Company
// Version: 1.0.0
// Description: A standalone chat widget inspired by Telegram's theme.

(function() {
// --- Configuration ---
// WebSocket URL for your chat server
const WEBSOCKET_URL = 'ws://localhost:2530/v1/chats/clients';
// Tailwind CSS CDN URL
const TAILWIND_CDN_URL = 'https://cdn.tailwindcss.com';
// Name for the global SDK object (optional, if you want to expose methods later)
const SDK_NAME = 'QuickConnectChat';
// Default chat title
const CHAT_TITLE = 'پشتیبانی آنلاین';
// Default welcome message
const WELCOME_MESSAGE = 'برای شروع گفتگو، پیام خود را تایپ کنید.';
// Default connection messages
const CONNECTING_MESSAGE = 'در حال اتصال به سرور چت...';
const CONNECTED_MESSAGE = 'به سرور چت متصل شدید.';
const CONNECTION_ERROR_MESSAGE = 'خطا در اتصال WebSocket.';
const CONNECTION_CLOSED_MESSAGE = 'اتصال WebSocket قطع شد.';
const RECONNECTING_MESSAGE = ' تلاش برای اتصال مجدد...';
const UNKNOWN_SERVER_MESSAGE_PREFIX = 'پیام ناشناخته از سرور: ';
const PROCESSING_ERROR_MESSAGE = 'خطا در پردازش پیام از سرور.';
const NOT_CONNECTED_MESSAGE = 'اتصال برقرار نیست. لطفاً منتظر بمانید.';
// Input placeholder
const INPUT_PLACEHOLDER = 'پیام خود را بنویسید...';


// --- Emojis ---
const emojis = [
'😀', '😃', '😄', '😁', '😆', '😅', '😂', '🤣', '😊', '😇',
'🙂', '🙃', '😉', '😌', '😍', '🥰', '😘', '😗', '😙', '😚',
'😋', '😛', '😜', '🤪', '🤨', '🧐', '🤓', '😎', '🤩', '🥳',
'😏', '😒', '😞', '😔', '😟', '😕', '🙁', '☹️', '😣', '😖',
'😫', '😩', '🥺', '😢', '😭', '😤', '😠', '😡', '🤬', '🤯',
'😳', '🥵', '🥶', '😱', '😨', '😰', '😥', '😓', '🤗', '🤔',
'🤭', '🤫', '🤥', '😶', '😐', '😑', '😬', '🙄', '😯', '😦',
'😧', '😮', '😲', '🥱', '😴', '🤤', '😪', '😵', '🤐', '🥴',
'🤢', '🤮', '🤧', '😷', '🤒', '🤕', '🤑', '🤠', '❤️', '👍',
'👌', '🙏', '👋', '🎉', '🔥', '⭐', '☀️', '☁️', '☂️'
];

// --- CSS Styles (Improved with Telegram Theme & Background Pattern) ---
const customCSS = `
:root {
--tg-bg-primary: #212121;
--tg-bg-secondary: #0F0F0F;
--tg-text-primary: #FFFFFF;
--tg-text-secondary: #AAAAAA;
--tg-text-system: #A3B6C7;
--tg-bubble-sent-bg: #766AC8;
--tg-bubble-sent-text: #FFFFFF;
--tg-bubble-sent-timestamp: rgba(255, 255, 255, 0.533);
--tg-bubble-received-bg: #182533;
--tg-bubble-received-text: #FFFFFF;
--tg-bubble-received-timestamp: #707C88;
--tg-icon-primary-bg: #8774E1;
--tg-icon-primary-stroke: #FFFFFF;
--tg-input-bg: #17212B;
--tg-input-text: #D1D5DB;
--tg-input-border: #303030;
--tg-input-placeholder: #5E6C7A;
--tg-input-icon-color: #AAAAAA;
--tg-button-send-bg: #8774E1;
--tg-button-send-icon: #FFFFFF;
--tg-scrollbar-thumb: #2A3F53;
--tg-scrollbar-track: var(--tg-bg-primary);
--tg-input-scrollbar-track: var(--tg-input-bg);
--tg-emoji-picker-bg: #1F2D3A;
--tg-emoji-picker-border: #303030;
--tg-emoji-item-hover-bg: #2A3F53;
--tg-pattern-color: rgba(170, 170, 170, 0.05);
--tg-pattern-color-darker: rgba(170, 170, 170, 0.03);

--border-radius-messages: 0.9375rem;
--border-radius-messages-small: 0.375rem;
--message-text-size: 0.95rem;
}

.chat-widget-container, .chat-widget-container * {
direction: rtl;
text-align: right;
font-family: "Roboto", "Vazirmatn", -apple-system, BlinkMacSystemFont, 'Segoe UI', Oxygen, Ubuntu, Cantarell, 'Fira Sans', 'Droid Sans', 'Helvetica Neue', sans-serif;
}

.telegram-chat-bg {
background-color: var(--tg-bg-primary);
position: relative;
}

.message-bubble {
max-width: 75%;
word-wrap: break-word;
padding: 8px 12px;
font-size: var(--message-text-size);
line-height: 1.375;
border-radius: var(--border-radius-messages);
display: flex;
flex-direction: column;
position: relative;
box-shadow: 0 1px 2px rgba(0,0,0,0.15);
margin-bottom: 0.375rem;
z-index: 1;
}
.Message.last-in-group {
margin-bottom: 0.625rem;
}

.message-bubble::after {
content: "";
position: absolute;
bottom: 0px;
width: 0;
height: 0;
border: 7px solid transparent;
}

.message-content {
margin-bottom: 4px;
white-space: pre-wrap;
}

.message-timestamp {
font-size: 0.7rem;
align-self: flex-end;
margin-right: 8px;
opacity: 0.8;
}

.message-sent {
background-color: var(--tg-bubble-sent-bg);
color: var(--tg-bubble-sent-text);
margin-left: auto;
border-bottom-right-radius: var(--border-radius-messages-small);
}
.message-sent::after {
right: -9px;
border-left-color: var(--tg-bubble-sent-bg);
border-right-width: 0;
border-top-width: 0;
}
.message-sent .message-timestamp { color: var(--tg-bubble-sent-timestamp); }

.message-received {
background-color: var(--tg-bubble-received-bg);
color: var(--tg-bubble-received-text);
margin-right: auto;
border-bottom-left-radius: var(--border-radius-messages-small);
}
.message-received::after {
left: -9px;
border-right-color: var(--tg-bubble-received-bg);
border-left-width: 0;
border-top-width: 0;
}
.message-received .message-timestamp { color: var(--tg-bubble-received-timestamp); }

.message-pending { opacity: 0.7; }
.message-confirmed { opacity: 1; }

.chat-icon-rtl {
position:fixed; left: 1.5rem; right: auto !important; bottom: 1.5rem;
background-color: var(--tg-icon-primary-bg);
transition: background-color 0.3s ease, transform 0.3s ease, opacity 0.3s ease, visibility 0.3s ease;
z-index: 2147483640; display: flex !important;
width: 56px;
height: 56px;
}
.chat-icon-rtl.qc-hidden {
opacity: 0; visibility: hidden; transform: scale(0.8);
}
.chat-icon-rtl svg { stroke: var(--tg-icon-primary-stroke); width: 28px; height: 28px; }
.chat-icon-rtl:hover { filter: brightness(110%); }

.chat-window-rtl {
position:fixed;
width: 24rem; height: 600px;
left: 1.5rem; right: auto !important; bottom: 6rem;
border: 1px solid var(--tg-input-border);
z-index: 2147483630;
border-radius: 0.75rem;
}

.chat-header {
background-color: var(--tg-bg-secondary);
color: var(--tg-text-primary);
border-bottom: 1px solid var(--tg-bg-primary);
border-top-left-radius: inherit;
border-top-right-radius: inherit;
}
.chat-header button:hover { background-color: rgba(255,255,255,0.1); }
.chat-header button svg { stroke: var(--tg-text-secondary); }

#qc-chat-messages {
scrollbar-width: thin;
scrollbar-color: var(--tg-scrollbar-thumb) var(--tg-scrollbar-track);
padding: 12px;
position: relative;
}
#qc-chat-messages::before {
content: '';
position: absolute;
top: 0;
left: 0;
width: 100%;
height: 100%;
background-image: url("data:image/svg+xml,%3Csvg width='20' height='20' viewBox='0 0 20 20' xmlns='http://www.w3.org/2000/svg'%3E%3Ccircle cx='3' cy='3' r='1' fill='rgba(170,170,170,0.07)'/%3E%3C/svg%3E");
background-repeat: repeat;
pointer-events: none;
z-index: 0;
}

#qc-chat-messages::-webkit-scrollbar { width: 6px; }
#qc-chat-messages::-webkit-scrollbar-track { background: var(--tg-scrollbar-track); border-radius: 3px; }
#qc-chat-messages::-webkit-scrollbar-thumb { background-color: var(--tg-scrollbar-thumb); border-radius: 3px; }
#qc-chat-messages::-webkit-scrollbar-thumb:hover { background-color: #375069; }

#qc-chat-input-textarea {
background-color: var(--tg-input-bg); color: var(--tg-input-text);
border: 1px solid var(--tg-input-border);
border-radius: 20px;
transition: box-shadow 0.2s ease, border-color 0.2s ease; resize: none;
overflow-y: auto; min-height: 42px;
max-height: 120px; padding: 10px 18px;
line-height: 1.45; flex-grow: 1;
scrollbar-width: thin; scrollbar-color: var(--tg-scrollbar-thumb) var(--tg-input-scrollbar-track);
font-size: 0.9rem;
}
#qc-chat-input-textarea::-webkit-scrollbar { width: 5px; }
#qc-chat-input-textarea::-webkit-scrollbar-track { background: var(--tg-input-scrollbar-track); border-radius: 2.5px;}
#qc-chat-input-textarea::-webkit-scrollbar-thumb { background-color: var(--tg-scrollbar-thumb); border-radius: 2.5px; }
#qc-chat-input-textarea::-webkit-scrollbar-thumb:hover { background-color: #375069; }
#qc-chat-input-textarea::placeholder { color: var(--tg-input-placeholder); }
#qc-chat-input-textarea:focus {
border-color: var(--tg-button-send-bg);
box-shadow: 0 0 0 1.5px var(--tg-button-send-bg);
}

#qc-chat-form button#qc-send-btn {
background-color: var(--tg-button-send-bg);
transition: background-color 0.3s ease, transform 0.2s ease;
align-self: flex-end; margin-bottom: 3px; padding: 0.6rem;
border-radius: 50%; width: 40px; height: 40px;
display: flex; align-items: center; justify-content: center;
}
#qc-chat-form button#qc-send-btn svg {
fill: var(--tg-button-send-icon);
width: 1.2rem; height: 1.2rem;
}
#qc-chat-form button#qc-send-btn:hover {
filter: brightness(115%); transform: scale(1.05);
}

.chat-footer {
background-color: var(--tg-bg-secondary);
border-top: 1px solid var(--tg-bg-primary);
padding: 10px 12px;
position: relative;
border-bottom-left-radius: inherit;
border-bottom-right-radius: inherit;
}
#qc-chat-form { align-items: flex-end; display: flex; gap: 8px; }
.system-message { color: var(--tg-text-system); font-size: 0.8rem; z-index: 1; position: relative; }

.emoji-picker-btn {
padding: 0.5rem; margin-left: 0; margin-right: 0.5rem;
align-self: flex-end; margin-bottom: 4px; cursor: pointer;
}
.emoji-picker-btn svg {
width: 1.5rem; height: 1.5rem;
fill: var(--tg-input-icon-color); transition: fill 0.2s ease, transform 0.2s ease;
}
.emoji-picker-btn:hover svg { fill: var(--tg-text-primary); transform: scale(1.1); }

#qc-emoji-picker-container {
position: absolute;
bottom: calc(100% + 8px);
right: 12px;
width: 300px;
max-height: 220px;
overflow-y: auto;
background-color: var(--tg-emoji-picker-bg);
border: 1px solid var(--tg-emoji-picker-border);
border-radius: 10px;
box-shadow: 0 6px 15px rgba(0,0,0,0.35);
padding: 10px;
display: grid;
grid-template-columns: repeat(auto-fill, minmax(36px, 1fr));
gap: 5px;
z-index: 2147483635;
scrollbar-width: thin;
scrollbar-color: var(--tg-scrollbar-thumb) var(--tg-emoji-picker-bg);
}
#qc-emoji-picker-container.hidden { display: none !important; }
#qc-emoji-picker-container::-webkit-scrollbar { width: 5px; }
#qc-emoji-picker-container::-webkit-scrollbar-track { background: var(--tg-emoji-picker-bg); }
#qc-emoji-picker-container::-webkit-scrollbar-thumb { background-color: var(--tg-scrollbar-thumb); border-radius: 3px; }

.emoji-item {
font-size: 1.3rem; padding: 5px; text-align: center; cursor: pointer;
border-radius: 6px;
transition: background-color 0.2s ease, transform 0.1s ease;
}
.emoji-item:hover {
background-color: var(--tg-emoji-item-hover-bg);
transform: scale(1.1);
}

@media (max-width: 640px) {
.chat-window-rtl {
width: calc(100% - 1rem);
height: calc(100% - 1rem);
top: 0.5rem;
left: 0.5rem;
right: 0.5rem !important;
bottom: 0.5rem;
border-radius: 0.75rem;
box-shadow: 0 4px 20px rgba(0,0,0,0.3);
}
.chat-icon-rtl {
bottom: 1rem; left: 1rem; transform: scale(0.9);
width: 52px; height: 52px;
}
.chat-icon-rtl svg { width: 26px; height: 26px; }
.chat-icon-rtl.qc-hidden-mobile {
opacity: 0 !important; visibility: hidden !important; transform: scale(0.8) !important;
}
.message-bubble { font-size: 0.9rem; max-width: 88%; }
#qc-chat-input-textarea {
min-height: 40px; max-height: 100px; padding: 9px 14px; font-size: 0.95rem;
}
#qc-emoji-picker-container {
width: calc(100% - 1rem); left: 0.5rem; right: 0.5rem; bottom: calc(100% + 4px);
max-height: 180px;
}
.chat-footer { padding: 8px 10px; }
#qc-chat-form button#qc-send-btn {
width: 36px; height: 36px; padding: 0.5rem;
}
#qc-chat-form button#qc-send-btn svg { width: 1.1rem; height: 1.1rem; }
.emoji-picker-btn svg { width: 1.4rem; height: 1.4rem; }
}

.chat-window-rtl {
transition: transform 0.3s cubic-bezier(0.33, 1, 0.68, 1), opacity 0.3s cubic-bezier(0.33, 1, 0.68, 1), visibility 0.3s cubic-bezier(0.33, 1, 0.68, 1);
}
body.is-ios .chat-window-rtl {
transition-duration: 0.4s;
}
`;

// --- DOM Elements (to be assigned after creation) ---
let chatIconEl, chatWindowEl, messagesContainer, chatForm, chatInputTextarea, closeChatBtn, emojiPickerBtn, emojiPickerContainer;
let socket;
let pendingMessages = {};
let sdkInitialized = false; // Flag to prevent multiple initializations

// --- Helper Functions ---
function injectStyles(cssString) {
if (document.getElementById('qc-chat-widget-styles')) return; // Prevent re-injecting styles
const styleElement = document.createElement('style');
styleElement.id = 'qc-chat-widget-styles';
styleElement.type = 'text/css';
styleElement.textContent = cssString;
document.head.appendChild(styleElement);
}

function createChatWidgetDOM() {
if (document.getElementById('qc-chat-icon') || document.getElementById('qc-chat-window')) {
// Widget already exists, perhaps from a previous script run or manual addition
console.warn(`${SDK_NAME}: Widget DOM elements already exist. Skipping creation.`);
// Attempt to re-assign elements if they exist
chatIconEl = document.getElementById('qc-chat-icon');
chatWindowEl = document.getElementById('qc-chat-window');
messagesContainer = document.getElementById('qc-chat-messages');
chatForm = document.getElementById('qc-chat-form');
chatInputTextarea = document.getElementById('qc-chat-input-textarea');
closeChatBtn = document.getElementById('qc-close-chat-btn');
emojiPickerBtn = document.getElementById('qc-emoji-picker-btn');
emojiPickerContainer = document.getElementById('qc-emoji-picker-container');
return;
}

// Create Chat Icon
chatIconEl = document.createElement('div');
chatIconEl.id = 'qc-chat-icon';
chatIconEl.className = 'chat-icon-rtl p-3.5 rounded-full shadow-xl cursor-pointer transition-colors flex items-center justify-center';
chatIconEl.innerHTML = `
<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2" class="w-7 h-7">
    <path stroke-linecap="round" stroke-linejoin="round" d="M21.75 6.75v10.5a2.25 2.25 0 0 1-2.25 2.25h-15a2.25 2.25 0 0 1-2.25-2.25V6.75m19.5 0A2.25 2.25 0 0 0 19.5 4.5h-15a2.25 2.25 0 0 0-2.25 2.25m19.5 0v.243a2.25 2.25 0 0 1-1.07 1.916l-7.5 4.615a2.25 2.25 0 0 1-2.36 0L3.32 8.91A2.25 2.25 0 0 1 2.25 6.993V6.75" />
</svg>
`;
document.body.appendChild(chatIconEl);

// Create Chat Window
chatWindowEl = document.createElement('div');
chatWindowEl.id = 'qc-chat-window';
chatWindowEl.className = 'chat-window-rtl chat-widget-container telegram-chat-bg border rounded-xl shadow-2xl flex flex-col transition-all duration-300 ease-in-out transform translate-y-full opacity-0 overflow-hidden';
chatWindowEl.innerHTML = `
<div class="chat-header flex items-center justify-between p-3.5 rounded-t-xl">
    <h3 class="font-semibold text-md">${CHAT_TITLE}</h3>
    <button id="qc-close-chat-btn" class="p-1.5 rounded-full">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="2.5" stroke="currentColor" class="w-5 h-5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
        </svg>
    </button>
</div>
<div id="qc-chat-messages" class="flex-1 p-3 space-y-3 overflow-y-auto">
    <div class="text-center text-xs py-2 system-message">${WELCOME_MESSAGE}</div>
</div>
<div class="chat-footer">
    <div id="qc-emoji-picker-container"></div>
    <form id="qc-chat-form" class="flex items-end space-x-2 space-x-reverse w-full">
        <button type="button" id="qc-emoji-picker-btn" class="emoji-picker-btn">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="w-6 h-6">
                <path fill-rule="evenodd" d="M12 2.25c-5.385 0-9.75 4.365-9.75 9.75s4.365 9.75 9.75 9.75 9.75-4.365 9.75-9.75S17.385 2.25 12 2.25Zm-2.625 6c-.54 0-.975.435-.975.975v.008c0 .54.435.975.975.975h.008c.54 0 .975-.435.975-.975V9.225A.975.975 0 0 0 9.375 8.25h-.008ZM12 8.25c-.54 0-.975.435-.975.975v.008c0 .54.435.975.975.975h.008c.54 0 .975-.435.975-.975V9.225A.975.975 0 0 0 12 8.25h-.008Zm2.625 0c-.54 0-.975.435-.975.975v.008c0 .54.435.975.975.975h.008c.54 0 .975-.435.975-.975V9.225a.975.975 0 0 0-.975-.975h-.008ZM12 18a6 6 0 0 0 4.682-2.144.75.75 0 0 0-1.034-1.086A4.5 4.5 0 0 1 12 15.75a4.5 4.5 0 0 1-3.648-1.48.75.75 0 0 0-1.034 1.086A6 6 0 0 0 12 18Z" clip-rule="evenodd" />
            </svg>
        </button>
        <textarea id="qc-chat-input-textarea" class="flex-1" placeholder="${INPUT_PLACEHOLDER}" rows="1"></textarea>
        <button type="submit" id="qc-send-btn" class="rounded-full">
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" class="w-5 h-5">
                <path d="M3.478 2.404a.75.75 0 0 0-.926.941l2.432 7.905H13.5a.75.75 0 0 1 0 1.5H4.984l-2.432 7.905a.75.75 0 0 0 .926.94 60.519 60.519 0 0 0 18.445-8.986.75.75 0 0 0 0-1.218A60.517 60.517 0 0 0 3.478 2.404Z" />
            </svg>
        </button>
    </form>
</div>
`;
document.body.appendChild(chatWindowEl);

messagesContainer = document.getElementById('qc-chat-messages');
chatForm = document.getElementById('qc-chat-form');
chatInputTextarea = document.getElementById('qc-chat-input-textarea');
closeChatBtn = document.getElementById('qc-close-chat-btn');
emojiPickerBtn = document.getElementById('qc-emoji-picker-btn');
emojiPickerContainer = document.getElementById('qc-emoji-picker-container');

if (emojiPickerContainer) {
emojiPickerContainer.classList.add('hidden');
}
}

function populateEmojiPicker() {
if (!emojiPickerContainer || emojiPickerContainer.children.length > 0) return; // Populate only once
emojis.forEach(emoji => {
const emojiElement = document.createElement('span');
emojiElement.className = 'emoji-item';
emojiElement.textContent = emoji;
emojiElement.addEventListener('click', () => {
const cursorPos = chatInputTextarea.selectionStart;
const textBefore = chatInputTextarea.value.substring(0, cursorPos);
const textAfter = chatInputTextarea.value.substring(cursorPos);
chatInputTextarea.value = textBefore + emoji + textAfter;

chatInputTextarea.focus();
chatInputTextarea.selectionStart = cursorPos + emoji.length;
chatInputTextarea.selectionEnd = cursorPos + emoji.length;

autoResizeTextarea();
if (emojiPickerContainer) {
emojiPickerContainer.classList.add('hidden');
}
});
emojiPickerContainer.appendChild(emojiElement);
});
}

function setupEventListeners() {
if (!chatIconEl || !chatWindowEl || !closeChatBtn || !emojiPickerBtn || !chatInputTextarea || !chatForm) {
console.error(`${SDK_NAME}: Could not find all required DOM elements to attach listeners.`);
return;
}

chatIconEl.addEventListener('click', () => {
if (!chatWindowEl || !chatIconEl) return;
const isChatWindowHidden = chatWindowEl.classList.contains('translate-y-full');

if (isChatWindowHidden) {
chatWindowEl.classList.remove('translate-y-full', 'opacity-0');
chatWindowEl.classList.add('translate-y-0', 'opacity-100');
chatIconEl.classList.add('qc-hidden');
if (window.innerWidth <= 640) {
chatIconEl.classList.add('qc-hidden-mobile');
}

if(chatInputTextarea) chatInputTextarea.focus();
autoResizeTextarea();
if (!socket || socket.readyState === WebSocket.CLOSED) {
connectWebSocket();
}
} else {
chatWindowEl.classList.add('translate-y-full', 'opacity-0');
chatWindowEl.classList.remove('translate-y-0', 'opacity-100');
chatIconEl.classList.remove('qc-hidden');
if (window.innerWidth <= 640) {
chatIconEl.classList.remove('qc-hidden-mobile');
}
if (emojiPickerContainer) emojiPickerContainer.classList.add('hidden');
}
});

closeChatBtn.addEventListener('click', () => {
if (!chatWindowEl || !chatIconEl) return;
chatWindowEl.classList.add('translate-y-full', 'opacity-0');
chatWindowEl.classList.remove('translate-y-0', 'opacity-100');
chatIconEl.classList.remove('qc-hidden');
if (window.innerWidth <= 640) {
chatIconEl.classList.remove('qc-hidden-mobile');
}
if (emojiPickerContainer) emojiPickerContainer.classList.add('hidden');
});

emojiPickerBtn.addEventListener('click', (event) => {
event.stopPropagation();
if (emojiPickerContainer) {
emojiPickerContainer.classList.toggle('hidden');
if (!emojiPickerContainer.classList.contains('hidden') && chatInputTextarea) {
chatInputTextarea.focus();
}
}
});

document.addEventListener('click', (event) => {
if (emojiPickerContainer && !emojiPickerContainer.classList.contains('hidden')) {
const isClickInsidePicker = emojiPickerContainer.contains(event.target);
const isClickOnPickerButton = emojiPickerBtn.contains(event.target) || event.target === emojiPickerBtn;

if (!isClickInsidePicker && !isClickOnPickerButton) {
emojiPickerContainer.classList.add('hidden');
}
}
});

if(chatInputTextarea) {
chatInputTextarea.addEventListener('input', autoResizeTextarea);
chatInputTextarea.addEventListener('keydown', (e) => {
if (e.key === 'Enter' && !e.shiftKey) {
e.preventDefault();
if(chatForm) chatForm.dispatchEvent(new Event('submit', { cancelable: true }));
}
});
}

if(chatForm) {
chatForm.addEventListener('submit', (e) => {
e.preventDefault();
if(!chatInputTextarea) return;
const messageText = chatInputTextarea.value.trim();
if (messageText && socket && socket.readyState === WebSocket.OPEN) {
const clientMessageId = `client_${Date.now()}_${Math.random().toString(36).substring(2,9)}`;
const messagePayload = { type: 'message', content: messageText, clientMsgId: clientMessageId };
socket.send(JSON.stringify(messagePayload));
addMessageToUI(messageText, 'sent', Date.now(), clientMessageId);
chatInputTextarea.value = '';
autoResizeTextarea();
if (emojiPickerContainer) emojiPickerContainer.classList.add('hidden');
} else if (!socket || socket.readyState !== WebSocket.OPEN) {
addSystemMessage(NOT_CONNECTED_MESSAGE);
if (!socket || socket.readyState === WebSocket.CLOSED) { connectWebSocket(); }
}
});
}
autoResizeTextarea();
}

function autoResizeTextarea() {
if (!chatInputTextarea) return;
chatInputTextarea.style.height = 'auto';
let scrollHeight = chatInputTextarea.scrollHeight;
const maxHeight = parseInt(getComputedStyle(chatInputTextarea).maxHeight);
if (scrollHeight > maxHeight) {
chatInputTextarea.style.height = maxHeight + 'px';
chatInputTextarea.style.overflowY = 'auto';
} else {
chatInputTextarea.style.height = scrollHeight + 'px';
chatInputTextarea.style.overflowY = 'hidden';
}
}

function formatTimestamp(dateValue) {
const date = new Date(dateValue);
return date.toLocaleTimeString('fa-IR', { hour: '2-digit', minute: '2-digit', hour12: false });
}

function addMessageToUI(text, type, timestamp, messageId = null) {
if (!messagesContainer) return;
const initialMessage = messagesContainer.querySelector('.text-center.system-message');
if (initialMessage && initialMessage.textContent.includes(WELCOME_MESSAGE.substring(0,10))) { // Check for part of welcome message
initialMessage.remove();
}
const messageElement = document.createElement('div');
messageElement.className = `message-bubble ${type === 'sent' ? 'message-sent' : 'message-received'}`;
const contentElement = document.createElement('div');
contentElement.className = 'message-content';
contentElement.textContent = text;
const timestampElement = document.createElement('div');
timestampElement.className = 'message-timestamp';
timestampElement.textContent = formatTimestamp(timestamp || Date.now());
messageElement.appendChild(contentElement);
messageElement.appendChild(timestampElement);
if (type === 'sent') {
messageElement.classList.add('message-pending');
if (messageId) {
pendingMessages[messageId] = { element: messageElement, content: text, timestampElement: timestampElement };
}
}
messagesContainer.appendChild(messageElement);
messagesContainer.scrollTop = messagesContainer.scrollHeight;
return messageElement;
}

function addSystemMessage(text) {
if (!messagesContainer) return;
const messageElement = document.createElement('div');
messageElement.className = 'text-center text-xs py-1.5 px-2 system-message';
messageElement.textContent = text;
messagesContainer.appendChild(messageElement);
messagesContainer.scrollTop = messagesContainer.scrollHeight;
}

function connectWebSocket() {
if (socket && (socket.readyState === WebSocket.OPEN || socket.readyState === WebSocket.CONNECTING)) {
return;
}
socket = new WebSocket(WEBSOCKET_URL);
addSystemMessage(CONNECTING_MESSAGE);
socket.onopen = () => { addSystemMessage(CONNECTED_MESSAGE); };
socket.onmessage = (event) => {
try {
const message = JSON.parse(event.data);
if (message.type === 'echo') {
const originalContent = message.content;
let confirmedMsgId = null;
if (message.clientMsgId && pendingMessages[message.clientMsgId]) {
confirmedMsgId = message.clientMsgId;
} else {
for (const msgId in pendingMessages) {
if (pendingMessages[msgId].content === originalContent) {
confirmedMsgId = msgId; break;
}
}
}
if (confirmedMsgId && pendingMessages[confirmedMsgId]) {
const pendingMsg = pendingMessages[confirmedMsgId];
if (pendingMsg.element) {
pendingMsg.element.classList.remove('message-pending');
pendingMsg.element.classList.add('message-confirmed');
}
delete pendingMessages[confirmedMsgId];
}
} else if (message.type === 'message') {
addMessageToUI(message.content, 'received', message.timestamp || Date.now());
} else {
addSystemMessage(`${UNKNOWN_SERVER_MESSAGE_PREFIX} ${message.content || JSON.stringify(message)}`);
}
} catch (error) { console.error('Error processing message from server:', error); addSystemMessage(PROCESSING_ERROR_MESSAGE); }
};
socket.onerror = (error) => { console.error('WebSocket Error:', error); addSystemMessage(CONNECTION_ERROR_MESSAGE); };
socket.onclose = (event) => {
let msgText = CONNECTION_CLOSED_MESSAGE;
if (event.code === 1006) { // Abnormal closure
msgText += RECONNECTING_MESSAGE;
setTimeout(connectWebSocket, 5000); // Reconnect after 5 seconds
}
addSystemMessage(msgText);
};
}

// --- Initialization ---
function initializeSDK() {
if (sdkInitialized) {
console.warn(`${SDK_NAME}: SDK already initialized.`);
return;
}
injectStyles(customCSS);
createChatWidgetDOM();
populateEmojiPicker();
setupEventListeners();
sdkInitialized = true;
console.log(`${SDK_NAME} initialized.`);
}

function loadTailwindAndInit() {
let tailwindLoaded = false;
try {
const testEl = document.createElement('div');
testEl.className = 'hidden';
document.body.appendChild(testEl);
tailwindLoaded = getComputedStyle(testEl).display === 'none';
document.body.removeChild(testEl);
} catch (e) { /* Tailwind not loaded */ }

if (tailwindLoaded || document.querySelector('script[src*="tailwindcss.com"]')) {
initializeSDK();
} else {
const tailwindScript = document.createElement('script');
tailwindScript.src = TAILWIND_CDN_URL;
tailwindScript.onload = initializeSDK;
tailwindScript.onerror = () => {
console.error(`${SDK_NAME}: Failed to load Tailwind CSS from CDN. Widget may not style correctly.`);
initializeSDK(); // Initialize anyway
};
document.head.appendChild(tailwindScript);
}
}

// --- Start Initialization when DOM is ready ---
if (document.readyState === 'loading') {
document.addEventListener('DOMContentLoaded', loadTailwindAndInit);
} else {
loadTailwindAndInit();
}

// Optional: Expose a global object for the SDK if needed later
// window[SDK_NAME] = {
// init: initializeSDK, // Example: allow re-init or passing config
// openChat: () => { if (chatIconEl) chatIconEl.click(); },
// closeChat: () => { if (closeChatBtn) closeChatBtn.click(); }
// };

})();
