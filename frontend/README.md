# Byte Payments Frontend

A modern, responsive frontend for the Byte Payments gateway built with Next.js, React, TypeScript, and shadcn/ui components.

## Features

- 🎨 **Modern UI**: Built with shadcn/ui components and Tailwind CSS
- 🌓 **Dark/Light Theme**: Automatic system theme detection with manual toggle
- 📱 **Responsive Design**: Mobile-first responsive design
- ⚡ **Real-time Updates**: Live payment status updates and countdown timer
- 🎯 **TypeScript**: Full type safety throughout the application
- 🔒 **Email Verification**: Secure OTP-based email verification flow
- 💳 **QR Code Payments**: Dynamic QR code generation for crypto payments
- 🏪 **Professional Gateway**: Payment processor-style UI with comprehensive information display

## Tech Stack

- **Framework**: Next.js 14 with App Router
- **UI Library**: shadcn/ui + Radix UI
- **Styling**: Tailwind CSS
- **Icons**: Lucide React
- **State Management**: React hooks
- **Theme**: next-themes
- **Notifications**: Sonner (toast notifications)
- **Language**: TypeScript

## Getting Started

### Prerequisites

- Node.js 18+ 
- Your Go backend server running on port 8080

### Installation

1. Install dependencies:
   ```bash
   npm install
   ```

2. Configure API endpoint:
   Create a `.env.local` file in the frontend directory:
   ```
   NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
   ```

3. Start the development server:
   ```bash
   npm run dev
   ```

4. Open [http://localhost:3000](http://localhost:3000) in your browser

## Project Structure

```
frontend/
├── src/
│   ├── app/                 # Next.js App Router pages
│   │   ├── page.tsx         # Plan selection page (/)
│   │   ├── pay/
│   │   │   └── page.tsx     # Payment page (/pay?id={paymentId})
│   │   ├── layout.tsx       # Root layout with theme provider
│   │   └── globals.css      # Global styles
│   ├── components/
│   │   ├── ui/             # shadcn/ui components
│   │   ├── theme-toggle.tsx # Theme switcher component
│   │   ├── theme-provider.tsx # Theme context provider
│   │   ├── plan-card.tsx   # Plan selection card
│   │   └── ...
│   └── lib/
│       ├── api.ts          # API client and types
│       └── utils.ts        # Utility functions
└── ...
```

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run start` - Start production server
- `npm run lint` - Run ESLint

## API Integration

The frontend communicates with your Go backend through the following endpoints:

### Core Payment APIs
- `GET /api/v1/plans` - Fetch available plans
- `GET /api/v1/currencies` - Fetch supported currencies
- `POST /api/v1/verification/send-code` - Send email verification code
- `POST /api/v1/verification/verify-code` - Verify email code
- `POST /api/v1/payments/create` - Create new payment
- `GET /api/v1/payments/{id}/status` - Get payment status
- `PATCH /api/v1/payments/{id}/cancel` - Cancel payment



## Pages

### Plan Selection (`/`)
- **Clean Professional Design**: Minimal layout focused on essential elements
- **Email Input**: Simple email validation with immediate feedback
- **List-Style Plan Selection**: Standard radio-button style plan selection
- **Email Verification**: Streamlined OTP verification flow
- **Responsive Layout**: Single-column design optimized for all devices
- **Clear Call-to-Action**: Full-width continue button with contextual messaging

### Payment Page (`/pay?id={paymentId}`)
- **Professional Gateway Design**: Multi-section layout with comprehensive payment information
- **Real-time Updates**: Live status monitoring with 10-second refresh intervals
- **Progress Indicators**: Visual countdown timer with color-coded urgency states
- **QR Code Integration**: High-quality QR codes with click-to-copy functionality
- **Security Information**: Trust indicators and security features display
- **Help & Support**: Built-in support links and transaction guides
- **Payment Management**: Cancel payment option with confirmation
- **Responsive Design**: Optimized for both desktop and mobile devices

## Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `NEXT_PUBLIC_API_BASE_URL` | Backend API base URL | `http://localhost:8080` | Yes |

### Optional Email Configuration
Create a `.env.local` file in the frontend directory with these variables if needed:

```env
# API Configuration
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080

# Optional: If your backend requires authentication
NEXT_PUBLIC_API_KEY=your_api_key_here

# Optional: Email service configuration (handled by backend)
NEXT_PUBLIC_EMAIL_NOTIFICATIONS_ENABLED=true
```

## Customization

### Theme Colors
The app uses a neutral color scheme by default. You can customize colors in `src/app/globals.css` by modifying the CSS custom properties.

### API Client
The API client is in `src/lib/api.ts`. You can modify request headers, error handling, or add new endpoints here.

## Deployment

1. Build the application:
   ```bash
   npm run build
   ```

2. Start the production server:
   ```bash
   npm start
   ```

Or deploy to platforms like Vercel, Netlify, or any Node.js hosting service.

## Email Features

The frontend includes comprehensive email functionality:

### Automatic Notifications
- **Status Change Alerts**: Automatically sends emails when payment status changes
- **Real-time Monitoring**: Detects status changes and triggers appropriate notifications
- **Smart Filtering**: Only sends relevant notifications (completed, expired, cancelled)

### Manual Email Actions
- **Resend Confirmation**: Users can request confirmation emails to be resent
- **Payment Receipts**: Generate and send detailed payment receipts
- **Email Preferences**: Full preference management with granular controls

### Email Types Supported
- `payment_created` - Sent when payment is initiated
- `payment_completed` - Sent when payment is successfully completed
- `payment_failed` - Sent when payment is cancelled
- `payment_expired` - Sent when payment time expires

### User Preferences
Users can control:
- Payment update notifications
- Status change alerts  
- Promotional emails
- Complete email opt-out options

## Notes

- The frontend is configured to work with the existing Go backend APIs
- QR codes are loaded from the backend response 
- Payment status is checked every 10 seconds for real-time updates
- Email notifications are sent automatically on status changes
- The app features a professional payment gateway design with comprehensive UX
