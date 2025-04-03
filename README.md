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

## Deploying to Vercel

This project has been configured to deploy on Vercel with Go serverless functions for the backend and static files for the frontend.

### Configuration Details

1. The project uses:

   - **Go Serverless Functions**: Located in `api/` directory
   - **Static Frontend**: Located in `public/` directory

2. The `vercel.json` configuration includes:
   ```json
   {
     "version": 2,
     "builds": [
       {
         "src": "api/**/*.go",
         "use": "@vercel/go"
       },
       {
         "src": "public/**",
         "use": "@vercel/static"
       }
     ],
     "routes": [
       {
         "src": "/convert",
         "dest": "/api/convert.go"
       },
       {
         "src": "/",
         "dest": "/api/index.go"
       },
       {
         "src": "/script.js",
         "dest": "/public/script.js"
       },
       {
         "src": "/demo.html",
         "dest": "/public/demo.html"
       },
       {
         "src": "/(.*)",
         "dest": "/public/$1"
       }
     ]
   }
   ```

### Deployment Steps

1. Install the Vercel CLI:

   ```
   npm install -g vercel
   ```

2. Login to Vercel:

   ```
   vercel login
   ```

3. Deploy the project:
   ```
   vercel
   ```
4. Deploy to production:
   ```
   vercel --prod
   ```

### Troubleshooting

If you encounter issues with the index.html file not being found, make sure the `api/index.go` handler includes the correct paths to search for the file in both `public/` and `static/` directories.

## License

MIT
