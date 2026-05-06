using Microsoft.AspNetCore.Components.Web;
using Microsoft.AspNetCore.Components.WebAssembly.Hosting;
using Microsoft.AspNetCore.Components.Authorization;
using note_nexus_client;
using note_nexus_client.Services;
using note_nexus_client.Models;
using System.Net.Http;
using System;

var builder = WebAssemblyHostBuilder.CreateDefault(args);
builder.RootComponents.Add<App>("#app");
builder.RootComponents.Add<HeadOutlet>("head::after");

string backendUrl = "https://note-nexus-anxm.onrender.com/";

// Auth and core services
builder.Services.AddScoped<LocalStorageService>();
builder.Services.AddScoped<AuthenticationStateProvider, CustomAuthStateProvider>();
builder.Services.AddScoped<CustomHttpHandler>();
builder.Services.AddSingleton<AppState>();

// Base HTTP Client for AuthService (No JWT attached)
builder.Services.AddScoped<AuthService>(sp => 
{
    var client = new HttpClient { BaseAddress = new Uri(backendUrl) };
    return new AuthService(client, sp.GetRequiredService<AuthenticationStateProvider>(), sp.GetRequiredService<LocalStorageService>());
});

// Authenticated HTTP Client for App Services
builder.Services.AddScoped<TodoService>(sp => 
{
    var handler = sp.GetRequiredService<CustomHttpHandler>();
    handler.InnerHandler = new HttpClientHandler();
    var client = new HttpClient(handler) { BaseAddress = new Uri(backendUrl) };
    return new TodoService(client, sp.GetRequiredService<AppState>());
});

builder.Services.AddScoped<NoteService>(sp => 
{
    var handler = sp.GetRequiredService<CustomHttpHandler>();
    handler.InnerHandler = new HttpClientHandler();
    var client = new HttpClient(handler) { BaseAddress = new Uri(backendUrl) };
    return new NoteService(client, sp.GetRequiredService<AppState>());
});

builder.Services.AddScoped<WorkspaceService>(sp => 
{
    var handler = sp.GetRequiredService<CustomHttpHandler>();
    handler.InnerHandler = new HttpClientHandler();
    var client = new HttpClient(handler) { BaseAddress = new Uri(backendUrl) };
    return new WorkspaceService(client);
});

builder.Services.AddScoped<ReminderService>(sp => 
{
    var handler = sp.GetRequiredService<CustomHttpHandler>();
    handler.InnerHandler = new HttpClientHandler();
    var client = new HttpClient(handler) { BaseAddress = new Uri(backendUrl) };
    return new ReminderService(client, sp.GetRequiredService<AppState>());
});

builder.Services.AddScoped<NotificationService>(sp => 
{
    var handler = sp.GetRequiredService<CustomHttpHandler>();
    handler.InnerHandler = new HttpClientHandler();
    var client = new HttpClient(handler) { BaseAddress = new Uri(backendUrl) };
    return new NotificationService(client);
});

// Enable Authorization
builder.Services.AddAuthorizationCore();

await builder.Build().RunAsync();