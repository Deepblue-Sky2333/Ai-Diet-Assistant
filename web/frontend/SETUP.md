# AI Diet Assistant Frontend Setup

This is the frontend application for the AI Diet Assistant. It connects to the backend API to provide a user-friendly interface for managing your diet and nutrition.

## Quick Start

### Option 1: Demo Mode (No Backend Required) ⭐ Recommended for Testing

The easiest way to explore the UI without setting up a backend:

1. The app is pre-configured for demo mode
2. Click "Login with Test Account" on the login page
3. Explore all features with sample data

**Perfect for:**
- Testing the UI in v0 preview
- Frontend development without backend dependency
- Quick demos and prototyping

### Option 2: Connect to Real Backend

1. Create a `.env.local` file in the root directory:

\`\`\`bash
cp .env.local.example .env.local
\`\`\`

2. Update the `NEXT_PUBLIC_API_URL` variable to point to your backend API:

\`\`\`bash
# For local development (backend running on port 9090)
NEXT_PUBLIC_API_URL=http://localhost:9090/api/v1

# For production deployment
NEXT_PUBLIC_API_URL=https://your-backend-domain.com/api/v1
\`\`\`

## Environment Configuration

### Demo Mode (Default)

By default, the app runs in demo mode for easy testing:

\`\`\`bash
NEXT_PUBLIC_DEMO_MODE=true
\`\`\`

This allows you to use the app without a backend connection. All features work with sample data.

### Production Mode

1. Create a `.env.local` file in the root directory:

\`\`\`bash
cp .env.local.example .env.local
\`\`\`

2. Update the environment variables in `.env.local`:

\`\`\`bash
# Disable demo mode to connect to real backend
NEXT_PUBLIC_DEMO_MODE=false

# Backend API URL
NEXT_PUBLIC_API_URL=http://localhost:9090/api/v1
\`\`\`

**For v0 Users:** Add environment variables in the **Vars section** (left sidebar) instead of creating `.env.local`.

## Installation

\`\`\`bash
# Install dependencies
npm install

# Run development server
npm run dev
\`\`\`

Open [http://localhost:3000](http://localhost:3000) in your browser.

## Test Account

For quick testing, use the "Login with Test Account" button on the login page.

- Username: `test`
- Password: `114514`

**Note:** In demo mode, any username/password will work. In production mode, you need valid credentials from your backend.

## Troubleshooting

### "Cannot connect to backend API" Error

This error occurs when the frontend cannot reach the backend server. Check:

1. **Backend is running**: Make sure your backend server is started and listening on the correct port
2. **Correct URL**: Verify `NEXT_PUBLIC_API_URL` in `.env.local` matches your backend URL
3. **CORS enabled**: The backend must allow requests from the frontend origin
4. **Firewall**: Check that no firewall is blocking the connection

**Quick Fix:** Enable demo mode by setting `NEXT_PUBLIC_DEMO_MODE=true` to bypass backend connection.

### Debug Mode

Open browser DevTools (F12) → Console tab to see detailed API logs prefixed with `[v0]`.

### Using in v0 Preview

When working in v0:

1. **Add environment variables** in the **Vars section** (left sidebar):
   - `NEXT_PUBLIC_DEMO_MODE` - Set to `true` for demo mode, `false` to connect to backend
   - `NEXT_PUBLIC_API_URL` - Your backend URL (only needed when demo mode is off)

2. **Refresh the preview** after changing environment variables

3. **Deploy to Vercel** - Click the "Publish" button to deploy your app

## Features

- Dashboard with nutrition tracking
- Food inventory management (Supermarket)
- Daily meal logging
- AI-powered diet plan generation
- Nutrition analysis and reports
- AI chat for dietary advice
- User settings and preferences

## Tech Stack

- Next.js 16 with App Router
- React 19
- TypeScript
- Tailwind CSS v4
- shadcn/ui components
