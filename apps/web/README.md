# ProvisionHub Web UI

The **Web UI** is the frontend interface for ProvisionHub, built with Next.js and Material-UI. It provides a user-friendly way to interact with the platform's provisioning capabilities.

---

## 🎯 Purpose

The Web UI serves as the primary interface for developers and platform teams to:

- **Browse & Create Systems** - Visual interface for system management
- **Compose Components** - Add and configure components to systems
- **Monitor Provisioning** - Track execution status and view logs
- **Manage Blueprints** - Configure provisioning templates
- **User Authentication** - OIDC-based login flow
- **Dashboard** - Overview of systems, components, and activity

---

## 🏗️ Architecture

### Technology Stack

- **Framework**: Next.js 16+ (App Router)
- **Language**: TypeScript
- **Styling**: Material-UI (MUI)
- **State Management**: React Context / Zustand / TanStack Query
- **Authentication**: NextAuth.js (OIDC)
- **HTTP Client**: Axios / Fetch API
- **UI Components**: Material-UI (MUI)

### Features

- **Server-Side Rendering** - Fast initial page loads
- **API Routes** - Proxy to Control Plane API
- **Real-time Updates** - WebSocket or polling for run status
- **Responsive Design** - Mobile-friendly interface
- **Dark Mode** - MUI theme support

---

## 🚀 Getting Started

### Prerequisites

- Node.js 18+ or Node.js 20+
- npm, yarn, or pnpm
- Control Plane API running

### Environment Variables

Create a `.env.local` file:

```bash
# Control Plane API
NEXT_PUBLIC_API_URL=http://localhost:8080

# Authentication (NextAuth.js)
NEXTAUTH_URL=http://localhost:3000
NEXTAUTH_SECRET=your-secret-key-here

# OIDC (Keycloak)
OIDC_ISSUER=http://localhost:8081/realms/provisionhub
OIDC_CLIENT_ID=provisionhub-web
OIDC_CLIENT_SECRET=your-keycloak-client-secret

# Optional
NEXT_PUBLIC_APP_NAME=ProvisionHub
NEXT_PUBLIC_APP_VERSION=0.1.0
```

### Installation

The project supports all three major package managers: **npm**, **yarn**, and **pnpm**.

```bash
cd apps/web

# Install dependencies with your preferred package manager
npm install
# or
yarn install
# or
pnpm install

# Run development server
npm run dev
# or
yarn dev
# or
pnpm dev
```

**Note**:

- A `.npmrc` file is included for pnpm compatibility with Next.js
- You can use any package manager, but stick to one throughout the project
- Lock files will be generated on first install (package-lock.json, yarn.lock, or pnpm-lock.yaml)

Open [http://localhost:3000](http://localhost:3000) in your browser.

### Build for Production

```bash
# Build
npm run build

# Start production server
npm run start
```

---

## 📦 Package Manager Support

This project supports **npm**, **yarn**, and **pnpm**. Choose the one that best fits your workflow:

### npm (Default)

```bash
npm install
npm run dev
npm run build
```

- ✅ Default Node.js package manager
- ✅ No additional setup required
- 📄 Creates `package-lock.json`

### Yarn

```bash
yarn install
yarn dev
yarn build
```

- ✅ Fast and reliable
- ✅ Works out of the box
- 📄 Creates `yarn.lock`

### pnpm (Recommended)

```bash
pnpm install
pnpm dev
pnpm build
```

- ✅ Fastest and most disk-efficient
- ✅ Pre-configured with `.npmrc` for Next.js compatibility
- ✅ Saves significant disk space with content-addressable storage
- 📄 Creates `pnpm-lock.yaml`

**Best Practices**:

- Choose one package manager and use it consistently throughout the project
- Don't mix package managers (avoid having multiple lock files)
- Commit your lock file to version control
- Team members should use the same package manager

---

## 📁 Project Structure

```
apps/web/
├── app/
│   ├── (auth)/
│   │   ├── login/
│   │   └── callback/
│   ├── (dashboard)/
│   │   ├── systems/
│   │   ├── components/
│   │   ├── runs/
│   │   └── settings/
│   ├── api/
│   │   └── auth/
│   ├── layout.tsx
│   └── page.tsx
├── components/
│   ├── ui/              # Reusable UI components
│   ├── systems/         # System-specific components
│   ├── components/      # Component-specific views
│   └── runs/            # Provisioning run views
├── lib/
│   ├── api.ts          # API client
│   ├── auth.ts         # Auth configuration
│   └── utils.ts        # Utilities
├── theme/
│   └── theme.ts        # MUI theme configuration
├── public/
├── next.config.js
├── tsconfig.json
└── README.md
```

---

## 🎨 Key Features

### Dashboard

- Overview of recent systems and components
- Active provisioning runs
- Quick actions and shortcuts

### System Management

- List all systems with search and filters
- Create new system with blueprint configuration
- View system details and associated components
- Edit and delete systems

### Component Creation

- Select component type (backend, frontend, database, etc.)
- Configure component blueprint
- Preview generated structure
- Submit for provisioning

### Provisioning Tracking

- Real-time status updates
- Execution logs streaming
- Retry failed runs
- Download artifacts

### User Settings

- Profile management
- API token generation
- Preferences and theme

---

## 🔐 Authentication

The Web UI uses **NextAuth.js** with OIDC provider (Keycloak).

Flow:

1. User clicks "Login"
2. Redirects to Keycloak login page
3. After authentication, returns to callback route
4. NextAuth.js creates session
5. Protected routes require active session
6. API calls include JWT token

---

## 🎨 Material-UI (MUI) Setup

### Theme Configuration

The app uses a custom MUI theme defined in `theme/theme.ts`:

```typescript
// theme/theme.ts
import { createTheme } from "@mui/material/styles";

export const theme = createTheme({
   palette: {
      mode: "light",
      primary: {
         main: "#1976d2",
      },
      secondary: {
         main: "#dc004e",
      },
   },
   typography: {
      fontFamily: "Inter, Roboto, sans-serif",
   },
   components: {
      MuiButton: {
         styleOverrides: {
            root: {
               textTransform: "none",
            },
         },
      },
   },
});
```

### Using MUI Components

```tsx
import { Button, Card, CardContent } from "@mui/material";

export default function MyComponent() {
   return (
      <Card>
         <CardContent>
            <Button variant="contained" color="primary">
               Create System
            </Button>
         </CardContent>
      </Card>
   );
}
```

### Dark Mode Support

Toggle between light and dark themes:

```typescript
const [mode, setMode] = useState<"light" | "dark">("light");

const theme = createTheme({
   palette: {
      mode,
   },
});
```

---

## 🧪 Testing

```bash
# Unit tests (Jest + React Testing Library)
npm run test

# E2E tests (Playwright)
npm run test:e2e

# Type checking
npm run type-check

# Linting
npm run lint
```

---

## 🐳 Docker

Build and run the Web UI in a container:

```bash
# Build image
docker build -t provisionhub-web .

# Run container
docker run -p 3000:3000 \
  --env-file .env.local \
  provisionhub-web
```

---

## 📦 Key Dependencies

```json
{
   "dependencies": {
      "next": "^15.0.0",
      "react": "^18.0.0",
      "@mui/material": "^5.0.0",
      "@mui/icons-material": "^5.0.0",
      "@emotion/react": "^11.0.0",
      "@emotion/styled": "^11.0.0",
      "next-auth": "^4.0.0",
      "axios": "^1.0.0",
      "@tanstack/react-query": "^5.0.0"
   }
}
```

---

## 🎯 MUI Best Practices

### Use the sx Prop

```tsx
<Box sx={{ p: 2, bgcolor: "primary.main" }}>Content</Box>
```

### Theme-aware Styling

```tsx
import { styled } from "@mui/material/styles";

const StyledCard = styled(Card)(({ theme }) => ({
   padding: theme.spacing(2),
   backgroundColor: theme.palette.background.paper,
}));
```

### Responsive Design

```tsx
<Box
   sx={{
      width: { xs: "100%", sm: "50%", md: "33%" },
      p: { xs: 1, sm: 2, md: 3 },
   }}
>
   Responsive content
</Box>
```

---

## 📊 Performance

- **Lighthouse Score Target**: 90+
- **First Contentful Paint**: < 1.5s
- **Time to Interactive**: < 3s
- **Bundle Size**: Monitored and optimized

### MUI Optimization

- Use tree-shaking imports
- Lazy load heavy components
- Optimize icon imports:

```tsx
// ✅ Good
import AddIcon from "@mui/icons-material/Add";

// ❌ Avoid
import { Add } from "@mui/icons-material";
```

---

## 🤝 Contributing

See the main [CONTRIBUTING.md](../../CONTRIBUTING.md) for contribution guidelines.

---

## 📜 License

Apache License 2.0 - See [LICENSE](../../LICENSE)

---

## 🔗 Related

- [Main Project Documentation](../../README.md)
- [Control Plane Documentation](../control-plane/README.md)
- [Worker Documentation](../worker/README.md)
