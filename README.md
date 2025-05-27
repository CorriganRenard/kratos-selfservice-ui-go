# kratos-selfservice-ui-go
ORY Kratos Self-Service UI written in golang 1.16.

A self service UI for [Kratos](https://www.ory.sh/kratos) based on the NodeJS version but written in go 1.16.

## Features

The application provides the following self service UI pages:

- Registration
- Login
- Logout
- Email Verification
- Recovery
- User settings
  - Update profile
  - Change password
- Dashboard (for logged in users)

## Configuration

### Required Environment Variables

The following environment variables are required for the application to function:

- `KRATOS_PUBLIC_URL`: The URL where ORY Kratos's Public API is located
- `KRATOS_ADMIN_URL`: The URL where ORY Kratos's Admin API is located  
- `KRATOS_BROWSER_URL`: The URL to build all of the kratos self service URLs
- `BASE_URL`: The base URL of this app (must be absolute)
- `COOKIE_STORE_KEY_PAIRS`: Pairs of authentication and encryption keys (generate with `--gen-cookie-store-key-pair`)

### Site Customization

You can customize the appearance of your identity portal using these environment variables:

#### Site Name
- **Environment Variable**: `SITE_NAME`
- **Command Line Flag**: `--site-name`
- **Default**: `"Kratos Selfservice UI"`
- **Description**: The name that appears in the browser tab title across all pages

#### Favicon
- **Environment Variable**: `FAVICON_URL` 
- **Command Line Flag**: `--favicon-url`
- **Default**: `"/static/images/favicon.svg"`
- **Description**: The URL to the favicon image. Can be a local path or external URL

### Configuration Examples

#### Docker Compose

```yaml
version: '3.8'

services:
  kratos-selfservice-ui:
    build: .
    ports:
      - "4455:4455"
    environment:
      # Required Kratos configuration
      KRATOS_PUBLIC_URL: "http://kratos:4433"
      KRATOS_ADMIN_URL: "http://kratos:4434"
      KRATOS_BROWSER_URL: "http://localhost:4433"
      BASE_URL: "http://localhost:4455"
      COOKIE_STORE_KEY_PAIRS: "your-auth-key your-encryption-key"
      
      # Site customization
      SITE_NAME: "My Company Identity Portal"
      FAVICON_URL: "https://mycompany.com/favicon.ico"
      
      # Optional
      PORT: "4455"
```

#### Command Line

```bash
# Using command line flags
./kratos-selfservice-ui-go \
  --kratos-public-url "http://localhost:4433" \
  --kratos-admin-url "http://localhost:4434" \
  --kratos-browser-url "http://localhost:4433" \
  --base-url "http://localhost:4455" \
  --site-name "My Custom Portal" \
  --favicon-url "/custom/favicon.ico" \
  --cookie-store-key-pairs "your-auth-key your-encryption-key"
```

#### Environment Variables Only

```bash
export KRATOS_PUBLIC_URL="http://localhost:4433"
export KRATOS_ADMIN_URL="http://localhost:4434"
export KRATOS_BROWSER_URL="http://localhost:4433"
export BASE_URL="http://localhost:4455"
export COOKIE_STORE_KEY_PAIRS="your-auth-key your-encryption-key"
export SITE_NAME="My Custom Portal"
export FAVICON_URL="/custom/favicon.ico"

./kratos-selfservice-ui-go
```

### Generating Cookie Store Keys

To generate secure cookie store key pairs:

```bash
./kratos-selfservice-ui-go --gen-cookie-store-key-pair
```

This will output base64-encoded keys that you can use for the `COOKIE_STORE_KEY_PAIRS` environment variable.

## Development

### Tailwind CSS

To create the initial CSS file:

```bash
task gen_css
```

### Building

```bash
go build -o kratos-selfservice-ui-go
```

### Static Assets

Static assets are served via [HashFS](https://github.com/benbjohnson/hashfs) that appends hashes to embedded static assets for aggressive HTTP caching.

## Help

To see all available command line options:

```bash
./kratos-selfservice-ui-go --help
```
