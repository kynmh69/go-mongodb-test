# Frontend - User Management System

This is the Next.js frontend for the Go MongoDB user management API.

## Features

- **User Management**: Create, read, update, and delete users
- **Search Functionality**: Search users by User ID or email
- **Modern UI**: Built with shadcn/ui components and Tailwind CSS
- **Responsive Design**: Works on desktop and mobile devices
- **Type Safety**: Full TypeScript implementation

## Tech Stack

- **Framework**: Next.js 15.1.3 with App Router
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **UI Components**: shadcn/ui (Radix UI primitives)
- **Icons**: Lucide React
- **Date Handling**: date-fns

## Getting Started

### Prerequisites

- Node.js 18+ and npm
- Go backend server running on port 8080

### Installation

1. Navigate to the frontend directory:
```bash
cd frontend
```

2. Install dependencies:
```bash
npm install
```

3. Start the development server:
```bash
npm run dev
```

4. Open [http://localhost:3000](http://localhost:3000) in your browser

### Build for Production

```bash
npm run build
npm start
```

## API Integration

The frontend connects to the Go backend API through:
- **Base URL**: `/api/v1` (proxied to `http://localhost:8080/api/v1`)
- **Endpoints**: All user management endpoints from the Go API

### API Proxy Configuration

The Next.js config includes a proxy setup that forwards `/api/*` requests to the Go backend:

```javascript
// next.config.js
async rewrites() {
  return [
    {
      source: '/api/:path*',
      destination: 'http://localhost:8080/api/:path*',
    },
  ]
}
```

## Project Structure

```
frontend/
├── src/
│   ├── app/                 # Next.js app directory
│   │   ├── globals.css      # Global styles
│   │   ├── layout.tsx       # Root layout
│   │   └── page.tsx         # Home page
│   ├── components/          # React components
│   │   ├── ui/              # shadcn/ui components
│   │   ├── CreateUserDialog.tsx
│   │   ├── UpdateUserDialog.tsx
│   │   ├── DeleteUserDialog.tsx
│   │   ├── SearchUsers.tsx
│   │   └── UserTable.tsx
│   └── lib/                 # Utilities and API
│       ├── api.ts           # API client
│       ├── types.ts         # TypeScript types
│       └── utils.ts         # Utility functions
├── package.json
├── tailwind.config.js
├── tsconfig.json
└── README.md
```

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm start` - Start production server
- `npm run lint` - Run ESLint

## Components

### UserTable
Displays users in a table format with edit and delete actions.

### CreateUserDialog
Modal dialog for creating new users with form validation.

### UpdateUserDialog
Modal dialog for editing existing users.

### DeleteUserDialog
Confirmation dialog for user deletion.

### SearchUsers
Component for searching users by User ID or email.

## API Error Handling

The frontend includes comprehensive error handling:
- Network errors
- HTTP status errors
- Validation errors
- User-friendly error messages

## Responsive Design

The application is fully responsive and works well on:
- Desktop computers
- Tablets
- Mobile phones