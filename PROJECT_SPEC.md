# WebSocket Chat Server - Product Specification

## Overview
A real-time chat server supporting multiple clients, chat rooms, and basic user management. Terminal-based chat application similar to a simplified Slack.

## Core Features

### 1. Connection Management
- Clients connect via WebSocket
- Server tracks all active connections
- Graceful handling of client disconnects and reconnects
- Basic authentication requiring username on connection

### 2. Chat Rooms
- Default "lobby" room - all users join automatically on connect
- Users can create new chat rooms
- Users can join existing rooms
- Users can leave rooms (must remain in at least one room)
- Ability to list all available rooms

### 3. Messaging
- Send messages to current room
- Messages broadcast to all users in that room only
- Message format includes: timestamp, username, room name, message content
- Server-side validation: no empty messages, enforce length limits

### 4. User Commands
- `/rooms` - List all available chat rooms
- `/join <room>` - Join a specific room
- `/leave` - Leave current room
- `/users` - Display all users in current room
- `/help` - Show available commands

### 5. Server Features
- Log all client connections and disconnections
- Log all messages for debugging purposes
- Display active connection count
- Configurable server port (via command-line flag or environment variable)

## Out of Scope
The following features are explicitly NOT included in this version:
- Message history or persistence
- Private/direct messages between users
- User roles or permission systems
- Message editing or deletion
- File uploads or media sharing
- Message encryption (plain WebSocket only)

## Client Interface

### Terminal Client (Primary)
- Command-line interface for connecting to server
- Type messages directly - they send on enter
- Incoming messages display as they arrive
- Commands work like regular messages (type `/join general`)
- Simple scrolling text output

### TUI Client (Optional Stretch Goal)
- Split-screen interface using Bubbletea
- Message display area (top section)
- Input area (bottom section)
- Current room shown in header
- User list displayed in sidebar

**Recommendation**: Build the simple terminal client first. Add TUI only after core functionality is solid.

## Success Criteria
- Support 3+ simultaneous client connections
- Messages in one room do not appear in other rooms
- Server handles client disconnects without crashing
- Room creation and joining functionality works reliably
- Can demonstrate working system without major bugs

## Stretch Goals
If time permits, consider adding:
- Rate limiting per user
- Automatic timeout for inactive users
- Room capacity limits
- Enhanced error messages for clients
- Reconnection logic with session recovery

## Technical Constraints
- Use WebSocket protocol (gorilla/websocket or golang.org/x/net/websocket recommended)
- Server must handle concurrent connections safely
- All room/user state managed in-memory (no database required)
