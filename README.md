# HTML2Go Converter WebUI

A web interface for the HTML2Go converter, which converts HTML to Go code.

## Getting Started

You can run the server locally using the provided script:

```bash
./start.sh
```

This will:

1. Create a static directory from the public directory if it doesn't exist
2. Start the server on an auto-selected port

If you want to specify a port, you can use:

```bash
./start.sh -port=8080
```

## Structure

- `api/` - Contains the API handlers for conversion
  - `convert.go` - The HTML to Go conversion API
  - `index.go` - Serves the main index.html file
- `cmd/server/` - Contains the main entry point for the server
  - `main.go` - Server implementation
- `public/` - Original static files
- `static/` - Copied static files served by the server

## API

The primary API endpoint is `/convert`, which accepts POST requests with JSON payloads.

Example request:

```json
{
  "html": "<div class=\"container\"><h1>Hello World</h1></div>",
  "packagePrefix": "h",
  "vuetifyPrefix": "v",
  "vuetifyXPrefix": "vx",
  "direction": "html2go",
  "childrenMode": false
}
```

Example response:

```json
{
  "code": "h.Div(\n\tClass(\"container\"),\n\th.H1(\"Hello World\"),\n)"
}
```

## Features

- Convert HTML to Go code
- Support for Vuetify and VuetifyX components
- Customizable prefixes
- Clean and responsive UI

## Vercel Deployment

This project is set up to be deployed on Vercel with Golang serverless functions.

### Local Development

1. Clone the repository
2. Set up the environment:
   ```
   # Create a local .env file
   cp .env.example .env
   ```
3. Run the development server:

   ```
   # Using Go directly
   go run main.go

   # Or using Vercel CLI
   pnpm i -g vercel
   vercel dev
   ```

### Production Deployment

1. Push your changes to GitHub
2. Connect your GitHub repository to Vercel
3. Configure the Environment Variables in Vercel:
   - `TLS`: Set to `false`
   - `CONFIG`: Set to your desired JSON configuration. Example:
     ```json
     {
       "APP_NAME": "HTML2GoConverter",
       "APP_ENV": "production",
       "APP_URL": "https://your-vercel-domain.vercel.app",
       "APP_PORT": 8080,
       "APP_PPROF": false,
       "HTTPS": 0,
       "ADDRESS_LIMIT": true
     }
     ```
4. Deploy your application

## License

MIT
