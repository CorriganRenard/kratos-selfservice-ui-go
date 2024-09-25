# kratos-selfservice-ui-go
ORY Kratos Self-Service UI written in golang 1.16.

A self service UI for [Kratos](https://www.ory.sh/kratos) based on the NodeJS version but written in go 1.16.

The application provides the following self service UI pages:

- Registration
- Login
- Logout
- Email Verification
- Recovery
- User settings
  - Update profile
  - Change password

A separate app will be respnosible for handling logged in users

- Dashboard

# Tailwind CSS

To create the initial css file:

```
task gen_css
```

Static assets served via [HashFS](https://github.com/benbjohnson/hashfs) that appends hashes to embedded static assets for aggresive HTTP caching.
