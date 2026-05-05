# ============================
# STAGE 1: Build the .NET App
# ============================
# Use the .NET 9 SDK image to build the app
FROM mcr.microsoft.com/dotnet/sdk:9.0 AS build
WORKDIR /src

# Copy the project file and restore dependencies
COPY ["note-nexus-client.csproj", "./"]
RUN dotnet restore "note-nexus-client.csproj"

# Copy the rest of the code and publish (compile) it in Release mode
COPY . .
RUN dotnet publish "note-nexus-client.csproj" -c Release -o /app/publish

# ============================
# STAGE 2: Serve with Nginx
# ============================
# Use the lightweight Nginx image
FROM nginx:alpine
WORKDIR /usr/share/nginx/html

# Remove the default Nginx welcome page files
RUN rm -rf ./*

# Copy the compiled static files from the build stage into Nginx's folder
COPY --from=build /app/publish/wwwroot .

# Copy our custom Nginx configuration file
COPY nginx.conf /etc/nginx/nginx.conf

# Expose port 80 (standard web port)
EXPOSE 80