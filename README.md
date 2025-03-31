# HTML to Go Converter

A web-based tool to convert HTML/Vuetify/VuetifyX components to Go code.

[![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https%3A%2F%2Fgithub.com%2Fdanni-cool%2Fhtmlgo-convert-webui&project-name=htmlgo-converter&description=HTML%20to%20Go%20code%20converter%20with%20Vuetify%20support)

## Features

- Convert HTML to Go code
- Support for Vuetify and VuetifyX components
- Customizable prefixes
- Clean and responsive UI

## Local Development

1. Clone this repository
2. Install dependencies:
   ```
   go mod tidy
   npm install
   ```
3. Start the server:
   ```
   go run main.go -port=8080
   ```
4. Open your browser and navigate to `http://localhost:8080`

## Vercel Deployment

This project is configured to be deployed on Vercel. You can deploy it with a single click using the "Deploy with Vercel" button at the top of this README.

Alternatively, follow these steps to deploy your own instance:

1. Install Vercel CLI:

   ```
   npm install -g vercel
   ```

2. Login to Vercel:

   ```
   vercel login
   ```

3. Configure your environment:

   - Create a `.env.local` file with your configuration
   - No special configuration is required for basic usage

4. Test locally using Vercel Dev:

   ```
   npm run dev
   ```

5. Deploy to Vercel:
   ```
   npm run deploy
   ```

## Structure

- `/api`: Serverless functions for Vercel deployment
- `/static`: Frontend files (HTML, CSS, JavaScript)
- `/test`: Test suites

## License

MIT
