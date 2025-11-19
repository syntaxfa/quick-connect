import React, { useState, useEffect, useRef } from 'react';
import { Send, X, Smile, User, MessageCircle, CheckCheck, Check } from 'lucide-react';

const QuickConnectSDK = () => {
    // State Management
    const [isOpen, setIsOpen] = useState(false);
    const [messages, setMessages] = useState([]);
    const [inputText, setInputText] = useState('');
    const [isConnected, setIsConnected] = useState(false);
    const [userState, setUserState] = useState(null); // 'guest' | 'client'
    const [token, setToken] = useState(null);
    const [conversationId] = useState('01KAAWW5V95WFFDERJZCTDE2QQ'); // Demo
    const [showEmojiPicker, setShowEmojiPicker] = useState(false);
    const [typingUsers, setTypingUsers] = useState(new Set());
    const [showProfile, setShowProfile] = useState(false);
    const [profileData, setProfileData] = useState({
        fullname: '',
        email: '',
        phone_number: ''
    });

    const wsRef = useRef(null);
    const messagesEndRef = useRef(null);
    const lastTypingSentRef = useRef(0);
    const typingTimeoutRef = useRef({});
    const reconnectTimeoutRef = useRef(null);
    const reconnectAttempts = useRef(0);

    // Emoji list
    const emojis = ['😊', '😂', '❤️', '👍', '🙏', '🎉', '🔥', '✨', '💯', '🚀', '💪', '😍', '🤔', '👌', '🎯', '⭐'];

    // Storage Management
    const storage = {
        get: (key) => localStorage.getItem(key),
        set: (key, value) => localStorage.setItem(key, value),
        remove: (key) => localStorage.removeItem(key)
    };

    // Initialize SDK
    useEffect(() => {
        initializeSDK();
        return () => {
            if (wsRef.current) {
                wsRef.current.close();
            }
        };
    }, []);

    const initializeSDK = async () => {
        const storedToken = storage.get('QC_TOKEN');
        const storedUserState = storage.get('QC_USER_STATE');

        if (storedToken && storedUserState) {
            setToken(storedToken);
            setUserState(storedUserState);
            connectWebSocket(storedToken);
        } else {
            await registerGuest();
        }
    };

    // Guest Registration
    const registerGuest = async () => {
        try {
            // Simulated API call
            const mockToken = 'mock_guest_token_' + Date.now();

            storage.set('QC_TOKEN', mockToken);
            storage.set('QC_USER_STATE', 'guest');

            setToken(mockToken);
            setUserState('guest');

            connectWebSocket(mockToken);
            setShowProfile(true); // Show profile completion for guests
        } catch (error) {
            console.error('Guest registration failed:', error);
        }
    };

    // Update Guest Profile
    const updateGuestProfile = async () => {
        if (userState !== 'guest') return;

        try {
            // Simulated API call
            console.log('Updating guest profile:', profileData);
            setShowProfile(false);
            addSystemMessage('پروفایل شما با موفقیت ثبت شد ✅');
        } catch (error) {
            console.error('Profile update failed:', error);
        }
    };

    // WebSocket Connection
    const connectWebSocket = (authToken) => {
        if (wsRef.current?.readyState === WebSocket.OPEN) return;

        try {
            // In production: ws://CHAT_URL/v1/chats/clients
            // For demo, we'll simulate WebSocket behavior
            simulateWebSocketConnection(authToken);
        } catch (error) {
            console.error('WebSocket connection failed:', error);
            scheduleReconnect();
        }
    };

    const simulateWebSocketConnection = (authToken) => {
        setIsConnected(true);
        addSystemMessage('به چت پشتیبانی خوش آمدید! 👋');
        reconnectAttempts.current = 0;

        // Simulate receiving a message after 2 seconds
        setTimeout(() => {
            receiveMessage({
                type: 'text',
                timestamp: new Date().toISOString(),
                payload: {
                    id: 'msg_' + Date.now(),
                    conversation_id: conversationId,
                    sender_id: 'support_agent',
                    content: 'سلام! چطور میتونم کمکتون کنم؟ 😊',
                    created_at: new Date().toISOString(),
                    read_at: null,
                    isSupport: true
                }
            });
        }, 2000);
    };

    const scheduleReconnect = () => {
        const delays = [1000, 2000, 5000, 10000];
        const delay = delays[Math.min(reconnectAttempts.current, delays.length - 1)];

        reconnectTimeoutRef.current = setTimeout(() => {
            reconnectAttempts.current++;
            connectWebSocket(token);
        }, delay);
    };

    // Send Message
    const sendMessage = () => {
        if (!inputText.trim() || !isConnected) return;

        const message = {
            type: 'text',
            conversation_id: conversationId,
            content: inputText.trim()
        };

        // Add to local state immediately
        const localMessage = {
            id: 'local_' + Date.now(),
            content: inputText.trim(),
            created_at: new Date().toISOString(),
            isOwn: true,
            status: 'sent'
        };

        setMessages(prev => [...prev, localMessage]);
        setInputText('');
        sendTypingStopped();

        // In production: wsRef.current.send(JSON.stringify(message));

        // Simulate echo response
        setTimeout(() => {
            receiveMessage({
                type: 'text',
                timestamp: new Date().toISOString(),
                payload: {
                    id: 'msg_' + Date.now(),
                    conversation_id: conversationId,
                    sender_id: 'support_agent',
                    content: 'پیام شما دریافت شد. یک لحظه صبر کنید...',
                    created_at: new Date().toISOString(),
                    isSupport: true
                }
            });
        }, 1000);
    };

    // Receive Message
    const receiveMessage = (data) => {
        if (data.type === 'text') {
            setMessages(prev => [...prev, {
                ...data.payload,
                isSupport: data.payload.isSupport || false,
                isOwn: false
            }]);
        } else if (data.type === 'system') {
            handleSystemMessage(data);
        }
    };

    // System Messages
    const handleSystemMessage = (data) => {
        const { sub_type, payload } = data;

        if (sub_type === 'typing_started') {
            setTypingUsers(prev => new Set(prev).add(payload.sender_id));

            if (typingTimeoutRef.current[payload.sender_id]) {
                clearTimeout(typingTimeoutRef.current[payload.sender_id]);
            }

            typingTimeoutRef.current[payload.sender_id] = setTimeout(() => {
                setTypingUsers(prev => {
                    const newSet = new Set(prev);
                    newSet.delete(payload.sender_id);
                    return newSet;
                });
            }, 6000);
        } else if (sub_type === 'typing_stopped') {
            setTypingUsers(prev => {
                const newSet = new Set(prev);
                newSet.delete(payload.sender_id);
                return newSet;
            });

            if (typingTimeoutRef.current[payload.sender_id]) {
                clearTimeout(typingTimeoutRef.current[payload.sender_id]);
            }
        }
    };

    // Typing Indicators
    const handleTyping = () => {
        const now = Date.now();
        if (now - lastTypingSentRef.current > 5000) {
            sendTypingStarted();
            lastTypingSentRef.current = now;
        }
    };

    const sendTypingStarted = () => {
        const message = {
            type: 'system',
            sub_type: 'typing_started',
            conversation_id: conversationId
        };
        // In production: wsRef.current.send(JSON.stringify(message));
    };

    const sendTypingStopped = () => {
        const message = {
            type: 'system',
            sub_type: 'typing_stopped',
            conversation_id: conversationId
        };
        // In production: wsRef.current.send(JSON.stringify(message));
    };

    const addSystemMessage = (text) => {
        setMessages(prev => [...prev, {
            id: 'system_' + Date.now(),
            content: text,
            created_at: new Date().toISOString(),
            isSystem: true
        }]);
    };

    // Scroll to bottom
    useEffect(() => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    }, [messages]);

    // Format time
    const formatTime = (timestamp) => {
        const date = new Date(timestamp);
        return date.toLocaleTimeString('fa-IR', { hour: '2-digit', minute: '2-digit' });
    };

    return (
        <div className="fixed bottom-6 right-6 z-50">
            {/* Chat Button */}
            {!isOpen && (
                <button
                    onClick={() => setIsOpen(true)}
                    className="group relative w-16 h-16 bg-gradient-to-br from-blue-500 via-purple-500 to-pink-500 rounded-full shadow-2xl hover:scale-110 transition-all duration-300 flex items-center justify-center animate-pulse"
                >
                    <MessageCircle className="w-8 h-8 text-white" />
                    <span className="absolute -top-1 -right-1 w-5 h-5 bg-red-500 rounded-full flex items-center justify-center text-xs text-white font-bold animate-bounce">
            3
          </span>
                </button>
            )}

            {/* Chat Window */}
            {isOpen && (
                <div className="w-96 h-[600px] bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900 rounded-2xl shadow-2xl flex flex-col overflow-hidden border border-purple-500/30 animate-in slide-in-from-bottom-5">
                    {/* Header */}
                    <div className="bg-gradient-to-r from-purple-600 to-blue-600 p-4 flex items-center justify-between">
                        <div className="flex items-center gap-3">
                            <div className="relative">
                                <div className="w-10 h-10 bg-white rounded-full flex items-center justify-center">
                                    <MessageCircle className="w-6 h-6 text-purple-600" />
                                </div>
                                {isConnected && (
                                    <span className="absolute bottom-0 right-0 w-3 h-3 bg-green-400 border-2 border-white rounded-full"></span>
                                )}
                            </div>
                            <div>
                                <h3 className="text-white font-bold text-lg">پشتیبانی آنلاین</h3>
                                <p className="text-purple-200 text-xs">
                                    {isConnected ? 'آنلاین' : 'در حال اتصال...'}
                                </p>
                            </div>
                        </div>
                        <button
                            onClick={() => setIsOpen(false)}
                            className="text-white hover:bg-white/20 rounded-lg p-2 transition-colors"
                        >
                            <X className="w-5 h-5" />
                        </button>
                    </div>

                    {/* Guest Profile Banner */}
                    {userState === 'guest' && showProfile && (
                        <div className="bg-gradient-to-r from-amber-500 to-orange-500 p-3">
                            <div className="flex items-center justify-between mb-2">
                                <p className="text-white text-sm font-medium">
                                    اطلاعات خود را تکمیل کنید 📝
                                </p>
                                <button
                                    onClick={() => setShowProfile(false)}
                                    className="text-white/80 hover:text-white"
                                >
                                    <X className="w-4 h-4" />
                                </button>
                            </div>
                            <div className="space-y-2">
                                <input
                                    type="text"
                                    placeholder="نام و نام خانوادگی"
                                    value={profileData.fullname}
                                    onChange={(e) => setProfileData(prev => ({ ...prev, fullname: e.target.value }))}
                                    className="w-full px-3 py-2 rounded-lg text-sm bg-white/90 border-0 focus:ring-2 focus:ring-white"
                                />
                                <div className="flex gap-2">
                                    <input
                                        type="email"
                                        placeholder="ایمیل"
                                        value={profileData.email}
                                        onChange={(e) => setProfileData(prev => ({ ...prev, email: e.target.value }))}
                                        className="flex-1 px-3 py-2 rounded-lg text-sm bg-white/90 border-0 focus:ring-2 focus:ring-white"
                                    />
                                    <button
                                        onClick={updateGuestProfile}
                                        className="px-4 py-2 bg-white text-orange-600 rounded-lg text-sm font-medium hover:bg-white/90 transition-colors"
                                    >
                                        ثبت
                                    </button>
                                </div>
                            </div>
                        </div>
                    )}

                    {/* Messages Area */}
                    <div className="flex-1 overflow-y-auto p-4 space-y-3 bg-slate-950/50">
                        {messages.map((msg) => (
                            <div key={msg.id}>
                                {msg.isSystem ? (
                                    <div className="flex justify-center">
                    <span className="bg-slate-800/50 text-slate-300 text-xs px-3 py-1 rounded-full">
                      {msg.content}
                    </span>
                                    </div>
                                ) : (
                                    <div className={`flex ${msg.isOwn ? 'justify-end' : 'justify-start'}`}>
                                        <div className={`max-w-[75%] ${msg.isOwn ? 'order-2' : 'order-1'}`}>
                                            <div
                                                className={`rounded-2xl px-4 py-3 shadow-lg ${
                                                    msg.isOwn
                                                        ? 'bg-gradient-to-br from-purple-600 to-blue-600 text-white rounded-br-sm'
                                                        : 'bg-slate-800 text-white rounded-bl-sm'
                                                }`}
                                            >
                                                <p className="text-sm leading-relaxed">{msg.content}</p>
                                                <div className="flex items-center justify-end gap-1 mt-1">
                          <span className="text-xs opacity-70">
                            {formatTime(msg.created_at)}
                          </span>
                                                    {msg.isOwn && (
                                                        msg.status === 'read' ? (
                                                            <CheckCheck className="w-4 h-4 text-blue-300" />
                                                        ) : (
                                                            <Check className="w-4 h-4 opacity-70" />
                                                        )
                                                    )}
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                )}
                            </div>
                        ))}

                        {/* Typing Indicator */}
                        {typingUsers.size > 0 && (
                            <div className="flex justify-start">
                                <div className="bg-slate-800 rounded-2xl px-4 py-3 rounded-bl-sm">
                                    <div className="flex gap-1">
                                        <span className="w-2 h-2 bg-slate-500 rounded-full animate-bounce"></span>
                                        <span className="w-2 h-2 bg-slate-500 rounded-full animate-bounce" style={{ animationDelay: '0.1s' }}></span>
                                        <span className="w-2 h-2 bg-slate-500 rounded-full animate-bounce" style={{ animationDelay: '0.2s' }}></span>
                                    </div>
                                </div>
                            </div>
                        )}

                        <div ref={messagesEndRef} />
                    </div>

                    {/* Emoji Picker */}
                    {showEmojiPicker && (
                        <div className="bg-slate-800 border-t border-slate-700 p-3">
                            <div className="grid grid-cols-8 gap-2">
                                {emojis.map((emoji, idx) => (
                                    <button
                                        key={idx}
                                        onClick={() => {
                                            setInputText(prev => prev + emoji);
                                            setShowEmojiPicker(false);
                                        }}
                                        className="text-2xl hover:bg-slate-700 rounded-lg p-1 transition-colors"
                                    >
                                        {emoji}
                                    </button>
                                ))}
                            </div>
                        </div>
                    )}

                    {/* Input Area */}
                    <div className="bg-slate-900 border-t border-slate-700 p-4">
                        <div className="flex items-end gap-2">
                            <button
                                onClick={() => setShowEmojiPicker(!showEmojiPicker)}
                                className="text-slate-400 hover:text-purple-400 transition-colors p-2 hover:bg-slate-800 rounded-lg"
                            >
                                <Smile className="w-5 h-5" />
                            </button>

                            <div className="flex-1 bg-slate-800 rounded-2xl px-4 py-2 border border-slate-700 focus-within:border-purple-500 transition-colors">
                <textarea
                    value={inputText}
                    onChange={(e) => {
                        setInputText(e.target.value);
                        handleTyping();
                    }}
                    onKeyDown={(e) => {
                        if (e.key === 'Enter' && !e.shiftKey) {
                            e.preventDefault();
                            sendMessage();
                        }
                    }}
                    placeholder="پیام خود را بنویسید..."
                    className="w-full bg-transparent text-white text-sm resize-none outline-none max-h-24"
                    rows="1"
                    style={{ direction: 'rtl' }}
                />
                            </div>

                            <button
                                onClick={sendMessage}
                                disabled={!inputText.trim() || !isConnected}
                                className="bg-gradient-to-br from-purple-600 to-blue-600 text-white p-3 rounded-xl hover:scale-105 active:scale-95 transition-all disabled:opacity-50 disabled:cursor-not-allowed shadow-lg hover:shadow-purple-500/50"
                            >
                                <Send className="w-5 h-5" />
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default QuickConnectSDK;